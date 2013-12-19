package backend

import (
	"github.com/heartszhang/feed"
)

const (
	ReadSourceSubscribeStateSubscribed = 1 << iota
	ReadSourceSubscribeStateUnsubscribed
	ReadSourceSubscribeStatusDisabled
)

type ReadSource struct {
	feed.FeedSource `json:",inline" bson:",inline"`
	EnableProxy     int      `json:"enable_proxy" bson:"enable_proxy"`
	SubscribeState  int      `json:"subscribe_state" bson:"subscribe_state"` //auto refresh enabled
	Categories      []string `json:"categories,omitempty" bson:"categories,omitempty"`
	LastTouch       int64    `json:"last_touch" bson:"last_touch"`
	NextTouch       int64    `json:"next_touch" bson:"next_touch"`
	LastUpdate      int64    `json:"last_update" bson:"last_update"`
}

func new_readsource(s feed.FeedSource) ReadSource {
	return ReadSource{FeedSource: s}
}

func new_readsources(sources []feed.FeedSource) []ReadSource {
	var v []ReadSource
	for _, s := range sources {
		v = append(v, new_readsource(s))
	}
	return v
}
