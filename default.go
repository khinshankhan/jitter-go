package jitter

// New is the default backoff retry strategy using [NewFull].
func New(cfg Config) Strategy {
	return NewFull(cfg)
}
