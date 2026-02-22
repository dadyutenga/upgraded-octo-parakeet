package audio

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	tempFile string
	once     sync.Once
	playing  bool
	mu       sync.Mutex
)

// extractToTemp writes the embedded MP3 to a temp file (once) and returns the path.
func extractToTemp(fs embed.FS, filename string) (string, error) {
	var extractErr error
	once.Do(func() {
		data, err := fs.ReadFile(filename)
		if err != nil {
			extractErr = fmt.Errorf("read embedded file: %w", err)
			return
		}
		tmp := filepath.Join(os.TempDir(), "azan_clock.mp3")
		if err := os.WriteFile(tmp, data, 0644); err != nil {
			extractErr = fmt.Errorf("write temp file: %w", err)
			return
		}
		tempFile = tmp
	})
	return tempFile, extractErr
}

// Play plays the azan MP3 from the embedded filesystem.
// It's non-blocking and prevents overlapping playback.
func Play(fs embed.FS, filename string) error {
	mu.Lock()
	if playing {
		mu.Unlock()
		return nil
	}
	playing = true
	mu.Unlock()

	path, err := extractToTemp(fs, filename)
	if err != nil {
		mu.Lock()
		playing = false
		mu.Unlock()
		return err
	}

	go func() {
		defer func() {
			mu.Lock()
			playing = false
			mu.Unlock()
		}()
		playFile(path)
	}()

	return nil
}

// Stop stops any currently playing audio.
func Stop() {
	mu.Lock()
	defer mu.Unlock()
	if !playing {
		return
	}
	stopPlayback()
	playing = false
}

// IsPlaying returns whether audio is currently playing.
func IsPlaying() bool {
	mu.Lock()
	defer mu.Unlock()
	return playing
}

// Cleanup removes the temp file.
func Cleanup() {
	if tempFile != "" {
		os.Remove(tempFile)
	}
}

func playFile(path string) {
	switch runtime.GOOS {
	case "windows":
		playWindows(path)
	case "darwin":
		cmd := exec.Command("afplay", path)
		cmd.Run()
	default:
		// Try common Linux players
		for _, player := range []string{"mpv", "ffplay", "aplay", "paplay"} {
			if p, err := exec.LookPath(player); err == nil {
				var cmd *exec.Cmd
				if player == "ffplay" {
					cmd = exec.Command(p, "-nodisp", "-autoexit", path)
				} else {
					cmd = exec.Command(p, path)
				}
				cmd.Run()
				return
			}
		}
	}
}
