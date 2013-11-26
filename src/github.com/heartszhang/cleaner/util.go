package cleaner

import (
	"code.google.com/p/go.net/html"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// 空节点： 文本长度0
// from/input/textarea
// 由空节点组成的节点
func node_is_not_empty(n *html.Node) bool {
	switch n.Type {
	case html.TextNode:
		return len(n.Data) > 0
	case html.ElementNode:
		switch n.Data {
		case "video", "audio", "object", "embed", "img", "a":
			return true
		case "form", "input", "textarea": // why we reserve these nodes?
			return false
		default:
			rtn := false
			foreach_child(n, func(child *html.Node) {
				cne := node_is_not_empty(child)
				rtn = cne || rtn
			})
			return rtn
		}
	default:
		return false
	}
}

func node_has_block_children(n *html.Node) bool {
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		v := node_is_block(child)
		if v {
			return true
		}
	}
	return false
}

func node_has_inline_children(n *html.Node) bool {
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		v := node_is_inline(child)
		if v {
			return true
		}
	}
	return false
}

func node_is_inline(n *html.Node) bool {
	rtn := n.Type == html.TextNode ||
		node_is_inline_element(n) ||
		node_is_object(n)
	return rtn
}

// li可能包含div节点，这个时候li节点是block节点
func node_is_inline_element(n *html.Node) bool {
	switch n.Data {
	case "a", "font", "small", "span", "strong", "em", "dt", "dd", "br", "cite":
		return true
	case "li":
		return li_is_inline_mode(n)
	default:
		return false
	}
}

//img 按照object计算

func node_is_object(n *html.Node) bool {
	return strings_find([]string{"img", "embed", "object", "video", "audio"}, n.Data)
}

func node_is_media(n *html.Node) bool {
	return strings_find([]string{"embed", "audio", "video"}, n.Data)
}

// 不考虑object节点
// ignorable: form
// block-level: div, p, h1-h6, body, html, object, embed, table, ol, ul, dl, video
// inline-level: a, span, strong, br, img, small, font, i
func node_is_block(n *html.Node) bool {
	if n.Type != html.ElementNode && n.Type != html.DocumentNode {
		return false
	}

	switch n.Data {
	case "div", "p", "pre", "h1", "h2", "h3", "h4", "h5", "h6",
		"body", "html", "article", "section", "head", "ol", "ul", "dl",
		"tbody", "td", "tr", "table", "form", "textarea", "input":
		return true
	case "li":
		return !li_is_inline_mode(n)
	default:
		return false
	}

}

var (
	continue_spaces                = regexp.MustCompile("[ \t]+$")
	lb_spaces                      = regexp.MustCompile("[ \t]*[\r\n]+[ \t]*")
	rex             *regexp.Regexp = regexp.MustCompile(`\w+|\d+|[\W\D\S]`)
)

// number
// word
// zh-char
// punc
func string_tokens(t string) []string {
	//	re := regexp.MustCompile(`\w+|\d+|[\W\D\S]`)
	rtn := rex.FindAllString(t, -1)
	return rtn
}

const (
	zh_stop_chars string = "。．？！，、；：“ ”﹃﹄‘ ’﹁﹂（）［］〔〕【】—…—-～·《》〈〉﹏＿."
)

// (any chinese char except puncs or letter)+
func string_is_word(t string) bool {
	rs := []rune(t)
	if len(rs) == 0 {
		return false
	}
	return ((rs[0] > unicode.MaxLatin1 && strings.ContainsRune(zh_stop_chars, rs[0])) == false) || unicode.IsLetter(rs[0])
}

func node_has_children(this *html.Node) bool {
	return this.FirstChild != nil
}

func node_get_attribute(n *html.Node, name string) string {
	for _, a := range n.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

func merge_tail_spaces(txt string) string {
	txt = continue_spaces.ReplaceAllString(txt, "")
	txt = lb_spaces.ReplaceAllString(txt, "\n")
	return txt
}

func node_inner_text_length(n *html.Node) int {
	if n.Type == html.TextNode {
		return len(n.Data)
	}
	if node_is_object(n) {
		w, h := media_get_dim(n)
		alt := node_get_attribute(n, "alt")
		if alt != "" || w*h > 320*240 { // 图片设置了大小，或者alt，可以认为图片为正文内容
			return 140
		}
	}
	// all comments has been removed
	var rtn int
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		rtn += node_inner_text_length(child)
	}
	return rtn
}

func foreach_child(n *html.Node, dof func(*html.Node)) {
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		dof(child)
	}
}

func string_count_words(txt string) (tokens int, words int, commas int) {

	for _, c := range txt {
		if unicode.IsPunct(c) {
			commas++
		}
	}

	tkns := string_tokens(txt)
	tokens = len(tkns)
	for _, token := range tkns {
		if string_is_word(token) {
			words++
		}
	}
	return
}

// without atom
func create_element(name string) (node *html.Node) {
	return &html.Node{Type: html.ElementNode, Data: name}
}

func create_text(txt string) (node *html.Node) {
	return &html.Node{Type: html.TextNode, Data: txt}
}

// 需要换行的li都认为不再是inline模式。这个函数主要使用来检查使用Li构造的menu
func li_is_inline_mode(li *html.Node) bool {
	if li.Parent == nil {
		return false
	}
	lis := 0
	var txtlen int
	foreach_child(li.Parent, func(n *html.Node) {
		txtlen += node_inner_text_length(n)
		lis++
	})
	return txtlen < 60
}

