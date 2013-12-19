package backend

import (
	"github.com/heartszhang/feed"
)

type ReadEntry struct {
	feed.FeedEntry `json:",inline" bson:",inline"`
	SummaryStatus  feed.FeedTextStatus `json:"summary_status" bson:"summary_status"`
	ContentStatus  feed.FeedTextStatus `json:"content_status" bson:"content_status"`
	Flags          uint                `json:"flags" bson:"flags"`   // unread, star, saved..
	Status         uint64              `json:"status" bson:"status"` // SummaryStatus | ContentStatus
	Readed         bool                `json:"readed" bson:"readed"`
}

func new_readentry(e feed.FeedEntry) ReadEntry {
	return ReadEntry{FeedEntry: e}
}

func new_readentries(entries []feed.FeedEntry) []ReadEntry {
	var v []ReadEntry
	for _, e := range entries {
		v = append(v, new_readentry(e))
	}
	return v
}
