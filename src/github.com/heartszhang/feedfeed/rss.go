package feedfeed

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"time"
)

func CreateFeedSourceRss2(rss2file string) (FeedSource, error) {
	return feed_source_create_rss2(rss2file)
}

func CreateFeedEntriesRss2(rss2file string) ([]FeedEntry, error) {
	return feed_entries_create_rss2(rss2file)
}

type rss_channel struct {
	Title           string     `xml:"title,omityempty"` // required  unique?
	Links           []rss_link `xml:"link,omitempty"`
	Description     string     `xml:"description,omitempty"`
	LastBuildDate   string     `xml:"lastBuildDate,omitemptty"`
	UpdatePeriod    string     `xml:"http://purl.org/rss/1.0/modules/syndication/ updatePeriod,omityempty"`
	UpdateFrequency uint       `xml:"http://purl.org/rss/1.0/modules/syndication/ updateFrequency,omityempty"`
	TTL             uint       `xml:"ttl"` // minitues
	Categories      []string   `xml:"category,omitempty"`
	Image           rss_image  `xml:"image"` // a img can be displayed
	Items           []rss_item `xml:"item"`
}
type rss_enclosure struct {
	Url    string `xml:"url,attr"`
	Length int64  `xml:"length,attr"` // bytes
	Type   string `xml:"type,attr"`   // mime-type
}

type rss_image struct {
	Url   string `xml:"url, omitempty"`
	Title string `xml:"title, omitempty"`
	Link  string `xml:"link, omitempty"` // may be same as channel's link
}

// http://search.yahoo.com/mrss/
type media_content struct {
	Url    string `xml:"url,attr,omitempty"`
	Medium string `xml:"medium,attr,omitempty"`
	Title  string `xml:"title,omitempty"`
}

type rss_item struct {
	Title       string          `xml:"title"`             // required
	Link        string          `xml:"link"`              // required, http://nytimes.com/2004/12/07FEST.html
	PubDate     string          `xml:"pubDate,omitempty"` // created or updated
	Categories  []string        `xml:"category"`
	Author      string          `xml:"author,omitempty"`    // email address of the author
	Description string          `xml:"description"`         // required
	Guid        string          `xml:"guid,omitempty"`      // dont care
	Comments    string          `xml:"comments, omitempty"` // comments url
	FullText    string          `xml:"http://purl.org/rss/1.0/modules/content/ encoded,omitmepty"`
	Enclosure   rss_enclosure   `xml:"enclosure, omitempty"` //attachment
	Medias      []media_content `xml:"http://search.yahoo.com/mrss/ content,omitempty"`
}

type rss_link struct {
	atom_link
	Link string `xml:",chardata"`
}

// private wrapper around the RssFeed which gives us the <rss>..</rss> xml
type rss struct {
	XMLName xml.Name    `xml:"rss"`
	Channel rss_channel `xml:"channel"`
	Version string      `xml:"version,attr"` // 2.0
}

// file has already converted to utf-8
func feed_source_create_rss2(filepath string) (FeedSource, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return FeedSource{}, err
	}
	defer f.Close()

	decoder := xml.NewDecoder(f)
	decoder.CharsetReader = charset_reader_passthrough

	var (
		v  rss
		fs FeedSource
	)
	err = decoder.Decode(&v)
	fs = v.Channel.to_source()
	return fs, err
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

func feed_entries_create_rss2(filepath string) ([]FeedEntry, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return []FeedEntry{}, err
	}
	defer f.Close()

	decoder := xml.NewDecoder(f)
	decoder.CharsetReader = charset_reader_passthrough

	var (
		v   rss
		fes []FeedEntry
	)
	err = decoder.Decode(&v)
	fes = v.Channel.to_entries()
	return fes, err
}

func (this rss_channel) to_entries() []FeedEntry {
	v := make([]FeedEntry, len(this.Items))
	for idx, i := range this.Items {
		v[idx] = i.to_feed_entry(this.self())
	}
	return v
}

func (this rss_channel) to_source() FeedSource {
	v := FeedSource{
		Name:        this.Title,
		Uri:         this.self(),
		Local:       "", // filled later
		Period:      this.ttl(),
		Deadline:    unixtime_nano_rfc822(this.LastBuildDate) + int64(this.ttl())*int64(time.Hour),
		Type:        Feed_type_rss,
		Disabled:    false,
		EnableProxy: false,
		Update:      unixtime_nano_rfc822(this.LastBuildDate),
		Tags:        this.Categories,
		Description: this.Description,
		WebSite:     this.website(),
		Media:       nil,
	}
	if this.Image.Url != "" {
		v.Media = &FeedMedia{
			media_type: Feed_media_type_image,
			Title:      this.Image.Title,
			Uri:        this.Image.Url,
		}
	}
	fmt.Println("rss:", this.Title)
	return v
}

const (
	minute  uint = 1
	hourly       = 60 * minute
	daily        = 24 * hourly
	weekly       = 7 * daily
	monthly      = 30 * daily
	year         = 365 * daily
)

var (
	sd_update_period = map[string]uint{
		"hourly":  hourly,
		"daily":   daily,
		"weekly":  weekly,
		"monthly": monthly,
		"year":    year,
	}
)

func (this rss_channel) ttl() uint {
	v := this.TTL
	if v == 0 {
		v = sd_update_period[this.UpdatePeriod] * this.UpdateFrequency
	}
	if v == 0 {
		v = 2 * hourly
	}
	return v
}

func (this rss_channel) self() string { // rel = self
	for _, l := range this.Links {
		if l.Rel == "self" {
			return l.Href
		}
	}
	return ""
}

func (this rss_channel) website() string { // rel = alternate
	for _, l := range this.Links {
		if l.Rel == "alternate" || l.Rel == "" {
			return l.Href
		}
	}
	return ""
}

func (this rss_item) save_content() *FeedContent {
	if this.FullText == "" {
		return nil
	}
	return &FeedContent{FullText: this.FullText}
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
	return Feed_media_type_none
}

func (this rss_item) to_feed_entry(feed_url string) FeedEntry {
	v := FeedEntry{
		Source:  feed_url, // rss link
		Type:    Feed_type_rss,
		Uri:     this.Link,
		PubDate: unixtime_nano_rfc822(this.PubDate),
		//		Author:   FeedAuthor{Email: this.Author},
		Summary: this.Description,
		Tags:    this.Categories,
		Title:   FeedTitle{Main: this.Title},
		Content: this.save_content(),
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
	fmt.Println(v.Title)
	return v
}

// file has been converted to utf-8, so we just ignore internal encoding-declaration
func charset_reader_passthrough(charset string, input io.Reader) (io.Reader, error) {
	return input, nil
}
