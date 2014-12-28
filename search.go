package wally

import (
	"strings"
	"time"

	rdb "github.com/dancannon/gorethink"
)

// Query is one result in a successful search.
type Query struct {
	Document
	Index
}

// Results contains the search results of a successful Search().
type Results struct {
	Count   int64
	Results []Query
	Time    float64
}

func parsePageNumber(page int) uint {
	if page < 1 {
		page = 1
	}
	return uint(page)
}

// Search returns a list of results along with the time taken to run and the
// number of results found.
func Search(query string, session *rdb.Session, currentPage int) (*Results, error) {
	start := time.Now()
	res := []Query{}
	keys := strings.Split(query, " ")

	r := new(Results)

	page := parsePageNumber(currentPage)

	lower := (page - 1) * ItemsPerPage
	upper := page * ItemsPerPage

	results, err := rdb.Db(Conf.Database.Name).Table(Conf.Tables.IndexTable).GetAllByIndex("word", rdb.Args(keys)).EqJoin("document_id", rdb.Db(Conf.Database.Name).Table(Conf.Tables.DocumentTable)).Zip().OrderBy(rdb.Desc("count")).Slice(lower, upper).Run(session)
	if err != nil {
		return nil, err
	}

	results.All(&res)

	r.Count = int64(len(res))
	r.Results = res
	t := time.Since(start).Seconds()
	r.Time = t

	return r, nil
}
