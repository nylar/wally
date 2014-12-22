package wally

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"code.google.com/p/go-uuid/uuid"
	rdb "github.com/dancannon/gorethink"
)

var stopWords = map[string]bool{
	"a":          true,
	"about":      true,
	"above":      true,
	"after":      true,
	"again":      true,
	"against":    true,
	"all":        true,
	"am":         true,
	"an":         true,
	"and":        true,
	"any":        true,
	"are":        true,
	"aren't":     true,
	"as":         true,
	"at":         true,
	"be":         true,
	"because":    true,
	"been":       true,
	"before":     true,
	"being":      true,
	"below":      true,
	"between":    true,
	"both":       true,
	"but":        true,
	"by":         true,
	"can't":      true,
	"cannot":     true,
	"could":      true,
	"couldn't":   true,
	"did":        true,
	"didn't":     true,
	"do":         true,
	"does":       true,
	"doesn't":    true,
	"doing":      true,
	"don't":      true,
	"down":       true,
	"during":     true,
	"each":       true,
	"few":        true,
	"for":        true,
	"from":       true,
	"further":    true,
	"had":        true,
	"hadn't":     true,
	"has":        true,
	"hasn't":     true,
	"have":       true,
	"haven't":    true,
	"having":     true,
	"he":         true,
	"he'd":       true,
	"he'll":      true,
	"he's":       true,
	"her":        true,
	"here":       true,
	"here's":     true,
	"hers":       true,
	"herself":    true,
	"him":        true,
	"himself":    true,
	"his":        true,
	"how":        true,
	"how's":      true,
	"i":          true,
	"i'd":        true,
	"i'll":       true,
	"i'm":        true,
	"i've":       true,
	"if":         true,
	"in":         true,
	"into":       true,
	"is":         true,
	"isn't":      true,
	"it":         true,
	"it's":       true,
	"its":        true,
	"itself":     true,
	"let's":      true,
	"me":         true,
	"more":       true,
	"most":       true,
	"mustn't":    true,
	"my":         true,
	"myself":     true,
	"no":         true,
	"nor":        true,
	"not":        true,
	"of":         true,
	"off":        true,
	"on":         true,
	"once":       true,
	"only":       true,
	"or":         true,
	"other":      true,
	"ought":      true,
	"our":        true,
	"ours":       true,
	"ourselves":  true,
	"out":        true,
	"over":       true,
	"own":        true,
	"same":       true,
	"shan't":     true,
	"she":        true,
	"she'd":      true,
	"she'll":     true,
	"she's":      true,
	"should":     true,
	"shouldn't":  true,
	"so":         true,
	"some":       true,
	"such":       true,
	"than":       true,
	"that":       true,
	"that's":     true,
	"the":        true,
	"their":      true,
	"theirs":     true,
	"them":       true,
	"themselves": true,
	"then":       true,
	"there":      true,
	"there's":    true,
	"these":      true,
	"they":       true,
	"they'd":     true,
	"they'll":    true,
	"they're":    true,
	"they've":    true,
	"this":       true,
	"those":      true,
	"through":    true,
	"to":         true,
	"too":        true,
	"under":      true,
	"until":      true,
	"up":         true,
	"very":       true,
	"was":        true,
	"wasn't":     true,
	"we":         true,
	"we'd":       true,
	"we'll":      true,
	"we're":      true,
	"we've":      true,
	"were":       true,
	"weren't":    true,
	"what":       true,
	"what's":     true,
	"when":       true,
	"when's":     true,
	"where":      true,
	"where's":    true,
	"which":      true,
	"while":      true,
	"who":        true,
	"who's":      true,
	"whom":       true,
	"why":        true,
	"why's":      true,
	"with":       true,
	"won't":      true,
	"would":      true,
	"wouldn't":   true,
	"you":        true,
	"you'd":      true,
	"you'll":     true,
	"you're":     true,
	"you've":     true,
	"your":       true,
	"yours":      true,
	"yourself":   true,
	"yourselves": true,
}

var (
	Database      = "wally"
	DocumentTable = "documents"
	IndexTable    = "indexes"
)

type Document struct {
	Id      string `gorethink:"id"`
	Source  string `gorethink:"source"`
	Title   string `gorethink:"title"`
	Author  string `gorethink:"author"`
	Content string `gorethink:"content"`
}

func (d *Document) String() string {
	return fmt.Sprintf("Document#%s", d.Id)
}

func (d *Document) Put(session *rdb.Session) error {
	if d.Id == "" {
		d.Id = uuid.New()
	}

	res, _ := rdb.Db(Database).Table(DocumentTable).Insert(d).RunWrite(session)
	if res.Errors > 0 {
		return errors.New(res.FirstError)
	}
	return nil
}

type Index struct {
	Id         string `gorethink:"id"`
	Word       string `gorethink:"word"`
	Count      int64  `gorethink:"count"`
	DocumentId string `gorethink:"document_id"`
}

func (i *Index) String() string {
	return fmt.Sprintf("Index#%s", i.Id)
}

func (i *Index) Put(session *rdb.Session) error {
	if i.Id == "" {
		i.Id = uuid.New()
	}
	_, err := rdb.Db(Database).Table(IndexTable).Insert(i).RunWrite(session)
	return err
}

func IndexBatchPut(session *rdb.Session, indexes []Index) error {
	// FIXME: This call should be returning an error, instead it is returning nil.
	res, _ := rdb.Db(Database).Table(IndexTable).Insert(indexes).RunWrite(session)
	if res.Errors > 0 {
		return errors.New(res.FirstError)
	}
	return nil
}

// Given a blob of text, as a string or slice of bytes, split each word separted
// by whitespace, extra whitespace should be removed, any other type returns
// an empty string and therefore will not be processed later.
func SplitTextIntoWords(text interface{}) []string {
	words := ToString(text)
	return strings.Fields(words)
}

// Compares a given word to a list of stopper words (words which are common
// and therefore should be ignored when indexing).
func Stopper(word string) string {
	if _, ok := stopWords[word]; ok {
		return ""
	}
	return word
}

func RemoveDuplicates(i []Index) []Index {
	result := []Index{}
	seen := map[string]int64{}
	for _, val := range i {
		if _, ok := seen[val.Word]; !ok {
			// val.Count = val.Count + 1
			result = append(result, val)
			seen[val.Word] = seen[val.Word] + 1
		} else {
			seen[val.Word] = seen[val.Word] + 1
		}
	}

	finalResults := []Index{}

	for _, res := range result {
		count := seen[res.Word]
		res.Count = count
		finalResults = append(finalResults, res)
	}
	return finalResults
}

func Indexer(text interface{}, documentId string) []Index {
	// Divide into individual words
	words := SplitTextIntoWords(text)

	var normalisedWords []Index
	var wg sync.WaitGroup

	wg.Add(len(words))
	for _, word := range words {
		go func(word string) {
			defer wg.Done()

			// Lowercase words
			word = strings.ToLower(word)

			// Remove stopper words
			word = Stopper(word)

			if word == "" || len(word) < 2 {
				return
			}

			// Apply stemming

			// Append to normalised word list
			normalisedWords = append(normalisedWords, Index{
				Id:         uuid.New(),
				Word:       word,
				DocumentId: documentId,
			})
		}(word)
	}

	wg.Wait()

	return RemoveDuplicates(normalisedWords)
}
