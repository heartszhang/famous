package backend

import (
	feed "github.com/heartszhang/feedfeed"
)

// /tick.json

func tick() (feed_tick, error) {
	s := backend_tick()
	return s, nil
}

func meta() (FeedsBackendConfig, error) {
	return backend_config(), nil
}

type feed_update struct {
	feed.FeedSource
	Entries []feed.FeedEntry `json:"entries,omitempty"`
}

func update_popup() (*feed_update, error) {
	fs, fes, err := feedentries_updated()
	if err == nil {
		return &feed_update{FeedSource: *fs, Entries: fes}, nil
	}
	return nil, err
}

func meta_cleanup() error {
	// clean temp files
	// entries
	// thumbnails
	// images
	return nil
}

func source_type_map(sourcetype string) uint {
	v, ok := feed.FeedSourceTypes[sourcetype]
	if !ok {
		v = feed.Feed_type_unknown
	}
	return v
}

func feedtag_all() ([]string, error) {
	fto := new_feedtag_operator()
	return fto.all()
}

const (
	refer = "http://iweizhi2.duapp.com"
)
