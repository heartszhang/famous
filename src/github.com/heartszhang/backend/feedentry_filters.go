package backend

import (
	//	"code.google.com/p/go.net/html"
	"github.com/heartszhang/cleaner"
	"github.com/heartszhang/curl"
	"github.com/heartszhang/feed"
	"github.com/heartszhang/markhtml"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

func feedentry_filter(v []feed.FeedEntry) []feed.FeedEntry {
	v = feed_entries_downloaded(v) // clean downloaded entries
	v = feed_entries_clean(v)
	v = feed_entries_clean_text(v)
	v = feed_entries_autotag(v) // convert category to tags, extract tags
	v = feed_entries_statis(v)  // setup flags
	v = feed_entries_backup(v)  // save to db
	return v
}

func feed_entries_downloaded(entries []feed.FeedEntry) []feed.FeedEntry {
	result := make(map[string]feed.FeedEntry)
	var hashs []string
	for _, entry := range entries {
		hash := entry.Uri + "|" + entry.Title.Main
		result[hash] = entry
		hashs = append(hashs, hash)
	}
	hashs, err := new_feedentrytouch_operator().touch(hashs)
	if err != nil {
		return entries
	}
	var v []feed.FeedEntry
	for _, hash := range hashs {
		v = append(v, result[hash])
	}
	return v
}

func feed_entries_clean(entries []feed.FeedEntry) []feed.FeedEntry {
	return entries
}

func feed_entries_clean_text(entries []feed.FeedEntry) []feed.FeedEntry {
	wg := &sync.WaitGroup{}
	for i := 0; i < len(entries); i++ {
		wg.Add(1)
		go feedentry_clean_text(&entries[i], wg)
	}
	wg.Wait()
	return entries
}

func feedentry_clean_text(entry *feed.FeedEntry, wg *sync.WaitGroup) {
	defer wg.Done()
	entry.Summary, entry.SummaryStatus = make_text_readable(entry, entry.Summary, true, false)
	entry.Content, entry.ContentStatus = make_text_readable(entry, entry.Content, false, true)
	// sum_empty := entry.SummaryStatus.Status&feed.Feed_status_content_empty != 0
	con_empty := entry.ContentStatus.Status&feed.Feed_status_content_empty != 0
	if !con_empty {
		backup := entry.Summary
		bs := entry.SummaryStatus
		entry.Summary = entry.Content
		entry.SummaryStatus = entry.ContentStatus
		entry.ContentStatus = bs
		entry.Content = backup
	}
	diff := feed.Feed_status_summary_empty / feed.Feed_status_content_empty
	entry.Status |= entry.SummaryStatus.Status * diff
}

func make_text_readable(entry *feed.FeedEntry, txt string, trans, insimg bool) (string, feed.FeedTextStatus) {
	var status feed.FeedTextStatus
	if txt == "" {
		status.Status = status.Status | feed.Feed_status_content_empty
		return empty_flowdocument, status
	}
	if trans {
		txt = markhtml.TransferText(txt)
	}
	redirector := func(uri string) string {
		return redirect_thumbnail(url_resolve(entry.Uri, uri))
	}
	imgurl_maker := func(uri string) string {
		u := url_resolve(entry.Uri, uri)
		return imageurl_from_video(u)
	}
	frag, _ := html_create_fragment(txt)
	frag, score, _ := cleaner.NewExtractor(backend_context.config.CleanFolder).MakeFragmentReadable(frag)
	entry.Images = append_unique(entry.Images, feedmedias_from_docsummary(score.Images, redirector)...)
	entry.Images = append_unique(entry.Images, feedmedias_from_docsummary(score.Medias, imgurl_maker)...)

	entry.Videos = append_unique(entry.Videos, feedmedias_from_docsummary(score.Medias, func(o string) string { return o })...)
	status.WordCount = score.WordCount
	status.LinkCount = score.LinkCount
	status.LinkWordCount = score.LinkWordCount
	if status.WordCount < backend_config().SummaryMinWords {
		if status.WordCount > 0 {
			entry.Title.Others = append(entry.Title.Others, score.Text)
		}
		status.Status = status.Status | feed.Feed_status_content_empty
	}
	if status.WordCount > 0 && feedentry_content_exists(score.Hash) {
		status.Status = status.Status | feed.Feed_status_content_duplicated
	}
	status.Status |= feed.Feed_status_content_ready
	imgs := entry.Images
	if insimg == false || len(imgs) > 1 ||
		len(entry.Videos) > 0 ||
		status.WordCount < backend_config().SummaryMinWords {
		imgs = nil
	} else if len(imgs) > 0 {
		status.Status |= feed.Feed_status_content_mediainline
	}
	return new_flowdoc_maker().make(frag, imgs), status
}

// extract tags from summary and fulltext
func feed_entries_autotag(entries []feed.FeedEntry) []feed.FeedEntry {
	return entries
}

func feed_entries_statis(entries []feed.FeedEntry) []feed.FeedEntry {
	count := len(entries)
	for i := 0; i < count; i++ {
		entry := &entries[i]
		if entry.Status&feed.Feed_status_summary_empty != 0 && entry.Status&feed.Feed_status_content_empty != 0 {
			entry.Status |= feed.Feed_status_text_empty
		}
		switch len(entry.Audios) + len(entry.Videos) {
		case 0:
			entry.Status |= feed.Feed_status_media_empty
		case 1:
			entry.Status |= feed.Feed_status_media_one
		default:
			entry.Status |= feed.Feed_status_media_many
		}
		switch len(entry.Images) {
		case 0:
			entry.Status |= feed.Feed_status_image_empty
		case 1:
			entry.Status |= feed.Feed_status_image_one
		default:
			entry.Status |= feed.Feed_status_image_many
		}
		if entry.Status&feed.Feed_status_text_empty != 0 {
			d := strings.Join(entry.Title.Others, "\n")
			if d == "" {
				d = entry.Title.Main
			}
			for i := len(entry.Videos) - 1; i >= 0; i = i - 1 {
				if entry.Videos[i].Description == "" {
					entry.Videos[i].Description = d
				}
			}
			for i := len(entry.Audios) - 1; i >= 0; i = i - 1 {
				if entry.Audios[i].Description == "" {
					entry.Audios[i].Description = d
				}
			}
			for i := len(entry.Images) - 1; i >= 0; i = i - 1 {
				if entry.Images[i].Description == "" {
					entry.Images[i].Description = d
				}
			}
		}
	}
	return entries
}

// save to mongo
func feed_entries_backup(entries []feed.FeedEntry) []feed.FeedEntry {
	if entries == nil || len(entries) == 0 {
		return entries
	}
	feo := new_feedentry_operator()
	feo.save(entries)
	return entries
}

const (
	empty_flowdocument = `<FlowDocument xmlns="http://schemas.microsoft.com/winfx/2006/xaml/presentation"/>`
)

func feedentry_content_exists(hash uint64) bool {
	co := new_feedcontent_operator()
	cnt, err := co.touch(int64(hash))
	if err != nil {
		log.Println("feedcontent-hashs failed", err)
		return false
	}
	return cnt > backend_context.config.SummaryDuplicateCount
}

func append_unique(set []feed.FeedMedia, v ...feed.FeedMedia) []feed.FeedMedia {
	hav := make(map[string]bool)
	for _, s := range set {
		hav[s.Uri] = true
	}
	for _, a := range v {
		if !hav[a.Uri] {
			hav[a.Uri] = true
			set = append(set, a)
		}
	}
	return set
}

func feedentry_log_fails(text string) {
	of, err := ioutil.TempFile(backend_context.config.FailedFolder, "text.")
	if err != nil {
		return
	}
	defer of.Close()
	of.WriteString(text)
}

const (
	desc_image_worker_count = 4
)

func describe_image(img_url_chan <-chan string, img_chan chan<- feed.FeedMedia) {
	for url := range img_url_chan {
		mt, w, h, l, _ := curl.DescribeImage(url)
		v := feed.FeedMedia{Width: w, Height: h, Length: l, Uri: url, Mime: mt}
		img_chan <- v
	}
}

func feedmedias_from_docsummary(medias []cleaner.MediaSummary, redirector func(string) string) []feed.FeedMedia {
	count := len(medias)
	v := make([]feed.FeedMedia, count)
	for idx, ms := range medias {
		v[idx].Uri = redirector(ms.Uri)
		v[idx].Width = int(ms.Width)
		v[idx].Height = int(ms.Height)
		v[idx].Description = ms.Alt
	}
	return v
}

func text_unique_append(set []string, txt string) []string {
	for _, s := range set {
		if s == txt {
			return set
		}
	}
	return append(set, txt)
}
