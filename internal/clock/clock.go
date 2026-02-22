package clock

import (
	"fmt"
	"strings"
	"time"
)

// Segments defines the 7-segment patterns for digits 0-9.
// Each digit is represented as 5 rows of 3-character wide strings.
var Segments = map[rune][5]string{
	'0': {" _ ", "| |", "   ", "| |", " _ "},
	'1': {"   ", "  |", "   ", "  |", "   "},
	'2': {" _ ", "  |", " _ ", "|  ", " _ "},
	'3': {" _ ", "  |", " _ ", "  |", " _ "},
	'4': {"   ", "| |", " _ ", "  |", "   "},
	'5': {" _ ", "|  ", " _ ", "  |", " _ "},
	'6': {" _ ", "|  ", " _ ", "| |", " _ "},
	'7': {" _ ", "  |", "   ", "  |", "   "},
	'8': {" _ ", "| |", " _ ", "| |", " _ "},
	'9': {" _ ", "| |", " _ ", "  |", " _ "},
}

// ColonOn is the colon separator when visible.
var ColonOn = [5]string{" ", "o", " ", "o", " "}

// ColonOff is the colon separator when hidden (for blinking).
var ColonOff = [5]string{" ", " ", " ", " ", " "}

// RenderDigit returns the 5-row representation of a single digit character.
func RenderDigit(ch rune) [5]string {
	if seg, ok := Segments[ch]; ok {
		return seg
	}
	return [5]string{"   ", "   ", "   ", "   ", "   "}
}

// RenderTime builds the full ASCII clock string for the given time.
// showColon controls whether the colon is displayed (for blinking effect).
func RenderTime(t time.Time, showColon bool) string {
	timeStr := fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())

	var parts [][5]string
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

	var lines [5]string
	for row := 0; row < 5; row++ {
		var segments []string
		for _, p := range parts {
			segments = append(segments, p[row])
		}
		lines[row] = strings.Join(segments, " ")
	}

	return strings.Join(lines[:], "\n")
}
