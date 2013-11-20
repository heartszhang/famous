package backend

import (
	"code.google.com/p/go.net/html"
	"github.com/heartszhang/cleaner"
	"github.com/heartszhang/curl"
	feed "github.com/heartszhang/feedfeed"
	"io/ioutil"
	"log"
	"sync"
)

type text_score struct {
	*cleaner.DocumentSummary
	flowdoc string
	status  uint64
}

func feedentry_eval_text(entry *feed.FeedEntry, text string,
	emptyflag,
	readyflag,
	dupflag,
	inlineflag uint64, disableinline bool) text_score {
	frag, _ := html_create_fragment(text)
	frag, score, _ := cleaner.NewExtractor(config.CleanFolder).MakeFragmentReadable(frag)
	entry.Images = feedmedia_append_unique(entry.Images, feedmedias_from_docsummary(score.Images)...)
	entry.Videos = feedmedia_append_unique(entry.Videos, feedmedias_from_docsummary(score.Medias)...)
	if len(entry.Videos) > 0 {
		fm := feed.FeedMedia{Uri: imageurl_from_video(entry.Videos[0].Uri)}
		entry.Images = feedmedia_append_unique(entry.Images, fm)
	}
	mc := len(entry.Images) + len(entry.Videos) + len(entry.Audios)
	imgs := len(entry.Images)
	ext := feedentry_content_exists(score.Hash) && text != ""

	if score.WordCount < config.SummaryMinWords && mc == 0 {
		log.Println("clean-failed", text, score.Text)
		entry.Title.Others = append(entry.Title.Others, score.Text)
	}
	flowdoc, s := feedentry_make_flowdoc(frag, score.WordCount, score.LinkWordCount, mc, imgs, ext, emptyflag, readyflag, dupflag, inlineflag, disableinline)
	return text_score{score, flowdoc, s}
}

func feedentry_make_flowdoc(frag *html.Node, words, linkwords int, medias, imgs int, dup bool, emptyflag, readyflag, dupflag, inlineflag uint64, disableinline bool) (v string, s uint64) {
	s = feed.Feed_status_format_flowdocument
	dh := words > 0 && linkwords*100/words > 50
	wm := words < config.SummaryMinWords
	switch {
	case dup == true:
		v = empty_flowdocument
		s |= emptyflag | dupflag
	case medias == 0 && wm == true: // no video/audio/image/text
		//v = node_extract_text(frag)
		v = empty_flowdocument
	case medias == 0 && wm == false: // no vi/audio/image, has text
		v = make_flowdocument(frag, true)
	case medias > 0 && wm == true: // only image
		v = empty_flowdocument
		s |= emptyflag
	case medias < 4: // a little images
		v = make_flowdocument(frag, true)
	case medias > 1 && imgs > 0 && medias > imgs: // has au/video
		v = make_flowdocument(frag, true)
	case medias > 1 && wm == false && dh == true: // has many images , text quality is low
		v = make_flowdocument(frag, true)
	case medias > 1 && wm == false && dh == false && !disableinline: //many images and text quality is high
		v = make_flowdocument(frag, false)
		s |= inlineflag
	case medias > 1 && wm == true: // text-quality is low
		v = make_flowdocument(frag, true)
	default:
		v = make_flowdocument(frag, true)
	}
	return
}

func feed_entries_unreaded(entries []feed.FeedEntry) []feed.FeedEntry {
	return entries
}

