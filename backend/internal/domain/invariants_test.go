package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestEnumsValid(t *testing.T) {
	for _, ok := range []bool{
		RoleParent.Valid(), RoleAdmin.Valid(), !Role("root").Valid(),
		VisibilityShared.Valid(), !Visibility("public").Valid(),
		EntryProgress.Valid(), !EntryType("x").Valid(),
		ModeInPerson.Valid(), !SessionMode("offline").Valid(),
		SessionCancelled.Valid(), !SessionStatus("canceled").Valid(),
		RelGuardian.Valid(), !Relationship("x").Valid(),
	} {
		if !ok {
			t.Fatal("enum validation mismatch")
		}
	}
}

func TestPasswordPolicy(t *testing.T) {
	cases := map[string]string{"short": "too_short", "Password123": "too_common", strings.Repeat("a", 300): "too_long", "correct horse battery": ""}
	for p, want := range cases {
		ok, code := ValidPassword(p)
		if code != want || ok != (want == "") {
			t.Errorf("%q: got %v %q", p, ok, code)
		}
	}
}

func TestEmailNameTextDate(t *testing.T) {
	now := time.Date(2026, 10, 24, 0, 0, 0, 0, time.UTC)
	if !ValidEmail("a@example.test") || ValidEmail("Bob <a@example.test>") || ValidEmail("nope") {
		t.Fatal("email")
	}
	if NormalizeEmail("  A@Example.Test ") != "a@example.test" {
		t.Fatal("normalize")
	}
	if ValidName("  ") || !ValidName("Aida") || ValidName(strings.Repeat("x", 101)) {
		t.Fatal("name")
	}
	if ValidText("", 5) || !ValidText("abc", 5) || ValidText("abcdef", 5) {
		t.Fatal("text")
	}
	if !ValidBirthDate(now.AddDate(-7, 0, 0), now) || ValidBirthDate(now.AddDate(1, 0, 0), now) || ValidBirthDate(now.AddDate(-40, 0, 0), now) {
		t.Fatal("birth date")
	}
	if !ParentVisibilityAllowed(VisibilityShared) || ParentVisibilityAllowed(VisibilitySpecialistsOnly) {
		t.Fatal("parent visibility invariant")
	}
}

func TestActivityAndAssignments(t *testing.T) {
	now := time.Now()
	past, future := now.Add(-time.Hour), now.Add(time.Hour)
	if !(User{Status: StatusActive}).Active(now) || (User{Status: StatusInactive}).Active(now) || (User{Status: StatusActive, LockedUntil: &future}).Active(now) {
		t.Fatal("user active")
	}
	if !(ChildSpecialist{AssignedAt: past}).ActiveAt(now) || (ChildSpecialist{AssignedAt: past, EndedAt: &past}).ActiveAt(now) || !(ChildSpecialist{AssignedAt: past, EndedAt: &future}).ActiveAt(now) {
		t.Fatal("assignment active")
	}
}

func TestValidator(t *testing.T) {
	var v Validator
	if v.Err() != nil {
		t.Fatal("empty validator must be nil")
	}
	v.Check(false, "email", "invalid")
	var ve *ValidationError
	if err := v.Err(); !errors.As(err, &ve) || ve.Error() != "invalid: email:invalid" {
		t.Fatalf("got %v", err)
	}
}
