package feed

import (
	"encoding/xml"
	"github.com/heartszhang/unixtime"
	"os"
)

type atom_link struct {
	Rel    string `xml:"rel,attr,omitempty"`
	Href   string `xml:"href,attr,omitempty"`
	Type   string `xml:"type,attr,omitempty"`
	Title  string `xml:"title,attr,omitempty"`
	Length uint64 `xml:"length"`
}
type gd_image struct {
	Width  int    `xml:"width"`
	Height int    `xml:"height"`
	Src    string `xml:"src,omitempty"`
	Rel    string `xml:"rel,omitempty"`
}
type atom_person struct {
	Name   string    `xml:"name,omitempty"`
	Uri    string    `xml:"uri,omitempty"`
	Email  string    `xml:"email,omitempty"`
	Avatar *gd_image `xml:"image,omitempty"`
}

type atom_text struct {
	Type string `xml:"type,attr,omitempty"`
	Body string `xml:",chardata"` // omitempty cannot be used
}

type atom_category struct {
	Term   string `xml:"term,attr,omitempty"`
	Scheme string `xml:"scheme,omitempty"`
	Label  string `xml:"label,omitempty"`
}

type atom_entry struct { // to feed_entry
	Title      string            `xml:"title"`
	Id         string            `xml:"id"`
	Updated    unixtime.UnixTime `xml:"updated"`
	Summary    atom_text         `xml:"summary,omitempty"`
	Content    atom_text         `xml:"content,omitempty"`
	Links      []atom_link       `xml:"link"`
	Authors    []atom_person     `xml:"author,omitempty"`
	Categories []atom_category   `xml:"category,omitempty"`
}

type atom_feed struct { // to feed_source
	XMLName    xml.Name          `xml:"http://www.w3.org/2005/Atom feed"`
	Title      string            `xml:"title"`
	Subtitle   string            `xml:"subtitle,omitempty"`
	Id         string            `xml:"id"`
	Updated    unixtime.UnixTime `xml:"updated"` // rfc-822
	Logo       string            `xml:"logo,omitempty"`
	Links      []atom_link       `xml:"link,omitempty"`
	Authors    []atom_person     `xml:"author,omitempty"`
	Entries    []atom_entry      `xml:"entry,omitempty"`
	Categories []atom_category   `xml:"category,omitempty"`
}

func feed_from_atom(filepath string) (FeedSource, []FeedEntry, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return FeedSource{Local: filepath}, nil, err
	}
	defer f.Close()

	var v atom_feed
	decoder := xml.NewDecoder(f)
	decoder.CharsetReader = charset_reader_passthrough

	err = decoder.Decode(&v)
	x := v.to_feed_source(filepath)
	x.Local = filepath
	fes := v.extract_entries()
	return x, fes, err
}

func atom_query_selector(links []atom_link, rel string) atom_link {
	for _, l := range links {
		if l.Rel == rel {
			return l
		}
	}
	return atom_link{}
}

const (
	_2hours = 2 * 60 // minutes
)

func (this atom_feed) to_feed_source(local string) FeedSource {
	f := FeedSource{
		FeedSourceMeta: FeedSourceMeta{
			Name:        this.Title,
			Uri:         this.self(),
			Period:      _2hours,
			Logo:        this.Logo,
			Type:        Feed_type_atom,
			WebSite:     this.website(),
			Description: this.Subtitle,
		},
		Update:         this.Updated,
		Local:          local,
		SubscribeState: FeedSourceSubscribeStateSubscribed,
		EnableProxy:    0,
	}
	f.Tags = make([]string, len(this.Categories))
	for i, c := range this.Categories {
		f.Tags[i] = c.Term
	}
	if len(this.Authors) > 0 && this.Logo == "" && this.Authors[0].Avatar != nil {
		f.Logo = this.Authors[0].Avatar.Src
	}
	if f.Logo == "" {
		f.Logo = favicon(f.WebSite)
	}
	return f
}

func (this atom_feed) extract_entries() []FeedEntry {
	v := make([]FeedEntry, len(this.Entries))
	for idx, e := range this.Entries {
		v[idx] = e.to_feed_entry(this.self())
	}
	return v
}

func (this atom_entry) to_feed_entry(source string) FeedEntry {
	e := FeedEntry{
		FeedEntryMeta: FeedEntryMeta{
			Parent:  source,
			Type:    Feed_type_atom,
			Uri:     this.website(),
			Title:   FeedTitle{Main: this.Title},
			PubDate: this.Updated,
			Summary: this.Summary.Body,
		},
		Flags: 0,
	}
	if this.Content.Body != "" {
		//		e.Content = &FeedContent{FullText: this.Content.Body}
		e.Content = this.Content.Body
	}
	if len(this.Authors) > 0 {
		auth := this.Authors[0].to_feedauthor()
		e.Author = &auth
	}
	return e
}

func (this atom_feed) self() string { // rel = self
	l := atom_query_selector(this.Links, link_rel_self)
	return l.Href
}

func (this atom_entry) website() string { // rel = alternate
	l := atom_query_selector(this.Links, link_rel_alternate).Href
	if l == "" {
		l = atom_query_selector(this.Links, "").Href
	}
	return l
}

func (this atom_person) to_feedauthor() FeedAuthor {
	return FeedAuthor{Name: this.Name, Email: this.Email}
}
func (this atom_feed) website() string {
	l := atom_query_selector(this.Links, "")
	if l.Href == "" {
		l = atom_query_selector(this.Links, link_rel_alternate)
	}
	return l.Href
}
