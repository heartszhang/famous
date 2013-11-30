package backend

import (
	"fmt"
	"github.com/heartszhang/baidu"
	"github.com/heartszhang/feedfeed"
	"github.com/heartszhang/pubsub"
	"github.com/heartszhang/unixtime"
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
	if (v.Status.StatusCode != 200 && v.Status.StatusCode != 0) || v.Status.Feed == "" {
		return nil, nil, fmt.Errorf("%d: %v", v.Status.StatusCode, v.Status.StatusReason)
	}

	fs := feedfeed.FeedSource{
		Name:        v.Title,
		Uri:         v.Status.Feed,
		Update:      v.Updated,
		Description: v.Subtitle,
		LastTouch:   unixtime.UnixTimeNow(),
		NextTouch:   unixtime.UnixTime(v.Status.Period) + unixtime.UnixTimeNow(),
		LastUpdate:  unixtime.UnixTimeNow(),
		Period:      v.Status.Period / 60,
	}
	if fs.Period == 0 {
		fs.Period = 120 // minutes
	}
	fes := make([]feedfeed.FeedEntry, len(v.Items))
	for idx, i := range v.Items {
		fes[idx] = feedfeed.FeedEntry{
			Uri:     i.Uri,
			Title:   feedfeed.FeedTitle{Main: i.Title},
			PubDate: i.Published,
			Summary: i.Summary,
			Content: i.Content,
			Tags:    i.Categories,
		}
		feedentry_init_from_standardlinks(i.StandardLinks, fes[idx])
		if fes[idx].Uri == "" {
			feedentry_init_from_links(i.Links, fes[idx])
		}
	}
	if err == nil {
		fes = feedentry_filter(fes)
		err = new_feedsource_operator().touch(fs.Uri, int64(fs.LastTouch), int64(fs.NextTouch), fs.Period)
	}
	return &fs, fes, err
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

func feedentry_init_from_links(links []pubsub.PubsubLink, fe feedfeed.FeedEntry) {
	if len(links) == 0 {
		return
	}

	for _, link := range links {
		if link.Rel == "alternate" {
			fe.Uri = link.Href
		}
	}
}
