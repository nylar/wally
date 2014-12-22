package main

import (
	"log"
	"strings"

	"github.com/nylar/wally"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
)

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
		if err := wally.Crawler(u, session); err != nil {
			color.Set(color.FgRed)
			log.Fatalln(err.Error())
			color.Unset()
		}
	}
}
