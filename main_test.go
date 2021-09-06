package main

import (
	"testing"
)

func TestWrapQuote(t *testing.T) {
	tests := []struct {
		in     string
		indent string
		out    string
	}{
		{
			in:     "hello",
			indent: "\t",
			out:    `"hello"`,
		},
		{
			in:     "one\ntwo",
			indent: "\t\t\t\t",
			out: `"one\n" +
				"two"`,
		},
		{
			in:     "one\ntwo\n",
			indent: "\t\t\t\t",
			out: `"one\n" +
				"two\n"`,
		},
		{
			in:     "one\ntwo\n\n",
			indent: "\t\t\t\t",
			out: `"one\n" +
				"two\n" +
				"\n"`,
		},
		{
			in:     "one\ntwo\n\nthree\n",
			indent: "\t\t\t\t",
			out: `"one\n" +
				"two\n" +
				"\n" +
				"three\n"`,
		},
	}
	for _, tt := range tests {
		got := WrapQuote(tt.indent, tt.in)
		if got != tt.out {
			t.Errorf("want:\n%s\ngot: \n%s\n", tt.out, got)
		}
	}
}
