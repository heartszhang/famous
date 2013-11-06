package cleaner

import (
	"code.google.com/p/go.net/html"
	"hash/fnv"
	"io"
)

type DocSummary struct {
	WordCount int      `json:"word_count" bson:"word_count"`
	LinkCount int      `json:"link_count" bson:"link_count"`
	Images    []string `json:"images,omitempty" bson:"images,omitempty"`
	Medias    []string `json:"medias,omitempty" bson:"images,omitempty"`
	Hash      uint64   `json:"hash" bson:"hash"`
}

func new_docsummary(n *html.Node) *DocSummary {
	rtn := &DocSummary{Images: []string{}}
	if n == nil {
		return rtn
	}
	f := fnv.New64()
	foreach_child(n, func(child *html.Node) {
		switch {
		case child.Type == html.CommentNode:
		case child.Type == html.DoctypeNode:
		case child.Type == html.TextNode:
			c, _ := io.WriteString(f, child.Data)
			//_, c, _ := string_count_words(child.Data)
			rtn.WordCount += c
		case child.Data == "img":
			rtn.Images = append(rtn.Images, node_get_attribute(child, "src"))
		case node_is_media(child):
			rtn.Medias = append(rtn.Medias, node_get_attribute(child, "src"))
		case child.Data == "a":
			rtn.LinkCount++
		default:
			sc := new_docsummary(child)
			rtn.add(sc)
		}
	})
	rtn.Hash = f.Sum64()
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
