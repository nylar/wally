![wally](http://i.imgur.com/MSny4Kj.png)

[![wercker status](https://img.shields.io/wercker/ci/544c0c84ea87f6374f000650.svg "wercker status")](https://app.wercker.com/project/bykey/ffa1468bc1ebe9c1dd7d0c2d00f4c76f)
[![godoc reference](https://img.shields.io/badge/godoc-reference-blue.svg "godoc reference")](https://godoc.org/github.com/nylar/wally)
[![license](http://img.shields.io/badge/license-unlicense-red.svg "license")](https://raw.githubusercontent.com/nylar/wally/master/LICENSE)
[![Coverage Status](https://coveralls.io/repos/nylar/wally/badge.png?branch=HEAD)](https://coveralls.io/r/nylar/wally?branch=HEAD)

A full-text search engine built on Go.

## Getting Started with Wally

Use go get to retrieve the latest version of Wally
```shell
go get -u github.com/nylar/wally
```
Wally's CLI can be built with go build and it will generate a binary for your machine.
```shell
cd cli/wally
go build
```
You can also install the compiled binary to your Go workspaces bin directory.
```shell
cd cli/wally
go install
```

### Dependencies

There is one external dependency and serveral Go packages that Wally relies on. RethinkDB, the external dependency, can be installed from [http://rethinkdb.com/](http://rethinkdb.com/) or through your favourite package manager.

Assuming you have you Go installed and configured on your system, you can grab all the Go dependencies with go get.
```shell
go get ./...
```

### Running Wally's Tests

Wally's tests can be run with the built in go tool, to see code coverage you need to have the cover tool for go installed.
```shell
go get golang.org/x/tools/cmd/cover
```

Conversely, Wally's coverage reports can be seen on [Coveralls](https://coveralls.io/r/nylar/wally). Wally is tested with continuous integration via [Wercker](https://app.wercker.com/#applications/544c0c84ea87f6374f000650/tab).

## Configuration

Wally depends on a [YAML](http://yaml.org/) configuration file, a sample configuration file can be found at cli/wally/config.yaml.
```yaml
database:
  host: localhost:28015
  name: wally
tables:
  document_table: documents
  index_table: indexes
```
To then use a configuration file in your project, you will need to do the following.
```go
package main

import (
  "io/ioutil"
  "log"
  
  "github.com/nylar/wally"
  rdb "github.com/dancannon/gorethink"
)

var session *rdb.Session

func main() {
  var err error
  confData, err := ioutil.ReadFile("config.yml")
  if err != nil {
    log.Fatalln(err.Error())
  }

  wally.Conf, err = wally.LoadConfig(confData)
  if err != nil {
    log.Fatalln(err.Error())
  }

  session, err = rdb.Connect(rdb.ConnectOpts{
    Address:  wally.Conf.Database.Host,
    Database: wally.Conf.Database.Name,
  })
  if err != nil {
    log.Fatalln(err.Error())
  }
}
```

## Demo

Wally demo app on [GitHub](https://github.com/nylar/wally-ui).

![screenshot](https://camo.githubusercontent.com/36c067e4d7c8e9b4d640f8357a0656d7f4ebce15/687474703a2f2f692e696d6775722e636f6d2f694f75394351542e706e67)

You can download the code for the Wally demo with the go get tool.
```shell
go get -u github.com/nylar/wally-ui
```

Then to run the application, use go run. You can then fire up your browser and point it to [http://localhost:8008](http://localhost:8008).

```shell
go run main.go",
```
Inside this directory you will find a config file that you modify to match your environment. If port 8008 is occupied on your machine, you can modify the main() function in main.go, like so.

```go
func main() {
  http.HandleFunc("/", handler)
  http.ListenAndServe(":8880", nil)
}
```
