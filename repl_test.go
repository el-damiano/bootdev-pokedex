package main

import (
	"fmt"
	"testing"
)

func TestInputClean(t *testing.T) {
	cases := map[string]struct {
		input    string
		expected []string
	}{
		"leading and trailing whitespaces + uppercase letters": {
			input:    "   heLLo    world    ",
			expected: []string{"hello", "world"},
		},
		"empty string": {
			input:    "",
			expected: []string{},
		},
		"newline and whitespaces": {
			input:    " live \n the questions now ",
			expected: []string{"live", "the", "questions", "now"},
		},
		"newline": {
			input:    "new\nline",
			expected: []string{"new", "line"},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			actual := inputClean(c.input)

			if len(c.expected) != len(actual) {
				t.Errorf("Error: expected %v length %d, got %d\n", actual, len(c.expected), len(actual))
			}

			for i := range actual {
				word := actual[i]
				wordExpected := c.expected[i]
				if wordExpected != word {
					t.Errorf("The expected word %s doesn't match with %s\n", wordExpected, word)
				}
			}
		})
	}
}
