package feed

import "io"

type opml struct { //<head><title/></head> is omitted
	Body struct {
		Outline []opml_outline `xml:"outline" bson:"outline,omitempty" json:"outline,omitempty"`
	} `xml:"body,omitempty" json:"-" bson:"-"`
}

type opml_outline struct {
	Text        string         `xml:"text,attr"` // same as title
	Title       string         `xml:"title,attr,omitempty"`
	Type        string         `xml:"type,attr,omitempty"` // type='rss' // or link
	Description string         `xml:"description,attr,omitempty"`
	Version     string         `xml:"version,attr,omitempty"`
	XmlUrl      string         `xml:"xmlUrl,attr"`
	HtmlUrl     string         `xml:"htmlUrl,attr,"`
	Category    string         `xml:"category,attr,omitempty"` // seped by ,
	Children    []opml_outline `xml:"outline,omitempty"`
}

func OpmlExportFeedSource(data io.Reader) ([]FeedSource, error) {
	return feeds_category_create_opml(data)
}

func feeds_category_create_opml(data io.Reader) ([]FeedSource, error) {
	var o opml
	err := new_xml_decoder(data).Decode(&o)
	return o.export_feedsources(), err
}

func (this opml_outline) name() string {
	if this.Title == "" {
		return this.Text
	}
	return this.Title
}

func (this opml_outline) export_feedsources(v []FeedSource, categories []string) []FeedSource {
	if this.XmlUrl != "" { // this is a feed-source, not a category
		x := FeedSource{
			Name:        this.name(),
			Uri:         this.XmlUrl,
			Period:      _2hours,
			Type:        FeedSourceType(this.Type), // may be atom?
			WebSite:     this.HtmlUrl,
			Description: this.Description,
			Tags:        categories,
		}
		v = append(v, x)
	} else {
		categories = append(categories, this.name())
	}
	for _, child := range this.Children {
		v = child.export_feedsources(v, categories)
	}
	return v
}

func (this opml) export_feedsources() []FeedSource {
	var v []FeedSource
	for _, outline := range this.Body.Outline {
		v = outline.export_feedsources(v, nil)
	}
	return v
}