func feed_entries_clean(entries []feed.FeedEntry) []feed.FeedEntry {
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

func feed_entries_clean_fulltext(entries []feed.FeedEntry) []feed.FeedEntry {
	wg := &sync.WaitGroup{}
	for i := 0; i < len(entries); i++ {
		wg.Add(1)
		go feedentry_content_clean(&entries[i], wg)
	}
	wg.Wait()
	return entries
}

// extract tags from summary and fulltext
func feed_entries_autotag(entries []feed.FeedEntry) []feed.FeedEntry {
	return entries
}

func feed_entries_statis(entries []feed.FeedEntry) []feed.FeedEntry {
	count := len(entries)
	for i := 0; i < count; i++ {
		entry := &entries[i]
		switch entry.Words < uint(config.SummaryMinWords) {
		case true:
			entry.Status |= feed.Feed_status_text_empty
		default:
			entry.Status |= feed.Feed_status_text_many
		}
		if entry.Status&feed.Feed_status_summary_empty != 0 && entry.Status&feed.Feed_status_content_empty != 0 {
			entry.Status |= feed.Feed_status_text_empty
		}
		if entry.Status&(feed.Feed_status_content_mediainline|feed.Feed_status_summary_mediainline) != 0 {
			entry.Status |= feed.Feed_status_media_inline
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
		switch entry.Density < config.LinkDensityThreshuld {
		case true:
			entry.Status |= feed.Feed_status_linkdensity_low
		default:
			entry.Status |= feed.Feed_status_linkdensity_high
		}
		if entry.Status&feed.Feed_status_text_empty != 0 {
			d := entry.Title.Main
			if entry.Status&feed.Feed_status_image_one != 0 && len(entry.Images[0].Description) == 0 {
				entry.Images[0].Description = d
			}
			if entry.Status&feed.Feed_status_media_one != 0 && len(entry.Videos[0].Description) == 0 {
				entry.Videos[0].Description = d
			}
			if entry.Status&feed.Feed_status_image_many != 0 && len(entry.Images[0].Description) == 0 {
				entry.Images[0].Description = d
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
	log.Println("touch-hash", cnt, hash, err)
	return cnt > config.SummaryDuplicateCount
}

func feedmedia_append_unique(set []feed.FeedMedia, v ...feed.FeedMedia) []feed.FeedMedia {
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

func feedentry_write_flowdoc(text string) {
	of, err := ioutil.TempFile(config.FlowDocumentFolder, "xaml.")
	if err != nil {
		return
	}
	defer of.Close()
	of.WriteString(text)
}
func feedentry_write_fails(text string) {
	of, err := ioutil.TempFile(config.FailedFolder, "text.")
	if err != nil {
		return
	}
	defer of.Close()
	of.WriteString(text)
}

func feed_entry_summary_clean(entry *feed.FeedEntry, wg *sync.WaitGroup) {
	defer wg.Done()
	if entry.Summary == "" {
		entry.Status |= feed.Feed_status_summary_empty
		return
	}
	score := feedentry_eval_text(entry, entry.Summary, feed.Feed_status_summary_empty, feed.Feed_status_summary_ready, feed.Feed_status_summary_duplicated, feed.Feed_status_summary_mediainline, true)
	entry.Summary = score.flowdoc
	entry.Status |= score.status
	if entry.Words < uint(score.WordCount) {
		entry.Words = uint(score.WordCount)
		entry.Density = uint(score.LinkWordCount * 100 / score.WordCount)
	}
}

func feedentry_content_clean(entry *feed.FeedEntry, wg *sync.WaitGroup) {
	defer wg.Done()
	if entry.Content == "" {
		entry.Status |= feed.Feed_status_content_empty
		return
	}
	backup := entry.Summary
	entry.Status |= feed.Feed_status_content_inline
	score := feedentry_eval_text(entry, entry.Content, feed.Feed_status_content_empty, feed.Feed_status_content_ready, feed.Feed_status_content_duplicated, feed.Feed_status_content_mediainline, false)
	entry.Status |= score.status
	if entry.Words < uint(score.WordCount) {
		entry.Summary = score.flowdoc
		entry.Words = uint(score.WordCount)
		entry.Density = uint(score.LinkWordCount * 100 / score.WordCount)
		entry.Content = backup
	} else {
		entry.Content = score.flowdoc
	}
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

func feedmedias_from_docsummary(medias []cleaner.MediaSummary) []feed.FeedMedia {
	count := len(medias)
	v := make([]feed.FeedMedia, count)
	for idx, ms := range medias {
		v[idx].Uri = ms.Uri
		v[idx].Width = int(ms.Width)
		v[idx].Height = int(ms.Height)
		v[idx].Description = ms.Alt
	}
	return v
}
