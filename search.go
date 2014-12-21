package wally

import (
	"strings"
	"time"

	rdb "github.com/dancannon/gorethink"
)

type Query struct {
	Document
	Index
}

type Results struct {
	Count   int64
	Results []Query
	Time    float64
}

func Search(query string, session *rdb.Session) (*Results, error) {
	start := time.Now()
	res := []Query{}
	keys := strings.Split(query, " ")

	r := new(Results)

	results, err := rdb.Db(Database).Table(IndexTable).GetAllByIndex("word", rdb.Args(keys)).EqJoin("document_id", rdb.Db(Database).Table(DocumentTable)).Zip().OrderBy(rdb.Desc("count")).Run(session)
	if err != nil {
		return nil, err
	}

	if err := results.All(&res); err != nil {
		return nil, err
	}

	r.Count = int64(len(res))
	r.Results = res
	t := time.Since(start).Seconds()
	r.Time = t

	return r, nil
}