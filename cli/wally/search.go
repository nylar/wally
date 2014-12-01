package main

import (
	"log"
	"fmt"
	"time"
	"strings"

	"github.com/nylar/wally"
	
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	rdb "github.com/dancannon/gorethink"
)

type Query struct {
	wally.Document
	wally.Index
}

func SearchCommand() cli.Command {
	return cli.Command{
		Name: "search",
		Usage: "search wally",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "query",
				Value: "",
				Usage: "query args",
			},
		},
		Action: func(c *cli.Context) {
			SearchFunc(c)
		},
	}
}

func SearchFunc(c *cli.Context) {
	start := time.Now()
	res := []Query{}
	query := c.String("query")
	
	keys := strings.Split(query, " ")
	
	results, err := rdb.Db(wally.Database).Table(wally.IndexTable).GetAllByIndex("word", rdb.Args(keys)).EqJoin("document_id", rdb.Db(wally.Database).Table(wally.DocumentTable)).Zip().OrderBy(rdb.Desc("count")).Run(session)
	if err != nil {
		color.Set(color.FgRed)
		log.Fatalln(err.Error())
		color.Unset()
	}
	
	if err := results.All(&res); err != nil {
		color.Set(color.FgRed)
		log.Fatalln(err.Error())
		color.Unset()
	}
	
	
	if len(res) == 0 {
		fmt.Println("No results found")
	} else {
		Std.Printf("\nFound %d results in %s\n\n", len(res), time.Since(start))
		for i, r := range res {
			content := r.Content
			Info.Printf("%d:\t", i+1)
			Success.Printf("%s\n\n", r.Source)
			if len(r.Content) > 150 {
				content = r.Content[:150]
			}
			fmt.Printf("\t%s\n\n", content)
		}	
	}
}
