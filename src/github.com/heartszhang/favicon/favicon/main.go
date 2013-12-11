package main

import (
	"code.google.com/p/go.net/html"
	"flag"
	"fmt"
	"github.com/heartszhang/curl"
	"net/url"
	"os"
	"strings"
)

var (
	uri = flag.String("uri", "", "favicon hint")
)

// candidate
// hint, "/", "/favicon.ico"
func main() {
	flag.Parse()
	u, err := url.ParseRequestURI(*uri)
	if err != nil {
		flag.PrintDefaults()
		return
	}
	candis := []string{u.RequestURI(), "/", "/favicon.ico"}
	var logourl string
	for _, candi := range candis {
		x, _ := u.Parse(candi)
		logourl = favicon_try_from_url(x.String())
		if logourl != "" {
			break
		}
	}
	fmt.Println("logo:", logourl)
}

func query_select(node *html.Node, name string) *html.Node {
	if node == nil {
		return nil
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == name {
			return child
		}
	}
	return nil
}
func query_selects(node *html.Node, name string) []*html.Node {
	if node == nil {
		return nil
	}
	var v []*html.Node
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == name {
			v = append(v, child)
		}
	}
	return v
}

func node_get_attribute(node *html.Node, name string) string {
	if node == nil {
		return ""
	}
	for _, attr := range node.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}
	return ""
}

func icon_from_link_rel(doc *html.Node) (string, bool) {
	htmln := query_select(doc, "html")
	head := query_select(htmln, "head")
	links := query_selects(head, "link")
	for _, link := range links {
		rel := node_get_attribute(link, "rel")
		if rel == "icon" || rel == "shortcut icon" || rel == "apple-touch-icon" {
			return node_get_attribute(link, "href"), true
		}
	}
	return "", false
}
func favicon_try_from_url(uri string) string {
	c := curl.NewCurl("")

	cache, err := c.Get(uri)
	fmt.Println(cache)
	if err != nil {
		return ""
	}
	// text/html, text/xml, image
	m := strings.Split(cache.Mime, "/")
	switch m[0] {
	case "text":
		f, err := os.Open(cache.Local)
		if err == nil {
			defer f.Close()
			n, err := html.Parse(f)
			if err == nil {
				if u, ok := icon_from_link_rel(n); ok {
					return u
				}
			}
		}

	case "image":
		return uri
	}
	return ""
}
