package jitter

// fullJitter implements backoff = min(cap, base * 2^attempt), return U[0, backoff].
type fullJitter struct {
	base   int64
	cap    int64
	random RandomFunc
}

// NewFull returns a Strategy using the "full jitter" algorithm.
// It panics if cfg.Base <= 0, cfg.Cap <= 0, or cfg.Random is nil.
func NewFull(cfg Config) Strategy {
	if cfg.Base <= 0 || cfg.Cap <= 0 || cfg.Random == nil {
		panic("jitter: invalid full jitter config")
	}

	return &fullJitter{
		base:   cfg.Base,
		cap:    cfg.Cap,
		random: cfg.Random,
	}
}

func (f *fullJitter) Next(attempt int) int64 {
	if attempt < 0 {
		return 0
	}

	// base * 2^attempt, clamped to cap
	// TODO: circle back and handle potential overflow
	max := f.base << attempt
	if max <= 0 || max > f.cap {
		max = f.cap
	}

	// U[0, max)
	return f.random(max)
}
