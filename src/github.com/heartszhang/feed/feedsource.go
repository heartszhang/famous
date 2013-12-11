package feed

import (
	"github.com/heartszhang/unixtime"
)

type FeedSourceMeta struct {
	Name        string   `json:"name,omitempty" bson:"name,omitempty"`
	Uri         string   `json:"uri,omitempty" bson:"uri,omitempty"`         // rss/atom url
	Period      int64    `json:"period" bson:"period"`                       // minutes
	Type        uint     `json:"type" bson:"type"`                           // feed_type...
	WebSite     string   `json:"website,omitempty" bson:"website,omitempty"` // home
	Description string   `json:"description,omitempty" bson:"description,omitempty"`
	Logo        string   `json:"logo,omitempty" bson:"logo,omitempty"`
	Hub         string   `json:"hub,omitempty" bson:"hub,omitempty"`
	Tags        []string `json:"tags,omitempty" bson:"tags,omitempty"`
}
type FeedSource struct {
	FeedSourceMeta `json:",inline" bson:",inline"`
	Local          string            `json:"local,omitempty" bson:"local,omitempty"`
	EnableProxy    int               `json:"enable_proxy" bson:"enable_proxy"`
	SubscribeState int               `json:"subscribe_state" bson:"subscribe_state"` //auto refresh enabled
	Categories     []string          `json:"categories,omitempty" bson:"categories,omitempty"`
	Update         unixtime.UnixTime `json:"update" bson:"update"`
	LastTouch      unixtime.UnixTime `json:"last_touch" bson:"last_touch"`
	NextTouch      unixtime.UnixTime `json:"next_touch" bson:"next_touch"`
	LastUpdate     unixtime.UnixTime `json:"last_update" bson:"last_update"`
}

const (
	FeedSourceSubscribeStateSubscribed = 1 << iota
	FeedSourceSubscribeStateUnsubscribed
	FeedSourceSubscribeStatusDisabled
)
