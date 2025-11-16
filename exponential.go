package jitter

// exponential implements backoff = min(cap, base * 2^attempt), return backoff.
type exponential struct {
	base int64
	cap  int64
}

// NewExponential returns a Strategy using plain exponential backoff with no jitter.
// It panics if cfg.Base <= 0 or cfg.Cap <= 0. cfg.Random is ignored.
func NewExponential(cfg Config) Strategy {
	if cfg.Base <= 0 || cfg.Cap <= 0 {
		panic("jitter: invalid exponential config")
	}

	return &exponential{
		base: cfg.Base,
		cap:  cfg.Cap,
	}
}

func (e *exponential) Next(attempt int) int64 {
	return expCap(e.base, e.cap, attempt)
}
