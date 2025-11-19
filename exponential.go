package jitter

// exponential implements backoff = min(cap, base * 2^attempt), return backoff.
type exponential struct {
	base int64
	cap  int64
}

func getExponentialConfigIssues(cfg Config) []string {
	var issues []string
	if cfg.Base <= 0 {
		issues = append(issues, "Base must be > 0")
	}
	if cfg.Cap <= 0 {
		issues = append(issues, "Cap must be > 0")
	}
	return issues
}

// NewExponential returns a Strategy using plain exponential backoff with no jitter.
// Returns an error if cfg.Base <= 0 or cfg.Cap <= 0.
func NewExponential(cfg Config) (Strategy, error) {
	if issues := getExponentialConfigIssues(cfg); len(issues) > 0 {
		return nil, &ConfigError{Issues: issues}
	}

	return &exponential{
		base: cfg.Base,
		cap:  cfg.Cap,
	}, nil
}

func (e *exponential) Next(attempt int) int64 {
	return expCap(e.base, e.cap, attempt)
}
