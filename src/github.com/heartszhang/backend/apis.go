package backend

import (
	"fmt"
	"github.com/heartszhang/curl"
	feed "github.com/heartszhang/feedfeed"
)

// since_unixtime , 0: from now
// categories, categories mask, every bit represent a category
// count: entries per page
// page: page no, start at 0
func feeds_entries_since(since_unixtime int64, categories uint64, count uint, page uint) ([]feed.FeedEntry, error) {
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
func feed_entry_fulltext(url string, entry_id string) (feed.FeedLink, error) {
	return feed.FeedLink{}, nil
}

// /feed/entry/image.json/{url}/{entry_id}
func feed_entry_image(url string, entry_id string) (feed.FeedImage, error) {
	return feed.FeedImage{}, nil
}

// /feed/entry/media.json/{url}/{entry_id}/{media_type:[0-9]+}

func feed_entry_media(url string, entry_id string, media_type uint) (feed.FeedMedia, error) {
	return feed.FeedMedia{}, nil
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

func Subscribe(url string) (feed.FeedSource, error) {
	return feed_source_subscribe(url, feed.Feed_type_feed, feed.Feed_category_root)
}

func feed_source_subscribe(url string, source_type uint, category uint64) (v feed.FeedSource, err error) {
	curler := curl.NewCurl(BackendConfig().FeedSourceDir)
	cache, err := curler.GetUtf8(url, curl.CurlProxyPolicyUseProxy)
	fmt.Println(cache, err)

	if cache.LocalUtf8 != "" {
		fstype := feed.DetectFeedSourceType(cache.LocalUtf8)
		switch fstype {
		case feed.Feed_type_rss:
			v, err = feed.CreateFeedSourceRss2(cache.LocalUtf8)
		case feed.Feed_type_atom:
			v, err = feed.CreateFeedSourceAtom(cache.LocalUtf8)
		}
	}
	return
}

func feed_source_unsubscribe(url string, source_type uint, category uint64) error {
	return nil
}

func meta_categories() ([]feed.FeedCategory, error) {
	return nil, nil
}

func meta() (FeedsBackendConfig, error) {
	return BackendConfig(), nil
}

func meta_cleanup() error {
	return nil
}

func source_type_map(sourcetype string) uint {
	v, ok := feed.FeedSourceTypes[sourcetype]
	if !ok {
		v = feed.Feed_type_unknown
	}
	return v
}
