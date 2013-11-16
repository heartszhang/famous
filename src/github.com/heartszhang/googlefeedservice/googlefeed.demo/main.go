package main

import (
	"fmt"
	g "github.com/heartszhang/googlefeedservice"
)

func main() {
	fi := g.NewGoogleFeedApi("http://iweizhi2.duapp.com", "")
	f, err := fi.Find("caoliu", "")
	fmt.Println(f, err)
	//	s, e, err := fi.Load("http://sex8.cc/rss-htm-fid-4.html", "", 10, false)
	//	fmt.Println(s.Name, e, err)
}
