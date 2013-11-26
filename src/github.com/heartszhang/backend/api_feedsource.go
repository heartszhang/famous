package backend

import (
	"fmt"
	feed "github.com/heartszhang/feedfeed"
	"github.com/heartszhang/curl"
	"github.com/heartszhang/google"
)
func feedsource_all() ([]feed.FeedSource, error) {
	dbo := new_feedsource_operator()
	fs, err := dbo.all()
	return fs, err
}

func feedsource_subscribe(uri string, source_type uint) (v feed.FeedSource, err error) {
	fso := new_feedsource_operator()
	fs, err := fso.find(uri)
	if err == nil {
		return *fs, nil
	}
	curler := curl.NewCurl(backend_config().FeedSourceFolder)
	cache, err := curler.GetUtf8(uri)
	ext := curl.MimeToExt(cache.Mime)
	if ext != "xml" && ext != "atom+xml" && ext != "rss+atom" {
		return v, fmt.Errorf("unsupported mime: %v", cache.Mime)
	}

	if cache.LocalUtf8 != "" {
		v, err = feed.MakeFeedSource(cache.LocalUtf8)
	}
	if v.Uri == "" {
		v.Uri = uri
	}
	if err == nil && v.Uri != "" {
		err = fso.upsert(&v)
	}
	return v, err
}

func feedsource_unsubscribe(url string) error {
	fso := new_feedsource_operator()
	err := fso.drop(url)
	return err
}

func feedsource_show(uri string) ([]feed.FeedSourceFindEntry, error) {
	fs, _, err := google.NewGoogleFeedApi(refer, config.FeedSourceFolder).Load(uri, config.Language, 4, false)
	if err != nil {
		return nil, err
	}
	v := make([]feed.FeedSourceFindEntry, 1)
	v[0].Summary = fs.Description
	v[0].Title = fs.Name
	v[0].Url = fs.Uri
	v[0].Website = fs.WebSite
	return v, err
}

func feedsource_find(q string) ([]feed.FeedSourceFindEntry, error) {
	svc := google.NewGoogleFeedApi(refer, config.FeedSourceFolder)
	v, err := svc.Find(q, config.Language)
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
