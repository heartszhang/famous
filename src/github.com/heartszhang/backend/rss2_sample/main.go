package main

import (
	"flag"
	"fmt"
	"github.com/heartszhang/backend"
)

var (
	file = flag.String("file", "sample.xml", "rss2 file")
)

func main() {
	flag.Parse()
	if *file == "" {
		flag.PrintDefaults()
		return
	}
	_, err := backend.CreateFeedSourceRss2(*file)
	if err != nil {
		fmt.Println(err)
	}
	_, err = backend.CreateFeedEntriesRss2(*file)
	if err != nil {
		fmt.Println(err)
	}
}
