package backend

import (
	"fmt"
	"github.com/heartszhang/curl"
	"github.com/heartszhang/feed"
	"github.com/heartszhang/google"
	"log"
	"net/url"
)

func feedsource_all() ([]feed.FeedSource, error) {
	dbo := new_feedsource_operator()
	fs, err := dbo.all()
	for i := 0; i < len(fs); i++ {
		f := &fs[i]
		f.Logo = resolve_logo(*f)
	}
	return fs, err
}
func feed_fetch(uri string) (v feed.FeedSource, fes []feed.FeedEntry, err error) {
	cache, err := curl.NewCurl(backend_config().FeedSourceFolder).GetUtf8(uri)
	if err != nil {
		return
	}
	ext := curl.MimeToExt(cache.Mime)
	if ext != "xml" && ext != "atom+xml" && ext != "rss+xml" {
		return v, nil, fmt.Errorf("unsupported mime: %v", cache.Mime)
	} else if cache.LocalUtf8 == "" {
		return v, nil, fmt.Errorf("unrecognized encoding %v", cache.Local)
	}
	v, fes, err = feed.NewFeedMaker(cache.LocalUtf8, uri).MakeFeed()
	if v.Uri == "" {
		v.Uri = uri
	}
	return
}
func feedsource_expired(beforeunx int64) ([]feed.FeedSource, error) {
	return new_feedsource_operator().expired(beforeunx)
}
func feedsource_save(fs feed.FeedSource) error {
	return new_feedsource_operator().save_one(fs)
}
func feedsource_subscribe(uri string, source_type uint) (v feed.FeedSource, err error) {
	fso := new_feedsource_operator()
	fs, err := fso.find(uri)
	if err == nil {
		return *fs, nil
	}
	v, _, err = feed_fetch(uri)
	v.Type = source_type

	if err == nil && v.Uri != "" {
		err = fso.upsert(v)
	}
	return v, err
}

func feedsource_unsubscribe(url string) error {
	fso := new_feedsource_operator()
	err := fso.drop(url)
	return err
}

const (
	refer = "http://iweizhi2.duapp.com"
)

func feedsource_show(uri string) (FeedEntity, error) {
	fs, entries, err := google.NewGoogleFeedApi(refer, backend_context.config.FeedSourceFolder).Load(uri, backend_context.config.Language, 4, false)
	if err != nil {
		return FeedEntity{}, err
	}
	return FeedEntity{fs, entries}, err
}

func feedsource_find(q string) ([]google.FeedSourceFindEntry, error) {
	svc := google.NewGoogleFeedApi(refer, backend_context.config.FeedSourceFolder)
	v, err := svc.Find(q, backend_context.config.Language)
	if err != nil {
		return v, err
	}
	uris := make([]string, 0)
	for _, ve := range v {
		uris = append(uris, ve.Url)
	}
	dbo := new_feedsource_operator()
	subed, err := dbo.findbatch(uris)
	if err != nil {
		return v, err
	}
	for _, fs := range subed {
		for i := 0; i < len(v); i++ {
			if v[i].Url == fs.Uri {
				v[i].Subscribed = true
			}
		}
	}
	return v, err
}

func resolve_logo(f feed.FeedSource) string {
	baseu, _ := url.Parse(f.WebSite)
	u, err := url.Parse(f.Logo)
	if err != nil {
		u, err = baseu.Parse("/favicon.ico")
	}
	if !u.IsAbs() {
		u, _ = baseu.Parse(f.Logo)
	}
	v := u.String()
	log.Println("web-site:", u.String())
	return v
}
