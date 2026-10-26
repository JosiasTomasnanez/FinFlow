package lab

import (
	"runtime"
	"time"
)

// Spin ocupa CPU de verdad (no sleep). El cgroup lo cuenta como uso y, si
// el pico supera el quota CFS, como throttle.
func Spin(d time.Duration) {
	deadline := time.Now().Add(d)
	var acc uint64
	for time.Now().Before(deadline) {
		acc++
	}
	// KeepAlive evita que el compilador descarte el trabajo del loop.
	runtime.KeepAlive(acc)
}
