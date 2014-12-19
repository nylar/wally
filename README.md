![wally](http://i.imgur.com/MSny4Kj.png)

[![wercker status](https://app.wercker.com/status/ffa1468bc1ebe9c1dd7d0c2d00f4c76f/s "wercker status")](https://app.wercker.com/project/bykey/ffa1468bc1ebe9c1dd7d0c2d00f4c76f)

A full-text search engine built on Go.

## How to install

You can grab that latest build of Wally using:

```go get -u github.com/nylar/wally```

## Using Wally

Parsing a block of text

```go
package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/nylar/wally"
)

func main() {
	file, err := ioutil.ReadFile("somefile.txt")
	if err != nil {
		panic("Could not read file")
	}

	parsedWords := wally.Parse(file)

	fmt.Printf("Processed:\t %d words.\n", len(strings.Fields(string(file))))
	fmt.Printf("Finished with:\t %d words.\n", len(parsedWords))
}
```
