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
	return entries
}

func feed_entry_summary_clean(entry *feed.FeedEntry, wg *sync.WaitGroup) {
	defer wg.Done()
	if (entry.Status & (feed.Feed_content_ready | feed.Feed_status_fulltext_ready)) != 0 {
		return
	}
	frag, _ := html_create_fragment(entry.Summary)
	frag, score, _ := cleaner.MakeFragmentReadable(frag)
	entry.Summary = html_encode_fragment(frag)
	entry.Words = uint(score.WordCount)
	entry.Images = append(entry.Images, feed_medias_from_docsummary(score)...)
	if score.WordCount < config.SummaryMinWords {
		entry.Status |= feed.Feed_content_summary_empty
	}
}

func feed_entry_content_clean(entry *feed.FeedEntry, wg *sync.WaitGroup) {
	defer wg.Done()
	if entry.Content != nil && entry.Content.FullText != "" {
		entry.Status |= feed.Feed_status_fulltext_inline
		frag, _ := html_create_fragment(entry.Content.FullText)
		frag, score, _ := cleaner.MakeFragmentReadable(frag)
		entry.Summary = html_encode_fragment(frag)
		entry.Density = uint(score.WordCount)

		entry.Images = append(entry.Images, feed_medias_from_docsummary(score)...)
		entry.Status |= feed.Feed_content_ready
		entry.Status |= feed.Feed_status_fulltext_ready
		log.Println("from-fulltext", entry.Title)
	}
}

func extract_imgsrc_attr(attrs []html.Attribute) []html.Attribute {
	for _, attr := range attrs {
		if attr.Key == "src" {
			return []html.Attribute{html.Attribute{Key: "Source", Val: redirect_thumbnail(attr.Val)}}
		}
	}
	return nil
}

func extract_ahref_attr(attrs []html.Attribute) []html.Attribute {
	for _, attr := range attrs {
		if attr.Key == "href" {
			return []html.Attribute{html.Attribute{Key: "NavigateUri", Val: redirect_link(attr.Val)}}
		}
	}
	return nil
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

//p, img, a, text
func html_to_flowdoc(frag *html.Node) {
	if frag == nil || frag.Type != html.ElementNode {
		return
	}
	switch frag.Data {
	case "img":
		frag.Data = "Figure"
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
