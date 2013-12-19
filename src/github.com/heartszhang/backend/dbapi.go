package backend

import (
	"github.com/heartszhang/feed"
)

type feedentry_operator interface {
	save([]ReadEntry) ([]ReadEntry, error)
	save_one(ReadEntry) (interface{}, error)
	topn(skip, limit int) ([]ReadEntry, error)
	topn_by_category(skip, limit int, category string) ([]ReadEntry, error)
	topn_by_feedsource(skip, limit int, feed string) ([]ReadEntry, error)
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
	save(feeds []ReadSource) ([]ReadSource, error)
	upsert(f ReadSource) error
	update(f ReadSource) error
	find(uri string) (ReadSource, error)
	all() ([]ReadSource, error)
	touch(uri string, last, next, period int64) error
	drop(uri string) error
	set_subscribe_state(uri string, s int) error
	save_one(f ReadSource) error
	findbatch(uris []string) ([]ReadSource, error)
	expired(beforeunxtime int64) ([]ReadSource, error)
}

type feedcontent_operator interface {
	touch(hash int64) (uint, error)
}
