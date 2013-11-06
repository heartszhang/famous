package cleaner

import (
	"code.google.com/p/go.net/html"
	"fmt"
)

type readabilitier struct {
	content    []readability_score
	candidates map[*html.Node]*readability_score
	article    *readability_score
	body       *html.Node
}

func readabilitier_make_readable(body *html.Node) (doc, article *html.Node) {
	reader := new_readabilitier(body)
	return reader.create_article()
}

func new_readabilitier(body *html.Node) *readabilitier {
	r := &readabilitier{
		content:    []readability_score{},
		candidates: make(map[*html.Node]*readability_score),
		body:       body}

	r.extract_paragraphs(body)

	var top_candi *readability_score = nil
	for _, candi := range r.candidates {
		candi.content_score = candi.content_score * (100 - candi.link_density()) / 100
		if top_candi == nil || candi.content_score > top_candi.content_score {
			top_candi = candi
		}
	}
	if top_candi != nil {
		r.article = top_candi
	}

	return r
}

/**
* Now that we have the top candidate, look through its siblings for content that might also be related.
* Things like preambles, content split by ads that we removed, etc.
**/
func (this *readabilitier) create_article() (*html.Node, *html.Node) {
	doc, _, article := create_html_sketch()
	if this.article == nil {
		return doc, article
	}
	threshold := max(10, this.article.content_score/5)

	class_name := node_get_attribute(this.article.element, "class")

	if this.article.element.Parent == nil {
		flatten_block_node(this.article.element, article, false, "")
		return doc, article
	}
	foreach_child(this.article.element.Parent, func(neib *html.Node) {
		append := false
		if neib == this.article.element {
			append = true
		} else if ext, ok := this.candidates[neib]; ok {
			cn := node_get_attribute(neib, "class")
			if len(cn) > 0 && cn == class_name {
				append = true
				//				log.Println("append same class", ext)
			}
			if ext.content_score > threshold {
				append = true
				//				log.Println("append high score neib", ext)
			}
		} else if neib.Type == html.ElementNode && neib.Data == "p" {
			sc := new_boilerpipe_score(neib)
			if sc.words > 65 && sc.link_density() < 22 {
				append = true
				//				log.Println("append high p", neib)
			}
		}
		if append {
			flatten_block_node(neib, article, false, "")
		}
	})
	return doc, article
}

func flat_html(body *html.Node) (doc *html.Node, article *html.Node) {
	doc, _, article = create_html_sketch()
	flatten_block_node(body, article, false, "")
	return
}

func (this *readabilitier) make_readability_score(n *html.Node) *readability_score {
	rtn := new_readability_score(n)
	var (
		pext     *readability_score = nil
		grandext *readability_score = nil
	)
	parent := n.Parent
	var grand *html.Node = nil

	if parent != nil {
		if i, ok := this.candidates[parent]; ok {
			pext = i
		} else {
			pext = new_readability_score(parent)
			this.candidates[parent] = pext
		}
		if parent != this.body {
			grand = parent.Parent
		}
	}
	if grand != nil {
		if i, ok := this.candidates[grand]; ok {
			grandext = i
		} else {
			grandext = new_readability_score(grand)
			this.candidates[grand] = grandext
		}
	}
	bc := new_boilerpipe_score(n)
	score := bc.commas + 1
	// wrap lines
	score += min(bc.lines(), 3)
	score += min(bc.imgs*3, 3)
	score += min(bc.img_score, 3)
	rtn.content_score += score

	if pext != nil {
		pext.content_score += score
	}
	if grandext != nil {
		grandext.content_score += score / 2
	}
	return rtn
}

func (this *readabilitier) extract_paragraphs(n *html.Node) {
	switch {
	case node_is_unflatten(n):
		// has only inlines here
	case node_has_inline_children(n):
		this.content = append(this.content, *this.make_readability_score(n))
	default:
		foreach_child(n, func(child *html.Node) {
			this.extract_paragraphs(child)
		})
	}
}

// text-node
// <a>
// <img> <object> <embed> <video> <audio>
// <ul> <ol> <form> <textarea> <input> will be reserved
func flatten_block_node(b *html.Node, article *html.Node, flatt bool, class string) {
	cur_class := node_cat_class(b, class)
	switch {
	case node_is_media(b):
		mp := create_p_with_child(b)
		article.AppendChild(mp)
	case flatt && node_is_unflatten(b): // make unflatten nodes flatted
		nb := create_element(b.Data)
		//		try_update_class_attr(nb, cur_class)
		flatten_block_node(b, nb, false, class)
		article.AppendChild(nb)
	case node_is_unflatten(b):
	case node_has_inline_children(b):
		p := create_p_with_child(b)
		//		try_update_class_attr(p, cur_class)
		article.AppendChild(p)
	default:
		foreach_child(b, func(child *html.Node) {
			flatten_block_node(child, article, true, cur_class)
		})
	}
}

func get_class_weight(n *html.Node, attname string) int {
	c := node_get_attribute(n, attname)

	weight := 0
	if negative.MatchString(c) {
		weight -= 25
	}
	if positive.MatchString(c) {
		weight += 25
	}
	return weight
}

func (this *readabilitier) String() string {
	return fmt.Sprint("readerabilitier content:", len(this.content), ", candidates:", len(this.candidates))
}
