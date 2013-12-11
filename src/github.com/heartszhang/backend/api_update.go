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

func feedentries_updated() (*feed.FeedSource, []feed.FeedEntry, error) {
	bcms := baidu.NewBcmsProxy(baiduq)
	var v pubsub.PubsubMessage
	err := bcms.FetchOneAsJson(&v)
	log.Println(v, err)
	if err != nil {
		return nil, nil, err
	}
	if (v.Status.StatusCode != 200 && v.Status.StatusCode != 0) || v.Status.Feed == "" {
		return nil, nil, new_backenderror(v.Status.StatusCode, v.Status.StatusReason)
	}

	fs := feed.FeedSource{
		FeedSourceMeta: feed.FeedSourceMeta{
			Name:        v.Title,
			Uri:         v.Status.Feed,
			Description: v.Subtitle,
			Period:      v.Status.Period / 60,
		},
		Update:     v.Updated,
		LastTouch:  unixtime.UnixTimeNow(),
		NextTouch:  unixtime.UnixTime(v.Status.Period) + unixtime.UnixTimeNow(),
		LastUpdate: unixtime.UnixTimeNow(),
	}
	if fs.Period == 0 {
		fs.Period = 120 // minutes
	}
	fes := make([]feed.FeedEntry, len(v.Items))
	for idx, i := range v.Items {
		fes[idx] = feed.FeedEntry{
			FeedEntryMeta: feed.FeedEntryMeta{
				Uri:     i.Uri,
				Title:   feed.FeedTitle{Main: i.Title},
				PubDate: i.Published,
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
		fes = feedentry_filter(fes)
		err = new_feedsource_operator().touch(fs.Uri, int64(fs.LastTouch), int64(fs.NextTouch), fs.Period)
		log.Println("updated", fs.Name, fs.Update, fs.Logo)
	}
	return &fs, fes, err
}

func feedentry_init_from_standardlinks(links *pubsub.PubsubStandardLink, fe feed.FeedEntry) {
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

func feedentry_init_from_links(links []pubsub.PubsubLink, fe feed.FeedEntry) {
	for _, link := range links {
		if link.Rel == "alternate" {
			fe.Uri = link.Href
		}
	}
}

func update_popup() (*feed.FeedEntity, error) {
	fs, fes, err := feedentries_updated()
	if err == nil {
		return &feed.FeedEntity{FeedSource: *fs, Entries: fes}, nil
	}
	return nil, err
}
