package prompt

import (
	"fmt"
	"strings"
)

type fixedMatcher struct {
	separator     string
	beforeMatcher stringMatcher
	afterMatcher  stringMatcher
}

func newFixedMatcher(cfg *FixedConfig) (lineMatcher, error) {
	beforeMatcher, err := newStringMatcher(cfg.Before)
	if err != nil {
		return nil, fmt.Errorf("fixed matcher parse before: %v", err)
	}

	afterMatcher, err := newStringMatcher(cfg.After)
	if err != nil {
		return nil, fmt.Errorf("fied matcher parse after: %v", err)
	}

	return fixedMatcher{cfg.Separator, beforeMatcher, afterMatcher}, nil
}

func (m fixedMatcher) match(line string) (string, bool) {
	before, contents, ok := strings.Cut(line, m.separator)
	if !ok {
		return "", false
	}

	if !m.beforeMatcher.match(before) {
		return "", false
	}

	if !m.afterMatcher.match(contents) {
		return "", false
	}

	return contents, true
}
