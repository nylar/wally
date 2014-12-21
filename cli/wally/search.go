package main

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
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

	results, err := wally.Search(query, session)
	if err != nil {
		color.Set(color.FgRed)
		log.Fatalln(err.Error())
		color.Unset()
	}

	if results.Count == 0 {
		fmt.Println("No results found")
	} else {
		Std.Printf("\nFound %d results in %s\n\n", results.Count, results.Time)
		for _, r := range results.Results {
			content := r.Content
			if r.Title != "" {
				Info.Printf("\n%s", r.Title)
				Success.Printf("\n%s\n", r.Source)
			} else {
				Success.Printf("\n%s\n", r.Source)
			}
			if len(r.Content) > 150 {
				content = r.Content[:150] + " ..."
			}
			fmt.Printf("%s\n\n", content)
		}
	}
}
