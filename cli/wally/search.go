package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/nylar/wally"
)

func SearchCommand() cli.Command {
	return cli.Command{
		Name:  "search",
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
	query := c.String("query")

	results, err := wally.Search(query, session, 1)
	if err != nil {
		logError(err)
	}

	if results.Count == 0 {
		fmt.Println("No results found")
	} else {
		wally.Std.Printf("\nFound %d results in %s\n\n", results.Count, results.Time)
		for _, r := range results.Results {
			content := r.Content
			if r.Title != "" {
				wally.Info.Printf("\n%s", r.Title)
				wally.Success.Printf("\n%s\n", r.Source)
			} else {
				wally.Success.Printf("\n%s\n", r.Source)
			}
			if len(r.Content) > 150 {
				content = r.Content[:150] + " ..."
			}
			fmt.Printf("%s\n\n", content)
		}
	}
}
