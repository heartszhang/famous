package backend

import (
	"github.com/heartszhang/curl"
	"github.com/heartszhang/feed"
	"github.com/heartszhang/google"
)

func feedsource_all() ([]feed.FeedSource, error) {
	dbo := new_feedsource_operator()
	fs, err := dbo.all()
	return fs, err
}
func feed_fetch(uri string) (v feed.FeedSource, fes []feed.FeedEntry, err error) {
	cache, err := curl.NewCurl(backend_config().FeedSourceFolder).GetUtf8(uri)
	if err != nil {
		return
	}
	ext := curl.MimeToExt(cache.Mime)
	if ext != "xml" && ext != "atom+xml" && ext != "rss+xml" {
		return v, nil, new_backenderror(-1, "unsupported mime: "+cache.Mime)
	} else if cache.LocalUtf8 == "" {
		return v, nil, new_backenderror(-1, "unrecognized encoding: "+cache.Local)
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
		return fs, nil
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

func feedsource_show(uri string) (feed.FeedEntity, error) {
	g := google.NewGoogleFeedApi(refer, backend_context.config.FeedSourceFolder)
	fs, err := g.Load(uri, backend_context.config.Language, 4, false)
	if err != nil {
		return feed.FeedEntity{}, err
	}
	if _, err := new_feedsource_operator().find(uri); err != nil {
		fs.SubscribeState = feed.FeedSourceSubscribeStateUnsubscribed
	}
	fs.Entries = feedentry_filter(fs.Entries)
	return fs, err
}

func feedsource_mark_subscribed(sources []feed.FeedSource) []feed.FeedSource {
	var uris []string
	u2fs := make(map[string]bool)
	for i := 0; i < len(sources); i++ {
		uri := sources[i].Uri
		uris = append(uris, uri)
	}
	subed, _ := new_feedsource_operator().findbatch(uris)
	for _, s := range subed {
		u2fs[s.Uri] = true
	}
	var v []feed.FeedSource
	for _, s := range sources {
		if _, ok := u2fs[s.Uri]; ok {
			s.SubscribeState = feed.FeedSourceSubscribeStateSubscribed
		}
		v = append(v, s)
	}
	return v
}
func feedsource_find(q string) ([]feed.FeedEntity, error) {
	svc := google.NewGoogleFeedApi(refer, backend_context.config.FeedSourceFolder)
	v, err := svc.Find(q, backend_context.config.Language)
	if err != nil {
		return v, err
	}
	var uris []string
	for _, ve := range v {
		uris = append(uris, ve.Uri)
	}
	dbo := new_feedsource_operator()
	subed, err := dbo.findbatch(uris)
	if err != nil {
		return v, err
	}
	urims := make(map[string]bool)
	for _, ve := range subed {
		urims[ve.Uri] = true
	}
	for i := 0; i < len(v); i++ {
		if _, ok := urims[v[i].Uri]; !ok {
			v[i].SubscribeState = feed.FeedSourceSubscribeStateUnsubscribed
		}
	}
	return v, err
}

/*
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
*/
