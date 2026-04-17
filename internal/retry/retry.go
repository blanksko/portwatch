// Package retry provides simple retry logic with configurable attempts and backoff.
package retry

import (
	"errors"
	"time"
)

// ErrMaxAttempts is returned when all retry attempts are exhausted.
var ErrMaxAttempts = errors.New("max retry attempts reached")

// Config holds retry configuration.
type Config struct {
	MaxAttempts int
	Delay       time.Duration
}

// Default returns a Config with sensible defaults.
func Default() Config {
	return Config{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
	}
}

// Doer executes a function with retry logic based on the provided Config.
type Doer struct {
	cfg   Config
	sleep func(time.Duration)
}

// New creates a new Doer with the given Config.
func New(cfg Config) *Doer {
	return &Doer{
		cfg:   cfg,
		sleep: time.Sleep,
	}
}

// Do calls fn up to MaxAttempts times, sleeping Delay between attempts.
// Returns nil on first success, or ErrMaxAttempts if all attempts fail.
func (d *Doer) Do(fn func() error) error {
	if d.cfg.MaxAttempts <= 0 {
		return ErrMaxAttempts
	}
	var last error
	for i := 0; i < d.cfg.MaxAttempts; i++ {
		if err := fn(); err == nil {
			return nil
		} else {
			last = err
		}
		if i < d.cfg.MaxAttempts-1 {
			d.sleep(d.cfg.Delay)
		}
	}
	_ = last
	return ErrMaxAttempts
}
