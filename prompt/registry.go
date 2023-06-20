package prompt

import (
	"fmt"
)

type Config struct {
	Start     MatcherConfig
	Continued MatcherConfig
}

type MatcherConfig struct {
	// Fixed is a matcher that uses a fixed string to identify matches.
	Fixed *FixedConfig
}

type FixedConfig struct {
	Separator string
	Before    *StringMatchConfig // optional, match anything before the separator.
	After     *StringMatchConfig // optional, match anything after the separator.
}

type StringMatchConfig struct {
	Either      []*StringMatchConfig
	All         []*StringMatchConfig
	Regex       string
	Const       string
	ConstPrefix string
	ConstSuffix string
}

var registryOrder []string

func reg(s string) string {
	registryOrder = append(registryOrder, s)
	return s
}

var registry = map[string]*Config{
	reg("zsh-escape-newline"): {
		Start: MatcherConfig{
			Fixed: &FixedConfig{
				Separator: "> ",
				Before: &StringMatchConfig{
					Const: "",
				},
				After: &StringMatchConfig{
					ConstSuffix: `\`,
				},
			},
		},
		Continued: MatcherConfig{
			Fixed: &FixedConfig{
				Separator: "> ",
				Before: &StringMatchConfig{
					Const: "",
				},
			},
		},
	},
	reg("zsh-prashant-cont1"): {
		Start: MatcherConfig{
			Fixed: &FixedConfig{
				Separator: "> ",
				Before: &StringMatchConfig{
					Const: "",
				},
			},
		},
		Continued: MatcherConfig{
			Fixed: &FixedConfig{
				Separator: "> ",
				Before: &StringMatchConfig{
					Either: []*StringMatchConfig{
						{
							Const: "quote",
						},
						{
							Const: "pipe",
						},
					},
				},
			},
		},
	},
	reg("mysql"): {
		Start: MatcherConfig{
			Fixed: &FixedConfig{
				Separator: "mysql> ",
				Before: &StringMatchConfig{
					Const: "",
				},
			},
		},
		Continued: MatcherConfig{
			Fixed: &FixedConfig{
				Separator: "> ",
				Before: &StringMatchConfig{
					Either: []*StringMatchConfig{
						{ConstSuffix: "-"},
						{ConstSuffix: "'"},
						{ConstSuffix: `"`},
					},
				},
			},
		},
	},
}

// NewMatchers creates a list of matchers with specific names.
func NewMatchers(names []string) (*Parser, error) {
	matchers := make([]Matcher, 0, len(names))
	for _, n := range names {
		cfg, ok := registry[n]
		if !ok {
			return nil, fmt.Errorf("unknown config %q, must be one of [TODO]", n)
		}

		m, err := newGenericMatcher(cfg)
		if err != nil {
			return nil, fmt.Errorf("invalid config %q: %v", n, err)
		}

		matchers = append(matchers, m)
	}

	return &Parser{
		Matchers: matchers,
	}, nil
}

func NewAll() (*Parser, error) {
	return NewMatchers(registryOrder)
}
