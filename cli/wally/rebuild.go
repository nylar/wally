package main

import (
	"log"
	
	"github.com/nylar/wally"
	
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	rdb "github.com/dancannon/gorethink"
)

func RebuildCommand() cli.Command {
	return cli.Command{
		Name: "rebuild",
		Usage: "rebuild database",
		Action: func(c *cli.Context) {
			RebuildFunc()
		},
	}
}

func RebuildFunc() {
	if err := rdb.Db(wally.Database).TableDrop("documents").Exec(session); err != nil {
		color.Set(color.FgRed)
		log.Fatalln(err.Error())
		color.Unset()
	}
	Success.Println("Dropped 'documents' table.")
	
	if err := rdb.Db(wally.Database).TableDrop("indexes").Exec(session); err != nil {
		color.Set(color.FgRed)
		log.Fatalln(err.Error())
		color.Unset()
	}
	Success.Println("Dropped 'indexes' table.")
	
	if err := rdb.Db(wally.Database).TableCreate("documents").Exec(session); err != nil {
		color.Set(color.FgRed)
		log.Fatalln(err.Error())
		color.Unset()
	}
	Success.Println("Created 'documents' table.")
	
	if err := rdb.Db(wally.Database).TableCreate("indexes").Exec(session); err != nil {
		color.Set(color.FgRed)
		log.Fatalln(err.Error())
		color.Unset()
	}
	Success.Println("Created 'indexes' table.")
	
	if err := rdb.Db(wally.Database).Table("indexes").IndexCreate("word").Exec(session); err != nil {
		color.Set(color.FgRed)
		log.Fatalln(err.Error())
		color.Unset()
	}
	Success.Println("Created 'indexes' secondary index.")
}
