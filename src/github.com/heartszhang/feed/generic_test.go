package feed

import (
	"testing"
)

func TestMakeFeedSource(t *testing.T) {
	maker := NewFeedMaker(`E:\sourcesafe\famous\src\github.com\heartszhang\backend\run\data\a.xml`, "http://juetuzhi.net/feed")
	source, _, err := maker.MakeFeed()
	t.Log(source, err)
}
