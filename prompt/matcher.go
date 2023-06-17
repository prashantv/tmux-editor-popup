package prompt

import "strings"

// Matcher is the interface for matchers to detect prompts.
type Matcher interface {
	// Start is called to determine if the current line starts a match, and returns
	// the "parsed" prompt.
	Start(line string) (string, bool)

	// Continued is called only if the previous line started, or continued a match,
	Continued(line string, prevLines []string) (string, bool)
}

// FixedPrefixMatcher uses fixed string prefixes for prompt matching.
type FixedPrefixMatcher struct {
	InitialPrefix   string
	ContinuedPrefix string
}

// Start implements Matcher.Start.
func (m *FixedPrefixMatcher) Start(line string) (string, bool) {
	if !strings.HasPrefix(line, m.InitialPrefix) {
		return "", false
	}

	return line[len(m.InitialPrefix):], true
}

// Continued implements Matcher.Continued.
func (m *FixedPrefixMatcher) Continued(line string, prevLines []string) (string, bool) {
	if !strings.HasPrefix(line, m.ContinuedPrefix) {
		return "", false
	}
	return line[len(m.ContinuedPrefix):], true
}
