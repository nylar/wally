package wally

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStopper(t *testing.T) {
	tests := []struct {
		Input  string
		Output string
	}{
		{
			"computer",
			"computer",
		},

		{
			"the",
			"",
		},

		{
			"technology",
			"technology",
		},

		{
			"ComPutEr",
			"computer",
		},

		{
			"Wasn't",
			"",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Output, Stopper(test.Input), "Should be equal.")
	}
}

func TestSplitTextIntoWords(t *testing.T) {
	tests := []struct {
		Input  interface{}
		Output []string
	}{
		{
			"I am a block of text",
			[]string{"I", "am", "a", "block", "of", "text"},
		},

		{
			"        superfluous    whitespace ",
			[]string{"superfluous", "whitespace"},
		},

		{
			"              ",
			[]string{},
		},

		{
			[]byte("fancy a byte?"),
			[]string{"fancy", "a", "byte?"},
		},

		{
			32,
			[]string{},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Output, SplitTextIntoWords(test.Input), "Should return a slice of strings.")
	}
}
