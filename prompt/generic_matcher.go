package prompt

import "fmt"

type genericMatcher struct {
	start     lineMatcher
	continued lineMatcher
}

func newGenericMatcher(cfg *Config) (Matcher, error) {
	start, err := newLineMatcher(cfg.Start)
	if err != nil {
		return nil, fmt.Errorf("parse start: %v", err)
	}

	continued, err := newLineMatcher(cfg.Continued)
	if err != nil {
		return nil, fmt.Errorf("parse continued: %v", err)
	}

	return genericMatcher{start, continued}, nil
}

func (m genericMatcher) Start(line string) (string, bool) {
	return m.start.match(line)
}

func (m genericMatcher) Continued(line string, prevLines []string) (string, bool) {
	return m.continued.match(line)
}
