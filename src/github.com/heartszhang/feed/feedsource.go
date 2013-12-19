package feed

type FeedSource struct {
	Name        string   `json:"name,omitempty" bson:"name,omitempty"`
	Uri         string   `json:"uri,omitempty" bson:"uri,omitempty"`         // rss/atom url
	Period      int64    `json:"period" bson:"period"`                       // minutes
	Update      int64    `json:"update" bson:"update"`                       // unixtime
	Type        uint     `json:"type" bson:"type"`                           // feed_type...
	WebSite     string   `json:"website,omitempty" bson:"website,omitempty"` // home
	Description string   `json:"description,omitempty" bson:"description,omitempty"`
	Logo        string   `json:"logo,omitempty" bson:"logo,omitempty"`
	Hub         string   `json:"hub,omitempty" bson:"hub,omitempty"`
	Tags        []string `json:"tags,omitempty" bson:"tags,omitempty"`
}
