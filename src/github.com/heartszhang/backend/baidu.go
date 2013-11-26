package backend

import (
	"fmt"
	"github.com/heartszhang/baidu"
	"github.com/heartszhang/feedfeed"
	"github.com/heartszhang/pubsub"
	"time"
)

const (
	baiduq = "http://iweizhi2.duapp.com/hub/popup.json"
)

func feedentries_updated() (*feedfeed.FeedSource, []feedfeed.FeedEntry, error) {
	bcms := baidu.NewBcmsProxy(baiduq)
	var v pubsub.PubsubMessage
	err := bcms.FetchOneAsJson(&v)
	if err != nil {
		return nil, nil, err
	}
	if v.Status.StatusCode != 200 {
		return nil, nil, fmt.Errorf("%d: %v", v.Status.StatusCode, v.Status.StatusReason)
	}

	fs := feedfeed.FeedSource{
		Name:        v.Title,
		Uri:         v.Status.Feed,
		Update:      int64(v.Updated),
		Description: v.Subtitle,
		LastTouch:   unix_current(),
		NextTouch:   int64(v.Status.Period) + unix_current(),
		LastUpdate:  unix_current(),
	}
	fes := make([]feedfeed.FeedEntry, len(v.Items))
	for idx, i := range v.Items {
		fes[idx] = feedfeed.FeedEntry{
			Uri:     i.PermalinkUrl,
			Title:   feedfeed.FeedTitle{Main: i.Title},
			PubDate: int64(i.Published),
			Summary: i.Summary,
			Content: i.Content,
			Tags:    i.Categories,
		}
		feedentry_init_from_standardlinks(i.StandardLinks, fes[idx])
	}
	if err == nil {
		fes = feedentry_filter(fes)
		err = new_feedsource_operator().save_one(fs)
	}
	return &fs, fes, err
}
func feedentry_filter(v []feedfeed.FeedEntry) []feedfeed.FeedEntry {
	v = feed_entries_unreaded(v) // clean readed entries
	v = feed_entries_clean(v)
	v = feed_entries_clean_summary(v)
	v = feed_entries_clean_fulltext(v)
	v = feed_entries_autotag(v)
	v = feed_entries_statis(v)
	v = feed_entries_backup(v)
	return v
}
func feedentry_init_from_standardlinks(links *pubsub.PubsubStandardLink, fe feedfeed.FeedEntry) {
	if links == nil {
		return
	}
	if len(links.Self) > 0 {
		fe.Uri = links.Self[0].Href
	}
	for _, img := range links.Picture {
		fe.Images = append(fe.Images, feedfeed.FeedMedia{
			Title: img.Title,
			Uri:   img.Href,
			Mime:  img.Mime,
		})
	}
}

func unix_current() int64 {
	return time.Now().Unix()
}
