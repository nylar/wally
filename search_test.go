package wally

import (
	"testing"

	rdb "github.com/dancannon/gorethink"
	"github.com/stretchr/testify/assert"
)

func SearchSetup() error {
	DatabaseRebuild(session)

	d1 := Document{
		Id:      "1",
		Source:  "http://example.com",
		Title:   "Examples, Examples Everywhere",
		Author:  "John Johnson",
		Content: "This is an example of some example content remember though it's just an example",
	}

	d2 := Document{
		Id:      "2",
		Source:  "http://example.org",
		Title:   "Help Abandoned Examples",
		Author:  "",
		Content: "Save the example",
	}

	if err := d1.Put(session); err != nil {
		return err
	}

	if err := d2.Put(session); err != nil {
		return err
	}

	i1 := Indexer(d1.Content, d1.Id)
	i2 := Indexer(d2.Content, d2.Id)

	if err := IndexBatchPut(session, i1); err != nil {
		return err
	}

	if err := IndexBatchPut(session, i2); err != nil {
		return err
	}

	return nil
}

func TestSearch_Search(t *testing.T) {
	if err := SearchSetup(); err != nil {
		t.Errorf(err.Error())
	}

	results, err := Search("example", session)
	assert.NoError(t, err)

	assert.Equal(t, len(results.Results), 2)
}

func TestSearch_SearchWithNoIndex(t *testing.T) {
	DatabaseRebuild(session)
	rdb.Db(Database).Table(IndexTable).IndexDrop("word").Exec(session)

	_, err := Search("hello", session)
	assert.Error(t, err)
}
