package httpapi

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/application"
	"github.com/AknurOE/sen-support-platform/internal/domain"
)

var errUnauthorized = domain.ErrUnauthorized

// Config holds HTTP-level settings.
type Config struct {
	AllowedOrigins []string
	MaxBodyBytes   int64
	CookieSecure   bool
	RefreshTTL     time.Duration
}

// API is the HTTP adapter.
type API struct {
	svc         *application.Service
	cfg         Config
	log         *slog.Logger
	apiLimiter  application.Limiter
	authLimiter application.Limiter
}

// New builds the handler with the full middleware chain (master prompt §3.1.7).
func New(svc *application.Service, cfg Config, log *slog.Logger, apiLimiter, authLimiter application.Limiter) http.Handler {
	a := &API{svc: svc, cfg: cfg, log: log, apiLimiter: apiLimiter, authLimiter: authLimiter}
	mux := http.NewServeMux()
	for _, rt := range a.routes() {
		h := rt.h
		if rt.auth {
			h = requireAuth(log, h)
		}
		mux.HandleFunc(rt.method+" "+rt.path, h)
	}
	mux.HandleFunc("GET /healthz", a.healthz)
	mux.HandleFunc("GET /readyz", a.readyz)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { writeProblem(w, r, http.StatusNotFound, nil) })

	var h http.Handler = mux
	for _, mw := range []func(http.Handler) http.Handler{a.authenticate, a.rateLimit, a.bodyLimit, a.cors, securityHeaders, a.logging, a.recoverer, a.requestID} {
		h = mw(h)
	}
	return h
}

func (a *API) healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) readyz(w http.ResponseWriter, r *http.Request) {
	if err := a.svc.Ready(r.Context()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// pathID parses a path UUID; malformed IDs are 404 (do not reveal anything).
func pathID(r *http.Request, name string) (uuid.UUID, error) {
	id, err := uuid.Parse(r.PathValue(name))
	if err != nil {
		return uuid.Nil, domain.ErrNotFound
	}
	return id, nil
}

func pathIDs(r *http.Request, names ...string) ([]uuid.UUID, error) {
	out := make([]uuid.UUID, 0, len(names))
	for _, n := range names {
		id, err := pathID(r, n)
		if err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, nil
}

func invalid(field string) error {
	return &domain.ValidationError{Fields: []domain.FieldError{{Field: field, Code: "invalid"}}}
}

func parseTime(field, s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return t, invalid(field)
	}
	return t.UTC(), nil
}

type route struct {
	method, path string
	h            http.HandlerFunc
	auth         bool
}

func (a *API) routes() []route {
	const v = "/api/v1"
	c := v + "/children/{childId}"
	return []route{
		{"GET", v + "/healthz", a.healthz, false},
		{"GET", v + "/readyz", a.readyz, false},
		{"POST", v + "/auth/register", a.register, false},
		{"POST", v + "/auth/login", a.login, false},
		{"POST", v + "/auth/refresh", a.csrf(a.refresh), false},
		{"POST", v + "/auth/logout", a.csrf(a.logout), false},
		{"GET", v + "/auth/me", a.me, true},
		{"GET", v + "/users", a.listUsers, true},
		{"POST", v + "/users", a.createUser, true},
		{"PATCH", v + "/users/{userId}", a.patchUser, true},
		{"GET", v + "/audit-events", a.systemAudit, true},
		{"GET", v + "/children", a.listChildren, true},
		{"POST", v + "/children", a.createChild, true},
		{"GET", c, a.getChild, true},
		{"PATCH", c, a.patchChild, true},
		{"DELETE", c, a.deleteChild, true},
		{"POST", c + "/guardians", a.linkGuardian, true},
		{"POST", c + "/specialist-assignments", a.assign, true},
		{"PATCH", c + "/specialist-assignments/{specialistId}", a.endAssignment, true},
		{"GET", c + "/sessions", a.listSessions, true},
		{"POST", c + "/sessions", a.createSession, true},
		{"PATCH", c + "/sessions/{sessionId}", a.patchSession, true},
		{"GET", c + "/diary-entries", a.listEntries, true},
		{"POST", c + "/diary-entries", a.createEntry, true},
		{"PUT", c + "/diary-entries/{entryId}", a.editEntry, true},
		{"DELETE", c + "/diary-entries/{entryId}", a.deleteEntry, true},
		{"GET", c + "/diary-entries/{entryId}/versions", a.listVersions, true},
		{"GET", c + "/audit-events", a.childAudit, true},
		{"POST", c + "/exports", a.export, true},
		{"GET", c + "/team", a.team, true},
		{"GET", v + "/specialists", a.specialists, true},
	}
}

// Routes lists "METHOD /path" for every API route (used by the contract test).
func Routes() []string {
	var out []string
	for _, r := range (&API{}).routes() {
		out = append(out, r.method+" "+r.path)
	}
	return out
}
