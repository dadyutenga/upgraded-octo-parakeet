package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	azanFS "github.com/dadyutenga/upgraded-octo-parakeet/cmd/audio"
	"github.com/dadyutenga/upgraded-octo-parakeet/internal/audio"
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
	nav += "  |  a: azan on/off  |  s: stop audio"
	nav += "\n\n"
	return nav
}

func main() {
	// Cap memory at 55 MB
	debug.SetMemoryLimit(55 * 1024 * 1024)

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
	azanTriggered := make(map[string]bool) // track which prayers already triggered azan today
	azanEnabled := true
	var lastDateStr string

	defer audio.Cleanup()

	fmt.Print("\033[2J\033[H\033[?25l") // clear screen, hide cursor
	render(currentMode, showColon, sw, azanEnabled)

	for {
		select {
		case <-sig:
			audio.Stop()
			fmt.Print("\033[?25h") // show cursor
			fmt.Println("\nGoodbye!")
			return
		case key := <-keysCh:
			switch key {
			case 'q', 'Q':
				audio.Stop()
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
			case 'a', 'A':
				azanEnabled = !azanEnabled
			case 's', 'S':
				audio.Stop()
			}
			fmt.Print("\033[2J\033[H")
			render(currentMode, showColon, sw, azanEnabled)
		case <-ticker.C:
			blinkTick++
			if blinkTick%5 == 0 { // blink every 500ms
				showColon = !showColon
			}

			// Check for azan trigger
			now := time.Now()
			dateStr := now.Format("2006-01-02")
			if dateStr != lastDateStr {
				azanTriggered = make(map[string]bool)
				lastDateStr = dateStr
			}
			if azanEnabled {
				checkAzan(now, azanTriggered)
			}

			fmt.Print("\033[H")
			render(currentMode, showColon, sw, azanEnabled)
		}
	}
}

// checkAzan triggers the azan if we're within 1 minute of a prayer time.
func checkAzan(now time.Time, triggered map[string]bool) {
	prayers, err := prayer.GetPrayerTimes(now)
	if err != nil || len(prayers) == 0 {
		return
	}
	// Only trigger for actual prayer times (skip Sunrise)
	for _, p := range prayers {
		if p.Name == "Sunrise" {
			continue
		}
		if triggered[p.Name] {
			continue
		}
		diff := now.Sub(p.Time)
		if diff >= 0 && diff < time.Minute {
			triggered[p.Name] = true
			audio.Play(azanFS.FS, azanFS.AzanFile)
			return
		}
	}
}

func render(mode int, showColon bool, sw *stopwatch.Stopwatch, azanEnabled bool) {
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
		if azanEnabled {
			fmt.Println("  \033[32m🔊 Azan: ON\033[0m")
		} else {
			fmt.Println("  \033[90m🔇 Azan: OFF\033[0m")
		}
		if audio.IsPlaying() {
			fmt.Println("  \033[33m♪ Playing azan... (press 's' to stop)\033[0m")
		}
	}
}

