package wally

import (
	"testing"

	rdb "github.com/dancannon/gorethink"
	"github.com/stretchr/testify/assert"
)

func TestUtils_DatabaseRebuild(t *testing.T) {
	DatabaseRebuild(session)

	res, err := rdb.DB(Conf.Database.Name).TableList().Run(session)
	if err != nil {
		t.Errorf(err.Error())
	}

	var response []interface{}
	err = res.All(&response)
	if err != nil {
		t.Errorf(err.Error())
	}

	assert.Equal(t, len(response), 2)
}

func TestUtils_ToString(t *testing.T) {
	tests := []struct {
		input  interface{}
		output string
	}{
		{
			"string",
			"string",
		},
		{
			[]byte("byte slice"),
			"byte slice",
		},
		{
			34324,
			"34324",
		},
		{
			true,
			"",
		},
	}

	for _, test := range tests {
		assert.Equal(t, ToString(test.input), test.output)
	}
}

func TestUtils_TruncateText(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{
			"hello world",
			"hello world",
		},
		{
			"Lorem ipsum dolor sit amet, natoque quis",
			"Lorem ipsum dolor sit ...",
		},
	}

	for _, test := range tests {
		assert.Equal(t, TruncateText(test.input, " ...", 25), test.output)
	}
}
