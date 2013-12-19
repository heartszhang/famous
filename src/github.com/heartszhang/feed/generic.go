package feed

import (
	"encoding/xml"
	"io"
)

type FeedMaker interface {
	MakeFeed() (FeedSource, []FeedEntry, error)
}

type feed_maker struct {
	io.ReadSeeker
	uri string
}

func NewFeedMaker(reader io.ReadSeeker, source string) FeedMaker {
	return feed_maker{ReadSeeker: reader, uri: source}
}

func (this feed_maker) MakeFeed() (FeedSource, []FeedEntry, error) {
	var (
		fs  FeedSource
		fes []FeedEntry
		err error
	)
	t := DetectFeedSourceType(this.ReadSeeker)
	this.ReadSeeker.Seek(0, 0)
	switch t {
	case Feed_type_atom:
		fs, fes, err = feed_from_atom(this.ReadSeeker, this.uri)
	case Feed_type_rss:
		fs, fes, err = feed_from_rss(this.ReadSeeker, this.uri)
	default:
		fs, fes, err = FeedSource{}, nil, feed_error(this.uri+": invalid format")
	}

	return fs, fes, err
}

type feed_sketch struct {
	XMLName xml.Name
	XML     string `xml:",innerxml"`
}

func DetectFeedSourceType(f io.Reader) uint {
	var v feed_sketch
	new_xml_decoder(f).Decode(&v)
	return FeedSourceType(v.XMLName.Local)
}
func new_xml_decoder(f io.Reader) *xml.Decoder {
	d := xml.NewDecoder(f)
	d.CharsetReader = charset_reader_passthrough
	return d
}

type feed_error string

func (this feed_error) Error() string {
	return string(this)
}
