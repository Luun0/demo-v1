package httpapi

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/AknurOE/sen-support-platform/internal/application"
	"github.com/AknurOE/sen-support-platform/internal/domain"
)

const (
	refreshCookie = "ssp_refresh"
	csrfCookie    = "ssp_csrf"
)

func (a *API) register(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Role      string `json:"role"`
	}
	if err := decode(r, &in); err != nil {
		fail(w, r, a.log, err)
		return
	}
	u, err := a.svc.Register(r.Context(), application.RegisterInput{Email: in.Email, Password: in.Password,
		FirstName: in.FirstName, LastName: in.LastName, Role: domain.Role(in.Role)})
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	writeJSON(w, http.StatusCreated, toUser(u))
}

const passwordResetMessage = "If the account exists, password reset instructions will be provided."

func (a *API) requestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email string `json:"email"`
	}
	if err := decode(r, &in); err != nil {
		fail(w, r, a.log, err)
		return
	}
	token, err := a.svc.RequestPasswordReset(r.Context(), in.Email)
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	if token != "" && a.cfg.LogPasswordResetLinks {
		link := strings.TrimRight(a.cfg.PasswordResetBaseURL, "/") + "/reset-password?token=" + url.QueryEscape(token)
		// This is the prototype's development delivery channel. The flag is
		// rejected when COOKIE_SECURE=true; API request logs never include it.
		a.log.Info("development password reset link: " + link)
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"message": passwordResetMessage})
}

func (a *API) confirmPasswordReset(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := decode(r, &in); err != nil {
		fail(w, r, a.log, err)
		return
	}
	if err := a.svc.ConfirmPasswordReset(r.Context(), in.Token, in.Password); err != nil {
		fail(w, r, a.log, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) setSessionCookies(w http.ResponseWriter, s application.Session) {
	maxAge := int(a.cfg.RefreshTTL.Seconds())
	//nolint:gosec // G124: Secure is config-driven (COOKIE_SECURE=false only for local http dev); HttpOnly+SameSite=Strict are set
	http.SetCookie(w, &http.Cookie{Name: refreshCookie, Value: s.RefreshToken, Path: "/api/v1/auth", MaxAge: maxAge,
		HttpOnly: true, Secure: a.cfg.CookieSecure, SameSite: http.SameSiteStrictMode})
	//nolint:gosec // G124: CSRF double-submit cookie must be readable by JS (ADR-0004); Secure is config-driven
	http.SetCookie(w, &http.Cookie{Name: csrfCookie, Value: randomHex(16), Path: "/", MaxAge: maxAge,
		HttpOnly: false, Secure: a.cfg.CookieSecure, SameSite: http.SameSiteStrictMode})
}

func (a *API) clearSessionCookies(w http.ResponseWriter) {
	for _, c := range []struct{ name, path string }{{refreshCookie, "/api/v1/auth"}, {csrfCookie, "/"}} {
		//nolint:gosec // G124: expiring cookies; attributes mirror setSessionCookies
		http.SetCookie(w, &http.Cookie{Name: c.name, Value: "", Path: c.path, MaxAge: -1, HttpOnly: c.name == refreshCookie,
			Secure: a.cfg.CookieSecure, SameSite: http.SameSiteStrictMode})
	}
}

func (a *API) writeSession(w http.ResponseWriter, s application.Session) {
	a.setSessionCookies(w, s)
	writeJSON(w, http.StatusOK, map[string]any{"access_token": s.AccessToken, "token_type": "Bearer", "expires_in": int(s.ExpiresIn.Seconds())})
}

func (a *API) login(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := decode(r, &in); err != nil {
		fail(w, r, a.log, err)
		return
	}
	s, err := a.svc.Login(r.Context(), in.Email, in.Password)
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	a.writeSession(w, s)
}

func refreshValue(r *http.Request) string {
	c, err := r.Cookie(refreshCookie)
	if err != nil {
		return ""
	}
	return c.Value
}

func (a *API) refresh(w http.ResponseWriter, r *http.Request) {
	s, err := a.svc.Refresh(r.Context(), refreshValue(r))
	if err != nil {
		a.clearSessionCookies(w)
		fail(w, r, a.log, err)
		return
	}
	a.writeSession(w, s)
}

func (a *API) logout(w http.ResponseWriter, r *http.Request) {
	if err := a.svc.Logout(r.Context(), refreshValue(r)); err != nil {
		fail(w, r, a.log, err)
		return
	}
	a.clearSessionCookies(w)
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) me(w http.ResponseWriter, r *http.Request) {
	u, err := a.svc.Me(r.Context(), actorFrom(r.Context()))
	if err != nil {
		fail(w, r, a.log, err)
		return
	}
	writeJSON(w, http.StatusOK, toUser(u))
}
