package backend

import (
	feed "github.com/heartszhang/feedfeed"
)

type feedentry_operator interface {
	save([]feed.FeedEntry) ([]feed.FeedEntry, error)
	save_one(feed.FeedEntry) (interface{}, error)
	topn(skip, limit int) ([]feed.FeedEntry, error)
	topn_by_category(skip, limit int, category string) ([]feed.FeedEntry, error)
	topn_by_feedsource(skip, limit int, feed string) ([]feed.FeedEntry, error)
	mark(link string, newmark uint) error
	umark(uri string, markbit uint) error
	setcontent(link string, filepath string, words int, imgs []feed.FeedMedia) error
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
	upsert(f *feed.FeedSource) error
	find(uri string) (*feed.FeedSource, error)
	expired() ([]feed.FeedSource, error)
	all() ([]feed.FeedSource, error)
	touch(uri string, ttl int) error
	drop(uri string) error
	disable(uri string, dis bool) error
	update(f *feed.FeedSource) error
}

type feedcontent_operator interface {
	touch(hash int64) (uint, error)
}
