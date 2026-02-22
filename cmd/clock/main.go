package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/dadyutenga/upgraded-octo-parakeet/internal/clock"
	"github.com/dadyutenga/upgraded-octo-parakeet/internal/prayer"
	"github.com/dadyutenga/upgraded-octo-parakeet/internal/stopwatch"
)

const (
	ModeClock     = 0
	ModeStopwatch = 1
	ModePrayer    = 2
	ModeCount     = 3
)

var modeNames = [ModeCount]string{"🕐 Clock", "⏱  Stopwatch", "🕌 Prayer Times"}

func renderNav(currentMode int) string {
	nav := "\033[1m"
	for i, name := range modeNames {
		if i == currentMode {
			nav += fmt.Sprintf(" \033[7m %s \033[27m ", name) // inverted for selected
		} else {
			nav += fmt.Sprintf(" \033[90m%s\033[0m\033[1m ", name)
		}
		if i < ModeCount-1 {
			nav += "│"
		}
	}
	nav += "\033[0m\n"
	nav += "  ← → switch modes"
	if currentMode == ModeStopwatch {
		nav += "  |  SPACE: start/stop  |  r: reset"
	}
	nav += "\n\n"
	return nav
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Set up raw terminal input
	keysCh := make(chan byte, 10)
	go readKeys(keysCh)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	showColon := true
	blinkTick := 0
	currentMode := ModeClock
	sw := stopwatch.New()

	fmt.Print("\033[2J\033[H\033[?25l") // clear screen, hide cursor
	render(currentMode, showColon, sw)

	for {
		select {
		case <-sig:
			fmt.Print("\033[?25h") // show cursor
			fmt.Println("\nGoodbye!")
			return
		case key := <-keysCh:
			switch key {
			case 'q', 'Q':
				fmt.Print("\033[?25h")
				fmt.Println("\nGoodbye!")
				return
			case 'D': // left arrow (escape seq handled in readKeys)
				currentMode = (currentMode - 1 + ModeCount) % ModeCount
			case 'C': // right arrow
				currentMode = (currentMode + 1) % ModeCount
			case '1':
				currentMode = ModeClock
			case '2':
				currentMode = ModeStopwatch
			case '3':
				currentMode = ModePrayer
			case ' ':
				if currentMode == ModeStopwatch {
					sw.Toggle()
				}
			case 'r', 'R':
				if currentMode == ModeStopwatch {
					sw.Reset()
				}
			}
			fmt.Print("\033[2J\033[H")
			render(currentMode, showColon, sw)
		case <-ticker.C:
			blinkTick++
			if blinkTick%5 == 0 { // blink every 500ms
				showColon = !showColon
			}
			fmt.Print("\033[H")
			render(currentMode, showColon, sw)
		}
	}
}

func render(mode int, showColon bool, sw *stopwatch.Stopwatch) {
	fmt.Print(renderNav(mode))

	switch mode {
	case ModeClock:
		fmt.Println(clock.RenderTime(time.Now(), showColon))
		fmt.Printf("\n  %s\n", time.Now().Format("Monday, 02 January 2006"))
	case ModeStopwatch:
		elapsed := sw.Elapsed()
		fmt.Println(clock.RenderDuration(elapsed, showColon))
		ms := elapsed.Milliseconds() % 1000
		status := "\033[31m⏸ Stopped\033[0m"
		if sw.IsRunning() {
			status = "\033[32m● Running\033[0m"
		}
		fmt.Printf("\n  .%03d   %s\n", ms, status)
	case ModePrayer:
		now := time.Now()
		prayers, err := prayer.GetPrayerTimes(now)
		fmt.Println(prayer.Render(prayers, now, err))
	}
}

