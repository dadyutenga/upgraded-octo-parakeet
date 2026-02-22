package animate

import (
	"fmt"
	"strings"
)

// ScrollFrame generates a single frame of scrolling text across a terminal
// of the given width. position is the current offset (can exceed width for
// the text to scroll fully across and off screen).
func ScrollFrame(text string, width, position int) string {
	if width <= 0 {
		return ""
	}

	// Build a line of spaces with the text placed at the offset
	line := make([]byte, width)
	for i := range line {
		line[i] = ' '
	}

	textStart := width - position
	for i, ch := range []byte(text) {
		pos := textStart + i
		if pos >= 0 && pos < width {
			line[pos] = ch
		}
	}

	return string(line)
}

// ScrollCycleLength returns the number of positions needed for the text
// to scroll fully from right to left across a terminal of the given width.
func ScrollCycleLength(text string, width int) int {
	return width + len(text)
}

// CountdownFrame returns a formatted countdown string for the given
// remaining seconds.
func CountdownFrame(secondsLeft int) string {
	if secondsLeft < 0 {
		secondsLeft = 0
	}
	minutes := secondsLeft / 60
	seconds := secondsLeft % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// CountdownBanner returns a large ASCII art representation of the countdown.
func CountdownBanner(secondsLeft int) string {
	text := CountdownFrame(secondsLeft)
	border := strings.Repeat("=", len(text)+4)
	return fmt.Sprintf("%s\n| %s |\n%s", border, text, border)
}
