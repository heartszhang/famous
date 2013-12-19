package feed

import (
	"encoding/xml"
	"io"
	"strings"

	"github.com/heartszhang/unixtime"
)

type rss struct {
	XMLName xml.Name    `xml:"rss"`
	Version string      `xml:"version,attr"` // 2.0
	Channel rss_channel `xml:"channel"`
}

type rss_channel struct {
	Title           string        `xml:"title,omityempty"` // required  unique?
	Links           []rss_link    `xml:"link,omitempty"`
	Description     string        `xml:"description,omitempty"`
	Language        string        `xml:"language,omitempty"`
	LastBuildDate   unixtime.Time `xml:"lastBuildDate,omitemptty"`
	PubDate         unixtime.Time `xml:"pubDate,omitempty"`
	Docs            string        `xml:"docs,omitempty"`
	Cloud           string        `xml:"cloud,omitempty"`
	UpdatePeriod    string        `xml:"http://purl.org/rss/1.0/modules/syndication/ updatePeriod,omityempty"`
	UpdateFrequency int64         `xml:"http://purl.org/rss/1.0/modules/syndication/ updateFrequency,omityempty"`
	TTL             int64         `xml:"ttl"` // minitues
	Categories      []string      `xml:"category,omitempty"`
	Image           *rss_image    `xml:"image,omitempty"` // a img can be displayed
	Items           []rss_item    `xml:"item,omitempty"`
}

type rss_item struct {
	Title       string          `xml:"title"`             // required
	Link        string          `xml:"link"`              // required, like http://nytimes.com/2004/12/07FEST.html
	PubDate     unixtime.Time   `xml:"pubDate,omitempty"` // created or updated
	Categories  []string        `xml:"category,omitempty"`
	Author      string          `xml:"author,omitempty"`   // email address of the author
	Description string          `xml:"description"`        // required
	Guid        string          `xml:"guid,omitempty"`     // dont care
	Comments    string          `xml:"comments,omitempty"` // comments url
	FullText    string          `xml:"http://purl.org/rss/1.0/modules/content/ encoded,omitmepty"`
	Enclosure   []rss_enclosure `xml:"enclosure,omitempty"` //attachment
	Source      string          `xml:"source,omitempty"`
	Medias      []media_content `xml:"http://search.yahoo.com/mrss/ content,omitempty"`
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

// refer http://search.yahoo.com/mrss/
type media_content struct {
	Url    string `xml:"url,attr,omitempty"`
	Medium string `xml:"medium,attr,omitempty"`
	Title  string `xml:"title,omitempty"`
}

type rss_link struct { // many rss files use atom's link
	atom_link `xml:",inline"`
	Link      string `xml:",chardata"`
}

func feed_from_rss(f io.Reader, uri string) (FeedSource, []FeedEntry, error) {
	var (
		v   rss
		fs  FeedSource
		fes []FeedEntry
	)
	err := new_xml_decoder(f).Decode(&v)
	fs = v.Channel.export_source(uri)
	fes = v.Channel.export_entries(fs.Uri)

	return fs, fes, err
}

func (this rss_channel) export_entries(feeduri string) []FeedEntry {
	var v []FeedEntry
	for _, i := range this.Items {
		v = append(v, i.to_feedentry(feeduri))
	}
	return v
}

func (this rss_channel) export_source(uri string) FeedSource {
	v := FeedSource{
		Name:        this.Title,
		Uri:         this.self(uri),
		Period:      this.ttl(),
		Update:      this.update(),
		Type:        Feed_type_rss,
		Tags:        this.Categories,
		Description: this.Description,
		WebSite:     this.website(),
		Logo:        this.logo(),
		Hub:         this.hub(),
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
func (this rss_channel) logo() string {
	if this.Image != nil {
		return this.Image.Url
	}
	return ""
}
func (this rss_channel) self(downloaduri string) string { // rel = self
	if downloaduri != "" {
		return downloaduri
	}
	return query_selector(this.Links, link_rel_self)
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

func query_selector(links []rss_link, rel string) string {
	var alinks []atom_link
	for _, l := range links {
		alinks = append(alinks, l.to_atomlink())
	}
	return atom_query_selector(alinks, rel)
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

func (this rss_item) to_feedentry(feed_url string) FeedEntry {
	v := FeedEntry{
		Parent:  feed_url, // rss link
		Title:   this.Title,
		Uri:     this.Link,
		Type:    Feed_type_rss,
		PubDate: int64(this.PubDate),
		Summary: this.Description,
		Content: this.FullText,
		Tags:    this.Categories,
		Author:  this.Author,
	}
	for _, enclosure := range this.Enclosure {
		switch mime_to_feedmediatype(enclosure.Type) {
		case Feed_media_type_image:
			v.Images = append(v.Images, FeedMedia(enclosure.to_feedmedia()))
		case Feed_media_type_video:
			v.Videos = append(v.Videos, enclosure.to_feedmedia())
		case Feed_media_type_audio:
			v.Audios = append(v.Audios, enclosure.to_feedmedia())
		}
	}

	for _, media := range this.Medias {
		switch mime_to_feedmediatype(media.Medium) {
		case Feed_media_type_image:
			v.Images = append(v.Images, FeedMedia(media.to_feedmedia()))
		case Feed_media_type_video:
			v.Videos = append(v.Videos, media.to_feedmedia())
		case Feed_media_type_audio:
			v.Audios = append(v.Audios, media.to_feedmedia())
		}
	}
	return v
}

func (this rss_link) to_atomlink() atom_link {
	var v atom_link = this.atom_link
	if this.Link != "" {
		v.Href = this.Link
	}
	return v
}

func (this rss_enclosure) to_feedmedia() FeedMedia {
	return FeedMedia{
		Length: this.Length,
		Uri:    this.Url,
		Mime:   this.Type,
	}
}

func (this media_content) to_feedmedia() FeedMedia {
	return FeedMedia{
		Uri:   this.Url,
		Mime:  this.Medium,
		Title: this.Title,
	}
}

func (this rss_channel) update() int64 {
	var v int64
	if this.PubDate != 0 {
		v = int64(this.PubDate)
	}
	v = int64(this.LastBuildDate)
	return v
}
