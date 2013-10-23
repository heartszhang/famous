package backend

import (
	"encoding/xml"
	"io/ioutil"
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
	Href   string `xml:"href,attr, omitempty"`
	Type   string `xml:"type, attr, omitempty"`
	Title  string `xml:"title, attr, omitempty"`
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
	Title    string          `xml:"title"`
	Id       string          `xml:"id"`
	Updated  string          `xml:"updated"`
	Summary  atom_text       `xml:"summary,omitempty"`
	Content  atom_text       `xml:"content,omitempty"`
	Link     []atom_link     `xml:"link"`
	Author   []atom_person   `xml:"author,omitempty"`
	Category []atom_category `xml:"category,omitempty"`
}

type atom_feed struct { // to feed_source
	XMLName  xml.Name        `xml:"http://www.w3.org/2005/Atom feed"`
	Title    string          `xml:"title"`
	Subtitle string          `xml:"subtitle"`
	Id       string          `xml:"id"`
	Updated  string          `xml:"updated"` // rfc-822
	Link     []atom_link     `xml:"link"`
	Author   []atom_person   `xml:"author"`
	Entry    []atom_entry    `xml:"entry"`
	Category []atom_category `xml:"category"`
}

func feed_source_create_atom(filepath string) (FeedSource, error) {
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
	for _, l := range this.Link {
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
		TouchAt:     unixtime_nano_rfc822(this.Updated),
		Category:    feed_category_root,
		Type:        feed_type_atom,
		Disabled:    false,
		EnableProxy: false,
		PubDate:     unixtime_nano_rfc822(this.Updated),
		WebSite:     this.link(),
		Media:       FeedMedia{},
		Description: this.Subtitle,
	}
	f.Tags = make([]string, len(this.Category))
	for i, c := range this.Category {
		f.Tags[i] = c.Term
	}
	return f
}

func (this atom_entry) to_feed_entry() FeedEntry {
	e := FeedEntry{
		Flags:    0,
		Source:   "",
		Type:     feed_type_atom,
		Uri:      this.link(),
		Title:    FeedTitle{Main: this.Title},
		Author:   this.Author[0].to_feedauthor(),
		PubDate:  unixtime_nano_rfc822(this.Updated),
		Summary:  this.Summary.Body,
		Content:  this.save_content(),
		Category: feed_category_root,
	}
	return e
}

func (this atom_feed) docs() string { // rel = self
	for _, l := range this.Link {
		if l.Rel == "self" {
			return l.Href
		}
	}
	return ""
}

func (this atom_entry) link() string { // rel = alternate
	for _, l := range this.Link {
		if l.Rel == "alternate" || l.Rel == "" {
			return l.Href
		}
	}
	return ""
}
func mime_to_ext(mime string) string {
	return "html" // tbe
}
func (this atom_entry) save_content() FeedContent {
	ext := mime_to_ext(this.Content.Type)
	dir := feedsprofile().content_dir()
	f, err := ioutil.TempFile(dir, ext+".")
	if err != nil {
		return FeedContent{}
	}
	defer f.Close()
	_, err = f.Write([]byte(this.Content.Body))
	return FeedContent{Uri: this.link(),
		Local: f.Name()}
}

func (this atom_person) to_feedauthor() FeedAuthor {
	return FeedAuthor{Name: this.Name, Email: this.Email}
}
