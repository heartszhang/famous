package backend

import (
	"fmt"
	"github.com/heartszhang/curl"
	"github.com/heartszhang/feed"
	"log"
)

// category
// count: entries per page, must large than 0
// page: page no, start at 0
func feeds_entries_since(category string, count uint, page uint) ([]feed.FeedEntry, error) {
	return []feed.FeedEntry{}, backend_error{"not impl", -1}
}

// source : feed(atom/rss) url
// count: must large than 0
// page: 0 based page index
// if page is 0, entries may be fetched online
func feedentry_unread(source string, count int, page int) ([]feed.FeedEntry, error, int) {
	if count <= 0 {
		panic("invalid arg count")
	}
	var sc int
	if page == 0 {
		log.Println("curl-get...")
		c := curl.NewCurlerDetail(backend_config().FeedEntryFolder, 0, 0, nil, backend_context.ruler)
		cache, err := c.GetUtf8(source)
		log.Println("curl-get", cache.LocalUtf8)
		if err != nil || cache.LocalUtf8 == "" {
			return nil, err, cache.StatusCode
		}

		ext := curl.MimeToExt(cache.Mime)
		if ext != "xml" && ext != "atom+xml" && ext != "rss+xml" {
			return nil, fmt.Errorf("unsupported mime: %v, %d", cache.Mime, cache.StatusCode), 0
		}
		fs, v, err := feed.NewFeedMaker(cache.LocalUtf8, source).MakeFeed()
		if err == nil {
			new_feedsource_operator().update(fs)
			log.Println("feed-update", fs.Name)
		}
		v = feedentry_filter(v)
		log.Println("feedentries-filter", len(v))
		sc = cache.StatusCode
	}
	v, err := new_feedentry_operator().topn_by_feedsource(count*page, count, source)
	log.Println("unread-return(uri, page, count)", source, page, count, len(v), err)
	return v, err, sc
}

func feedentry_mark(uri string, flags uint) (uint, error) {
	dbo := new_feedentry_operator()
	err := dbo.mark(uri, flags)
	return flags, err
}

func feedentry_category_mark(cate string, flags uint) error {
	dbo := new_feedentry_operator()
	err := dbo.mark_category(cate, flags)
	return err
}

func feedentry_source_mark(src string, flags uint) error {
	err := new_feedentry_operator().mark_source(src, flags)
	return err
}

// /feed/entry/umark.json/{id}/{flags}
func feedentry_umark(uri string, flags uint) (uint, error) {
	dbo := new_feedentry_operator()
	err := dbo.mark(uri, flags)
	return flags, err
}

func feedentry_category_umark(category string, flags uint) error {
	err := new_feedentry_operator().umark_category(category, flags)
	return err
}

/*
func feedentry_fulltext(uri string, entry_uri string) (v feed.FeedContent, err error) {
	c := curl.NewCurl(backend_context.config.DocumentFolder)
	cache, err := c.GetUtf8(uri)
	v.Uri = uri
	v.Local = cache.LocalUtf8
	//	v.Length = cache.LengthUtf8
	if err != nil {
		return v, err
	}
	if curl.MimeToExt(cache.Mime) != "html" {
		return v, fmt.Errorf("unsupported mime %v", cache.Mime)
	}
	doc, err := html_create_from_file(cache.LocalUtf8)
	if err != nil {
		return v, err
	}
	article, sum, err := cleaner.NewExtractor(backend_context.config.CleanFolder).MakeHtmlReadable(doc, uri)
	v.Images = make([]feed.FeedMedia, len(sum.Images))
	for idx, img := range sum.Images {
		v.Images[idx].Uri = img.Uri
		v.Images[idx].Width = int(img.Width)
		v.Images[idx].Height = int(img.Height)
		v.Images[idx].Description = img.Alt
	}
	v.Words = uint(sum.WordCount)
	v.Links = uint(sum.LinkCount)
	v.Local, err = html_write_file(article, backend_context.config.DocumentFolder)
	return v, err
}
*/
// /feed/entry/media.json/{url}/{entry_id}/{media_type:[0-9]+}

func feedentry_media(url string, entry_id string, media_type uint) (feed.FeedMedia, error) {
	return feed.FeedMedia{}, nil
}

// /feed/entry/drop.json/{id}

// id is mongo's _id
func feedentry_drop(id string) error {
	return nil
}
