package feedfeed

import (
	"encoding/xml"
	"github.com/heartszhang/unixtime"
	"os"
)

const (
	atom_link_rel_self      = "self"
	atom_link_rel_related   = "related"
	atom_link_rel_alternate = "alternate"
	atom_link_rel_enclosure = "enclosure"
	atom_link_rel_via       = "via"
)

type atom_link struct {
	Rel    string `xml:"rel,attr,omitempty"`
	Href   string `xml:"href,attr,omitempty"`
	Type   string `xml:"type,attr,omitempty"`
	Title  string `xml:"title,attr,omitempty"`
	Length uint64 `xml:"length"`
}

type atom_person struct {
	Name  string `xml:"name,omitempty"`
	Uri   string `xml:"uri,omitempty"`
	Email string `xml:"email,omitempty"`
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
	Subtitle   string            `xml:"subtitle"`
	Id         string            `xml:"id"`
	Updated    unixtime.UnixTime `xml:"updated"` // rfc-822
	Logo       string            `xml:"logo,omitempty"`
	Links      []atom_link       `xml:"link"`
	Authors    []atom_person     `xml:"author"`
	Entries    []atom_entry      `xml:"entry"`
	Categories []atom_category   `xml:"category"`
}

/*
func feedsource_from_atom(filepath string) (FeedSource, error) {
	x, _, err := feed_from_atom(filepath)
	return x, err
}
*/
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
	x := v.to_feed_source()
	x.Local = filepath
	fes := v.extract_entries()
	return x, fes, err
}

func (this atom_feed) link() string {
	for _, l := range this.Links {
		if l.Rel == "alternate" || l.Rel == "" {
			return l.Href
		}
	}
	return ""
}

const (
	_2hours = 2 * 60 // minutes
)

func (this atom_feed) to_feed_source() FeedSource {
	f := FeedSource{
		Name:        this.Title,
		Uri:         this.docs(),
		Period:      _2hours,
		Logo:        this.Logo,
		Type:        Feed_type_atom,
		Disabled:    false,
		EnableProxy: false,
		Update:      this.Updated,
		WebSite:     this.link(),
		Description: this.Subtitle,
	}
	f.Tags = make([]string, len(this.Categories))
	for i, c := range this.Categories {
		f.Tags[i] = c.Term
	}
	return f
}

/*
func feedentries_from_atom(filepath string) ([]FeedEntry, error) {
	_, fes, err := feed_from_atom(filepath)
	return fes, err
}
*/
func (this atom_feed) extract_entries() []FeedEntry {
	v := make([]FeedEntry, len(this.Entries))
	for idx, e := range this.Entries {
		v[idx] = e.to_feed_entry(this.docs())
	}
	return v
}

func (this atom_entry) to_feed_entry(source string) FeedEntry {
	e := FeedEntry{
		Flags:   0,
		Parent:  source,
		Type:    Feed_type_atom,
		Uri:     this.link(),
		Title:   FeedTitle{Main: this.Title},
		PubDate: this.Updated,
		Summary: this.Summary.Body,
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

func (this atom_feed) docs() string { // rel = self
	for _, l := range this.Links {
		if l.Rel == "self" {
			return l.Href
		}
	}
	return ""
}

func (this atom_entry) link() string { // rel = alternate
	for _, l := range this.Links {
		if l.Rel == "alternate" || l.Rel == "" {
			return l.Href
		}
	}
	return ""
}

func (this atom_person) to_feedauthor() FeedAuthor {
	return FeedAuthor{Name: this.Name, Email: this.Email}
}
