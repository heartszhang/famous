package backend

import (
	"time"
)

// since_unixtime , 0: from now
// categories, categories mask, every bit represent a category
// count: entries per page
// page: page no, start at 0
func feeds_entries_since(since_unixtime int64, categories uint64, count uint, page uint) ([]FeedEntry, error) {
	return nil, nil
}

func feed_entry_mark(id string, flags int) (int, error) {
	return 0, nil
}

// /feed/entry/umark.json/{id}/{flags}
func feed_entry_umark(id string, flags int) (int, error) {
	return 0, nil
}

// /feed/entry/full_text.json/{url}/{entry_id}
func feed_entry_fulltext(url string, entry_id string) (FeedLink, error) {
	return FeedLink{}, nil
}

// /feed/entry/image.json/{url}/{entry_id}
func feed_entry_image(url string, entry_id string) (FeedImage, error) {
	return FeedImage{}, nil
}

// /feed/entry/media.json/{url}/{entry_id}/{media_type:[0-9]+}

func feed_entry_media(url string, entry_id string, media_type uint) (FeedMedia, error) {
	return FeedMedia{}, nil
}

// /feed/entry/drop.json/{id}

// id is mongo's _id
func feed_entry_drop(id string) error {
	return nil
}

// select a idle category_id, assigned to category
func feed_category_create(name string) (uint64, error) {
	return 0, nil
}

// id : isn't root or all, drop the category whoes name is name
// id : other, drop categories
// name : can be empty. if id is root or all, name cann't be empty
func feed_category_drop(name string) error {
	return nil
}

// /tick.json

func tick() (FeedsStatus, error) {
	s := BackendStatus()
	return s, nil
}

func feed_source_subscribe(url string, source_type uint, category uint64) (FeedSource, error) {
	return FeedSource{}, nil
}

var (
	_source_types = map[string]uint{
		"":        feed_type_unknown,
		"rss":     feed_type_rss,
		"atom":    feed_type_atom,
		"blog":    feed_type_blog,
		"tweet":   feed_type_tweet,
		"weibo":   feed_type_sina_weibo,
		"qqweibo": feed_type_qq_weibo}
)

func feed_source_unsubscribe(url string, source_type uint, category uint64) error {
	return nil
}

func meta_categories() ([]FeedCategory, error) {
	return nil, nil
}

func meta() (FeedsProfile, error) {
	return feedsprofile(), nil
}

func meta_cleanup() error {
	return nil
}

func source_type_map(sourcetype string) uint {
	v, ok := _source_types[sourcetype]
	if !ok {
		v = feed_type_unknown
	}
	return v
}

func unixtime_now() int64 {
	return time.Now().UnixNano()
}

func unixtime_nano(nano int64) time.Time {
	sec := time.Duration(nano) / time.Second
	n := time.Duration(nano) % time.Second
	return time.Unix(int64(sec), int64(n))
}
