package cleaner

import (
	"code.google.com/p/go.net/html"
	"fmt"
)

type boilerpipe_score struct {
	element      *html.Node
	words        int // word or zh-char
	tokens       int // number, punc, zh-char or word
	anchor_words int
	imgs         int
	anchor_imgs  int
	objects      int
	forms        int
	anchors      int
	commas       int
	inner_text   string
	img_score    int
	is_content   bool
}

func new_boilerpipe_score_omit_table(n *html.Node, omit bool, omit_form bool) boilerpipe_score {
	p := boilerpipe_score{element: n}
	switch {
	case n.Type == html.TextNode:
		p.inner_text += n.Data
		t, w, c := string_count_words(n.Data)
		p.tokens += t
		p.words += w
		p.commas += c
	case n.Data == "a":
		foreach_child(n, func(child *html.Node) {
			np := new_boilerpipe_score_omit_table(child, omit, omit_form)
			p.add(np)
			p.anchor_words += np.words
			p.anchor_imgs += np.imgs
		})
		p.anchors++
	case n.Data == "img":
		width, height := get_image_dim(n)
		if width > 320 || height > 320 || (width < 0 && height < 0) {
			p.imgs++
			p.img_score = int_min(p.img_score+int((width/21)*(height/21)/30), 140)
		}
	case node_is_media(n):
		mw, wh := get_image_dim(n)
		if mw > 400 {
			p.objects++
			p.img_score = int_min(p.img_score+int((mw/21)*(wh/21)/11), 140)
		}
	case omit_form && n.Data == "form":
		p.forms++
	case n.Data == "input" || n.Data == "textarea":
		p.forms++
	case omit && n.Data == "table":
	default:
		foreach_child(n, func(child *html.Node) {
			np := new_boilerpipe_score_omit_table(child, omit, omit_form)
			p.add(np)
		})
	}
	p.words += int_min(p.img_score, 40)
	return p
}

//包含n的子孙的评分
func new_boilerpipe_score(n *html.Node) boilerpipe_score {
	return new_boilerpipe_score_omit_table(n, false, true)
}

func (this *boilerpipe_score) add(rhs boilerpipe_score) {
	this.anchors += rhs.anchors
	this.anchor_words += rhs.anchor_words
	this.inner_text += rhs.inner_text
	this.tokens += rhs.tokens
	this.words += rhs.words
	this.anchor_imgs += rhs.anchor_imgs
	this.imgs += rhs.imgs
	this.img_score = int_min(this.img_score+rhs.img_score, 140)
	this.objects += rhs.objects
	//  this.forms += rhs.forms
}

func int_min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//有链接链接文字的情况，认为全部是图片链接
func (this *boilerpipe_score) link_density() int {
	switch {
	case this.words == 0 && this.anchors > 0:
		return 100
	case this.words == 0 && this.anchors == 0:
		return 0
	default:
		return (this.anchor_words + this.anchor_imgs*4) * 100 / (this.words + this.anchor_imgs*4)
	}
}

const (
	wordwrap = 65
)

func (this boilerpipe_score) lines() int {
	return (this.words + wordwrap - 1) / wordwrap
}
func (this boilerpipe_score) wrapped_words() int {
	return this.words - (this.words % wordwrap)
}

func (this boilerpipe_score) table_score() int {
	return this.words*(100-this.link_density())/100 + (this.imgs-this.anchor_imgs)*8
}

func (this boilerpipe_score) String() string {
	return fmt.Sprint("boilerpipe-score node-tag:", this.element.Data,
		", words:", this.words, ", anchor_words:", this.anchor_words,
		", imgs:", this.imgs, ", aimgs", this.anchor_imgs)
}
