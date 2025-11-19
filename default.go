package jitter

// New is the default backoff retry strategy using [NewFull].
func New(cfg Config) (Strategy, error) {
	return NewFull(cfg)
}
