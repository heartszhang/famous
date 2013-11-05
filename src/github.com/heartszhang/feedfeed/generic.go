package feedfeed

import (
	"encoding/xml"
	"fmt"
	"os"
)

func MakeFeedSource(filepath string) (FeedSource, error) {
	t := DetectFeedSourceType(filepath)
	switch t {
	case Feed_type_atom:
		return feedsource_from_atom(filepath)
	case Feed_type_rss:
		return feedsource_from_rss(filepath)
	default:
		return FeedSource{}, fmt.Errorf("invalid format")
	}
}
func MakeFeedEntries(filepath string) ([]FeedEntry, error) {
	t := DetectFeedSourceType(filepath)
	switch t {
	case Feed_type_atom:
		return feedentries_from_atom(filepath)
	case Feed_type_rss:
		return feedentries_from_rss(filepath)
	default:
		return []FeedEntry{}, fmt.Errorf("invalid format")
	}
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
