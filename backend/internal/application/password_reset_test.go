package application

import (
	"encoding/base64"
	"testing"
)

func TestPasswordResetTokenIsOpaqueAndHashed(t *testing.T) {
	raw, hash, err := newPasswordResetToken()
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil || len(decoded) != 32 {
		t.Fatalf("reset token should be 256 random bits: bytes=%d err=%v", len(decoded), err)
	}
	if string(hash) == raw || string(hash) != string(HashToken(raw)) {
		t.Fatal("persistence value must be the hash, not the raw reset token")
	}
}
