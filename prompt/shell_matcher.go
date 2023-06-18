package prompt

type shellMatcher struct{}

var _ Matcher = shellMatcher{}

func (shellMatcher) Start(line string) (string, bool) {
	if trimmed, ok := cutPrefixes(line,
		"$ ",
		"% ",
	); ok {
		return trimmed, true
	}

	if _, right, ok := cuts(line,
		"$ ",
		"% "); ok {
		return right, ok
	}

	return "", false
}

func (shellMatcher) Continued(line string, prevLines []string) (string, bool) {
	if trimmed, ok := cutPrefixes(line,
		"> ",
	); ok {
		return trimmed, true
	}

	if _, trimmed, ok := cuts(line,
		"> ",
	); ok {
		return trimmed, true
	}

	return "", false
}
