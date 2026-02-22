package clock

import (
	"fmt"
	"strings"
	"time"
)

const (
	DigitWidth = 8
	DigitRows  = 7
	ColonWidth = 4
)

// Segments defines large block-style patterns for digits 0-9.
// Each digit is 7 rows tall and 8 characters wide using █ and spaces.
var Segments = map[rune][DigitRows]string{
	'0': {
		" ██████ ",
		"██    ██",
		"██    ██",
		"██    ██",
		"██    ██",
		"██    ██",
		" ██████ ",
	},
	'1': {
		"   ██   ",
		" ████   ",
		"   ██   ",
		"   ██   ",
		"   ██   ",
		"   ██   ",
		" ██████ ",
	},
	'2': {
		" ██████ ",
		"██    ██",
		"      ██",
		"  ████  ",
		"██      ",
		"██      ",
		"████████",
	},
	'3': {
		" ██████ ",
		"██    ██",
		"      ██",
		"  ████  ",
		"      ██",
		"██    ██",
		" ██████ ",
	},
	'4': {
		"██    ██",
		"██    ██",
		"██    ██",
		"████████",
		"      ██",
		"      ██",
		"      ██",
	},
	'5': {
		"████████",
		"██      ",
		"██      ",
		"███████ ",
		"      ██",
		"██    ██",
		" ██████ ",
	},
	'6': {
		" ██████ ",
		"██      ",
		"██      ",
		"███████ ",
		"██    ██",
		"██    ██",
		" ██████ ",
	},
	'7': {
		"████████",
		"      ██",
		"     ██ ",
		"    ██  ",
		"   ██   ",
		"   ██   ",
		"   ██   ",
	},
	'8': {
		" ██████ ",
		"██    ██",
		"██    ██",
		" ██████ ",
		"██    ██",
		"██    ██",
		" ██████ ",
	},
	'9': {
		" ██████ ",
		"██    ██",
		"██    ██",
		" ███████",
		"      ██",
		"      ██",
		" ██████ ",
	},
}

// ColonOn is the colon separator when visible.
var ColonOn = [DigitRows]string{
	"    ",
	" ██ ",
	" ██ ",
	"    ",
	" ██ ",
	" ██ ",
	"    ",
}

// ColonOff is the colon separator when hidden (for blinking).
var ColonOff = [DigitRows]string{
	"    ",
	"    ",
	"    ",
	"    ",
	"    ",
	"    ",
	"    ",
}

// RenderDigit returns the row representation of a single digit character.
func RenderDigit(ch rune) [DigitRows]string {
	if seg, ok := Segments[ch]; ok {
		return seg
	}
	blank := strings.Repeat(" ", DigitWidth)
	var result [DigitRows]string
	for i := range result {
		result[i] = blank
	}
	return result
}

// RenderTime builds the full ASCII clock string for the given time with seconds.
// showColon controls whether the colon is displayed (for blinking effect).
func RenderTime(t time.Time, showColon bool) string {
	timeStr := fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())

	var parts [][DigitRows]string
	for _, ch := range timeStr {
		if ch == ':' {
			if showColon {
				parts = append(parts, ColonOn)
			} else {
				parts = append(parts, ColonOff)
			}
		} else {
			parts = append(parts, RenderDigit(ch))
		}
	}

	var lines [DigitRows]string
	for row := 0; row < DigitRows; row++ {
		var segments []string
		for _, p := range parts {
			segments = append(segments, p[row])
		}
		lines[row] = strings.Join(segments, "  ")
	}

	return strings.Join(lines[:], "\n")
}

// RenderDuration builds the ASCII clock string for a duration (used by stopwatch).
func RenderDuration(d time.Duration, showColon bool) string {
	total := int(d.Seconds())
	minutes := total / 60
	seconds := total % 60
	timeStr := fmt.Sprintf("%02d:%02d", minutes, seconds)

	var parts [][DigitRows]string
	for _, ch := range timeStr {
		if ch == ':' {
			if showColon {
				parts = append(parts, ColonOn)
			} else {
				parts = append(parts, ColonOff)
			}
		} else {
			parts = append(parts, RenderDigit(ch))
		}
	}

	var lines [DigitRows]string
	for row := 0; row < DigitRows; row++ {
		var segments []string
		for _, p := range parts {
			segments = append(segments, p[row])
		}
		lines[row] = strings.Join(segments, "  ")
	}

	return strings.Join(lines[:], "\n")
}
