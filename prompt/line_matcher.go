package prompt

import "fmt"

type lineMatcher interface {
	match(line string) (string, bool)
}

type neverMatch struct{}

func (neverMatch) match(line string) (string, bool) {
	return "", false
}

func newLineMatcher(cfg MatcherConfig) (lineMatcher, error) {
	if cfg.Fixed == nil {
		return neverMatch{}, nil
	}

	m, err := newFixedMatcher(cfg.Fixed)
	if err != nil {
		return nil, fmt.Errorf("create fixed matcher: %v", err)
	}

	return m, nil
}
