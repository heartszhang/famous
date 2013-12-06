package main

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"flag"
	"fmt"
	"github.com/heartszhang/cleaner"
	"github.com/heartszhang/curl"
	"os"
)

var (
	uri = flag.String("uri", "", "html doc link")
)

func main() {
	flag.Parse()
	if *uri == "" {
		flag.PrintDefaults()
		return
	}
	c := curl.NewCurl("")
	cache, err := c.GetUtf8(*uri)
	if err != nil {
		panic(err)
	}
	f, err := os.Open(cache.LocalUtf8)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	doc, err := html.Parse(f)
	if err != nil {
		panic(err)
	}
	ex := cleaner.NewExtractor("")
	article, _, err := ex.MakeHtmlReadable(doc, *uri)
	if err != nil {
		panic(err)
	}
	print_html_doc(article)
}
func print_html_doc(node *html.Node) {
	var buffer bytes.Buffer
	html.Render(&buffer, node) // ignore return error
	fmt.Println(buffer.String())
}
