package wally

import (
	"io/ioutil"
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

func BenchmarkSplitTextIntoWords(b *testing.B) {
	file, err := ioutil.ReadFile("test.txt") // 30654 words
	if err != nil {
		b.Error("Could not load test data")
	}

	for n := 0; n < b.N; n++ {
		SplitTextIntoWords(file)
	}
}

func BenchmarkSplitTextIntoWords_two(b *testing.B) {
	file, err := ioutil.ReadFile("test 2.txt") // 1891 words
	if err != nil {
		b.Error("Could not load test data")
	}

	for n := 0; n < b.N; n++ {
		SplitTextIntoWords(file)
	}
}
