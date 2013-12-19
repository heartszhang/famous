package backend

import (
	"net/url"
	"os"

	"github.com/heartszhang/cleaner"
	"github.com/heartszhang/curl"
	"github.com/heartszhang/feed"
	"github.com/qiniu/log"
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
func feedentry_unread(source string, count int, page int) ([]ReadEntry, error, int) {
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
			return nil, new_backenderror(cache.StatusCode, "unsupported mime: "+cache.Mime), 0
		}
		f, err := os.Open(cache.LocalUtf8)
		if err != nil {
			return nil, err, cache.StatusCode
		}

		fs, v, err := feed.NewFeedMaker(f, source).MakeFeed()
		f.Close()
		rs := new_readsource(fs)
		if err == nil {
			new_feedsource_operator().update(rs)
			log.Println("feed-update", fs.Name)
		}
		rv := readentry_filter(new_readentries(v))
		log.Println("feedentries-filter", len(rv))
		sc = cache.StatusCode
	}
	rv, err := new_feedentry_operator().topn_by_feedsource(count*page, count, source)
	log.Println("unread-return(uri, page, count)", source, page, count, len(rv), err)
	return rv, err, sc
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

func feedentry_fulldoc(uri string) (v feed.FeedContent, err error) {
	c := curl.NewCurlerDetail(backend_context.config.DocumentFolder, 0, 0, nil, backend_context.ruler)
	cache, err := c.GetUtf8(uri)
	if err != nil {
		return v, err
	}
	v.Uri = uri
	v.Local = cache.LocalUtf8
	if curl.MimeToExt(cache.Mime) != "html" {
		return v, new_backenderror(-1, "unsupported mime: "+cache.Mime)
	}
	doc, err := html_create_from_file(cache.LocalUtf8)
	if err != nil {
		return v, err
	}
	article, sum, err := cleaner.NewExtractor(backend_context.config.CleanFolder).MakeHtmlReadable(doc, uri)
	if err != nil {
		return v, err
	}
	v.Local, err = html_write_file(article, backend_context.config.DocumentFolder)
	redirector := func(turi string) string {
		return redirect_thumbnail(url_resolve(uri, turi))
	}
	imgurl_maker := func(reluri string) string {
		u := url_resolve(uri, reluri)
		return imageurl_from_video(u)
	}
	v.Images = append_unique(v.Images, feedmedias_from_docsummary(sum.Images, redirector)...)
	v.Images = append_unique(v.Images, feedmedias_from_docsummary(sum.Medias, imgurl_maker)...)
	v.Words = uint(sum.WordCount)
	v.Links = uint(sum.LinkCount)
	v.FlowDoc = new_flowdoc_maker().make(article, v.Images)
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

func url_resolve(root, ref string) string {
	u, e := url.ParseRequestURI(root)
	u, e = u.Parse(ref)
	if e == nil {
		return u.String()
	}
	return ref
}
