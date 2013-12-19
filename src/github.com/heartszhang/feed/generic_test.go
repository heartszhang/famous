package feed

import (
	"os"
	"testing"
)

func TestAtomSource(t *testing.T) {
	f, _ := os.Open("atom.xml")
	maker := NewFeedMaker(f, "http://juetuzhi.net/feed")
	source, _, err := maker.MakeFeed()
	f.Close()
	t.Log(source, err)
}

func TestRssSource(t *testing.T) {
	f, _ := os.Open("rss.xml")
	maker := NewFeedMaker(f, "http://juetuzhi.net/feed")
	source, _, err := maker.MakeFeed()
	f.Close()
	t.Log(source, err)
}

func TestOpmlSource(t *testing.T) {
	f, _ := os.Open("feedly.opml")
	source, err := OpmlExportFeedSource(f)

	f.Close()
	t.Log(source, err)
}
