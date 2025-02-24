package main

import (
	"html/template"
	"testing"
)

func TestEmphasize(t *testing.T) {
	tests := map[string]struct {
		input    string
		col      int
		expected template.HTML
	}{
		"in bounds": {
			input:    "foo",
			col:      2,
			expected: template.HTML(`<span style="position: relative;">f<span class="emphasize">o</span>o</span>`),
		},
		"out of bounds": {
			input:    "foo",
			col:      10,
			expected: template.HTML(`<span style="position: relative;">fo<span class="emphasize">o</span></span>`),
		},
		"out of bounds with HTML": {
			input:    "<span>foo</span>",
			col:      101,
			expected: template.HTML(`<span style="position: relative;">&lt;span&gt;foo&lt;/span<span class="emphasize">&gt;</span></span>`),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := emphasize(test.col, test.input)
			if actual != test.expected {
				t.Errorf("expected: %s, actual: %s", test.expected, actual)
			}
		})
	}
}

func TestEscapeDelimiter(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"empty": {
			input:    "",
			expected: "",
		},
		"no style": {
			input:    "foo",
			expected: "foo",
		},
		"open delimiter": {
			input:    "foo {{",
			expected: "foo \\{\\{",
		},
		"close delimiter": {
			input:    "foo }}",
			expected: "foo \\}\\}",
		},
		"both delimiter": {
			input:    "{{foo}}",
			expected: "\\{\\{foo\\}\\}",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := escapeDelimiter(test.input)
			if actual != test.expected {
				t.Errorf("expected: %s, actual: %s", test.expected, actual)
			}
		})
	}
}

func TestUnescapeDelimiter(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"empty": {
			input:    "",
			expected: "",
		},
		"no style": {
			input:    "foo",
			expected: "foo",
		},
		"open delimiter": {
			input:    "foo \\{\\{",
			expected: "foo {{",
		},
		"close delimiter": {
			input:    "foo \\}\\}",
			expected: "foo }}",
		},
		"both delimiter": {
			input:    "\\{\\{foo\\}\\}",
			expected: "{{foo}}",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := unescapeDelimiter(test.input)
			if actual != test.expected {
				t.Errorf("expected: %s, actual: %s", test.expected, actual)
			}
		})
	}
}
