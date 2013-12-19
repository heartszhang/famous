package feed

type FeedEntry struct {
	Parent     string      `json:"src,omitempty" bson:"src,omitempty"` // feed source's uri
	Uri        string      `json:"uri,omitempty" bson:"uri,omitempty"`
	Summary    string      `json:"summary,omitempty" bson:"summary,omitempty"`
	Content    string      `json:"content,omitempty" bson:"content,omitempty"`
	Type       uint        `json:"type" bson:"type"` // feed_type...
	Title      string      `json:"title,omitempty" bson:"title,omitempty"`
	Author     string      `json:"author,omitempty" bson:"author,omitempty"`
	PubDate    int64       `json:"pubdate" bson:"pubdate"` // unix time
	Tags       []string    `json:"tags,omitempty" bson:"tags,omitempty"`
	Categories []string    `json:"categories,omitempty" bson:"category,omitempty"`
	Images     []FeedMedia `json:"images,omitempty" bson:"images,omitempty"`
	Videos     []FeedMedia `json:"videos,omitempty" bson:"videos,omitempty"`
	Audios     []FeedMedia `json:"audios,omitempty" bson:"audios,omitempty"`
	//	Links      []FeedLink    `json:"links,omitempty" bson:"links,omitempty"`
}

type FeedEntity struct {
	FeedSource `json:",inline"`
	Entries    []FeedEntry `json:"entries,omitempty"`
}
