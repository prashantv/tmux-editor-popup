package prompt

var defaultMatchers = []Matcher{
	&FixedPrefixMatcher{"> ", "quote> "},
}

// NewDefault returns a parser with the default matchers.
func NewDefault() *Parser {
	return &Parser{
		Matchers: defaultMatchers,
	}
}
