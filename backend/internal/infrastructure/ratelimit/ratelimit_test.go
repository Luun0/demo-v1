package ratelimit

import (
	"testing"
	"time"
)

func TestFixedWindow(t *testing.T) {
	l := New(2, time.Minute)
	now := time.Now()
	first, second, third := l.Allow("a", now), l.Allow("a", now), l.Allow("a", now)
	if !first || !second || third {
		t.Fatal("limit 2 not enforced")
	}
	if !l.Allow("b", now) {
		t.Fatal("keys must be independent")
	}
	if !l.Allow("a", now.Add(time.Minute)) {
		t.Fatal("window must reset")
	}
}