func node_get_element_by_tag_name2(n *html.Node, tag string, set []*html.Node) []*html.Node {
	foreach_child(n, func(child *html.Node) {
		if child.Type == html.ElementNode && child.Data == tag {
			set = append(set, child)
		} else if child.Type == html.ElementNode {
			set = node_get_element_by_tag_name2(child, tag, set)
		}
	})
	return set
}

/*
func clean_element_before_header(body *html.Node, name string) {
	child := body.FirstChild
	for child != nil {
		if child.Type == html.ElementNode && child.Data != name {
			next := child.NextSibling
			body.RemoveChild(child)
			child = next
		} else {
			break
		}
	}
}
*/
// 通过header节点向上查找文本内容超过141的节点
// 如果正文以图片为主，这里就有可能找不到
func find_article_via_header_i(h *html.Node) *html.Node {
	parent := h.Parent
	pcl := 0
	if parent != nil {
		pcl = node_inner_text_length(parent)
	} else {
		return nil
	}
	// 内容超过3行才行，每行大概47个字符
	if pcl > 141 {
		return parent
	}
	return find_article_via_header_i(parent)
}

func node_is_unflatten(b *html.Node) bool {
	return strings_find([]string{"form", "textarea", "input", "embed", "audio", "video"}, b.Data)
}

func deep_clone_element(n *html.Node) (inline *html.Node) {
	inline = shallow_clone_element(n)
	foreach_child(n, func(child *html.Node) {
		i := deep_clone_element(child)
		inline.AppendChild(i)
	})
	return
}

// shallow copy
func shallow_clone_element(n *html.Node) (inline *html.Node) {
	inline = &html.Node{Type: n.Type, Data: n.Data}
	inline.Attr = make([]html.Attribute, len(n.Attr))
	copy(inline.Attr, n.Attr)
	return
}

func node_append_children(src *html.Node, target *html.Node) {
	foreach_child(src, func(child *html.Node) {
		switch {
		case child.Type == html.TextNode:
			target.AppendChild(create_text(child.Data))
		case child.Data == "a" || node_is_object(child):
			// ommit all children elements
			a := shallow_clone_element(child)
			node_append_children(child, a)
			target.AppendChild(a)
		default:
			node_append_children(child, target)
		}
	})
}

func create_p_with_child(n *html.Node) (p *html.Node) {
	p = create_element("p")
	node_append_children(n, p)
	return
}

/*
func node_update_class_attr(b *html.Node, class string) {
	if len(class) > 0 {
		ca := make([]html.Attribute, len(b.Attr)+1)
		copy(ca, b.Attr)
		ca[len(b.Attr)] = html.Attribute{Key: "class", Val: class}
		b.Attr = ca
	}
}
*/
func node_cat_class(b *html.Node, class string) (rtn string) {
	c := node_get_attribute(b, "class")
	id := node_get_attribute(b, "id")
	rtn = class
	if len(c) > 0 {
		rtn = class + "/" + c
	}
	if len(id) > 0 {
		rtn = class + "#" + id
	}
	return
}

func create_html_sketch() (doc *html.Node, body *html.Node, article *html.Node) {
	doc = &html.Node{Type: html.DocumentNode}
	dt := &html.Node{Type: html.DoctypeNode, Data: "html"}
	root := create_element("html")
	body = create_element("body")
	article = create_element("article")
	doc.AppendChild(dt)
	doc.AppendChild(root)
	root.AppendChild(body)
	body.AppendChild(article)
	return
}

/*
func node_clear_decendant_by_type(n *html.Node, tag string) {
	child := n.FirstChild
	for child != nil {
		if child.Type == html.ElementNode && child.Data == tag {
			next := child.NextSibling
			n.RemoveChild(child)
			child = next
		} else {
			node_clear_decendant_by_type(child, tag)
			child = child.NextSibling
		}
	}
}
*/
func int_max(l int, r int) int {
	if l > r {
		return l
	}
	return r
}

func int_min(l int, r int) int {
	if l < r {
		return l
	}
	return r
}

// update only, attribute must exist already
func node_update_attribute(n *html.Node, key string, val string) {
	for idx, attr := range n.Attr {
		if attr.Key == key {
			n.Attr[idx].Val = val
		}
	}
}

func media_get_dim(img *html.Node) (w, h int64) {
	ws := node_get_attribute(img, "width")
	ws = strings.TrimSuffix(ws, "px")
	hs := node_get_attribute(img, "height")
	hs = strings.TrimSuffix(hs, "px")
	var err error
	if w, err = strconv.ParseInt(ws, 0, 0); err != nil {
		w = -1
	}
	if h, err = strconv.ParseInt(hs, 0, 0); err != nil {
		h = -1
	}

	return
}

func node_is_in_a(n *html.Node) bool {
	for p := n.Parent; p != nil; p = p.Parent {
		if p.Type == html.ElementNode && p.Data == "a" {
			return true
		}
	}
	return false
}

/*
func node_clear_children(a *html.Node) {
	for a.FirstChild != nil {
		a.RemoveChild(a.FirstChild)
	}
}
*/
func trim_display_none(n *html.Node) {
	st := node_get_attribute(n, "style")
	if strings.Contains(st, "display") && (strings.Contains(st, "none")) {
		n.Data = "input"
	}
}

func strings_find(set []string, val string) bool {
	for _, i := range set {
		if i == val {
			return true
		}
	}
	return false
}
