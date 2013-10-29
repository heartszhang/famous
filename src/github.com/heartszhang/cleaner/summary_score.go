package cleaner

import (
	"code.google.com/p/go.net/html"
)

type DocSummary struct {
	WordCount int      `json:"word_count" bson:"word_count"`
	LinkCount int      `json:"link_count" bson:"link_count"`
	Images    []string `json:"image,omitempty" bson:"image,omitempty"`
}

func new_docsummary(n *html.Node) *DocSummary {
	rtn := &DocSummary{Images: []string{}}
	if n == nil {
		return rtn
	}
	foreach_child(n, func(child *html.Node) {
		switch {
		case child.Type == html.CommentNode:
		case child.Type == html.DoctypeNode:
		case child.Type == html.TextNode:
			_, c, _ := string_count_words(child.Data)
			rtn.WordCount += c
		case child.Data == "img":
			rtn.Images = append(rtn.Images, node_get_attribute(child, "src"))
		case child.Data == "a":
			rtn.LinkCount++
		default:
			sc := new_docsummary(child)
			rtn.add(sc)
		}
	})
	return rtn
}

func (this *DocSummary) add(l *DocSummary) {
	if l == nil {
		return
	}
	this.WordCount += l.WordCount
	this.LinkCount += l.LinkCount
	imgs := make([]string, len(this.Images)+len(l.Images))
	copy(imgs, this.Images)
	copy(imgs[len(this.Images):], l.Images)
	this.Images = imgs
}
