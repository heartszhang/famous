package feedfeed

import (
	"encoding/xml"
	"github.com/heartszhang/unixtime"
	"os"
)

type feed_error string

func (this feed_error) Error() string {
	return string(this)
}

type FeedMaker interface {
	MakeFeed() (FeedSource, []FeedEntry, error)
	MakeFeedSource() (FeedSource, error)
	MakeFeedEntries() ([]FeedEntry, error)
}

type feed_maker struct {
	cache  string
	source string
}

func NewFeedMaker(filepath, source string) FeedMaker {
	return feed_maker{cache: filepath, source: source}
}
func (this feed_maker) MakeFeed() (FeedSource, []FeedEntry, error) {
	var (
		fs  FeedSource
		fes []FeedEntry
		err error
	)
	t := DetectFeedSourceType(this.cache)
	switch t {
	case Feed_type_atom:
		fs, fes, err = feed_from_atom(this.cache)
	case Feed_type_rss:
		fs, fes, err = feed_from_rss(this.cache)
	default:
		fs, fes, err = FeedSource{Uri: this.source}, nil, feed_error("invalid format")
	}
	if this.source == "" {
		this.source = fs.Uri
	}
	fs.Uri = this.source

	fs.LastTouch = unixtime.UnixTimeNow()
	fs.LastUpdate = fs.Update
	fs.NextTouch = unixtime.UnixTime(int64(fs.Period)*60 + int64(fs.LastTouch))
	count := len(fes)
	for i := 0; i < count; i++ {
		fes[i].Parent = this.source
	}
	return fs, fes, err
}

func (this feed_maker) MakeFeedSource() (FeedSource, error) {
	fs, _, err := this.MakeFeed()
	return fs, err
}
func (this feed_maker) MakeFeedEntries() ([]FeedEntry, error) {
	_, fes, err := this.MakeFeed()

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
