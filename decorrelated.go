package jitter

// decorrelatedJitter implements backoff = min(cap, U[base, prev * 3]), return backoff,
// where prev is the previous sleep value.
type decorrelatedJitter struct {
	base   int64
	cap    int64
	random RandomFunc

	sleep int64 // last computed delay
}

// NewDecorrelated returns a Strategy using the "decorrelated jitter" algorithm.
// It panics if cfg.Base <= 0, cfg.Cap <= 0, or cfg.Random is nil.
//
// The returned Strategy is stateful and not safe for concurrent use from
// multiple goroutines. Callers should create a new Strategy for each
// independent retry loop.
func NewDecorrelated(cfg Config) Strategy {
	if cfg.Base <= 0 || cfg.Cap <= 0 || cfg.Random == nil {
		panic("jitter: invalid decorrelated jitter config")
	}

	return &decorrelatedJitter{
		base:   cfg.Base,
		cap:    cfg.Cap,
		random: cfg.Random,
		sleep:  cfg.Base,
	}
}

// The attempt parameter is used only as a reset signal:
// - attempt < 1: reset internal state to base and start a new sequence
// - attempt  >= 1: continue from the previous sleep value
func (d *decorrelatedJitter) Next(attempt int) int64 {
	// reset on first attempt
	if attempt < 1 {
		d.sleep = d.base
	}

	// safeguard, ensure we have a valid previous sleep value
	if d.sleep <= 0 {
		d.sleep = d.base
	}

	min := d.base
	max := d.sleep * 3

	if max < min {
		max = min
	}
	if max > d.cap {
		max = d.cap
	}
	if max <= min {
		// no room to jitter
		d.sleep = min
		return d.sleep
	}

	span := max - min
	if span <= 0 {
		d.sleep = min
		return d.sleep
	}

	// U[base, prev*3] clamped to cap == base + U[0, span]
	next := min + d.random(span)
	if next <= 0 {
		next = min
	}
	if next > d.cap {
		next = d.cap
	}

	d.sleep = next
	return d.sleep
}
