package application

import (
	"strings"
	"testing"
	"time"

	"github.com/AknurOE/sen-support-platform/internal/domain"
)

func TestValidateEntryBoundaries(t *testing.T) {
	valid := EntryInput{
		Type: domain.EntryObservation, Visibility: domain.VisibilityShared,
		Content: strings.Repeat("ж", domain.MaxContentLen), ObservedAt: time.Date(2026, 10, 26, 9, 0, 0, 0, time.UTC),
	}
	cases := []struct {
		name  string
		input EntryInput
		valid bool
	}{
		{name: "maximum unicode rune count", input: valid, valid: true},
		{name: "invalid entry type", input: func() EntryInput { v := valid; v.Type = "unknown"; return v }()},
		{name: "invalid visibility", input: func() EntryInput { v := valid; v.Visibility = "public"; return v }()},
		{name: "blank content", input: func() EntryInput { v := valid; v.Content = " \t\n"; return v }()},
		{name: "over maximum unicode rune count", input: func() EntryInput { v := valid; v.Content += "ж"; return v }()},
		{name: "missing observation time", input: func() EntryInput { v := valid; v.ObservedAt = time.Time{}; return v }()},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateEntry(tc.input)
			if (err == nil) != tc.valid {
				t.Fatalf("validateEntry() error = %v, valid = %v", err, tc.valid)
			}
		})
	}
}
