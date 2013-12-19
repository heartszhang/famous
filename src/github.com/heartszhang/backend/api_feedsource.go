package backend

import (
	"os"

	"github.com/heartszhang/curl"
	"github.com/heartszhang/feed"
	"github.com/heartszhang/google"
)

func feedsource_all() ([]ReadSource, error) {
	dbo := new_feedsource_operator()
	fs, err := dbo.all()
	return fs, err
}
func feed_fetch(uri string) (v ReadSource, res []ReadEntry, err error) {
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
	f, err := os.Open(cache.LocalUtf8)
	if err != nil {
		return
	}
	var fv feed.FeedSource
	var fes []feed.FeedEntry
	fv, fes, err = feed.NewFeedMaker(f, uri).MakeFeed()
	f.Close()
	v = new_readsource(fv)
	res = new_readentries(fes)
	return
}
func feedsource_expired(beforeunx int64) ([]ReadSource, error) {
	return new_feedsource_operator().expired(beforeunx)
}
func feedsource_save(fs ReadSource) error {
	return new_feedsource_operator().save_one(fs)
}
func feedsource_subscribe(uri string, source_type uint) (v ReadSource, err error) {
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

func feedsource_show(uri string) (ReadEntity, error) {
	g := google.NewGoogleFeedApi(refer, backend_context.config.FeedSourceFolder)
	fs, err := g.Load(uri, backend_context.config.Language, 4, false)
	if err != nil {
		return ReadEntity{}, err
	}
	rs := new_readsource(fs.FeedSource)
	if _, err := new_feedsource_operator().find(uri); err != nil {
		rs.SubscribeState = ReadSourceSubscribeStateUnsubscribed
	}
	entries := readentry_filter(new_readentries(fs.Entries))
	return new_readentity(rs, entries), err
}

func feedsource_mark_subscribed(sources []ReadSource) []ReadSource {
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
	var v []ReadSource
	for _, s := range sources {
		if _, ok := u2fs[s.Uri]; ok {
			s.SubscribeState = ReadSourceSubscribeStateSubscribed
		} else {
			s.SubscribeState = ReadSourceSubscribeStateUnsubscribed
		}
		v = append(v, s)
	}
	return v
}
func feedsource_find(q string) ([]ReadEntity, error) {
	svc := google.NewGoogleFeedApi(refer, backend_context.config.FeedSourceFolder)
	v, err := svc.Find(q, backend_context.config.Language)
	if err != nil {
		return nil, err
	}
	var uris []string
	for _, ve := range v {
		uris = append(uris, ve.Uri)
	}
	dbo := new_feedsource_operator()
	subed, err := dbo.findbatch(uris)
	if err != nil {
		return nil, err
	}
	urims := make(map[string]bool)
	for _, ve := range subed {
		urims[ve.Uri] = true
	}
	rv := from_feedentities(v)
	for i := 0; i < len(rv); i++ {
		if _, ok := urims[rv[i].Uri]; !ok {
			rv[i].SubscribeState = ReadSourceSubscribeStateUnsubscribed
		}
	}
	return rv, err
}
