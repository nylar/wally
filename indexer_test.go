package wally

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	rdb "github.com/dancannon/gorethink"
	"github.com/stretchr/testify/assert"
)

var (
	session *rdb.Session
)

func init() {
	Database = "testing"
	DocumentTable = "documents"
	IndexTable = "indexes"

	var err error
	session, err = rdb.Connect(rdb.ConnectOpts{
		Address:  os.Getenv("RETHINKDB_URL"),
		Database: "test",
	})

	if err != nil {
		fmt.Errorf(err.Error())
	}

	// Reset database
	rdb.DbDrop(Database).Exec(session)
	rdb.DbCreate(Database).Exec(session)
}

func DbBootstrap() {
	// Drop tables
	rdb.Db(Database).TableDrop(DocumentTable).Run(session)
	rdb.Db(Database).TableDrop(IndexTable).Run(session)

	// Create tables
	rdb.Db(Database).TableCreate(DocumentTable).Run(session)
	rdb.Db(Database).TableCreate(IndexTable).Run(session)
}

func TestIndexer_Stopper(t *testing.T) {
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
			"wasn't",
			"",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Output, Stopper(test.Input))
	}
}

func TestIndexer_SplitTextIntoWords(t *testing.T) {
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
			[]string{"32"},
		},

		{
			true,
			[]string{},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Output, SplitTextIntoWords(test.Input))
	}
}

func TestIndexer_Indexer(t *testing.T) {
	docId := "12345-67890-ABCDE"

	tests := []struct {
		Input  interface{}
		Output []Index
	}{
		{
			"I am a block of text and I am going to be parsed",
			[]Index{
				Index{Word: "block", DocumentId: docId, Count: 0},
				Index{Word: "text", DocumentId: docId, Count: 0},
				Index{Word: "going", DocumentId: docId, Count: 0},
				Index{Word: "parsed", DocumentId: docId, Count: 0},
			},
		},

		{
			[]byte("I am another block of text but now I am in bytes"),
			[]Index{
				Index{Word: "another", DocumentId: docId, Count: 0},
				Index{Word: "block", DocumentId: docId, Count: 0},
				Index{Word: "text", DocumentId: docId, Count: 0},
				Index{Word: "now", DocumentId: docId, Count: 0},
				Index{Word: "bytes", DocumentId: docId, Count: 0},
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Output, Indexer(test.Input, docId))
	}
}

func TestIndexer_IndexString(t *testing.T) {
	indexId := "world"

	index := Index{Id: indexId}

	assert.Equal(t, "Index#world", index.String())
}

func TestIndexer_IndexPut(t *testing.T) {
	DbBootstrap()

	index := Index{Word: "hello", Count: 5, DocumentId: "12345-67890-ABCDE"}

	err := index.Put(session)
	assert.Nil(t, err)
	

	res, err := rdb.Db(Database).Table(IndexTable).Get(index.Id).Run(session)
	assert.Nil(t, err)

	var i Index
	err = res.One(&i)
	assert.Nil(t, err)

	assert.NotEqual(t, i.Id, "")
	assert.Equal(t, i.Word, "hello")
	assert.Equal(t, i.Count, 5)
	assert.Equal(t, i.DocumentId, "12345-67890-ABCDE")
}

func TestIndexer_DocumentString(t *testing.T) {
	docId := "12345-67890-ABCDE"

	doc := Document{Id: docId}

	assert.Equal(t, "Document#12345-67890-ABCDE", doc.String())
}

func TestIndexer_DocumentPut(t *testing.T) {
	DbBootstrap()

	doc := Document{
		Source:  "www.google.com",
		Content: "Lorem ipsum dolor sit amet.",
	}

	err := doc.Put(session)
	assert.Nil(t, err)

	res, err := rdb.Db(Database).Table(DocumentTable).Run(session)
	assert.Nil(t, err)

	var d Document
	err = res.One(&d)
	assert.Nil(t, err)

	assert.NotEqual(t, d.Id, "")
	assert.Equal(t, d.Source, "www.google.com")
	assert.Equal(t, d.Content, "Lorem ipsum dolor sit amet.")
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
		for _, word := range data {
			Stopper(word)
		}
	}
}

func BenchmarkStopper_two(b *testing.B) {
	file, err := ioutil.ReadFile("test_data/test_2.txt") // 1891 words
	if err != nil {
		b.Error("Could not load test data")
	}

	data := strings.Fields(string(file))

	for n := 0; n < b.N; n++ {
		for _, word := range data {
			Stopper(word)
		}
	}
}

func BenchmarkIndexer(b *testing.B) {
	docId := "12345-67890-ABCDE"
	file, err := ioutil.ReadFile("test_data/test.txt") // 30654 words
	if err != nil {
		b.Error("Could not load test data")
	}

	for n := 0; n < b.N; n++ {
		Indexer(file, docId)
	}
}

func BenchmarkIndexer_two(b *testing.B) {
	docId := "12345-67890-ABCDE"
	file, err := ioutil.ReadFile("test_data/test_2.txt") // 30654 words
	if err != nil {
		b.Error("Could not load test data")
	}

	for n := 0; n < b.N; n++ {
		Indexer(file, docId)
	}
}
