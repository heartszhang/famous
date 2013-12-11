package main

import (
	"flag"
	"fmt"
	"github.com/heartszhang/pubsub"
)

var (
	//	uri  = "http://feeds.feedburner.com/ftchina"
	//	uri2 = "http://www.bbc.co.uk/zhongwen/trad/index.xml"
	uri         = flag.String("uri", "http://feed.feedsky.com/199it", "feed rss url")
	method      = flag.String("method", "subscribe", "subscribe/unsubscribe/retrieve")
	provider    = flag.String("provider", "superfeedr", "superfeedr/google")
	verify_mode = flag.String("verify", "async", "async/sync")
)

func main() {
	flag.Parse()
	var sub pubsub.PubSubscriber
	if *provider == "google" {
		sub = pubsub.NewGooglePubSubscriber()
	} else {
		sub = pubsub.NewSuperFeedrPubSubscriber(*verify_mode, "Hearts", "Refresh")
	}
	switch *method {
	case "subscribe":
		code, err := sub.Subscribe(*uri)
		fmt.Println(code, err)
	case "unsubscribe":
		code, err := sub.Unsubscribe(*uri)
		fmt.Println(code, err)
	case "retrieve":
		s, err := sub.Retrieve(*uri, 1)
		fmt.Println(s, err)
	}
}
