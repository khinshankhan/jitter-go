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
	// base * 2^attempt, clamped to cap
	max := expCap(f.base, f.cap, attempt)

	if max < 0 {
		return 0
	}

	// U[0, max)
	return f.random(max)
}
