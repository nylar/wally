package wally

import (
	"testing"

	rdb "github.com/dancannon/gorethink"
	"github.com/stretchr/testify/assert"
)

var originalItemsPerPage = ItemsPerPage

func setUp(pages int) {
	ItemsPerPage = pages
}

func tearDown() {
	ItemsPerPage = originalItemsPerPage
}

func SearchSetup() error {
	DatabaseRebuild(session)

	d1 := Document{
		ID:      "1",
		Source:  "http://example.com",
		Title:   "Examples, Examples Everywhere",
		Author:  "John Johnson",
		Content: "This is an example of some example content remember though it's just an example",
	}

	d2 := Document{
		ID:      "2",
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

	i1 := Indexer(d1.Content, d1.ID)
	i2 := Indexer(d2.Content, d2.ID)

	if err := IndexBatchPut(session, i1); err != nil {
		return err
	}

	if err := IndexBatchPut(session, i2); err != nil {
		return err
	}

	return nil
}

func TestSearch_Search(t *testing.T) {
	setUp(1)
	if err := SearchSetup(); err != nil {
		t.Errorf(err.Error())
	}

	for i := 1; i < 3; i++ {
		results, err := Search("example", session, i)
		assert.NoError(t, err)

		assert.Equal(t, len(results.Results), 1)
	}
	tearDown()
}

func TestSearch_SearchNumberOfResults(t *testing.T) {
	if err := SearchSetup(); err != nil {
		t.Errorf(err.Error())
	}

	r := new(Results)
	err := r.NumberOfResults([]string{"example"}, session)

	assert.Equal(t, r.Count, 2)
	assert.NoError(t, err)
}

func TestSearch_SearchNumberOfResultsNoIndex(t *testing.T) {
	if err := SearchSetup(); err != nil {
		t.Errorf(err.Error())
	}
	rdb.Db(Conf.Database.Name).Table(Conf.Tables.IndexTable).IndexDrop("word").Exec(session)

	r := new(Results)
	err := r.NumberOfResults([]string{"example"}, session)

	assert.Equal(t, r.Count, 0)
	assert.Error(t, err)
}

func TestSearch_SearchWithNoIndex(t *testing.T) {
	DatabaseRebuild(session)
	rdb.Db(Conf.Database.Name).Table(Conf.Tables.IndexTable).IndexDrop("word").Exec(session)

	_, err := Search("hello", session, 1)
	assert.Error(t, err)
}

func TestSearch_parsePageNumber(t *testing.T) {
	tests := []struct {
		input  int
		output int
	}{
		{
			0,
			1,
		},
		{
			1,
			1,
		},
		{
			2,
			2,
		},
		{
			192983,
			192983,
		},
	}

	for _, test := range tests {
		assert.Equal(t, parsePageNumber(test.input), test.output)
	}
}

func TestSearch_PaginationNumberOfPages(t *testing.T) {
	tests := []struct {
		input  int
		output int
	}{
		{
			45,
			5,
		},
		{
			109,
			11,
		},
		{
			0,
			1,
		},
		{
			53,
			6,
		},
	}

	p := new(Pagination)

	for _, test := range tests {
		p.PageCount = test.input
		assert.Equal(t, test.output, p.NumberOfPages())
	}

}

func TestSearch_PaginationPrevious(t *testing.T) {
	tests := []struct {
		input  int
		output int
	}{
		{
			1,
			1,
		},
		{
			20,
			19,
		},
		{
			2,
			1,
		},
	}

	for _, test := range tests {
		p := new(Pagination)
		p.CurrentPage = test.input
		assert.Equal(t, test.output, p.Previous())
	}

}

func TestSearch_PaginationNext(t *testing.T) {
	tests := []struct {
		input  int
		output int
	}{
		{
			6,
			6,
		},
		{
			5,
			6,
		},
		{
			1,
			2,
		},
	}

	for _, test := range tests {
		p := new(Pagination)
		p.PageCount = 56
		p.CurrentPage = test.input
		assert.Equal(t, test.output, p.Next())
	}

}
