package backoff

import (
	"math/rand"
	"time"
)

// Config contains the minimum and maximum backoff times
type Config struct {
	MinBackoff int
	MaxBackoff int
}

// NewBackoff instantiates a new backoff struct
// min: the minimum backoff time
// max: the maximum backoff time
func NewBackoff(min int, max int) *Config {
	return &Config{
		MinBackoff: min,
		MaxBackoff: max,
	}
}

// SetBackoff sets a randomized backoff timer and returns a pointer to the timer
func (b *Config) SetBackoff() *time.Timer {
	//backoff duration
	backoffMS := rand.Intn(b.MaxBackoff-b.MinBackoff) + b.MinBackoff
	return time.NewTimer(time.Duration(backoffMS) * time.Millisecond)

}

// ResetBackoff stops an existing timer and resets the backoff to a new randomized time
// t: pointer to the timer to be reset
func (b *Config) ResetBackoff(t *time.Timer) *time.Timer {
	if !t.Stop() {
		<-t.C
	}
	backoffMS := rand.Intn(b.MaxBackoff-b.MinBackoff) + b.MinBackoff
	t.Reset(time.Duration(backoffMS) * time.Millisecond)
	return t
}
