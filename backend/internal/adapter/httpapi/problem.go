// Package httpapi is the HTTP adapter: routing, middleware, DTOs and
// RFC 9457 problem responses. It makes no access decisions (ADR-0001).
package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/AknurOE/sen-support-platform/internal/application"
	"github.com/AknurOE/sen-support-platform/internal/domain"
)

type fieldErr struct {
	Field string `json:"field"`
	Code  string `json:"code"`
}

type problem struct {
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Status    int        `json:"status"`
	Detail    string     `json:"detail,omitempty"`
	RequestID string     `json:"request_id"`
	Errors    []fieldErr `json:"errors,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeProblem(w http.ResponseWriter, r *http.Request, status int, fields []fieldErr) {
	if status == http.StatusTooManyRequests {
		w.Header().Set("Retry-After", "60")
	}
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(problem{Type: "about:blank", Title: http.StatusText(status), Status: status,
		RequestID: application.RequestID(r.Context()), Errors: fields})
}

// fail maps domain errors to status codes; unknown errors become 500
// without details (SEC-23).
func fail(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	var ve *domain.ValidationError
	var mbe *http.MaxBytesError
	switch {
	case errors.As(err, &ve):
		fs := make([]fieldErr, 0, len(ve.Fields))
		for _, f := range ve.Fields {
			fs = append(fs, fieldErr(f))
		}
		writeProblem(w, r, http.StatusUnprocessableEntity, fs)
	case errors.As(err, &mbe):
		writeProblem(w, r, http.StatusRequestEntityTooLarge, nil)
	case errors.Is(err, domain.ErrNotFound):
		writeProblem(w, r, http.StatusNotFound, nil)
	case errors.Is(err, domain.ErrForbidden):
		writeProblem(w, r, http.StatusForbidden, nil)
	case errors.Is(err, domain.ErrUnauthorized):
		w.Header().Set("WWW-Authenticate", `Bearer realm="api"`)
		writeProblem(w, r, http.StatusUnauthorized, nil)
	case errors.Is(err, domain.ErrConflict):
		writeProblem(w, r, http.StatusConflict, nil)
	case errors.Is(err, domain.ErrRateLimited):
		writeProblem(w, r, http.StatusTooManyRequests, nil)
	default:
		log.ErrorContext(r.Context(), "internal error", "request_id", application.RequestID(r.Context()), "err_type", errType(err))
		writeProblem(w, r, http.StatusInternalServerError, nil)
	}
}

func errType(err error) string { return fmt.Sprintf("%T", err) }

// decode reads strict JSON: unknown fields (mass assignment) are rejected.
func decode(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		var mbe *http.MaxBytesError
		if errors.As(err, &mbe) {
			return err
		}
		return &domain.ValidationError{Fields: []domain.FieldError{{Field: "body", Code: "malformed_or_unknown_field"}}}
	}
	if dec.More() {
		return &domain.ValidationError{Fields: []domain.FieldError{{Field: "body", Code: "trailing_data"}}}
	}
	return nil
}
