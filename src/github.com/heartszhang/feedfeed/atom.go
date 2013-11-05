package feedfeed

import (
	"encoding/xml"
	"os"
	"time"
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
	Title      string          `xml:"title"`
	Id         string          `xml:"id"`
	Updated    string          `xml:"updated"`
	Summary    atom_text       `xml:"summary,omitempty"`
	Content    atom_text       `xml:"content,omitempty"`
	Links      []atom_link     `xml:"link"`
	Authors    []atom_person   `xml:"author,omitempty"`
	Categories []atom_category `xml:"category,omitempty"`
}

type atom_feed struct { // to feed_source
	XMLName    xml.Name        `xml:"http://www.w3.org/2005/Atom feed"`
	Title      string          `xml:"title"`
	Subtitle   string          `xml:"subtitle"`
	Id         string          `xml:"id"`
	Updated    string          `xml:"updated"` // rfc-822
	Links      []atom_link     `xml:"link"`
	Authors    []atom_person   `xml:"author"`
	Entries    []atom_entry    `xml:"entry"`
	Categories []atom_category `xml:"category"`
}

/*
func CreateFeedSourceAtom(filepath string) (FeedSource, error) {
	return feed_source_create_atom(filepath)
}

func CreateFeedEntriesAtom(filepath string) ([]FeedEntry, error) {
	return feed_entries_create_atom(filepath)
}
*/
func feedsource_from_atom(filepath string) (FeedSource, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return FeedSource{}, err
	}
	defer f.Close()

	var v atom_feed
	decoder := xml.NewDecoder(f)
	decoder.CharsetReader = charset_reader_passthrough

	err = decoder.Decode(&v)
	x := v.to_feed_soruce()
	return x, err
}

func (this atom_feed) link() string {
	for _, l := range this.Links {
		if l.Rel == "alternate" || l.Rel == "" {
			return l.Href
		}
	}
	return ""
}
func (this atom_feed) to_feed_soruce() FeedSource {
	f := FeedSource{
		Name:        this.Title,
		Uri:         this.docs(),
		Local:       "",
		Period:      _2hours,
		Deadline:    unixtime_nano_rfc822(this.Updated) + int64(_2hours*time.Hour),
		Type:        Feed_type_atom,
		Disabled:    false,
		EnableProxy: false,
		Update:      unixtime_nano_rfc822(this.Updated),
		WebSite:     this.link(),
		Media:       nil,
		Description: this.Subtitle,
	}
	f.Tags = make([]string, len(this.Categories))
	for i, c := range this.Categories {
		f.Tags[i] = c.Term
	}
	return f
}

func feedentries_from_atom(filepath string) ([]FeedEntry, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return []FeedEntry{}, err
	}
	defer f.Close()
	decoder := xml.NewDecoder(f)
	decoder.CharsetReader = charset_reader_passthrough
	var (
		v   atom_feed
		fes []FeedEntry
	)
	err = decoder.Decode(&v)
	if err == nil {
		fes = v.extract_entries()
	}
	return fes, err
}

func (this atom_feed) extract_entries() []FeedEntry {
	v := make([]FeedEntry, len(this.Entries))
	for idx, e := range this.Entries {
		v[idx] = e.to_feed_entry()
	}
	return v
}

func (this atom_entry) to_feed_entry() FeedEntry {
	e := FeedEntry{
		Flags:   0,
		Source:  "",
		Type:    Feed_type_atom,
		Uri:     this.link(),
		Title:   FeedTitle{Main: this.Title},
		PubDate: unixtime_nano_rfc822(this.Updated),
		Summary: this.Summary.Body,
	}
	if this.Content.Body != "" {
		e.Content = &FeedContent{FullText: this.Content.Body}
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
