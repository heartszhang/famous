package backend

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"github.com/heartszhang/cleaner"
	"github.com/heartszhang/curl"
	feed "github.com/heartszhang/feedfeed"
	"io/ioutil"
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
		go feed_entry_content_clean(&entries[i], wg)
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
	frag, _ := html_create_fragment(entry.Summary)
	frag, score, _ := cleaner.MakeFragmentReadable(frag)
	entry.Summary = html_encode_fragment(frag)
	entry.Words = uint(score.WordCount)
	entry.Images = append(entry.Images, feed_medias_from_docsummary(score)...)
	if feedentry_content_exists(score.Hash) {
		entry.Summary = empty_flowdocument
	}
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
	cnt, _ := co.touch(hash)
	return cnt > config.SummaryDuplicateCount
}
func feed_entry_content_clean(entry *feed.FeedEntry, wg *sync.WaitGroup) {
	defer wg.Done()
	if entry.Content != nil && entry.Content.FullText != "" {
		entry.Status |= feed.Feed_content_inline
		frag, _ := html_create_fragment(entry.Content.FullText)
		frag, score, _ := cleaner.MakeFragmentReadable(frag)
		entry.Summary = html_encode_fragment(frag)
		entry.Density = uint(score.WordCount)

		entry.Images = append(entry.Images, feed_medias_from_docsummary(score)...)
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

//p, img, a, text
func html_to_flowdoc(frag *html.Node) {
	if frag == nil || frag.Type != html.ElementNode {
		return
	}
	switch frag.Data {
	case "img":
		frag.Data = "Figure"
		node_clear_children(frag)
		frag.AppendChild(make_image_node(frag))
		frag.Attr = nil
		return
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
		frag.Data = "Hyperlink"
		node_clear_children(frag)
		frag.AppendChild(make_run_node(frag))
		frag.Attr = node_convert_attr(frag.Attr, "src", "NavigateUri", func(v string) string { return v })
		return
	case "p":
		fallthrough
	default:
		frag.Data = "Paragraph"
		frag.Attr = nil
	}
	for child := frag.FirstChild; child != nil; child = child.NextSibling {
		html_to_flowdoc(child)
	}
}

func html_encode_fragment(frag *html.Node) string {
	html_to_flowdoc(frag)
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
	v := make([]feed.FeedMedia, count)
	for idx, ms := range score.Images {
		v[idx].Uri = ms
	}
	return v
}
