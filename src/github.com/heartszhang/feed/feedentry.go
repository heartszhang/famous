package feed

import (
	"github.com/heartszhang/unixtime"
)

type FeedEntryMeta struct {
	Parent     string            `json:"src,omitempty" bson:"src,omitempty"` // source's uri
	Uri        string            `json:"uri,omitempty" bson:"uri,omitempty"`
	Summary    string            `json:"summary,omitempty" bson:"summary,omitempty"`
	Content    string            `json:"content,omitempty" bson:"content,omitempty"`
	Type       uint              `json:"type" bson:"type"` // feed_type...
	Title      FeedTitle         `json:"title,omitempty" bson:"title,omitempty"`
	Author     *FeedAuthor       `json:"author,omitempty" bson:"author,omitempty"`
	PubDate    unixtime.UnixTime `json:"pubdate" bson:"pubdate"` // unix time
	Tags       []string          `json:"tags,omitempty" bson:"tags,omitempty"`
	Categories []string          `json:"categories,omitempty" bson:"category,omitempty"`
	Links      []FeedLink        `json:"links,omitempty" bson:"links,omitempty"`
}
type FeedEntry struct {
	FeedEntryMeta `json:",inline" bson:",inline"`
	Flags         uint           `json:"flags" bson:"flags"`
	Images        []FeedMedia    `json:"images,omitempty" bson:"images,omitempty"`
	Videos        []FeedMedia    `json:"videos,omitempty" bson:"videos,omitempty"`
	Audios        []FeedMedia    `json:"audios,omitempty" bson:"audios,omitempty"`
	SummaryStatus FeedTextStatus `json:"summary_status" bson:"summary_status"`
	ContentStatus FeedTextStatus `json:"content_status" bson:"content_status"`
	Status        uint64         `json:"status" bson:"status"`
	Readed        bool           `json:"readed" bson:"readed"`
}
