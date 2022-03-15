package util_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/outblocks/outblocks-plugin-go/util"
)

func TestExpand(t *testing.T) {
	tests := []struct {
		content  string
		vars     map[string]interface{}
		expected string
	}{
		{
			content:  "abc: ${var.abc}",
			vars:     map[string]interface{}{"var": map[string]interface{}{"abc": 1}},
			expected: `abc: 1`,
		},
		{
			content:  "abc: ${var}\na: true",
			vars:     map[string]interface{}{"var": map[string]interface{}{"abc": 1, "cba": []interface{}{"1", 2, true, 1.5}}},
			expected: "abc: map[abc:1 cba:[1 2 true 1.5]]\na: true",
		},
		{
			content:  "abc: I am ${var}",
			vars:     map[string]interface{}{"var": "cornholio"},
			expected: "abc: I am cornholio",
		},
		{
			content:  "abc: ${var.nested.val.y}",
			vars:     map[string]interface{}{"var": map[string]interface{}{"nested": map[string]interface{}{"val": map[string]interface{}{"y": 1}}}},
			expected: `abc: 1`,
		},
		{
			content:  "val: ${var.base_url}/func1",
			vars:     map[string]interface{}{"var": map[string]interface{}{"base_url": "test"}},
			expected: "val: test/func1",
		},
		{
			content:  `val: "*.${var.base_url}"`,
			vars:     map[string]interface{}{"var": map[string]interface{}{"base_url": "test"}},
			expected: `val: "*.test"`,
		},
	}

	for _, test := range tests {
		out, params, err := util.NewBaseVarEvaluator(test.vars).ExpandRaw([]byte(test.content))
		outF := fmt.Sprintf(string(out), params...)

		if err != nil {
			t.Fatalf(`Expand(%q) for %q = (%q, %q), expected non error`, test.content, test.vars, outF, err)
		}
		if test.expected != outF {
			t.Fatalf(`Expand(%q) for %q = (%q, %q), expected: %q`, test.content, test.vars, outF, err, test.expected)
		}
	}
}

func TestExpand_Invalid(t *testing.T) {
	tests := []struct {
		content  string
		vars     map[string]interface{}
		expected string
	}{
		{
			content:  "\nabc ${}",
			vars:     nil,
			expected: "[2:5] empty expansion found",
		},
		{
			content:  "abc ${var.abc}",
			vars:     nil,
			expected: "[1:5] expansion value for 'var.abc' could not be evaluated",
		},
	}

	for _, test := range tests {
		out, params, err := util.NewBaseVarEvaluator(test.vars).ExpandRaw([]byte(test.content))
		outF := fmt.Sprintf(string(out), params...)

		if err == nil {
			t.Fatalf(`Expand(%q) for %q = (%q, %q), expected error`, test.content, test.vars, outF, err)
		}
		if !strings.Contains(err.Error(), test.expected) {
			t.Fatalf(`Expand(%q) for %q = (%q, %q), expected error: %q`, test.content, test.vars, outF, err, test.expected)
		}
	}
}
