package wally

import (
	"io/ioutil"
	"net/http"

	rdb "github.com/dancannon/gorethink"
	"github.com/nylar/odlaw"
)

type Sourcer interface {
	Grab(resource string) ([]byte, error)
}

type WebSource struct{}

func (ws *WebSource) Grab(resource string) ([]byte, error) {
	resp, err := http.Get(resource)
	if err != nil {
		return []byte{}, err
	}

	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return data, nil
}

func GrabResource(source Sourcer, resource string) ([]byte, error) {
	return source.Grab(resource)
}

// Crawler grabs the contents of a URL and passes the data to Odlaw for
// processing, it is then written in bulk to the database.
func Crawler(url string, session *rdb.Session, source Sourcer) error {
	data, err := GrabResource(source, url)
	if err != nil {
		return err
	}

	doc := odlaw.NewDocument(string(data))
	title := odlaw.ExtractTitle(doc)
	author := odlaw.ExtractAuthor(doc)
	content := odlaw.ExtractText(doc)

	d := Document{
		Source:  url,
		Title:   title,
		Author:  author,
		Content: content,
	}
	_ = d.Put(session)

	indexes := Indexer(content, d.ID)

	_, _ = rdb.Db(Conf.Database.Name).Table(Conf.Tables.IndexTable).Insert(indexes).RunWrite(session)
	return nil
}
