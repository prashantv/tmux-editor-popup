package prompt

import "strings"

func cutPrefixes(line string, prefixes ...string) (trimmed string, ok bool) {
	for _, p := range prefixes {
		if trimmed, ok := strings.CutPrefix(line, p); ok {
			return trimmed, ok
		}
	}

	return "", false
}

func cuts(line string, cuts ...string) (left, right string, ok bool) {
	for _, c := range cuts {
		if l, r, ok := strings.Cut(line, c); ok {
			return l, r, ok
		}
	}

	return "", "", false
}
