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

const (
	TestConfig = `
database:
  host: localhost:28015
  name: testing

tables:
  document_table: documents
  index_table: indexes
`
)

var (
	session *rdb.Session
)

func init() {
	var err error
	Conf, err = LoadConfig([]byte(TestConfig))
	if err != nil {
		fmt.Errorf(err.Error())
	}

	session, err = rdb.Connect(rdb.ConnectOpts{
		Address:  os.Getenv("RETHINKDB_URL"),
		Database: "test",
	})

	if err != nil {
		fmt.Errorf(err.Error())
	}

	// Reset database
	rdb.DbDrop(Conf.Database.Name).Exec(session)
	rdb.DbCreate(Conf.Database.Name).Exec(session)
	DatabaseRebuild(session)
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
	docID := "12345-67890-ABCDE"

	tests := []struct {
		Input  interface{}
		Output []Index
	}{
		{
			"I am a block of text and I am going to be parsed",
			[]Index{
				Index{Word: "block", DocumentID: docID, Count: 0},
				Index{Word: "text", DocumentID: docID, Count: 0},
				Index{Word: "going", DocumentID: docID, Count: 0},
				Index{Word: "parsed", DocumentID: docID, Count: 0},
			},
		},

		{
			[]byte("I am another block of text but now I am in bytes"),
			[]Index{
				Index{Word: "another", DocumentID: docID, Count: 0},
				Index{Word: "block", DocumentID: docID, Count: 0},
				Index{Word: "text", DocumentID: docID, Count: 0},
				Index{Word: "now", DocumentID: docID, Count: 0},
				Index{Word: "bytes", DocumentID: docID, Count: 0},
			},
		},
	}

	for x, test := range tests {
		i := Indexer(test.Input, docID)
		assert.NotEqual(t, "", i[x].ID)
	}
}

func TestIndexer_IndexString(t *testing.T) {
	indexID := "world"

	index := Index{ID: indexID}

	assert.Equal(t, "Index#world", index.String())
}

func TestIndexer_IndexPut(t *testing.T) {
	index := Index{Word: "hello", Count: 5, DocumentID: "12345-67890-ABCDE"}

	err := index.Put(session)
	assert.Nil(t, err)

	res, err := rdb.Db(Conf.Database.Name).Table(Conf.Tables.IndexTable).Get(index.ID).Run(session)
	assert.Nil(t, err)

	var i Index
	err = res.One(&i)
	assert.Nil(t, err)

	assert.NotEqual(t, i.ID, "")
	assert.Equal(t, i.Word, "hello")
	assert.Equal(t, i.Count, 5)
	assert.Equal(t, i.DocumentID, "12345-67890-ABCDE")
}

func TestIndexer_IndexPutInvalid(t *testing.T) {

	i := Index{ID: "1"}
	i2 := Index{ID: "1"}

	err := i.Put(session)
	assert.NoError(t, err)

	err = i2.Put(session)
	assert.Error(t, err)
}

func TestIndexer_IndexBatchPut(t *testing.T) {
	indexes := []Index{
		{
			ID:   "1",
			Word: "dupe",
		},
		{
			ID:   "1",
			Word: "dupe",
		},
	}
	err := IndexBatchPut(session, indexes)
	assert.Error(t, err)
}

func TestIndexer_DocumentString(t *testing.T) {
	docID := "12345-67890-ABCDE"

	doc := Document{ID: docID}

	assert.Equal(t, "Document#12345-67890-ABCDE", doc.String())
}

func TestIndexer_DocumentPut(t *testing.T) {
	doc := Document{
		Source:  "www.google.com",
		Content: "Lorem ipsum dolor sit amet.",
	}

	err := doc.Put(session)
	assert.Nil(t, err)

	res, err := rdb.Db(Conf.Database.Name).Table(Conf.Tables.DocumentTable).Run(session)
	assert.Nil(t, err)

	var d Document
	err = res.One(&d)
	assert.Nil(t, err)

	assert.NotEqual(t, d.ID, "")
	assert.Equal(t, d.Source, "www.google.com")
	assert.Equal(t, d.Content, "Lorem ipsum dolor sit amet.")
}

func TestIndexer_DocumentPutDupeDocs(t *testing.T) {
	doc1 := Document{ID: "1"}
	doc2 := Document{ID: "1"}

	err := doc1.Put(session)
	assert.NoError(t, err)

	err = doc2.Put(session)
	assert.Error(t, err)
}

func TestIndexer_RemoveDuplicates(t *testing.T) {
	indexes := []Index{
		Index{Word: "hello"},
		Index{Word: "world"},
		Index{Word: "hello"},
		Index{Word: "hello"},
		Index{Word: "world"},
		Index{Word: "cruel"},
	}

	indexes = RemoveDuplicates(indexes)
	assert.Equal(t, len(indexes), 3)
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
	docID := "12345-67890-ABCDE"
	file, err := ioutil.ReadFile("test_data/test.txt") // 30654 words
	if err != nil {
		b.Error("Could not load test data")
	}

	for n := 0; n < b.N; n++ {
		Indexer(file, docID)
	}
}

func BenchmarkIndexer_two(b *testing.B) {
	docID := "12345-67890-ABCDE"
	file, err := ioutil.ReadFile("test_data/test_2.txt") // 30654 words
	if err != nil {
		b.Error("Could not load test data")
	}

	for n := 0; n < b.N; n++ {
		Indexer(file, docID)
	}
}
