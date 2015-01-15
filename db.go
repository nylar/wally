package wally

import rdb "github.com/dancannon/gorethink"

type Store interface {
	Put(session *rdb.Session) error
}

func NewDocument(source string) *Document {
	return &Document{ID: source}
}

func NewIndex(word, documentID string) *Index {
	return &Index{Word: word, DocumentID: documentID}
}
