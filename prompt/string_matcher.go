package prompt

import (
	"fmt"
	"regexp"
	"strings"
)

type stringMatcher interface {
	match(s string) bool
}

type allMatcher struct {
	matchers []stringMatcher
}

func newAllMatcher(all []*StringMatchConfig) (stringMatcher, error) {
	matchers := make([]stringMatcher, 0, len(all))
	for _, c := range all {
		m, err := newStringMatcher(c)
		if err != nil {
			return nil, fmt.Errorf("parse %+v: %v", c, err)
		}
		matchers = append(matchers, m)
	}
	return allMatcher{matchers}, nil
}

func (m allMatcher) match(s string) bool {
	for _, m := range m.matchers {
		if !m.match(s) {
			return false
		}
	}
	return true
}

type orMatcher struct {
	matchers []stringMatcher
}

func newOrMatcher(all []*StringMatchConfig) (stringMatcher, error) {
	matchers := make([]stringMatcher, 0, len(all))
	for _, c := range all {
		m, err := newStringMatcher(c)
		if err != nil {
			return nil, fmt.Errorf("parse %+v: %v", c, err)
		}
		matchers = append(matchers, m)
	}
	return orMatcher{matchers}, nil
}

func (m orMatcher) match(s string) bool {
	for _, m := range m.matchers {
		if m.match(s) {
			return true
		}
	}
	return false
}

type stringMatcherFunc func(s string) bool

func (m stringMatcherFunc) match(s string) bool {
	return m(s)
}

type regexMatcher struct {
	r *regexp.Regexp
}

func newRegexMatcher(regexStr string) (stringMatcher, error) {
	r, err := regexp.Compile(regexStr)
	if err != nil {
		return nil, fmt.Errorf("parse regex: %v", err)
	}

	return regexMatcher{r}, nil
}

func (m regexMatcher) match(s string) bool {
	return m.r.MatchString(s)
}

func newStringMatcher(cfg *StringMatchConfig) (stringMatcher, error) {
	if cfg == nil {
		// We want a matcher that always returns true, which allMatcher with no matchers does.
		return allMatcher{}, nil
	}
	if len(cfg.All) > 0 {
		return newAllMatcher(cfg.All)
	}
	if len(cfg.Either) > 0 {
		return newOrMatcher(cfg.Either)
	}
	if len(cfg.Regex) > 0 {
		return newRegexMatcher(cfg.Regex)
	}
	if len(cfg.ConstPrefix) > 0 {
		return stringMatcherFunc(func(s string) bool {
			return strings.HasPrefix(s, cfg.ConstPrefix)
		}), nil
	}
	if len(cfg.ConstSuffix) > 0 {
		return stringMatcherFunc(func(s string) bool {
			return strings.HasSuffix(s, cfg.ConstSuffix)
		}), nil
	}

	return stringMatcherFunc(func(s string) bool {
		return s == cfg.Const
	}), nil
}
