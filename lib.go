package jitter

// Strategy returns a non-negative integer "delay" for a given attempt.
// Units are arbitrary (eg milliseconds), the caller decides.
type Strategy interface {
	Next(attempt int) int64
}

// RandomFunc is the minimal RNG contract we need.
// Implementations must handle max > 0; callers should not pass max <= 0.
// This is to remove dependency on `math/rand` in this package.
type RandomFunc func(max int64) int64

// Config configures a jitter strategy. Base and Cap must be > 0.
type Config struct {
	Base   int64      // base delay, eg 100 (ms, but time unit is caller's responsibility)
	Cap    int64      // maximum delay (should be the same unit as Base)
	Random RandomFunc // random function to produce U[0, max)
}

// expCap computes base * 2^attempt, clamped to cap.
// Returns 0 if base or cap are invalid or attempt < 0.
func expCap(base, cap int64, attempt int) int64 {
	if attempt < 0 || base <= 0 || cap <= 0 {
		return 0
	}

	// grow exponentially, but stop before we overflow or exceed cap
	max := base
	for i := 0; i < attempt && max < cap; i++ {
		// if doubling would exceed cap, clamp to cap
		if max > cap/2 {
			max = cap
			break
		}
		max <<= 1
	}

	if max > cap {
		return cap
	}
	return max
}
