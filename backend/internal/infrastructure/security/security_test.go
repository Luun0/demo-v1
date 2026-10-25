package security

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestArgon2HashVerifyAndRehash(t *testing.T) {
	a := Argon2{Memory: 8192, Time: 1, Threads: 1}
	h, err := a.Hash("correct horse battery")
	if err != nil || !strings.HasPrefix(h, "$argon2id$v=19$m=8192,t=1,p=1$") {
		t.Fatalf("%v %s", err, h)
	}
	h2, _ := a.Hash("correct horse battery")
	if h == h2 {
		t.Fatal("salt must be random")
	}
	if ok, re, err := a.Verify(h, "correct horse battery"); !ok || re || err != nil {
		t.Fatal("verify")
	}
	if ok, _, _ := a.Verify(h, "wrong"); ok {
		t.Fatal("wrong password accepted")
	}
	stronger := Argon2{Memory: 16384, Time: 2, Threads: 1}
	if ok, re, _ := stronger.Verify(h, "correct horse battery"); !ok || !re {
		t.Fatal("outdated params must request rehash")
	}
	for _, bad := range []string{"", "$bcrypt$x", "$argon2id$v=1$m=1,t=1,p=1$a$b", "$argon2id$v=19$m=x$a$b", "$argon2id$v=19$m=1,t=1,p=1$!!$b"} {
		if _, _, err := a.Verify(bad, "x"); err == nil {
			t.Errorf("malformed hash accepted: %q", bad)
		}
	}
}

func TestJWTPinsAlgIssuerAudienceExpiry(t *testing.T) {
	j := JWT{Secret: []byte(strings.Repeat("s", 32)), Issuer: "iss", Audience: "aud", TTL: 15 * time.Minute}
	now := time.Now()
	id := uuid.New()
	tok, ttl, err := j.Issue(id, now)
	if err != nil || ttl != 15*time.Minute {
		t.Fatal(err)
	}
	if got, err := j.Verify(tok, now); err != nil || got != id {
		t.Fatalf("%v %v", got, err)
	}
	if _, err := j.Verify(tok, now.Add(16*time.Minute)); err == nil {
		t.Fatal("expired accepted")
	}
	other := j
	other.Audience = "x"
	if _, err := other.Verify(tok, now); err == nil {
		t.Fatal("wrong audience accepted")
	}
	other = j
	other.Secret = []byte(strings.Repeat("t", 32))
	if _, err := other.Verify(tok, now); err == nil {
		t.Fatal("wrong key accepted")
	}
}
