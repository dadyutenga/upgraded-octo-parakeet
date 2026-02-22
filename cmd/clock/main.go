package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/dadyutenga/upgraded-octo-parakeet/internal/clock"
)

func main() {
	// Handle Ctrl+C gracefully
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	showColon := true

	// Initial render
	fmt.Print("\033[2J\033[H") // clear screen
	fmt.Println(clock.RenderTime(time.Now(), showColon))

	for {
		select {
		case <-sig:
			fmt.Print("\033[?25h") // show cursor
			fmt.Println("\nGoodbye!")
			return
		case <-ticker.C:
			showColon = !showColon
			fmt.Print("\033[H") // move cursor to top-left
			fmt.Println(clock.RenderTime(time.Now(), showColon))
		}
	}
}
