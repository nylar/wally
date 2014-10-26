package main

import (
	"github.com/codegangsta/cli"
	"os"
)

var version = "0.1"

func main() {
	app := cli.NewApp()
	app.Name = "wally"
	app.Version = version
	app.Usage = "command line utility"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "p, port",
			Value: "",
			Usage: "wally server port",
			EnvVar: "WALLY_PORT",
		},
		cli.StringFlag{
			Name: "c, config",
			Value: "",
			Usage: "wally config file",
			EnvVar: "WALLY_CONFIG",
		},
	}
	
	app.Commands = []cli.Command{}
	
	app.Run(os.Args)
}
