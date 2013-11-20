package backend

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"io/ioutil"
	"os"
	"strings"
)

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

//p, img, a, text
func make_flowdocument(frag *html.Node, excludeimg bool) string {
	if frag == nil || frag.Type != html.ElementNode {
		return empty_flowdocument
	}
	node_convert_flowdocument(frag, excludeimg)
	node_clean_empty(frag)
	var buffer bytes.Buffer
	html.Render(&buffer, frag) // ignore return error
	body := buffer.String()

	return body
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
		} else if node_is_hyperlink_decendant(frag) == false {
			frag.Data = "Figure"
			node_clear_children(frag)
			frag.AppendChild(make_image_node(frag))
			frag.Attr = nil
		} else {
			frag.Data = "Image"
			frag.Attr = extract_imgsrc_attr(frag.Attr)
		}
		ignore_children = true
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

func node_is_hyperlink_decendant(frag *html.Node) bool {
	for p := frag.Parent; p != nil; p = p.Parent {
		if p.Type == html.ElementNode && (p.Data == "a" || p.Data == "Hyperlink") {
			return true
		}
	}
	return false
}
