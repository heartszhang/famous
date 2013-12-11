package feed

import (
	"encoding/xml"
	"github.com/heartszhang/unixtime"
	"io"
	"net/url"
	"os"
	"strings"
)

type rss_channel struct {
	Title           string            `xml:"title,omityempty"` // required  unique?
	Links           []rss_link        `xml:"link,omitempty"`
	Description     string            `xml:"description,omitempty"`
	LastBuildDate   unixtime.UnixTime `xml:"lastBuildDate,omitemptty"`
	UpdatePeriod    string            `xml:"http://purl.org/rss/1.0/modules/syndication/ updatePeriod,omityempty"`
	UpdateFrequency int64             `xml:"http://purl.org/rss/1.0/modules/syndication/ updateFrequency,omityempty"`
	TTL             int64             `xml:"ttl"` // minitues
	Categories      []string          `xml:"category,omitempty"`
	Image           *rss_image        `xml:"image,omitempty"` // a img can be displayed
	Items           []rss_item        `xml:"item"`
}
type rss_enclosure struct {
	Url    string `xml:"url,attr,omitempty"`
	Length int64  `xml:"length,attr"`         // bytes
	Type   string `xml:"type,attr,omitempty"` // mime-type
}

type rss_image struct {
	Url   string `xml:"url,omitempty"`
	Title string `xml:"title,omitempty"`
	Link  string `xml:"link,omitempty"` // may be same as channel's link
}

// http://search.yahoo.com/mrss/
type media_content struct {
	Url    string `xml:"url,attr,omitempty"`
	Medium string `xml:"medium,attr,omitempty"`
	Title  string `xml:"title,omitempty"`
}

type rss_item struct {
	Title       string            `xml:"title"`             // required
	Link        string            `xml:"link"`              // required, http://nytimes.com/2004/12/07FEST.html
	PubDate     unixtime.UnixTime `xml:"pubDate,omitempty"` // created or updated
	Categories  []string          `xml:"category,omitempty"`
	Author      string            `xml:"author,omitempty"`   // email address of the author
	Description string            `xml:"description"`        // required
	Guid        string            `xml:"guid,omitempty"`     // dont care
	Comments    string            `xml:"comments,omitempty"` // comments url
	FullText    string            `xml:"http://purl.org/rss/1.0/modules/content/ encoded,omitmepty"`
	Enclosure   rss_enclosure     `xml:"enclosure,omitempty"` //attachment
	Medias      []media_content   `xml:"http://search.yahoo.com/mrss/ content,omitempty"`
}

type rss_link struct {
	atom_link `xml:",inline"`
	Link      string `xml:",chardata"`
}

// private wrapper around the RssFeed which gives us the <rss>..</rss> xml
type rss struct {
	XMLName xml.Name    `xml:"rss"`
	Channel rss_channel `xml:"channel"`
	Version string      `xml:"version,attr"` // 2.0
}

// file has already converted to utf-8
func feed_from_rss(filepath string) (FeedSource, []FeedEntry, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return FeedSource{}, nil, err
	}
	defer f.Close()

	decoder := xml.NewDecoder(f)
	decoder.CharsetReader = charset_reader_passthrough

	var (
		v   rss
		fs  FeedSource
		fes []FeedEntry
	)
	err = decoder.Decode(&v)
	fs = v.Channel.to_source(filepath)
	fes = v.Channel.to_entries()

	return fs, fes, err
}

func (this rss_channel) to_entries() []FeedEntry {
	v := make([]FeedEntry, len(this.Items))
	for idx, i := range this.Items {
		v[idx] = i.to_feed_entry(this.self())
	}
	return v
}

