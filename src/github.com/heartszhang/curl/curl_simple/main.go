package main

import (
	"flag"
	"fmt"
	"github.com/heartszhang/curl"
)

var (
	uri = flag.String("uri", "http://www.voachinese.com/img/assignedIcons/icon-blank.gif", "testing url")
)

func main() {
	flag.Parse()
	curler := curl.NewCurl("")
	cache, err := curler.Get(*uri)
	fmt.Println(cache, err)
	fp, m, w, h, err := curl.NewThumbnail(cache.Local, "", 400, 0)
	fmt.Println(fp, m, w, h, err)
}
