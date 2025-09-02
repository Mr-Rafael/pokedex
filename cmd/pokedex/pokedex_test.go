package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
    cases := []struct {
		input string
		expected []string
	}{
		{
			input: "  hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "HELLO",
			expected: []string{"hello"}
		},
		{
			input: "hello   world. Yo  ",
			expected: []string{"hello", "world.", "yo"}
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Expected %v words in the output, but got %v.", len(c.expected), len(actual))
		}
		for i := range actual {
			if c.expected[i] != actual[i] {
				t.Errorf("Expected word %v to be %v, but got %v", i, c.expected[i], actual[i])
			}
		}
	}
}