func (this rss_channel) to_source(local string) FeedSource {
	v := FeedSource{
		FeedSourceMeta: FeedSourceMeta{
			Name:        this.Title,
			Uri:         this.self(),
			Period:      this.ttl(),
			Type:        Feed_type_rss,
			Tags:        this.Categories,
			Description: this.Description,
			WebSite:     this.website(),
		},
		SubscribeState: FeedSourceSubscribeStateSubscribed,
		EnableProxy:    0,
		Local:          local, // filled later
		Update:         this.LastBuildDate,
	}

	if this.Image != nil {
		v.Logo = this.Image.Url // this may be relative
	}
	if v.Logo == "" {
		v.Logo = favicon(this.website())
	}
	return v
}

const (
	minute  int64 = 1
	hourly        = 60 * minute
	daily         = 24 * hourly
	weekly        = 7 * daily
	monthly       = 30 * daily
	year          = 365 * daily
)

var (
	sd_update_period = map[string]int64{
		"hourly":  hourly,
		"daily":   daily,
		"weekly":  weekly,
		"monthly": monthly,
		"year":    year,
	}
)

func (this rss_channel) ttl() int64 {
	v := this.TTL
	if v == 0 {
		v = sd_update_period[this.UpdatePeriod] * this.UpdateFrequency
	}
	if v == 0 {
		v = _2hours
	}
	return v
}

func (this rss_channel) self() string { // rel = self
	return query_selector(this.Links, link_rel_self)
}
func query_selector(links []rss_link, rel string) string {
	for _, l := range links {
		if l.Rel == rel {
			x := l.Href
			if x == "" {
				x = l.Link
			}
			return x
		}
	}
	return ""
}

func (this rss_channel) website() string { // rel = alternate
	l := query_selector(this.Links, "")
	if l == "" {
		l = query_selector(this.Links, link_rel_alternate)
	}
	return l
}
func (this rss_channel) hub() string {
	return query_selector(this.Links, link_rel_hub)
}

func (this rss_enclosure) to_feed_media(mt uint) FeedMedia {
	return FeedMedia{
		media_type: Feed_media_type_image,
		Length:     this.Length,
		Uri:        this.Url,
		Mime:       this.Type,
	}
}

func mime_to_feedmediatype(mime string) uint {
	t := Feed_media_type_none
	types := strings.Split(mime, "/")
	switch types[0] {
	case "image":
		t = Feed_media_type_image
	case "video":
		t = Feed_media_type_video
	case "audio":
		t = Feed_media_type_audio
	}
	return t
}

func (this rss_item) to_feed_entry(feed_url string) FeedEntry {
	v := FeedEntry{
		FeedEntryMeta: FeedEntryMeta{
			Parent:  feed_url, // rss link
			Type:    Feed_type_rss,
			Uri:     this.Link,
			PubDate: this.PubDate,
			Summary: this.Description,
			Tags:    this.Categories,
			Title:   FeedTitle{Main: this.Title},
			Content: this.FullText,
		},
	}
	if this.Author != "" {
		v.Author = &FeedAuthor{Email: this.Author}
	}

	switch mime_to_feedmediatype(this.Enclosure.Type) {
	case Feed_media_type_image:
		v.Images = append(v.Images, FeedMedia(this.Enclosure.to_feed_media(Feed_media_type_image)))
	case Feed_media_type_video:
		v.Videos = append(v.Videos, this.Enclosure.to_feed_media(Feed_media_type_video))
	case Feed_media_type_audio:
		v.Audios = append(v.Audios, this.Enclosure.to_feed_media(Feed_media_type_audio))
	case Feed_media_type_url:
		v.Links = append(v.Links,
			FeedLink{
				media_type: Feed_media_type_url,
				Uri:        this.Enclosure.Url,
				Length:     this.Enclosure.Length,
			})

	default:
	}
	return v
}

// file has been converted to utf-8, so we just ignore internal encoding-declaration
func charset_reader_passthrough(charset string, input io.Reader) (io.Reader, error) {
	return input, nil
}

func favicon(uri string) string {
	u, err := url.Parse(uri)
	if err == nil {
		u, err = u.Parse("/favicon.ico")
	}
	if err == nil {
		return u.String()
	}
	return ""
}
