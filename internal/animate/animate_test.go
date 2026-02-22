package animate

import (
	"strings"
	"testing"
)

func TestScrollFrame_TextVisible(t *testing.T) {
	text := "Hello"
	width := 20

	// At position = width, text should start at position 0
	frame := ScrollFrame(text, width, width)
	if !strings.Contains(frame, "Hello") {
		t.Errorf("expected text visible at position=width, got %q", frame)
	}
}

func TestScrollFrame_EmptyWidth(t *testing.T) {
	frame := ScrollFrame("test", 0, 5)
	if frame != "" {
		t.Errorf("expected empty frame for width=0, got %q", frame)
	}
}

func TestScrollFrame_WidthPreserved(t *testing.T) {
	width := 40
	frame := ScrollFrame("Hi", width, 10)
	if len(frame) != width {
		t.Errorf("expected frame length %d, got %d", width, len(frame))
	}
}

func TestScrollCycleLength(t *testing.T) {
	text := "Hello"
	width := 20
	expected := width + len(text)

	got := ScrollCycleLength(text, width)
	if got != expected {
		t.Errorf("expected cycle length %d, got %d", expected, got)
	}
}

func TestCountdownFrame(t *testing.T) {
	tests := []struct {
		seconds  int
		expected string
	}{
		{0, "00:00"},
		{59, "00:59"},
		{60, "01:00"},
		{125, "02:05"},
		{-1, "00:00"}, // negative clamped to 0
	}

	for _, tc := range tests {
		got := CountdownFrame(tc.seconds)
		if got != tc.expected {
			t.Errorf("CountdownFrame(%d) = %q, want %q", tc.seconds, got, tc.expected)
		}
	}
}

func TestCountdownBanner(t *testing.T) {
	banner := CountdownBanner(90)
	if !strings.Contains(banner, "01:30") {
		t.Errorf("expected banner to contain '01:30', got:\n%s", banner)
	}
	// Should have border characters
	if !strings.Contains(banner, "=") {
		t.Error("expected banner to contain border '='")
	}
}
