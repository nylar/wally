package wally

import (
	"io/ioutil"
	"testing"
	"strings"

	"github.com/stretchr/testify/assert"
)

func TestStopper(t *testing.T) {
	tests := []struct {
		Input  []string
		Output []string
	}{
		{
			[]string{"computer"},
			[]string{"computer"},
		},

		{
			[]string{"the"},
			[]string{},
		},

		{
			[]string{"technology"},
			[]string{"technology"},
		},

		{
			[]string{"ComPutEr"},
			[]string{"computer"},
		},

		{
			[]string{"Wasn't"},
			[]string{},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Output, Stopper(test.Input))
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
		assert.Equal(t, test.Output, SplitTextIntoWords(test.Input))
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		Input  interface{}
		Output []string
	}{
		{
			"I am a block of text and I am going to be parsed",
			[]string{"block", "text", "going", "parsed"},
		},
		
		{
			[]byte("I am another block of text but now I am in bytes"),
			[]string{"another", "block", "text", "now", "bytes"},
		},
	}
	
	for _, test := range tests {
		assert.Equal(t, test.Output, Parse(test.Input))
	}
}

func BenchmarkSplitTextIntoWords(b *testing.B) {
	file, err := ioutil.ReadFile("test_data/test.txt") // 30654 words
	if err != nil {
		b.Error("Could not load test data")
	}

	for n := 0; n < b.N; n++ {
		SplitTextIntoWords(file)
	}
}

func BenchmarkSplitTextIntoWords_two(b *testing.B) {
	file, err := ioutil.ReadFile("test_data/test_2.txt") // 1891 words
	if err != nil {
		b.Error("Could not load test data")
	}

	for n := 0; n < b.N; n++ {
		SplitTextIntoWords(file)
	}
}

func BenchmarkStopper(b *testing.B) {
	file, err := ioutil.ReadFile("test_data/test.txt") // 30654 words
	if err != nil {
		b.Error("Could not load test data")
	}
	
	data := strings.Fields(string(file))
	
	for n := 0; n < b.N; n++ {
		Stopper(data)
	}
}

func BenchmarkStopper_two(b *testing.B) {
	file, err := ioutil.ReadFile("test_data/test_2.txt") // 1891 words
	if err != nil {
		b.Error("Could not load test data")
	}

	data := strings.Fields(string(file))

	for n := 0; n < b.N; n++ {
		Stopper(data)
	}
}

func BenchmarkParse(b *testing.B) {
	file, err := ioutil.ReadFile("test_data/test.txt") // 30654 words
	if err != nil {
		b.Error("Could not load test data")
	}
	
	for n := 0; n < b.N; n++ {
		Parse(file)
	}
}

func BenchmarkParse_two(b *testing.B) {
	file, err := ioutil.ReadFile("test_data/test_2.txt") // 30654 words
	if err != nil {
		b.Error("Could not load test data")
	}

	for n := 0; n < b.N; n++ {
		Parse(file)
	}
}
