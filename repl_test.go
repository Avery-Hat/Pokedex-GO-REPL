// repl_test.go
package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{input: "  hello  world  ", expected: []string{"hello", "world"}},
		{input: "Charmander Bulbasaur PIKACHU", expected: []string{"charmander", "bulbasaur", "pikachu"}},
		{input: "\t  foo\tbar   baz  \n", expected: []string{"foo", "bar", "baz"}},
		{input: "", expected: []string{}},
		{input: "   \n\t", expected: []string{}},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Fatalf("for input %q: expected len=%d, got %d (%v)", c.input, len(c.expected), len(actual), actual)
		}

		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("for input %q at index %d: expected %q, got %q",
					c.input, i, c.expected[i], actual[i])
			}
		}
	}
}
