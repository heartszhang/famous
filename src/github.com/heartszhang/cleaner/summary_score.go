package cleaner

import (
	"code.google.com/p/go.net/html"
	"hash"
	"hash/fnv"
)

type MediaSummary struct {
	Uri    string
	Alt    string
	Width  int64
	Height int64
}
type DocumentSummary struct {
	WordCount     int
	LinkCount     int
	LinkWordCount int
	Images        []MediaSummary
	Medias        []MediaSummary
	Hash          uint64
	Text          string
}

func make_mediasummary(m *html.Node) MediaSummary {
	v := MediaSummary{
		Uri: node_get_attribute(m, "src"),
		Alt: node_get_attribute(m, "alt"),
	}
	v.Width, v.Height = media_get_dim(m)
	return v
}
func new_docsummary_internal(n *html.Node, f hash.Hash64) *DocumentSummary {
	rtn := &DocumentSummary{}
	if n == nil {
		return rtn
	}
	foreach_child(n, func(child *html.Node) {
		switch {
		case child.Type == html.CommentNode:
		case child.Type == html.DoctypeNode:
		case child.Type == html.TextNode:
			c, _ := f.Write([]byte(child.Data))
			rtn.WordCount += c
			if node_is_in_a(child) {
				rtn.LinkWordCount += c
			}
			rtn.Text += child.Data
		case child.Data == "img":
			rtn.Images = append(rtn.Images, make_mediasummary(child))
		case node_is_media(child):
			rtn.Medias = append(rtn.Medias, make_mediasummary(child))
		case child.Data == "a":
			rtn.LinkCount++
			ac := new_docsummary_internal(child, f)
			rtn.Images = append(rtn.Images, ac.Images...)
			rtn.Medias = append(rtn.Medias, ac.Medias...)
		default:
			sc := new_docsummary_internal(child, f)
			rtn.add(sc)
		}
	})
	return rtn
}
func new_docsummary(n *html.Node, images []boiler_image) *DocumentSummary {
	f := fnv.New64()
	rtn := new_docsummary_internal(n, f)
	rtn.Hash = f.Sum64()
	for _, img := range images {
		rtn.Images = append(rtn.Images, MediaSummary{img.url, img.alt, img.width, img.height})
	}
	return rtn
}

func (this *DocumentSummary) add(l *DocumentSummary) {
	if l == nil {
		return
	}
	this.WordCount += l.WordCount
	this.LinkCount += l.LinkCount
	this.LinkWordCount += l.LinkWordCount
	this.Images = append(this.Images, l.Images...)
	this.Medias = append(this.Medias, l.Medias...)
	this.Text += l.Text
}
