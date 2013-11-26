// 通过链接密度，判定链接群及前后修辞，将正文中内嵌的广告群去除
// 此算法处理图片能力欠妥，特别是对以图片为主的网页，会将图片去除
package cleaner

import (
	"code.google.com/p/go.net/html"
)

type boilerpiper struct {
	titles                   []string
	authors                  []string
	keywords                 []string
	description              string
	content                  []*boilerpipe_score
	images                   []boiler_image
	words                    int
	lines                    int
	chars                    int
	inner_chars              int
	outer_chars              int
	parags                   int
	words_p_boilerpipe_score float64
	quality                  float64
	body                     *html.Node
	//经过初次处理后得到的html/body节点，或者更准确的article节点
}

func new_boilerpiper(article *html.Node) *boilerpiper {
	rtn := boilerpiper{titles: []string{},
		authors:  []string{},
		keywords: []string{},
		content:  []*boilerpipe_score{}}
	rtn.evaluate_score(article)
	return &rtn
}

func boiler_clean_by_link_density(article *html.Node) (*html.Node, []boiler_image) {
	boiler := new_boilerpiper(article)
	boiler.clean_by_link_density()
	return article, boiler.images
}

//http://www.l3s.de/~kohlschuetter/boilerplate/
//implement
func (this *boilerpiper) clean_by_link_density() {
	for idx, current := range this.content {
		var (
			prev *boilerpipe_score = &boilerpipe_score{}
			next *boilerpipe_score = &boilerpipe_score{}
		)
		if idx != 0 {
			prev = this.content[idx-1]
		}
		if idx < len(this.content)-1 {
			next = this.content[idx+1]
		}
		this.classify(prev, current, next)
	}
	for _, p := range this.content {
		if !p.is_content {
			p.element.Parent.RemoveChild(p.element)
		}
	}
}

// 清除表单前的提示行
func boiler_clean_form_prefix(article *html.Node) *html.Node {
	boiler := new_boilerpiper(article)
	//表单前面的一段内容如果很短，基本上可以认定是form或者menu的标题，可以进行清除
	for idx, current := range boiler.content {
		var next = &boilerpipe_score{}
		if idx < len(boiler.content)-1 {
			next = boiler.content[idx+1]
		}
		if current.is_content && next.forms > 0 && current.words < 16 {
			current.is_content = false
			current.element.Parent.RemoveChild(current.element)
		}
	}
	return article
}

func (this *boilerpiper) evaluate_score(n *html.Node) {
	switch {
	case n.Data == "form" || n.Data == "input" || n.Data == "textarea":
		//form中的内容，仍然需要进行统计
		bs := new_boilerpipe_score(n)
		this.content = append(this.content, &bs)
		//经过前面的整理，如果节点包含inline-nodes，则所有子节点必然都是inline-nodes
		// 说明这是一个段落
	case node_has_inline_children(n):
		bs := new_boilerpipe_score(n)
		this.content = append(this.content, &bs)
	default:
		foreach_child(n, func(child *html.Node) {
			this.evaluate_score(child)
		})
	}
}

const (
	ld_link_para_t        = 13
	ld_link_group_t       = 33
	ld_link_group_title_t = 55
	w_current_line_l      = 20
	w_next_line_l         = 15
	w_prev_line_l         = 8
)

// 链接密度高于0.33的段落，直接认为不是正文
// 第一段是一个图片链接，图片大小确定或者有alt属性的情况下，这个图片是正文内容
// 当前段不多于一定字符，后续段落的链接密度很高，认为这一段落是后续段落的标题，可以进行清除
// form组成的段落，直接抛弃
// many magic numbers in this function
func (this *boilerpiper) classify(prev *boilerpipe_score,
	current *boilerpipe_score,
	next *boilerpipe_score) {
	// doc's picture
	tagged_imgs := len(current.tagged_imgs)
	ifpi := prev.element == nil && current.link_density() > 90 && tagged_imgs == 1
	imgbtx := prev.is_content && current.link_density() > 90 && tagged_imgs >= 1 && tagged_imgs == current.imgs && current.imgs == current.anchor_imgs
	if current.link_density() > 33 && !ifpi && !imgbtx {
		current.is_content = false
	} else {
		c := (prev.link_density() <= 55 &&
			(current.words > 20 || next.words > 15 || prev.words > 8)) ||
			(prev.link_density() > 55 && (current.words > 40 || next.words > 17))
		current.is_content = current.is_content || c

		//images between content paragraphs
		if prev.link_density() <= 33 && next.link_density() <= 33 &&
			current.words == 0 && current.imgs > 0 && current.anchor_imgs == 0 {
			current.is_content = true
		}
		// short paragraphs
		if prev.link_density() < ld_link_para_t && next.link_density() < ld_link_para_t &&
			current.link_density() < ld_link_para_t && current.words < 40 {
			current.is_content = true
		}
	}
	//链接群标题
	if current.words < 15 && next.link_density() > 55 {
		current.is_content = false
	}
	if current.forms > 0 && current.words == 0 {
		current.is_content = false
	}
	if current.is_content == false {
		this.images = append(this.images, current.tagged_imgs...)
	}
}
