package domain

import (
	"errors"
	"strings"
)

// Sentinel errors mapped to HTTP status codes by the adapter (ADR-0006).
var (
	ErrNotFound     = errors.New("not found")
	ErrForbidden    = errors.New("forbidden")
	ErrUnauthorized = errors.New("unauthorized")
	ErrConflict     = errors.New("conflict")
	ErrRateLimited  = errors.New("rate limited")
)

// FieldError names an invalid field and a machine code; never the value.
type FieldError struct {
	Field string
	Code  string
}

// ValidationError aggregates field errors (HTTP 422).
type ValidationError struct{ Fields []FieldError }

func (e *ValidationError) Error() string {
	parts := make([]string, 0, len(e.Fields))
	for _, f := range e.Fields {
		parts = append(parts, f.Field+":"+f.Code)
	}
	return "invalid: " + strings.Join(parts, ",")
}

// Validator collects field errors.
type Validator struct{ fields []FieldError }

// Check records code for field when ok is false.
func (v *Validator) Check(ok bool, field, code string) {
	if !ok {
		v.fields = append(v.fields, FieldError{Field: field, Code: code})
	}
}

// Err returns a *ValidationError or nil.
func (v *Validator) Err() error {
	if len(v.fields) == 0 {
		return nil
	}
	return &ValidationError{Fields: v.fields}
}
