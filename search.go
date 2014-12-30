package wally

import (
	"math"
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

type Pagination struct {
	CurrentPage, PageCount       int
	HasPreviousPage, HasNextPage bool
}

func (p *Pagination) NumberOfPages() int {
	if p.PageCount == 0 {
		return 1
	}
	count := float64(p.PageCount) / float64(ItemsPerPage)
	pages := math.Ceil(count)
	return int(pages)
}

func (p *Pagination) Previous() int {
	if p.CurrentPage > 1 {
		return p.CurrentPage - 1
	}

	return 1
}

func (p *Pagination) Next() int {
	if p.CurrentPage < p.NumberOfPages() {
		return p.CurrentPage + 1
	}

	return p.CurrentPage
}

func (r *Results) NumberOfResults(keys []string, session *rdb.Session) error {
	var count int
	res, err := rdb.Db(Conf.Database.Name).Table(Conf.Tables.IndexTable).GetAllByIndex("word", rdb.Args(keys)).Count().Run(session)
	if err != nil {
		return err
	}

	res.One(&count)
	r.Count = int64(count)
	return nil
}

func parsePageNumber(page int) int {
	if page < 1 {
		page = 1
	}
	return int(page)
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

	r.NumberOfResults(keys, session)

	r.Results = res
	t := time.Since(start).Seconds()
	r.Time = t

	return r, nil
}
