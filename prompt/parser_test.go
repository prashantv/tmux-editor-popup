package prompt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	m := &FixedPrefixMatcher{
		InitialPrefix:   "> ",
		ContinuedPrefix: "> ",
	}
	p := &Parser{Matchers: []Matcher{m}}

	tests := []struct {
		msg string
		in  string

		want    string
		isFirst bool
		isLast  bool
		wantErr string
	}{
		{
			msg:     "empty",
			in:      ``,
			wantErr: errNotFound.Error(),
		},
		{
			msg:     "single line no prompt",
			in:      `testing`,
			wantErr: errNotFound.Error(),
		},
		{
			msg:     "single line with prompt",
			in:      `> testing`,
			want:    "testing",
			isFirst: true,
			isLast:  true,
		},
		{
			msg: "multiple lines all prompt",
			in: `
> hello
> world
`,
			want:    "hello\nworld",
			isFirst: true,
			isLast:  true,
		},
		{
			msg: "ignore empty prompt at the end",
			in: `
> echo 'hello
> there';
hello
there
> 
`,
			want:    "echo 'hello\nthere';",
			isFirst: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			r := strings.NewReader(
				strings.TrimPrefix(tt.in, "\n"), // multi-line strings in tests have an empty first line.
			)
			prompt, err := p.Find(r)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, strings.Join(prompt.Parsed, "\n"))
			assert.Equal(t, tt.isFirst, prompt.IsFirst, "IsFirst")
			assert.Equal(t, tt.isLast, prompt.IsLast, "IsLast")
		})
	}
}
