package jitter

// equalJitter implements backoff = min(cap, base * 2^attempt), return backoff/2 + U[0, backoff/2].
type equalJitter struct {
	base   int64
	cap    int64
	random RandomFunc
}

// NewEqual returns a Strategy using the "equal jitter" algorithm.
// It panics if cfg.Base <= 0, cfg.Cap <= 0, or cfg.Random is nil.
func NewEqual(cfg Config) Strategy {
	if cfg.Base <= 0 || cfg.Cap <= 0 || cfg.Random == nil {
		panic("jitter: invalid equal jitter config")
	}

	return &equalJitter{
		base:   cfg.Base,
		cap:    cfg.Cap,
		random: cfg.Random,
	}
}

func (e *equalJitter) Next(attempt int) int64 {
	backoff := expCap(e.base, e.cap, attempt)
	if backoff <= 0 {
		return 0
	}

	half := backoff / 2
	span := backoff - half // >= 0

	if span <= 0 {
		// no room to jitter
		return half
	}

	// backoff/2 + U[0, backoff/2]
	return half + e.random(span)
}
