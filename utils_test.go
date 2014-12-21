package wally

import (
	"testing"

	rdb "github.com/dancannon/gorethink"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseRebuild(t *testing.T) {
	DatabaseRebuild(session)

	res, err := rdb.Db(Database).TableList().Run(session)
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

func TestToString(t *testing.T) {
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
