package backend

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"github.com/heartszhang/cleaner"
	"github.com/heartszhang/curl"
	feed "github.com/heartszhang/feedfeed"
	"strings"
	"sync"
)

func feed_entries_unreaded(entries []feed.FeedEntry) []feed.FeedEntry {
	return entries
}

func feed_entries_statis(entries []feed.FeedEntry) []feed.FeedEntry {
	count := len(entries)
	for i := 0; i < count; i++ {
		entry := &entries[i]
		if entry.Content != nil && entry.Content.FullText != "" {
			entry.Status |= feed.Feed_content_ready
		}
		if len(entry.Images) > 0 {
			entry.Status |= feed.Feed_status_has_image
		}
	}
	return entries
}

func feed_entry_content_clean(entry *feed.FeedEntry, wg *sync.WaitGroup) {
	defer wg.Done()
	if entry.Content != nil && entry.Content.FullText != "" {
		entry.Status |= feed.Feed_status_fulltext_inline
		frag, _ := html_create_fragment(entry.Content.FullText)
		frag, score, _ := cleaner.MakeFragmentReadable(frag)
		entry.Content.FullText = html_encode_fragment(frag)
		entry.Content.Words = uint(score.WordCount)
		entry.Content.Medias = feed_medias_from_docsummary(score)
		entry.Images = append(entry.Images, entry.Content.Medias...)

		if len(entry.Content.FullText) > len(entry.Summary) {
			entry.Status |= feed.Feed_content_ready
		}
		if entry.Words > backend_config().SummaryThreshold {
			entry.Status |= feed.Feed_content_summary
		}
		entry.Status |= feed.Feed_status_fulltext_ready
	}
}

func feed_entry_summary_clean(entry *feed.FeedEntry, wg *sync.WaitGroup) {
	defer wg.Done()

	frag, _ := html_create_fragment(entry.Summary)
	frag, score, _ := cleaner.MakeFragmentReadable(frag)
	entry.Summary = html_encode_fragment(frag)
	entry.Words = uint(score.WordCount)
	entry.Images = append(entry.Images, feed_medias_from_docsummary(score)...)
}

func feed_entries_clean_fulltext(entries []feed.FeedEntry) []feed.FeedEntry {
	wg := &sync.WaitGroup{}
	for i := 0; i < len(entries); i++ {
		wg.Add(1)
		go feed_entry_content_clean(&entries[i], wg)
	}
	wg.Wait()
	return entries
}

func feed_entries_clean_summary(entries []feed.FeedEntry) []feed.FeedEntry {
	wg := &sync.WaitGroup{}
	for i := 0; i < len(entries); i++ {
		wg.Add(1)
		go feed_entry_summary_clean(&entries[i], wg)
	}
	wg.Wait()
	return entries
}

// extract tags from summary and fulltext
func feed_entries_autotag(entries []feed.FeedEntry) []feed.FeedEntry {
	return entries
}

// save to mongo
func feed_entries_backup(entries []feed.FeedEntry) []feed.FeedEntry {
	return entries
}

func feed_entries_clean(entries []feed.FeedEntry) []feed.FeedEntry {
	return entries
}

func html_encode_fragment(frag *html.Node) string {
	var buffer bytes.Buffer
	html.Render(&buffer, frag) // ignore return error
	return buffer.String()
}

func html_create_fragment(fulltext string) (*html.Node, error) {
	reader := strings.NewReader(fulltext)

	v := &html.Node{Type: html.ElementNode, Data: "article", DataAtom: atom.Article}
	frags, err := html.ParseFragment(reader, v)
	if err != nil {
		return v, err
	}
	for _, frag := range frags {
		v.AppendChild(frag)
	}
	return v, err
}

const (
	desc_image_worker_count = 16
)

func describe_image(img_url_chan <-chan string, img_chan chan<- feed.FeedMedia) {
	for url := range img_url_chan {
		mt, w, h, l, _ := curl.DescribeImage(url)
		v := feed.FeedMedia{Width: w, Height: h, Length: l, Uri: url, Mime: mt}
		img_chan <- v
	}
}

func feed_medias_from_docsummary(score *cleaner.DocSummary) []feed.FeedMedia {
	count := len(score.Images)
	img_url_chan := make(chan string)

	v := make([]feed.FeedMedia, count)
	img_chan := make(chan feed.FeedMedia, len(score.Images))

	worker_count := desc_image_worker_count
	if count < worker_count {
		worker_count = count
	}
	for i := 0; i < worker_count; i++ {
		go describe_image(img_url_chan, img_chan)
	}

	for _, link := range score.Images {
		img_url_chan <- link
	}
	close(img_url_chan)
	for i := 0; i < count; i++ {
		x := <-img_chan
		v[i] = x
	}
	return v
}
