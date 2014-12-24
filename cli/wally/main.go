package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/nylar/wally"

	"github.com/codegangsta/cli"
	rdb "github.com/dancannon/gorethink"
	"github.com/fatih/color"
)

var (
	version = "0.1"
	session *rdb.Session
)

func logError(err error) {
	color.Set(color.FgRed)
	log.Fatalln(err.Error())
	color.Unset()
}

func init() {
	var err error
	confData, err := ioutil.ReadFile("config.yml")
	if err != nil {
		logError(err)
	}

	wally.Conf, err = wally.LoadConfig(confData)
	if err != nil {
		logError(err)
	}

	session, err = rdb.Connect(rdb.ConnectOpts{
		Address:  wally.Conf.Database.Host,
		Database: "test",
	})
	if err != nil {
		logError(err)
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
