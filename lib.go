package jitter

import (
	"errors"
	"strings"
)

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

func getJitterConfigIssues(cfg Config) []string {
	var issues []string
	if cfg.Base <= 0 {
		issues = append(issues, "Base must be > 0")
	}
	if cfg.Cap <= 0 {
		issues = append(issues, "Cap must be > 0")
	}
	if cfg.Random == nil {
		issues = append(issues, "Random function must be provided")
	}
	return issues
}

// ErrInvalidConfig is returned when a jitter config is invalid.
var ErrInvalidConfig = errors.New("jitter: invalid config")

// ConfigError describes an invalid jitter config.
// Use errors.Is(err, ErrInvalidConfig) to check for this error.
type ConfigError struct {
	// List of details why given config is invalid.
	Issues []string
}

func (e *ConfigError) Error() string {
	return "jitter: invalid config: " + strings.Join(e.Issues, ", ")
}

func (e *ConfigError) Is(target error) bool {
	return target == ErrInvalidConfig
}
