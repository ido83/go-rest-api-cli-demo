package command

// headers_test.go - tests header flag parsing in internal/command.
//
// How to run tests (from repo root):
//
//	go test ./...
//
// or just this folder:
//
//	go test ./test

import (
	"strings"
	"testing"
)

func TestHeaderFlagSetAndString(t *testing.T) {
	var h HeaderFlag

	if err := h.Set("X-Test: value"); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}

	if len(h) != 1 {
		t.Fatalf("expected 1 header, got %d", len(h))
	}

	if got := h["X-Test"]; got != "value" {
		t.Errorf("expected X-Test == %q, got %q", "value", got)
	}

	// Add another header
	if err := h.Set("X-Other: 123"); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}

	if len(h) != 2 {
		t.Fatalf("expected 2 headers, got %d", len(h))
	}

	s := h.String()
	if !strings.Contains(s, "X-Test: value") || !strings.Contains(s, "X-Other: 123") {
		t.Errorf("String() did not contain expected headers, got: %q", s)
	}
}

func TestHeaderFlagSetInvalidFormat(t *testing.T) {
	var h HeaderFlag

	if err := h.Set("NoColonHere"); err == nil {
		t.Fatalf("expected error for header without colon, got nil")
	}
}

func TestHeaderFlagSetEmptyKey(t *testing.T) {
	var h HeaderFlag

	if err := h.Set(": value"); err == nil {
		t.Fatalf("expected error for empty header key, got nil")
	}
}
