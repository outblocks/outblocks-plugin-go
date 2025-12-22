package fields_test

import (
	"testing"

	"github.com/outblocks/outblocks-plugin-go/registry/fields"
)

func TestExpandWithPercent(t *testing.T) {
	tests := []struct {
		content  string
		vars     map[string]any
		expected string
	}{
		{
			content:  "abc: %{var.abc}",
			vars:     map[string]any{"var": map[string]any{"abc": 1}},
			expected: `abc: 1`,
		},
		{
			content:  "abc: %d%{var.abc}",
			vars:     map[string]any{"var": map[string]any{"abc": 1}},
			expected: `abc: %d1`,
		},
		{
			content:  "abc: %",
			expected: `abc: %`,
		},
		{
			content:  "abc: %{}",
			expected: `abc: %{}`,
		},
		{
			content:  "abc: %{",
			expected: `abc: %{`,
		},
		{
			content:  "abc: %{abra.kadabra!}",
			expected: `abc: %{abra.kadabra!}`,
		},
	}

	for _, test := range tests {
		out, err := fields.NewFieldVarEvaluator(test.vars).Expand(test.content)
		if err != nil {
			t.Fatalf(`Expand(%q) for %q = (%q, %q), expected non error`, test.content, test.vars, out, err)
		}

		outF := out.Any()
		if test.expected != outF {
			t.Fatalf(`Expand(%q) for %q = (%q, %q), expected: %q`, test.content, test.vars, outF, err, test.expected)
		}
	}
}
