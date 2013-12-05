package feed

import (
	"encoding/xml"
	"os"
)

//<head><title/></head> is omitted
type opml struct {
	Body struct {
		Outline []opml_outline `xml:"outline" bson:"outline,omitempty" json:"outline,omitempty"`
	} `xml:"body,omitempty" json:"-" bson:"-"`
}

type opml_outline struct {
	Text        string         `xml:"text,attr" bson:"-" json:"-"` // same as title
	Title       string         `xml:"title,attr,omitempty" bson:"title" json:"title"`
	Type        string         `xml:"type,attr,omitempty" bson:"type" json:"type"` // type='rss' // or link
	Description string         `xml:"description,attr,omitempty"`
	Version     string         `xml:"version,attr,omitempty"`
	Docs        string         `xml:"xmlUrl,attr" bson:"link" json:"link"`
	Link        string         `xml:"htmlUrl,attr," bson:"htmlurl" json:"htmlurl"`
	Category    string         `xml:"category,attr,omitempty"` // seped by ,
	Children    []opml_outline `xml:"outline,omitempty" bson:"children,omitempty" json:"omitempty"`
}

func CreateFeedsCategoryOpml(opmlfile string) ([]FeedSource, error) {
	return feeds_category_create_opml(opmlfile)
}

func feeds_category_create_opml(filepath string) ([]FeedSource, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var o opml
	d := xml.NewDecoder(f)
	d.CharsetReader = charset_reader_passthrough

	err = d.Decode(&o)
	return o.to_feedscategory(), err
}

func (this opml_outline) name() string {
	if this.Title == "" {
		return this.Text
	}
	return this.Title
}

func (this opml_outline) export_feedsource(v []FeedSource) []FeedSource {
	if this.Docs != "" {
		x := FeedSource{
			FeedSourceMeta: FeedSourceMeta{
				Name:        this.name(),
				Uri:         this.Docs,
				Period:      _2hours,
				Type:        Feed_type_feed, // may be atom?
				WebSite:     this.Link,
				Description: this.Description,
			},
			Disabled:    false,
			EnableProxy: 0,
		}
		v = append(v, x)
	}
	for _, child := range this.Children {
		v = child.export_feedsource(v)
	}
	return v
}

func (this opml) to_feedscategory() []FeedSource {
	v := []FeedSource{}
	if len(this.Body.Outline) == 0 {
		return v
	}
	for _, outline := range this.Body.Outline {
		v = outline.export_feedsource(v)
	}
	return v
}
