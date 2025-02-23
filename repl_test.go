package main

import (
	"testing"
)

func TestInputClean(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "   heLLo    world    ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    " live \n the questions now ",
			expected: []string{"live", "the", "questions", "now"},
		},
		{
			input:    "sentinel string",
			expected: []string{"sentinel", "string"},
		},
		{
			input:    "new\nline",
			expected: []string{"new", "line"},
		},
	}

	for _, c := range cases {
		actual := inputClean(c.input)

		if len(c.expected) != len(actual) {
			t.Errorf("Error: expected %v length %d, got %d\n", actual, len(c.expected), len(actual))
		}

		for i := range actual {
			word := actual[i]
			wordExpected := c.expected[i]
			// t.Logf("comparing %s to %s\n", word, wordExpected)
			if wordExpected != word {
				t.Errorf("The expected word %s doesn't match with %s\n", wordExpected, word)
			}
		}

	}
}
