package backend

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"github.com/heartszhang/cleaner"
	"github.com/heartszhang/curl"
	feed "github.com/heartszhang/feedfeed"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

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
		switch entry.Density == 0 {
		case true:
			entry.Status |= feed.Feed_status_linkdensity_low
		default:
			d := entry.Density * 100 / entry.Words
			switch d < 33 {
			case true:
				entry.Status |= feed.Feed_status_linkdensity_low
			default:
				entry.Status |= feed.Feed_status_linkdensity_high
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

func feed_entry_summary_clean(entry *feed.FeedEntry, wg *sync.WaitGroup) {
	defer wg.Done()
	if (entry.Status & (feed.Feed_content_ready | feed.Feed_content_external_ready)) != 0 {
		return
	}
	score := feedentry_fill_summary(entry, entry.Summary)

	if score.WordCount < config.SummaryMinWords {
		entry.Status |= feed.Feed_summary_empty
	} else {
		entry.Status |= feed.Feed_summary_ready
	}
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

func feedentry_fill_summary(entry *feed.FeedEntry, text string) *cleaner.DocumentSummary {
	frag, _ := html_create_fragment(text)
	frag, score, _ := cleaner.MakeFragmentReadable(frag)
	entry.Words = uint(score.WordCount)
	entry.Density = uint(score.LinkWordCount)
	entry.Images = append(entry.Images, feedmedias_from_docsummary(score.Images)...)
	entry.Videos = append(entry.Videos, feedmedias_from_docsummary(score.Medias)...)
	//	log.Println("text-imgs:", len(entry.Images), len(entry.Videos))
	mc := len(entry.Images) + len(entry.Videos) + len(entry.Audios)
	ext := feedentry_content_exists(score.Hash)
	summary, s := feedentry_make_summary(frag, entry.Words, entry.Density, mc, ext)
	entry.Summary = summary
	entry.Status |= s
	return score
}

func feedentry_content_clean(entry *feed.FeedEntry, wg *sync.WaitGroup) {
	defer wg.Done()
	if entry.Content != nil && entry.Content.FullText != "" {
		entry.Status |= feed.Feed_content_inline
		score := feedentry_fill_summary(entry, entry.Content.FullText)

		if score.WordCount < config.SummaryMinWords {
			entry.Status |= feed.Feed_content_empty
		} else {
			entry.Status |= feed.Feed_content_ready
		}
	}
}

func node_convert_attr(attrs []html.Attribute, origin, updated string, converter func(string) string) []html.Attribute {
	for _, attr := range attrs {
		if attr.Key == origin {
			return []html.Attribute{html.Attribute{Key: updated, Val: converter(attr.Val)}}
		}
	}
	return nil
}
func extract_imgsrc_attr(attrs []html.Attribute) []html.Attribute {
	return node_convert_attr(attrs, "src", "Source", redirect_thumbnail)
}

func extract_ahref_attr(attrs []html.Attribute) []html.Attribute {
	return node_convert_attr(attrs, "href", "NavigateUri", redirect_link)
}

const (
	fdocns = "http://schemas.microsoft.com/winfx/2006/xaml/presentation"
)

func make_image_node(n *html.Node) *html.Node {
	c := &html.Node{Type: html.ElementNode, Data: "BlockUIContainer", DataAtom: n.DataAtom}
	v := &html.Node{Type: html.ElementNode, Data: "Image", DataAtom: n.DataAtom}
	v.Attr = extract_imgsrc_attr(n.Attr)
	c.AppendChild(v)
	return c
}
func make_run_node(n *html.Node) *html.Node {
	v := &html.Node{Type: html.TextNode, Data: "VIDEO", DataAtom: n.DataAtom}
	return v
}
func node_clear_children(frag *html.Node) {
	for child := frag.FirstChild; child != nil; {
		next := child.NextSibling
		frag.RemoveChild(child)
		child = next
	}
}

func node_convert_flowdocument(frag *html.Node, excludeimg bool) {
	if frag.Type == html.TextNode {
		return
	}
	ignore_children := false
	switch frag.Data {
	case "img":
		if excludeimg == true {
			frag.Type = html.CommentNode
			node_clear_children(frag)
			frag.Attr = nil
		} else {
			frag.Data = "Figure"
			node_clear_children(frag)
			frag.AppendChild(make_image_node(frag))
			frag.Attr = nil
		}
		ignore_children = true
		//		frag.Data = "Image"
		//		frag.Attr = extract_imgsrc_attr(frag.Attr)
	case "a":
		frag.Data = "Hyperlink"
		frag.Attr = extract_ahref_attr(frag.Attr)
	case "article":
		frag.Data = "FlowDocument"
		// set namespace dont work
		frag.Attr = []html.Attribute{html.Attribute{Key: "xmlns", Val: fdocns}}
	case "object":
		fallthrough
	case "video":
		fallthrough
	case "audio":
		fallthrough
	case "embed":
		frag.Type = html.CommentNode
		node_clear_children(frag)
		frag.Attr = nil
		ignore_children = true
	case "p":
		fallthrough
	default:
		frag.Data = "Paragraph"
		frag.Attr = nil
	}
	for child := frag.FirstChild; ignore_children == false && child != nil; child = child.NextSibling {
		node_convert_flowdocument(child, excludeimg)
	}
}
func node_is_empty(n *html.Node) bool {
	return n.Type == html.CommentNode ||
		(n.Type == html.ElementNode && (n.Data == "Paragraph" || n.Data == "Hyperlink") && n.FirstChild == nil) ||
		(n.Type == html.TextNode && n.Data == "")
}

func node_clean_empty(n *html.Node) {
	child := n.FirstChild
	for child != nil {
		next := child.NextSibling
		node_clean_empty(child)
		child = next
	}

	if node_is_empty(n) && n.Parent != nil {
		parent := n.Parent
		parent.RemoveChild(n)
	}
}

//p, img, a, text
func make_flowdocument(frag *html.Node, excludeimg bool) string {
	if frag == nil || frag.Type != html.ElementNode {
		return empty_flowdocument
	}
	node_convert_flowdocument(frag, excludeimg)
	node_clean_empty(frag)
	var buffer bytes.Buffer
	html.Render(&buffer, frag) // ignore return error
	return buffer.String()
}

func feedentry_make_summary(frag *html.Node, words, linkwords uint, medias int, dup bool) (v string, s uint64) {
	s = feed.Feed_status_format_flowdocument
	dh := words > 0 && linkwords*100/words > 50
	wm := words < uint(config.SummaryMinWords)
	switch {
	case dup == true:
		v = empty_flowdocument
	case medias == 0 && wm == true:
		//v = node_extract_text(frag)
		v = empty_flowdocument
	case medias == 0 && wm == false:
		v = make_flowdocument(frag, true)
	case medias == 1:
		v = make_flowdocument(frag, true)
	case medias > 1 && wm == false && dh == true:
		v = make_flowdocument(frag, true)
	case medias > 1 && wm == false && dh == false:
		v = make_flowdocument(frag, false)
		s |= feed.Feed_status_media_inline
	case medias > 1 && wm == true:
		v = make_flowdocument(frag, true)
	}
	return
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

func html_create_from_file(filepath string) (doc *html.Node, err error) {
	f, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer f.Close()
	doc, err = html.Parse(f)
	return
}

func html_write_file(article *html.Node, dir string) (string, error) {
	f, err := ioutil.TempFile(dir, "html.")
	if err != nil {
		return "", err
	}
	defer f.Close()
	err = html.Render(f, article)
	return f.Name(), err
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
