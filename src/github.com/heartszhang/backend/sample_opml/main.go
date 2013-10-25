//
package main

import (
	"flag"
	"fmt"
	"github.com/heartszhang/backend"
)

var (
	file = flag.String("file", "feedly.opml", "local opml file path")
)

func main() {
	flag.Parse()
	_, err := backend.CreateFeedsCategoryOpml(*file)
	fmt.Println(err)
}
