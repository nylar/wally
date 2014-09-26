package wally

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStopper(t *testing.T) {
	tests := map[string]string{
		"computer":   "computer",
		"the":        "",
		"technology": "technology",
		"ComPutEr":   "computer",
		"Wasn't":     "",
	}

	for k, v := range tests {
		assert.Equal(t, v, Stopper(k), "Should be equal.")
	}
}
