package integration

import (
	"context"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestRefreshRotationReuseDetectionAndLogout(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	p := w.parent1
	csrf := func() []string { return []string{"X-CSRF-Token", p.csrf()} }

	oldCookie := ""
	for _, ck := range p.http.Jar.Cookies(mustURL(e.srv.URL + "/api/v1/auth/")) {
		if ck.Name == "ssp_refresh" {
			oldCookie = ck.Value
		}
	}
	if oldCookie == "" {
		t.Fatal("refresh cookie not set")
	}
	// CSRF required
	must(t, p.do("POST", "/api/v1/auth/refresh", nil), 403)
	must(t, p.do("POST", "/api/v1/auth/refresh", nil, "X-CSRF-Token", "wrong"), 403)
	must(t, p.do("POST", "/api/v1/auth/refresh", nil, append(csrf(), "Origin", "https://evil.example")...), 403)
	r := must(t, p.do("POST", "/api/v1/auth/refresh", nil, append(csrf(), "Origin", testOrigin)...), 200)
	if !strings.Contains(r.body, `"token_type":"Bearer"`) || strings.Contains(r.body, "refresh") {
		t.Fatalf("refresh token must not be in the body: %s", r.body)
	}
	// reuse of the rotated cookie: 401 and the whole family is revoked
	attacker := e.client()
	attacker.http.Jar.SetCookies(mustURL(e.srv.URL+"/api/v1/auth/"), p.http.Jar.Cookies(mustURL(e.srv.URL+"/api/v1/auth/")))
	attacker.http.Jar.SetCookies(mustURL(e.srv.URL+"/"), p.http.Jar.Cookies(mustURL(e.srv.URL+"/")))
	replay := e.client()
	replay.http.Jar.SetCookies(mustURL(e.srv.URL+"/"), p.http.Jar.Cookies(mustURL(e.srv.URL+"/")))
	replay.http.Jar.SetCookies(mustURL(e.srv.URL+"/api/v1/auth/"), []*httpCookie{{Name: "ssp_refresh", Value: oldCookie, Path: "/api/v1/auth"}})
	must(t, replay.do("POST", "/api/v1/auth/refresh", nil, "X-CSRF-Token", replay.csrf()), 401)
	must(t, p.do("POST", "/api/v1/auth/refresh", nil, csrf()...), 401) // legitimate new token is revoked too
	if auditCount(t, e, "auth:refresh_reuse", "denied") != 1 {
		t.Fatal("reuse must be audited once")
	}
	// logout revokes
	q := w.parent2
	must(t, q.do("POST", "/api/v1/auth/logout", nil, "X-CSRF-Token", q.csrf()), 204)
	must(t, q.do("POST", "/api/v1/auth/refresh", nil, "X-CSRF-Token", q.csrf()), 403) // csrf cookie cleared
}

func TestLoginUniformResponseAndTiming(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	_ = w
	g := e.client()
	norm := regexp.MustCompile(`"request_id":"[^"]*"`)
	var bodies [2]string
	var durs [2]time.Duration
	const n = 6
	for i := 0; i < n; i++ {
		for j, email := range []string{"nobody@example.test", "p1@example.test"} {
			start := time.Now()
			r := g.do("POST", "/api/v1/auth/login", map[string]string{"email": email, "password": "Wrong-password-1"})
			durs[j] += time.Since(start)
			must(t, r, 401)
			bodies[j] = norm.ReplaceAllString(r.body, "")
		}
	}
	if bodies[0] != bodies[1] {
		t.Fatalf("responses differ: %q vs %q", bodies[0], bodies[1])
	}
	diff := (durs[0] - durs[1]) / n
	if diff < 0 {
		diff = -diff
	}
	t.Logf("mean timing difference unknown-email vs wrong-password: %v", diff)
	if diff > 50*time.Millisecond {
		t.Fatalf("timing difference %v exceeds 50ms", diff)
	}
}

func TestLockoutAfterRepeatedFailures(t *testing.T) {
	e := newEnv(t, map[string]string{"MAX_FAILED_LOGINS": "3"})
	e.world()
	g := e.client()
	for i := 0; i < 3; i++ {
		must(t, g.login("p1@example.test", "Wrong-password-1"), 401)
	}
	must(t, g.login("p1@example.test", pw), 401) // locked even with the right password
	if auditCount(t, e, "auth:login", "failure") < 4 {
		t.Fatal("failed logins must be audited")
	}
}

func TestRateLimitPerIdentifierAndIP(t *testing.T) {
	e := newEnv(t, map[string]string{"LOGIN_RATE_PER_MINUTE": "3"})
	g := e.client()
	codes := []int{}
	for i := 0; i < 5; i++ {
		codes = append(codes, g.login("ghost@example.test", "whatever-123").code)
	}
	if codes[4] != 429 {
		t.Fatalf("expected 429 for non-existent identifier, got %v", codes)
	}
	r := g.login("ghost@example.test", "whatever-123")
	if r.hdr.Get("Retry-After") == "" || !strings.Contains(r.hdr.Get("Content-Type"), "application/problem+json") {
		t.Fatalf("429 must carry Retry-After and problem+json: %v", r.hdr)
	}
}

func TestDeactivationAndTamperedTokens(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	must(t, w.parent2.do("GET", "/api/v1/auth/me", nil), 200)
	must(t, w.admin.do("PATCH", "/api/v1/users/"+w.parent2.id.String(), map[string]any{"status": "inactive"}), 200)
	must(t, w.parent2.do("GET", "/api/v1/auth/me", nil), 401) // immediate
	must(t, w.parent2.do("POST", "/api/v1/auth/refresh", nil, "X-CSRF-Token", w.parent2.csrf()), 401)
	must(t, e.client().login("p2@example.test", pw), 401)

	tamper := e.client()
	for _, tok := range []string{
		jwtNone(w.parent1.id.String()),
		signed(w.parent1.id.String(), strings.Repeat("x", 48), "sen-support-platform", time.Hour),
		signed(w.parent1.id.String(), strings.Repeat("k", 48), "other-issuer", time.Hour),
		signed(w.parent1.id.String(), strings.Repeat("k", 48), "sen-support-platform", -time.Hour),
		w.parent1.token[:len(w.parent1.token)-2] + "xx",
	} {
		tamper.token = tok
		must(t, tamper.do("GET", "/api/v1/auth/me", nil), 401)
	}
	tamper.token = ""
	must(t, tamper.do("GET", "/api/v1/auth/me", nil, "Authorization", "Basic abc"), 401)
	_ = context.Background()
}

func jwtNone(sub string) string {
	s, _ := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{"sub": sub, "iss": "sen-support-platform", "aud": "sen-support-platform-web", "exp": time.Now().Add(time.Hour).Unix()}).
		SignedString(jwt.UnsafeAllowNoneSignatureType)
	return s
}

func signed(sub, secret, iss string, ttl time.Duration) string {
	s, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": sub, "iss": iss, "aud": "sen-support-platform-web", "exp": time.Now().Add(ttl).Unix()}).
		SignedString([]byte(secret))
	return s
}
