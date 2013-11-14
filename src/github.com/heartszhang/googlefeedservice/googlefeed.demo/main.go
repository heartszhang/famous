package main

import (
	"fmt"
	g "github.com/heartszhang/googlefeedservice"
)

func main() {
	fi := g.NewGoogleFeedApi("http://iweizhi2.duapp.com", "")
	//	r, err := fi.Find("caoliu", "")
	r, err := fi.Load("http://caoliuok.blogspot.com/feeds/posts/default", "", 1, false)
	fmt.Println(r.ResponseStatus, r.ResponseData.Feed, err)
}
