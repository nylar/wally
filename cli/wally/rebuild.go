package main

import (
	"github.com/nylar/wally"

	"github.com/codegangsta/cli"
)

func RebuildCommand() cli.Command {
	return cli.Command{
		Name:  "rebuild",
		Usage: "rebuild database",
		Action: func(c *cli.Context) {
			RebuildFunc()
		},
	}
}

func RebuildFunc() {
	wally.DatabaseRebuild(session)
	wally.Success.Println("Rebuilt database")
}
