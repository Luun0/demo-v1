// Package security provides Argon2id hashing and JWT access tokens.
package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2 hashes passwords as PHC strings.
type Argon2 struct {
	Memory  uint32
	Time    uint32
	Threads uint8
}

const keyLen = 32

// Hash returns $argon2id$v=19$m=..,t=..,p=..$salt$hash.
func (a Argon2) Hash(pw string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := argon2.IDKey([]byte(pw), salt, a.Time, a.Memory, a.Threads, keyLen)
	b := base64.RawStdEncoding
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, a.Memory, a.Time, a.Threads,
		b.EncodeToString(salt), b.EncodeToString(key)), nil
}

var errFormat = errors.New("invalid hash format")

// Verify checks pw in constant time and reports outdated parameters.
func (a Argon2) Verify(encoded, pw string) (bool, bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, false, errFormat
	}
	var v int
	var m, t uint32
	var p uint8
	if _, err := fmt.Sscanf(parts[2], "v=%d", &v); err != nil || v != argon2.Version {
		return false, false, errFormat
	}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &m, &t, &p); err != nil {
		return false, false, errFormat
	}
	b := base64.RawStdEncoding
	salt, err1 := b.DecodeString(parts[4])
	want, err2 := b.DecodeString(parts[5])
	if err1 != nil || err2 != nil || len(want) == 0 {
		return false, false, errFormat
	}
	got := argon2.IDKey([]byte(pw), salt, t, m, p, uint32(len(want))) //nolint:gosec // len(want) is small
	ok := subtle.ConstantTimeCompare(got, want) == 1
	return ok, m != a.Memory || t != a.Time || p != a.Threads, nil
}
