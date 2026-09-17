package lab

import "time"

const (
	DefaultDelay     = 400 * time.Millisecond
	MaxDelay         = 5 * time.Second
	DefaultCPU       = 5 * time.Second
	MaxCPU           = 10 * time.Second
	DefaultBurst     = 80 * time.Millisecond
	MaxBurst         = 200 * time.Millisecond
	DefaultMemoryMiB = 64
	MaxMemoryMiB     = 128
)

// DurationFromMS convierte milisegundos del JSON a Duration, con default y tope.
func DurationFromMS(ms int, fallback, max time.Duration) time.Duration {
	if ms <= 0 {
		return fallback
	}
	d := time.Duration(ms) * time.Millisecond
	if d > max {
		return max
	}
	return d
}

func ClampInt(n, fallback, max int) int {
	if n <= 0 {
		return fallback
	}
	if n > max {
		return max
	}
	return n
}
