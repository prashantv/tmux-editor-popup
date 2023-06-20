package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode"
)

var errNotFound = errors.New("no prompt found")

type Prompt struct {
	Lines   []string
	Parsed  []string
	IsFirst bool
	IsLast  bool
}

type Parser struct {
	Matchers []Matcher
}

func (p Parser) Find(r io.Reader) (Prompt, error) {
	m, totalLines, err := p.parse(r)
	if err != nil {
		return Prompt{}, err
	}
	if m == nil {
		return Prompt{}, errNotFound
	}

	return Prompt{
		Lines:   m.lines,
		Parsed:  m.parsed,
		IsFirst: m.start == 0,
		IsLast:  m.last == totalLines-1, // TODO: missing case where there's a lot of blank lines.
	}, nil
}

type match struct {
	matcher Matcher
	start   int
	last    int
	lines   []string
	parsed  []string
}

func (m *match) isEmpty() bool {
	for _, p := range m.parsed {
		for _, c := range p {
			if !unicode.IsSpace(c) {
				return false
			}
		}
	}
	return true
}

func (p *Parser) parse(r io.Reader) (_ *match, lines int, _ error) {
	var (
		last    *match
		current *match
	)

	var lineIdx int
	endCurrent := func() {
		current.last = lineIdx - 1

		// Only add the match if it's non-empty
		if !current.isEmpty() {
			last = current
		}
		current = nil
	}

	scanner := bufio.NewScanner(r)
	for lineIdx = 0; scanner.Scan(); lineIdx++ {
		line := scanner.Text()

		if current != nil {
			parsed, ok := current.matcher.Continued(line, current.lines)
			fmt.Printf("  and continued? %v on %q\n", ok, line)
			if ok {
				current.lines = append(current.lines, line)
				current.parsed = append(current.parsed, parsed)
				continue
			}

			// No match on a continued line, so end the match.
			endCurrent()
		}

		// Check for a new match start.
		m, parsed, ok := p.checkStartMatch(line)
		if !ok {
			continue
		}

		current = &match{
			matcher: m,
			start:   lineIdx,
			lines:   []string{line},
			parsed:  []string{parsed},
		}
	}

	if current != nil {
		endCurrent()
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}

	return last, lineIdx, nil
}

func (p *Parser) checkStartMatch(line string) (Matcher, string, bool) {
	for i, m := range p.Matchers {
		parsed, ok := m.Start(line)
		if !ok {
			continue
		}

		fmt.Printf("matcher %d started on %q\n", i, line)
		return m, parsed, true
	}

	return nil, "", false
}
