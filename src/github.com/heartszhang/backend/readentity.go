package backend

import (
	"github.com/heartszhang/feed"
)

type ReadEntity struct {
	ReadSource `json:",inline"`
	Entries    []ReadEntry `json:"entries,omitempty"`
}

func new_readentity(s ReadSource, entries []ReadEntry) ReadEntity {
	return ReadEntity{s, entries}
}

func from_feedentity(e feed.FeedEntity) ReadEntity {
	return ReadEntity{new_readsource(e.FeedSource), new_readentries(e.Entries)}
}

func from_feedentities(entities []feed.FeedEntity) []ReadEntity {
	var v []ReadEntity
	for _, e := range entities {
		v = append(v, from_feedentity(e))
	}
	return v
}
