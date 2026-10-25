package domain

import (
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"
)

// Limits for free-text fields.
const (
	MaxNameLen     = 100
	MaxTitleLen    = 200
	MaxContentLen  = 10000
	MinPasswordLen = 10
	MaxPasswordLen = 256
)

var commonPasswords = map[string]struct{}{
	"password123": {}, "1234567890": {}, "qwertyuiop": {}, "password1234": {},
	"iloveyou123": {}, "12345678910": {}, "qwerty12345": {}, "admin12345": {},
	"letmein1234": {}, "welcome123": {}, "0987654321": {}, "1111111111": {},
	"example-password": {},
}

// NormalizeEmail lowercases and trims an email address.
func NormalizeEmail(e string) string { return strings.ToLower(strings.TrimSpace(e)) }

// ValidEmail reports whether e is a plain RFC 5322 address.
func ValidEmail(e string) bool {
	a, err := mail.ParseAddress(e)
	return err == nil && a.Address == e && len(e) <= 254
}

// ValidPassword enforces SEC-02: length and a common-password deny list.
func ValidPassword(p string) (bool, string) {
	n := utf8.RuneCountInString(p)
	if n < MinPasswordLen {
		return false, "too_short"
	}
	if n > MaxPasswordLen {
		return false, "too_long"
	}
	if _, bad := commonPasswords[strings.ToLower(p)]; bad {
		return false, "too_common"
	}
	return true, ""
}

// ValidName reports whether s is a non-empty name within limits.
func ValidName(s string) bool {
	s = strings.TrimSpace(s)
	return s != "" && utf8.RuneCountInString(s) <= MaxNameLen
}

// ValidText reports whether s is non-empty and at most max runes.
func ValidText(s string, max int) bool {
	return strings.TrimSpace(s) != "" && utf8.RuneCountInString(s) <= max
}

// ValidBirthDate requires a date in the past, not older than 30 years.
func ValidBirthDate(d, now time.Time) bool {
	return d.Before(now) && d.After(now.AddDate(-30, 0, 0))
}

// ParentVisibilityAllowed is the invariant "parent diary = shared only".
func ParentVisibilityAllowed(v Visibility) bool { return v == VisibilityShared }
