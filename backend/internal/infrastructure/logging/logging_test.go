package logging

import (
	"bytes"
	"strings"
	"testing"
)

func TestRedactsSensitiveKeys(t *testing.T) {
	var b bytes.Buffer
	New(&b).Info("x", "email", "a@b.test", "password", "p", "refresh_token", "t", "content", "diary", "route", "/api")
	out := b.String()
	for _, leak := range []string{"a@b.test", `"p"`, `"t"`, "diary"} {
		if strings.Contains(out, leak) {
			t.Errorf("leaked %s: %s", leak, out)
		}
	}
	if !strings.Contains(out, "/api") {
		t.Fatal("non-sensitive fields must stay")
	}
}
