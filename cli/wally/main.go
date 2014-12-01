package main

import (
	"log"
	"os"

	"github.com/nylar/wally"

	"github.com/codegangsta/cli"
	rdb "github.com/dancannon/gorethink"
	"github.com/fatih/color"
)

var version = "0.1"

func init() {
	var err error
	session, err = rdb.Connect(rdb.ConnectOpts{
		Address:  os.Getenv("RETHINKDB_URL"),
		Database: wally.Database,
	})
	if err != nil {
		color.Set(color.FgRed)
		log.Fatalln(err.Error())
		color.Unset()
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "wally"
	app.Version = version
	app.Usage = "command line utility"
	app.Flags = []cli.Flag{
		// cli.StringFlag{
		// 	Name:   "p, port",
		// 	Value:  "",
		// 	Usage:  "wally server port",
		// 	EnvVar: "WALLY_PORT",
		// },
		// cli.StringFlag{
		// 	Name:   "c, config",
		// 	Value:  "",
		// 	Usage:  "wally config file",
		// 	EnvVar: "WALLY_CONFIG",
		// },
	}

	app.Commands = []cli.Command{
		CrawlCommand(),
		RebuildCommand(),
		SearchCommand(),
	}

	app.Run(os.Args)
}
