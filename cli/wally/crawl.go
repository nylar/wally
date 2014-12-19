package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/nylar/odlaw"
	"github.com/nylar/wally"

	"github.com/codegangsta/cli"
	rdb "github.com/dancannon/gorethink"
	"github.com/fatih/color"
)

var (
	session *rdb.Session
	Success = color.New(color.FgGreen)
	Info    = color.New(color.FgBlue)
	Warning = color.New(color.FgYellow)
	Std     = color.New(color.FgMagenta)
)

func Crawler(url string) {
	start := time.Now()

	Std.Printf("\nGrabbing: %s\n\n", url)

	resp, err := http.Get(url)
	if err != nil {
		color.Set(color.FgRed)
		log.Fatalln(err.Error())
		color.Unset()
	}

	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		color.Set(color.FgRed)
		log.Fatalln(err.Error())
		color.Unset()
	}

	doc := odlaw.NewDocument(string(data))
	title := odlaw.ExtractTitle(doc)
	author := odlaw.ExtractAuthor(doc)
	content := odlaw.ExtractText(doc)

	d := wally.Document{
		Source: url,
		Title:  title,
		Author: author,
	}
	if err := d.Put(session); err != nil {
		color.Set(color.FgRed)
		log.Fatalln(err.Error())
		color.Unset()
	}

	Success.Printf("Created document: %s.\n\n", d.String())

	Info.Printf("Processing %d words.\n", len(strings.Fields(content)))

	indexes := wally.Indexer(content, d.Id)

	if _, err := rdb.Db(wally.Database).Table(wally.IndexTable).Insert(indexes).RunWrite(session); err != nil {
		color.Set(color.FgRed)
		log.Fatalln(err.Error())
		color.Unset()
	}

	rdb.Db(wally.Database).Table(wally.DocumentTable).Get(d.Id).Update(
		map[string]interface{}{"content": content},
	).RunWrite(session)

	Success.Printf("\nIndexing complete. Completed in %s.\n", time.Since(start))
}

func CrawlCommand() cli.Command {
	return cli.Command{
		Name:  "crawl",
		Usage: "crawls resource",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "url",
				Value: "",
				Usage: "crawls url",
			},
		},
		Action: func(c *cli.Context) {
			CrawlFunc(c)
		},
	}
}

func CrawlFunc(c *cli.Context) {
	url := c.String("url")

	urls := strings.Split(url, "|")

	for _, u := range urls {
		Crawler(u)
	}
}
