package tasks

import "testing"

func TestIsAllowedTransition(t *testing.T) {
	if !IsAllowedTransition("new", "received", false) {
		t.Fatal("expected transition new -> received to be allowed")
	}

	if IsAllowedTransition("completed", "in_progress", false) {
		t.Fatal("expected transition from final to working to be denied for non-admin")
	}

	if !IsAllowedTransition("completed", "in_progress", true) {
		t.Fatal("expected transition from final to working to be allowed for admin")
	}
}
