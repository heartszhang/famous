package main

import (
	"flag"
	//	"fmt"
	"github.com/heartszhang/curl"
)

var (
	uri = flag.String("uri", "http://www.baidu.com/", "testing url")
)

func main() {
	flag.Parse()
	curler := curl.NewCurl("")
	curler.Get(*uri)
}
