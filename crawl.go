package wally

import (
	"io/ioutil"
	"net/http"

	rdb "github.com/dancannon/gorethink"
	"github.com/nylar/odlaw"
)

func GrabUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return data, nil
}

func Crawler(url string, session *rdb.Session) error {
	data, err := GrabUrl(url)
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
	if err := d.Put(session); err != nil {
		return err
	}

	indexes := Indexer(content, d.Id)

	if _, err := rdb.Db(Database).Table(IndexTable).Insert(indexes).RunWrite(session); err != nil {
		return err
	}

	return nil
}
