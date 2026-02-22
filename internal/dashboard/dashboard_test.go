package dashboard

import (
	"strings"
	"testing"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tc := range tests {
		got := FormatBytes(tc.input)
		if got != tc.expected {
			t.Errorf("FormatBytes(%d) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestColorForPercent(t *testing.T) {
	tests := []struct {
		pct      float64
		expected string
	}{
		{95, Red},
		{75, Yellow},
		{55, Magenta},
		{30, Green},
		{0, Green},
	}

	for _, tc := range tests {
		got := ColorForPercent(tc.pct)
		if got != tc.expected {
			t.Errorf("ColorForPercent(%.0f) = %q, want %q", tc.pct, got, tc.expected)
		}
	}
}

func TestProgressBar(t *testing.T) {
	bar := ProgressBar(50, 10)
	if !strings.HasPrefix(bar, "[") || !strings.HasSuffix(bar, "]") {
		t.Errorf("expected bar with brackets, got %q", bar)
	}
	// 50% of 10 = 5 filled
	if !strings.Contains(bar, "█████") {
		t.Errorf("expected 5 filled blocks at 50%%, got %q", bar)
	}
}

func TestProgressBar_EdgeCases(t *testing.T) {
	bar0 := ProgressBar(0, 10)
	if strings.Contains(bar0, "█") {
		t.Errorf("expected no filled blocks at 0%%, got %q", bar0)
	}

	bar100 := ProgressBar(100, 10)
	if strings.Contains(bar100, "░") {
		t.Errorf("expected all filled blocks at 100%%, got %q", bar100)
	}

	barSmall := ProgressBar(50, 1)
	if !strings.HasPrefix(barSmall, "[") {
		t.Errorf("expected valid bar even with width=1, got %q", barSmall)
	}
}

func TestRender(t *testing.T) {
	info := SystemInfo{
		DateTime:   "2025-01-01 12:00:00",
		CPUUsage:   45.5,
		MemTotal:   8589934592,
		MemUsed:    4294967296,
		MemPercent: 50.0,
		GoRoutines: 1,
		NumCPU:     4,
		OS:         "linux",
		Arch:       "amd64",
	}

	output := Render(info)

	if !strings.Contains(output, "2025-01-01 12:00:00") {
		t.Error("expected date/time in output")
	}
	if !strings.Contains(output, "linux/amd64") {
		t.Error("expected platform in output")
	}
	if !strings.Contains(output, "45.5%") {
		t.Error("expected CPU percentage in output")
	}
	if !strings.Contains(output, "50.0%") {
		t.Error("expected memory percentage in output")
	}
}
