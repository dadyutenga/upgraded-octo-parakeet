package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/dadyutenga/upgraded-octo-parakeet/internal/dashboard"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	fmt.Print("\033[2J")    // clear screen
	fmt.Print("\033[?25l")  // hide cursor
	defer fmt.Print("\033[?25h\n")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Initial render
	render()

	for {
		select {
		case <-sig:
			fmt.Println("\nGoodbye!")
			return
		case <-ticker.C:
			render()
		}
	}
}

func render() {
	info := dashboard.Collect(200 * time.Millisecond)
	fmt.Print("\033[H") // move cursor to top
	fmt.Print(dashboard.Render(info))
}
