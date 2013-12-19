package backend

import (
	"github.com/heartszhang/baidu"
	"github.com/heartszhang/feed"
	"github.com/heartszhang/pubsub"
	"github.com/heartszhang/unixtime"
	"github.com/qiniu/log"
)

const (
	baiduq = "http://iweizhi2.duapp.com/hub/popup.json"
)

func feedentries_updated() (ReadSource, []ReadEntry, error) {
	bcms := baidu.NewBcmsProxy(baiduq)
	var v pubsub.PubsubMessage
	err := bcms.FetchOneAsJson(&v)
	if err != nil {
		return ReadSource{}, nil, err
	}
	if (v.Status.StatusCode != 200 && v.Status.StatusCode != 0) || v.Status.Feed == "" {
		return ReadSource{}, nil, new_backenderror(v.Status.StatusCode, v.Status.StatusReason)
	}
	fs := ReadSource{
		FeedSource: feed.FeedSource{
			Name:        v.Title,
			Uri:         v.Status.Feed,
			Description: v.Subtitle,
			Period:      v.Status.Period / 60,
			Update:      int64(v.Updated),
		},
		LastTouch:  int64(unixtime.TimeNow()),
		NextTouch:  int64(unixtime.Time(v.Status.Period) + unixtime.TimeNow()),
		LastUpdate: int64(unixtime.TimeNow()),
	}
	if fs.Period == 0 {
		fs.Period = 120 // minutes
	}
	fes := make([]ReadEntry, len(v.Items))
	for idx, i := range v.Items {
		fes[idx] = ReadEntry{
			FeedEntry: feed.FeedEntry{
				Uri:     i.Uri,
				Title:   i.Title,
				PubDate: int64(i.Published),
				Summary: i.Summary,
				Content: i.Content,
				Tags:    i.Categories,
			},
		}
		feedentry_init_from_standardlinks(i.StandardLinks, fes[idx])
		if fes[idx].Uri == "" {
			feedentry_init_from_links(i.Links, fes[idx])
		}
	}
	if err == nil {
		fes = readentry_filter(fes)
		fo := new_feedsource_operator()
		// ignore touch error, because, source may not be subscribed
		fo.touch(fs.Uri, int64(fs.LastTouch), int64(fs.NextTouch), fs.Period)
	}
	log.Println("updated", fs.Name, fs.Update)
	return fs, fes, err
}

func feedentry_init_from_standardlinks(links *pubsub.PubsubStandardLink, fe ReadEntry) {
	if links == nil {
		return
	}
	if len(links.Self) > 0 {
		fe.Uri = links.Self[0].Href
	}
	for _, img := range links.Picture {
		fe.Images = append(fe.Images, feed.FeedMedia{
			Title: img.Title,
			Uri:   img.Href,
			Mime:  img.Mime,
		})
	}
}

func feedentry_init_from_links(links []pubsub.PubsubLink, fe ReadEntry) {
	for _, link := range links {
		if link.Rel == "alternate" {
			fe.Uri = link.Href
		}
	}
}

func update_popup() (ReadEntity, error) {
	fs, fes, err := feedentries_updated()
	if err == nil {
		return ReadEntity{ReadSource: fs, Entries: fes}, nil
	}
	log.Println("update-popup", err)
	return ReadEntity{}, err
}
