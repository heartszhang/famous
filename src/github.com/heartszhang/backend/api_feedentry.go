package backend

import (
	"fmt"
	"github.com/heartszhang/cleaner"
	"github.com/heartszhang/curl"
	feed "github.com/heartszhang/feedfeed"
)

// since_unixtime , 0: from now
// categories, categories mask, every bit represent a category
// count: entries per page
// page: page no, start at 0
func feeds_entries_since(since_unixtime int64, category string, count uint, page uint) ([]feed.FeedEntry, error) {
	return []feed.FeedEntry{}, nil
}

func feedentry_unread(source string, count uint, page uint) ([]feed.FeedEntry, error, int) {
	c := curl.NewCurl(backend_config().FeedEntryFolder)
	cache, err := c.GetUtf8(source)

	if err != nil || cache.LocalUtf8 == "" {
		return nil, err, cache.StatusCode
	}
	//atom+xml;xml;html
	ext := curl.MimeToExt(cache.Mime)
	if ext != "xml" && ext != "atom+xml" && ext != "rss+xml" {
		return nil, fmt.Errorf("unsupported mime: %v, %d", cache.Mime, cache.StatusCode), 0
	}
	v, err := feed.MakeFeedEntries(cache.LocalUtf8)
	v = feedentry_filter(v)
	return v, err, cache.StatusCode
}

func feedentry_mark(uri string, flags uint) (uint, error) {
	dbo := new_feedentry_operator()
	err := dbo.mark(uri, flags)
	return flags, err
}

// /feed/entry/umark.json/{id}/{flags}
func feedentry_umark(uri string, flags uint) (uint, error) {
	dbo := new_feedentry_operator()
	err := dbo.mark(uri, flags)
	return flags, err
}

// /feed/entry/full_text.json/{url}/{entry_id}
func feedentry_fulltext(uri string, entry_uri string) (v feed.FeedContent, err error) {
	c := curl.NewCurl(config.DocumentFolder)
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
	article, sum, err := cleaner.NewExtractor(config.CleanFolder).MakeHtmlReadable(doc, uri)
	v.Images = make([]feed.FeedMedia, len(sum.Images))
	for idx, img := range sum.Images {
		v.Images[idx].Uri = img.Uri
		v.Images[idx].Width = int(img.Width)
		v.Images[idx].Height = int(img.Height)
		v.Images[idx].Description = img.Alt
	}
	v.Words = uint(sum.WordCount)
	v.Links = uint(sum.LinkCount)
	v.Local, err = html_write_file(article, config.DocumentFolder)
	return v, err
}

// /feed/entry/media.json/{url}/{entry_id}/{media_type:[0-9]+}

func feedentry_media(url string, entry_id string, media_type uint) (feed.FeedMedia, error) {
	return feed.FeedMedia{}, nil
}

// /feed/entry/drop.json/{id}

// id is mongo's _id
func feedentry_drop(id string) error {
	return nil
}
