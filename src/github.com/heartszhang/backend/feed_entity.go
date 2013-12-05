package backend

import (
	"github.com/heartszhang/feed"
)

type FeedEntity struct {
	feed.FeedSource
	Entries []feed.FeedEntry `json:"entries,omitempty"`
}
