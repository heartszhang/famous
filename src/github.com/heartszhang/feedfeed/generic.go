package feedfeed

import (
	"encoding/xml"
	//	"fmt"
	"os"
	"time"
)

type feed_error string

func (this feed_error) Error() string {
	return string(this)
}

func MakeFeed(filepath string) (FeedSource, []FeedEntry, error) {
	var (
		fs  FeedSource
		fes []FeedEntry
		err error
	)
	t := DetectFeedSourceType(filepath)
	switch t {
	case Feed_type_atom:
		fs, fes, err = feed_from_atom(filepath)
	case Feed_type_rss:
		fs, fes, err = feed_from_rss(filepath)
	default:
		fs, fes, err = FeedSource{}, nil, feed_error("invalid format")
	}
	fs.LastTouch = time.Now().Unix()
	fs.LastUpdate = fs.Update
	fs.NextTouch = int64(fs.Period)*60 + fs.LastTouch
	return fs, fes, err
}
func MakeFeedSource(filepath string) (FeedSource, error) {
	fs, _, err := MakeFeed(filepath)
	return fs, err
}
func MakeFeedEntries(filepath string) ([]FeedEntry, error) {
	_, fes, err := MakeFeed(filepath)
	return fes, err
}

type feed_sketch struct {
	XMLName xml.Name
	XML     string `xml:",innerxml"`
}

func DetectFeedSourceType(filepath string) uint {
	f, err := os.Open(filepath)
	if err != nil {
		return Feed_type_unknown
	}
	defer f.Close()

	decoder := xml.NewDecoder(f)
	decoder.CharsetReader = charset_reader_passthrough

	var (
		v feed_sketch
	)
	err = decoder.Decode(&v)
	return FeedSourceTypes[v.XMLName.Local]
}
