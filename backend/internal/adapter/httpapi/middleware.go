package httpapi

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"log/slog"
	"net"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/AknurOE/sen-support-platform/internal/application"
)

type actorKey struct{}

func actorFrom(ctx context.Context) *application.Actor {
	a, _ := ctx.Value(actorKey{}).(*application.Actor)
	return a
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusWriter) WriteHeader(c int) { s.status = c; s.ResponseWriter.WriteHeader(c) }

func (a *API) requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if len(id) == 0 || len(id) > 64 || strings.ContainsFunc(id, func(c rune) bool { return c < 33 || c > 126 }) {
			id = randomHex(8)
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(application.WithRequestID(r.Context(), id)))
	})
}

func (a *API) recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				a.log.ErrorContext(r.Context(), "panic", "request_id", application.RequestID(r.Context()))
				writeProblem(w, r, http.StatusInternalServerError, nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// logging writes one JSON line per request without PII: no query string,
// body, IP, email or tokens.
func (a *API) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(sw, r)
		a.log.LogAttrs(r.Context(), slog.LevelInfo, "request", slog.String("request_id", application.RequestID(r.Context())),
			slog.String("method", r.Method), slog.String("route", r.Pattern), slog.Int("status", sw.status),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()))
	})
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; base-uri 'none'")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		h.Set("Cross-Origin-Opener-Policy", "same-origin")
		h.Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (a *API) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowed := origin != "" && slices.Contains(a.cfg.AllowedOrigins, origin)
		if allowed {
			h := w.Header()
			h.Set("Access-Control-Allow-Origin", origin)
			h.Set("Access-Control-Allow-Credentials", "true")
			h.Set("Access-Control-Expose-Headers", "X-Request-ID, X-Next-Cursor, X-Total-Count, Retry-After")
			h.Add("Vary", "Origin")
		}
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			if !allowed {
				writeProblem(w, r, http.StatusForbidden, nil)
				return
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-CSRF-Token, X-Request-ID")
			w.Header().Set("Access-Control-Max-Age", "600")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *API) bodyLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > a.cfg.MaxBodyBytes {
			writeProblem(w, r, http.StatusRequestEntityTooLarge, nil)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, a.cfg.MaxBodyBytes)
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func (a *API) rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		ip := clientIP(r)
		if !a.apiLimiter.Allow(ip, now) {
			writeProblem(w, r, http.StatusTooManyRequests, nil)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/v1/auth/") && r.Method == http.MethodPost && !a.authLimiter.Allow(ip, now) {
			writeProblem(w, r, http.StatusTooManyRequests, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// authenticate attaches the actor when a valid bearer token is present.
func (a *API) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if tok, ok := strings.CutPrefix(h, "Bearer "); ok && tok != "" {
			actor, err := a.svc.Authenticate(r.Context(), tok)
			if err != nil {
				fail(w, r, a.log, err)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), actorKey{}, actor))
		} else if h != "" {
			fail(w, r, a.log, errUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// requireAuth rejects guests with 401.
func requireAuth(log *slog.Logger, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if actorFrom(r.Context()) == nil {
			fail(w, r, log, errUnauthorized)
			return
		}
		h(w, r)
	}
}

// csrf protects cookie-authenticated endpoints (double submit + Origin).
func (a *API) csrf(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if o := r.Header.Get("Origin"); o != "" && !slices.Contains(a.cfg.AllowedOrigins, o) {
			writeProblem(w, r, http.StatusForbidden, nil)
			return
		}
		c, err := r.Cookie(csrfCookie)
		hdr := r.Header.Get("X-CSRF-Token")
		if err != nil || hdr == "" || subtle.ConstantTimeCompare([]byte(c.Value), []byte(hdr)) != 1 {
			writeProblem(w, r, http.StatusForbidden, nil)
			return
		}
		h(w, r)
	}
}
