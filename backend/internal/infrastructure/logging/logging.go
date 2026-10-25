// Package logging builds a JSON slog logger that redacts sensitive keys.
package logging

import (
	"io"
	"log/slog"
	"strings"
)

var sensitive = []string{"password", "token", "secret", "authorization", "cookie", "content", "email"}

// New returns a JSON logger writing to w with redaction of sensitive keys.
func New(w io.Writer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
		k := strings.ToLower(a.Key)
		for _, s := range sensitive {
			if strings.Contains(k, s) {
				return slog.String(a.Key, "[REDACTED]")
			}
		}
		return a
	}}))
}
