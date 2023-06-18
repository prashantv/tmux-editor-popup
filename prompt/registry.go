package prompt

var defaultMatchers = []Matcher{
	shellMatcher{},
}

// NewDefault returns a parser with the default matchers.
func NewDefault() *Parser {
	return &Parser{
		Matchers: defaultMatchers,
	}
}
