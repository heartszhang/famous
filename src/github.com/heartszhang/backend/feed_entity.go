package backend

import (
	feed "github.com/heartszhang/feedfeed"
)

type FeedEntity struct {
	feed.FeedSource
	Entries []feed.FeedEntry `json:"entries,omitempty"`
}
