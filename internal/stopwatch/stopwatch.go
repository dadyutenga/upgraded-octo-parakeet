package stopwatch

import "time"

// Stopwatch tracks elapsed time with start/stop/reset functionality.
type Stopwatch struct {
	startTime time.Time
	elapsed   time.Duration
	running   bool
}

// New creates a new Stopwatch.
func New() *Stopwatch {
	return &Stopwatch{}
}

// Toggle starts or stops the stopwatch.
func (s *Stopwatch) Toggle() {
	if s.running {
		s.elapsed += time.Since(s.startTime)
		s.running = false
	} else {
		s.startTime = time.Now()
		s.running = true
	}
}

// Reset resets the stopwatch to zero.
func (s *Stopwatch) Reset() {
	s.elapsed = 0
	s.running = false
}

// Elapsed returns the current elapsed duration.
func (s *Stopwatch) Elapsed() time.Duration {
	if s.running {
		return s.elapsed + time.Since(s.startTime)
	}
	return s.elapsed
}

// IsRunning returns whether the stopwatch is currently running.
func (s *Stopwatch) IsRunning() bool {
	return s.running
}
