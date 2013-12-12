package backend

import (
	"github.com/heartszhang/feed"
)

type feedentry_operator interface {
	save([]feed.FeedEntry) ([]feed.FeedEntry, error)
	save_one(feed.FeedEntry) (interface{}, error)
	topn(skip, limit int) ([]feed.FeedEntry, error)
	topn_by_category(skip, limit int, category string) ([]feed.FeedEntry, error)
	topn_by_feedsource(skip, limit int, feed string) ([]feed.FeedEntry, error)
	mark(link string, newmark uint) error
	umark(uri string, markbit uint) error
	umark_category(category string, markbit uint) error
	mark_category(category string, newmark uint) error
	umark_source(category string, markbit uint) error
	mark_source(category string, newmark uint) error
	//	setcontent(link string, filepath string, words int, imgs []feed.FeedMedia) error
	unread_count(uri string) (int, error)
	unread_count_sources() ([]feedentry_unreadcount, error)
	unread_count_category(category string) (int, error)
	unread_count_categories() ([]feedentry_unreadcount, error)
}

type feedentrytouch_operator interface {
	touch(hashes []string) ([]string, error)
}
type feedentry_unreadcount struct {
	Source string `json:"source" bson:"_id"`
	Count  int    `json:"count" bson:"value"`
}

type feedentrytag_operator interface {
}

type imagecache_operator interface {
	find(uri string) (feed.FeedImage, error)
	save(uri string, v feed.FeedImage) error
}
type feedcategory_operator interface {
	save(cate string) (interface{}, error)
	all() ([]string, error)
	drop(category string) error
}
type feedtag_operator feedcategory_operator

type feedsource_operator interface {
	save(feeds []feed.FeedSource) ([]feed.FeedSource, error)
	upsert(f feed.FeedSource) error
	update(f feed.FeedSource) error
	find(uri string) (feed.FeedSource, error)
	all() ([]feed.FeedSource, error)
	touch(uri string, last, next, period int64) error
	drop(uri string) error
	set_subscribe_state(uri string, s int) error
	save_one(f feed.FeedSource) error
	findbatch(uris []string) ([]feed.FeedSource, error)
	expired(beforeunxtime int64) ([]feed.FeedSource, error)
}

type feedcontent_operator interface {
	touch(hash int64) (uint, error)
}
