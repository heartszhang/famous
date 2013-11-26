package backend

import (
//	"github.com/heartszhang/cleaner"
//	"github.com/heartszhang/curl"
	feed "github.com/heartszhang/feedfeed"
//	"github.com/heartszhang/google"
)


// /tick.json

func tick() (FeedsStatus, error) {
	s := backend_status()
	return s, nil
}

func meta() (FeedsBackendConfig, error) {
	return backend_config(), nil
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
