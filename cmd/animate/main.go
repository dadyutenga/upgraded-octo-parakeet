package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/dadyutenga/upgraded-octo-parakeet/internal/animate"
)

func main() {
	mode := flag.String("mode", "scroll", "Animation mode: scroll or countdown")
	text := flag.String("text", "Good Morning Dadi", "Text to scroll across terminal")
	width := flag.Int("width", 60, "Terminal width for scrolling")
	seconds := flag.Int("seconds", 10, "Countdown duration in seconds")
	flag.Parse()

	switch *mode {
	case "scroll":
		runScroll(*text, *width)
	case "countdown":
		runCountdown(*seconds)
	default:
		fmt.Fprintf(os.Stderr, "Unknown mode: %s (use 'scroll' or 'countdown')\n", *mode)
		os.Exit(1)
	}
}

func runScroll(text string, width int) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	cycleLen := animate.ScrollCycleLength(text, width)
	pos := 0

	fmt.Print("\033[?25l") // hide cursor
	defer fmt.Print("\033[?25h\n")

	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-sig:
			return
		case <-ticker.C:
			frame := animate.ScrollFrame(text, width, pos)
			fmt.Printf("\r\033[36m%s\033[0m", frame) // cyan color
			pos++
			if pos > cycleLen {
				pos = 0
			}
		}
	}
}

func runCountdown(totalSeconds int) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	fmt.Print("\033[2J") // clear screen

	for remaining := totalSeconds; remaining >= 0; remaining-- {
		select {
		case <-sig:
			fmt.Print("\033[?25h")
			fmt.Println("\nCancelled!")
			return
		default:
		}

		fmt.Print("\033[H") // move to top
		banner := animate.CountdownBanner(remaining)
		fmt.Printf("\033[33m%s\033[0m\n", banner) // yellow color

		if remaining > 0 {
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Println("\n\033[32m🎉 Time's up!\033[0m")
}
