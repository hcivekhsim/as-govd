package tasks

import "testing"

func TestOptionalDate(t *testing.T) {
	date, err := optionalDate("2026-11-14")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if date == nil {
		t.Fatal("expected parsed date")
	}

	if _, err := optionalDate("14-11-2026"); err == nil {
		t.Fatal("expected invalid date error")
	}
}

func TestOptionalString(t *testing.T) {
	if optionalString("") != nil {
		t.Fatal("expected nil for empty string")
	}
	if value := optionalString("x"); value == nil || *value != "x" {
		t.Fatal("expected pointer to x")
	}
}
