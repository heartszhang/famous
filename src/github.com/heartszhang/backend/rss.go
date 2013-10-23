package backend

import (
	"encoding/xml"
	"io"
	"os"
)

// private wrapper around the RssFeed which gives us the <rss>..</rss> xml
type rss struct {
	XMLName xml.Name    `xml:"rss"`
	Channel rss_channel `xml:"channel" json:"channel,omitempty"`
	Version string      `xml:"version,attr"` // 2.0
}

type rss_channel struct {
	Title         string     `xml:"title, omityempty"` // required  unique?
	Link          string     `xml:"link, omityempty"`  // site's url
	Docs          string     `xml:"docs, omitempty"`   // rss link, be this rss url
	Description   string     `xml:"description, omitempty"`
	PubDate       string     `xml:"pubDate,omitempty"` // created or updated
	LastBuildDate string     `xml:"lastBuildDate"`
	TTL           uint       `xml:"ttl"` // minitues
	Category      []string   `xml:"category,omitempty"`
	Items         []rss_item `xml:"item"`
	Image         rss_image  `xml:"image"` // a img can be displayed
}

type rss_image struct {
	Url   string `xml:"url, omitempty"`
	Title string `xml:"title, omitempty"`
	Link  string `xml:"link, omitempty"` // may be same as channel's link
}

type rss_enclosure struct {
	Url    string `xml:"url,attr"`
	Length int64  `xml:"length,attr"` // bytes
	Type   string `xml:"type,attr"`   // mime-type
}

type rss_item struct {
	Title       string        `xml:"title"`             // required
	Link        string        `xml:"link"`              // required, http://nytimes.com/2004/12/07FEST.html
	PubDate     string        `xml:"pubDate,omitempty"` // created or updated
	Category    []string      `xml:"category"`
	Author      string        `xml:"creator,omitempty"`    // email address of the author
	Description string        `xml:"description"`          // required
	Guid        string        `xml:"guid,omitempty"`       // dont care
	Comments    string        `xml:"comments, omitempty"`  // comments url
	Enclosure   rss_enclosure `xml:"enclosure, omitempty"` //attachment
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

func feed_entries_from_rss2(filepath string) ([]FeedEntry, error) {
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
		v[idx] = i.to_feed_entry(this.Docs)
	}
	return v
}

func (this rss_channel) to_source() FeedSource {
	v := FeedSource{
		Name:        this.Title,
		Uri:         this.Docs,
		Local:       "", // filled later
		Period:      this.TTL,
		TouchAt:     unixtime_nano_rfc822(this.LastBuildDate),
		Category:    feed_category_root,
		Type:        feed_type_rss,
		Disabled:    false,
		EnableProxy: false,
		PubDate:     unixtime_nano_rfc822(this.PubDate),
		WebSite:     this.Link,
		Tags:        this.Category,
		Description: this.Description,
		Media: FeedMedia{
			media_type: feed_media_type_image,
			Title:      this.Image.Title,
			Uri:        this.Image.Url,
		},
	}
	return v
}
func (this rss_enclosure) to_feed_media(mt uint) FeedMedia {
	return FeedMedia{
		media_type: feed_media_type_image,
		Length:     this.Length,
		Uri:        this.Url,
		Mime:       this.Type,
	}
}

func mime_to_feedmediatype(mime string) uint {
	return feed_media_type_none
}

func (this rss_item) to_feed_entry(feed_url string) FeedEntry {
	v := FeedEntry{
		Source:   feed_url, // rss link
		Type:     feed_type_rss,
		Uri:      this.Link,
		PubDate:  unixtime_nano_rfc822(this.PubDate),
		Author:   FeedAuthor{Email: this.Author},
		Summary:  this.Description,
		Tags:     this.Category,
		Title:    FeedTitle{Main: this.Title},
		Category: feed_category_root,
	}
	switch mime_to_feedmediatype(this.Enclosure.Type) {
	case feed_media_type_image:
		v.Images = append(v.Images, FeedImage(this.Enclosure.to_feed_media(feed_media_type_image)))
	case feed_media_type_video:
		v.Video = append(v.Video, this.Enclosure.to_feed_media(feed_media_type_video))
	case feed_media_type_audio:
		v.Audio = append(v.Audio, this.Enclosure.to_feed_media(feed_media_type_audio))
	case feed_media_type_url:
		fallthrough
	default:
		v.Links = append(v.Links,
			FeedLink{
				media_type: feed_media_type_url,
				Uri:        this.Enclosure.Url,
				Length:     this.Enclosure.Length,
			})
	}
	return v
}

// file has been converted to utf-8, so we just ignore internal encoding-declaration
func charset_reader_passthrough(charset string, input io.Reader) (io.Reader, error) {
	return input, nil
}
