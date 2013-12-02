package backend

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	feed "github.com/heartszhang/feedfeed"
	"io/ioutil"
	"os"
	"strings"
)

type flowdocument_maker struct {
	first_paragraph *html.Node
}

func new_flowdoc_maker() *flowdocument_maker {
	return &flowdocument_maker{}
}

func (this *flowdocument_maker) make(frag *html.Node, imgs []feed.FeedMedia) string {
	if frag == nil || frag.Type != html.ElementNode {
		return empty_flowdocument
	}

	this.convert_flowdocument(frag)
	this.insert_images(imgs)
	node_clean_empty(frag)
	var buffer bytes.Buffer
	html.Render(&buffer, frag) // ignore return error
	body := buffer.String()

	return body
}

func (this *flowdocument_maker) convert_flowdocument(frag *html.Node) {
	if frag.Type == html.TextNode {
		return
	}
	ignore_children := false
	switch frag.Data {
	case "img":
		frag.Type = html.CommentNode
		node_clear_children(frag)
		frag.Attr = nil
	case "a":
		frag.Data = "Hyperlink"
		frag.Attr = extract_ahref_attr(frag.Attr)
	case "article":
		frag.Data = "FlowDocument"
		// set namespace dont work
		frag.Attr = []html.Attribute{html.Attribute{Key: "xmlns", Val: fdocns}}
	case "object", "video", "audio", "embed":
		frag.Type = html.CommentNode
		node_clear_children(frag)
		frag.Attr = nil
	case "p":
		fallthrough
	default:
		frag.Data = "Paragraph"
		frag.Attr = nil
		if this.first_paragraph == nil {
			this.first_paragraph = frag
		}
	}
	for child := frag.FirstChild; ignore_children == false && child != nil; child = child.NextSibling {
		this.convert_flowdocument(child)
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

func extract_imgsrc_attr(img feed.FeedMedia) []html.Attribute {
	return []html.Attribute{html.Attribute{Key: "Source", Val: redirect_thumbnail(img.Uri)}}
}

func extract_ahref_attr(attrs []html.Attribute) []html.Attribute {
	return node_convert_attr(attrs, "href", "NavigateUri", redirect_link)
}

const (
	fdocns = "http://schemas.microsoft.com/winfx/2006/xaml/presentation"
)

func (this *flowdocument_maker) insert_images(imgs []feed.FeedMedia) {
	if this.first_paragraph == nil || len(imgs) == 0 {
		return
	}
	neib := this.first_paragraph.FirstChild

	f := &html.Node{Type: html.ElementNode, Data: "Figure", DataAtom: atom.Div}
	for _, img := range imgs {
		c := &html.Node{Type: html.ElementNode, Data: "BlockUIContainer", DataAtom: atom.Div}
		v := &html.Node{Type: html.ElementNode, Data: "Image", DataAtom: atom.Div}
		v.Attr = []html.Attribute{html.Attribute{Key: "Source", Val: redirect_thumbnail(img.Uri)}}
		// should add image-desc?
		c.AppendChild(v)
		f.AppendChild(c)
	}
	this.first_paragraph.InsertBefore(f, neib)
}

func node_clear_children(frag *html.Node) {
	for child := frag.FirstChild; child != nil; {
		next := child.NextSibling
		frag.RemoveChild(child)
		child = next
	}
}

func node_is_hyperlink_decendant(frag *html.Node) bool {
	for p := frag.Parent; p != nil; p = p.Parent {
		if p.Type == html.ElementNode && (p.Data == "a" || p.Data == "Hyperlink") {
			return true
		}
	}
	return false
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
