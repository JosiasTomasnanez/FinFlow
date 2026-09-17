package lab

import (
	"testing"
	"time"
)

func TestDurationFromMS(t *testing.T) {
	if got := DurationFromMS(0, DefaultDelay, MaxDelay); got != DefaultDelay {
		t.Fatalf("default: got %v", got)
	}
	if got := DurationFromMS(50, DefaultDelay, MaxDelay); got != 50*time.Millisecond {
		t.Fatalf("value: got %v", got)
	}
	if got := DurationFromMS(99999, DefaultCPU, MaxCPU); got != MaxCPU {
		t.Fatalf("clamp: got %v", got)
	}
}

func TestHoldRetainAndRelease(t *testing.T) {
	var h Hold
	if got := h.Retain(2); got != 2 {
		t.Fatalf("retain: %d", got)
	}
	if h.MiB() != 2 {
		t.Fatalf("held: %d", h.MiB())
	}
	h.Release()
	if h.MiB() != 0 {
		t.Fatalf("after release: %d", h.MiB())
	}
}

func TestHoldRetainReusesOrDropsBeforeGrow(t *testing.T) {
	var h Hold
	if h.Retain(3) != 3 || h.MiB() != 3 {
		t.Fatalf("first hold %d", h.MiB())
	}
	if h.Retain(1) != 1 || h.MiB() != 1 {
		t.Fatalf("shrink reuse %d", h.MiB())
	}
	if h.Retain(4) != 4 || h.MiB() != 4 {
		t.Fatalf("grow after drop %d", h.MiB())
	}
}

func TestAllowedErrorStatus(t *testing.T) {
	if AllowedErrorStatus(404) != 500 {
		t.Fatal("4xx must not pass")
	}
	if AllowedErrorStatus(503) != 503 {
		t.Fatal("503")
	}
}
