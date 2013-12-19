package feed

import (
	"encoding/xml"
	"io"

	"github.com/heartszhang/unixtime"
)

const atom_ns = "http://www.w3.org/2005/Atom"

type atom_feed struct { // to feed_source
	XMLName    xml.Name        `xml:"http://www.w3.org/2005/Atom feed"`
	Title      string          `xml:"title"`
	Subtitle   string          `xml:"subtitle,omitempty"`
	Id         string          `xml:"id"`
	Updated    unixtime.Time   `xml:"updated"` // rfc-822
	Logo       string          `xml:"logo,omitempty"`
	Icon       string          `xml:"icon,omitempty`
	Links      []atom_link     `xml:"link,omitempty"`
	Authors    []atom_person   `xml:"author,omitempty"`
	Categories []atom_category `xml:"category,omitempty"`
	Entries    []atom_entry    `xml:"entry,omitempty"`
}

type atom_entry struct { // to feed_entry
	Title      string          `xml:"title"`
	Id         string          `xml:"id"`
	Updated    unixtime.Time   `xml:"updated"`
	Published  unixtime.Time   `xml:"published"`
	Summary    atom_text       `xml:"summary,omitempty"`
	Content    atom_text       `xml:"content,omitempty"`
	Source     string          `xml:"source,omitempty"`
	Links      []atom_link     `xml:"link"`
	Authors    []atom_person   `xml:"author,omitempty"`
	Categories []atom_category `xml:"category,omitempty"`
}

type atom_link struct {
	Rel    string `xml:"rel,attr,omitempty"`
	Href   string `xml:"href,attr,omitempty"`
	Type   string `xml:"type,attr,omitempty"`
	Title  string `xml:"title,attr,omitempty"`
	Length int64  `xml:"length"`
}

type gd_image struct { // google ajax feed api
	Width  int    `xml:"width"`
	Height int    `xml:"height"`
	Src    string `xml:"src,omitempty"`
	Rel    string `xml:"rel,omitempty"`
}
type atom_person struct {
	Name   string    `xml:"name,omitempty"`
	Uri    string    `xml:"uri,omitempty"`
	Email  string    `xml:"email,omitempty"`
	Avatar *gd_image `xml:"image,omitempty"` // no specified
}

type atom_text struct {
	Type    string `xml:"type,attr,omitempty"` // mime
	Content string `xml:",chardata"`
}

type atom_category struct {
	Term   string `xml:"term,attr,omitempty"`
	Scheme string `xml:"scheme,omitempty"`
	Label  string `xml:"label,omitempty"`
}

func feed_from_atom(f io.Reader, uri string) (FeedSource, []FeedEntry, error) {
	var v atom_feed
	err := new_xml_decoder(f).Decode(&v)
	x := v.to_feedsource(uri)
	fes := v.extract_feedentries(x.Uri)
	return x, fes, err
}

const (
	_2hours = 2 * 60 // minutes
)

func (this atom_feed) to_feedsource(uri string) FeedSource {
	return FeedSource{
		Name:        this.Title,
		Uri:         this.self(uri),
		Period:      _2hours,
		Logo:        this.logo(),
		Type:        Feed_type_atom,
		WebSite:     this.website(),
		Description: this.Subtitle,
		Hub:         this.hub(),
		Tags:        convert_categories(this.Categories),
		Update:      int64(this.Updated),
	}
}

func convert_categories(categories []atom_category) []string {
	var v []string
	for _, category := range categories {
		v = append(v, category.Term)
	}
	return v
}
func (this atom_feed) extract_feedentries(feeduri string) []FeedEntry {
	var v []FeedEntry
	for _, e := range this.Entries {
		v = append(v, e.to_feedentry(feeduri))
	}
	return v
}

func (this atom_entry) to_feedentry(feeduri string) FeedEntry {
	e := FeedEntry{
		Parent:  feeduri,
		Type:    Feed_type_atom,
		Uri:     this.website(),
		Title:   this.Title,
		PubDate: this.updated(),
		Summary: this.Summary.Content,
		Content: this.Content.Content,
		Author:  this.select_author(),
		Tags:    convert_categories(this.Categories),
	}
	for _, link := range this.Links {
		switch mime_to_feedmediatype(link.Type) {
		case Feed_media_type_image:
			e.Images = append(e.Images, FeedMedia(link.to_feedmedia()))
		case Feed_media_type_video:
			e.Videos = append(e.Videos, link.to_feedmedia())
		case Feed_media_type_audio:
			e.Audios = append(e.Audios, link.to_feedmedia())
		}
	}
	return e
}

func (this atom_entry) select_author() string {
	var name, avatar_name string
	for _, a := range this.Authors {
		if a.Avatar != nil {
			avatar_name = a.name()
		}
		name = a.name()
	}
	if avatar_name != "" {
		return avatar_name
	}
	return name
}

func (this atom_feed) self(downloadeduri string) string { // rel = self
	if downloadeduri != "" {
		return downloadeduri
	}
	return atom_query_selector(this.Links, link_rel_self)
}

func (this atom_entry) website() string { // rel = alternate
	l := atom_query_selector(this.Links, link_rel_alternate)
	if l == "" {
		l = atom_query_selector(this.Links, "")
	}
	return l
}

func (this atom_feed) website() string {
	l := atom_query_selector(this.Links, "")
	if l == "" {
		l = atom_query_selector(this.Links, link_rel_alternate)
	}
	return l
}

func (this atom_feed) hub() string {
	return atom_query_selector(this.Links, link_rel_hub)
}

func (this atom_person) name() string {
	if this.Name != "" {
		return this.Name
	}
	if this.Email != "" {
		return this.Email
	}
	return this.Uri
}

func atom_query_selector(links []atom_link, rel string) string {
	for _, l := range links {
		if l.Rel == rel {
			return l.Href
		}
	}
	return ""
}

func (this atom_feed) logo() string {
	if this.Logo != "" {
		return this.Logo
	}
	return this.Icon
}

func (this atom_link) to_feedmedia() FeedMedia {
	return FeedMedia{
		Uri:    this.Href,
		Mime:   this.Type,
		Title:  this.Title,
		Length: this.Length,
	}
}

func (this atom_entry) updated() int64 {
	var v int64
	if this.Updated != 0 {
		v = int64(this.Updated)
	}
	v = int64(this.Published)
	return v
}
