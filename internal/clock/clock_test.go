package clock

import (
	"strings"
	"testing"
	"time"
)

func TestRenderDigit_ValidDigits(t *testing.T) {
	for _, ch := range "0123456789" {
		seg := RenderDigit(ch)
		if len(seg) != DigitRows {
			t.Errorf("digit %c: expected %d rows, got %d", ch, DigitRows, len(seg))
		}
		for row, line := range seg {
			if len([]rune(line)) != DigitWidth {
				t.Errorf("digit %c row %d: expected width %d, got %d (%q)", ch, row, DigitWidth, len([]rune(line)), line)
			}
		}
	}
}

func TestRenderDigit_Invalid(t *testing.T) {
	seg := RenderDigit('X')
	for row, line := range seg {
		if strings.TrimSpace(line) != "" {
			t.Errorf("invalid char row %d: expected blank, got %q", row, line)
		}
	}
}

func TestRenderTime_Format(t *testing.T) {
	// Use a fixed time: 13:45
	tm := time.Date(2025, 1, 1, 13, 45, 0, 0, time.UTC)
	output := RenderTime(tm, true)

	lines := strings.Split(output, "\n")
	if len(lines) != DigitRows {
		t.Fatalf("expected %d lines, got %d", DigitRows, len(lines))
	}

	// All lines should have the same rune width
	width := len([]rune(lines[0]))
	for i, line := range lines {
		if len([]rune(line)) != width {
			t.Errorf("line %d width %d != expected %d", i, len([]rune(line)), width)
		}
	}
}

func TestRenderTime_ColonBlink(t *testing.T) {
	tm := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	withColon := RenderTime(tm, true)
	withoutColon := RenderTime(tm, false)

	if withColon == withoutColon {
		t.Error("expected colon on/off to produce different output")
	}

	// "██" should appear in the colon-on version but not colon-off
	if !strings.Contains(withColon, "██") {
		t.Error("expected '██' in colon-on output")
	}
}

func TestRenderTime_AllDigits(t *testing.T) {
	// Test edge cases: 00:00 and 23:59
	tests := []struct {
		hour, min int
	}{
		{0, 0},
		{23, 59},
		{12, 30},
		{9, 5},
	}

	for _, tc := range tests {
		tm := time.Date(2025, 1, 1, tc.hour, tc.min, 0, 0, time.UTC)
		output := RenderTime(tm, true)
		if output == "" {
			t.Errorf("empty output for %02d:%02d", tc.hour, tc.min)
		}
	}
}
