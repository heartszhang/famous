package cleaner

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type html_cleaner struct {
	may_be_html5 bool
	current_url  *url.URL
	article      *html.Node // body or article or a table's body
	head         *html.Node
	header1s     []*html.Node
	header2s     []*html.Node
	header3s     []*html.Node
	header4s     []*html.Node
	uls          []*html.Node
	ols          []*html.Node
	forms        []*html.Node
	tables       []*html.Node
	tds          []boilerpipe_score
	pages        []string
	titles       []string
	keywords     []string
	author       []string
	text_words   int
	anchor_words int
	table_words  int
	links        int
	imgs         int
	link_imgs    int
	lis          int
	description  string
}

/*
茅于轼 | 中国是个忘恩负义的国家吗？ - 中国数字时代
中国抗议日内阁成员参拜靖国神社 - BBC中文网 - 两岸
媒体札记：胜利者姿态 - 评论 - FT中文网
发现已知最大的贼兽: 金氏树贼兽(图)  - 阿波罗新闻网
GFW BLOG（功夫网与翻墙）: 通过 ToyVPN 网站获取 5 个免费的 PPTP VPN 帐号
\r\t[导入]VK Cup 2012 Qualification Round 1    E. Phone Talks - ACM博客_kuangbin - C++博客
译言网 | 南非零售销售额六月份缓慢增长
南方周末 - 广州公安局原副局长受贿600余万被起诉
*/
func (cleaner *html_cleaner) grab_title(title *html.Node) {

}

func html_clean_fragment(root *html.Node) *html.Node {
	cleaner := &html_cleaner{}
	var (
		dropping []*html.Node = []*html.Node{}
	)
	cleaner.clean_unprintable_element(&dropping, root)

	for _, drop := range dropping {
		p := drop.Parent
		p.RemoveChild(drop)
	}

	if cleaner.head != nil {
		cleaner.head.Parent.RemoveChild(cleaner.head)
	}
	cleaner.article = root
	cleaner.clean_body()

	cleaner.clean_empty_nodes(cleaner.article)
	cleaner.clean_attributes(cleaner.article)

	return cleaner.article
}

func html_clean_root(root *html.Node, uribase string) *html.Node {
	cleaner := &html_cleaner{}
	cleaner.current_url, _ = url.Parse(uribase)

	var (
		dropping []*html.Node = []*html.Node{}
	)
	cleaner.clean_unprintable_element(&dropping, root)

	for _, drop := range dropping {
		p := drop.Parent
		p.RemoveChild(drop)
	}

	if cleaner.head != nil {
		cleaner.head.Parent.RemoveChild(cleaner.head)
	}

	//文档中如果只有一个h1,通常这个h1所在的div就是文档内容
	if len(cleaner.header1s) == 1 { // only one h1
		ab := find_article_via_header_i(cleaner.header1s[0])
		cleaner.try_update_article(ab)
	}
	//如果文档中只有一个h2，这时又没有h1，h2就是其中的标题，所在的div就是文档内容
	if len(cleaner.header1s) == 0 && len(cleaner.header2s) == 1 {
		ab := find_article_via_header_i(cleaner.header2s[0])
		cleaner.try_update_article(ab)
	}

	if cleaner.article == nil {

		cleaner.article = &html.Node{Type: html.ElementNode,
			DataAtom: atom.Body,
			Data:     "body"}
		root.AppendChild(cleaner.article)
	}
	cleaner.try_catch_phpwnd()
	cleaner.fix_forms()

	cleaner.clean_body()

	cleaner.clean_empty_nodes(cleaner.article)
	cleaner.clean_attributes(cleaner.article)

	return cleaner.article
}

func (this *html_cleaner) try_catch_phpwnd() {
	// have not table, or some  content not in table
	if len(this.tds) == 0 || this.table_words*100/(this.text_words+1) < 33 {
		return
	}
	top := boilerpipe_score{}
	for _, td_score := range this.tds {
		if top.element == nil || top.table_score() < td_score.table_score() {
			top = td_score
		}
	}
	if top.element == nil {
		return
	}
	//remove_decentant(top.element, "table")
	this.article = top.element
}

var (
	unlikely *regexp.Regexp = regexp.MustCompile(`combx|comment|community|disqus|extra|foot|header|menu|remark|rss|shoutbox|sidebar|sponsor|ad-break|agegate|pagination|pager|popup|tweet|twitter`)
)

//清除所有的脚本，css和Link等等不能显示的内容
//多文档结构进行统计
func (cleaner *html_cleaner) clean_unprintable_element(dropping *[]*html.Node, n *html.Node) {
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.CommentNode {
			*dropping = append(*dropping, child)
		} else if child.Type == html.ElementNode {
			drop := false
			child.Data = strings.ToLower(child.Data)
			idc := node_get_attribute(child, "class") + node_get_attribute(child, "id")

			if unlikely.MatchString(idc) {
				drop = true
				*dropping = append(*dropping, child)
				//				log.Println("dropping by class-id", idc, ", of ", child.Data)
			} else {
				switch child.Data {
				case "script", "link", "iframe", "nav", "aside", "noscript", "style", "input", "textarea", "marquee":
					*dropping = append(*dropping, child)
					drop = true
				case "meta":
					cleaner.grab_keywords(child)
					cleaner.grab_description(child)
				case "title":
					cleaner.grab_title(child)
				case "head":
					cleaner.head = child
				case "body":
					cleaner.article = child
				case "br":
					child.Data = "p"
				case "article":
					// a html may have more article nodes,
				case "h1":
					cleaner.header1s = append(cleaner.header1s, child)
				case "h2":
					cleaner.header2s = append(cleaner.header2s, child)
				case "h3":
					cleaner.header3s = append(cleaner.header3s, child)
				case "h4":
					cleaner.header4s = append(cleaner.header4s, child)
				case "form":
					cleaner.forms = append(cleaner.forms, child)
				case "ul":
					cleaner.uls = append(cleaner.uls, child)
				case "ol":
					cleaner.ols = append(cleaner.ols, child)
				case "table":
					cleaner.tables = append(cleaner.tables, child)
				case "td":
					ts := new_boilerpipe_score_omit_table(child, true, true)
					cleaner.tds = append(cleaner.tds, ts)
					cleaner.table_words += ts.words
				case "th":
					ts := new_boilerpipe_score_omit_table(child, true, true)
					cleaner.tds = append(cleaner.tds, ts)
					cleaner.table_words += ts.words
				case "option":
					child.Data = "a"
				case "img":
					cleaner.imgs++
					if node_is_ownered_by_a(child) {
						cleaner.link_imgs++
					}
					//					log.Println(get_attribute(child, "src"))
					drop = trim_small_image(child)
					if !drop {
						drop = trim_invisible_image(child)
					}
					if drop {
						*dropping = append(*dropping, child)
					}
				case "a":
					cleaner.links++
					cleaner.fix_a_href(child)
				case "li":
					cleaner.lis++
					trim_display_none(child)
				default:
					/* 有些菜单使用了这个属性，如果直接去除，菜单头会被保留下来*/
					trim_display_none(child)
				}
			}
			if !drop {
				cleaner.clean_unprintable_element(dropping, child)
			}
		} else if child.Type == html.TextNode {
			child.Data = merge_tail_spaces(child.Data)
			l := new_boilerpipe_score(child).words
			cleaner.text_words += l
			if node_is_ownered_by_a(child) {
				cleaner.anchor_words += l
			}
		}
	}

	return
}

func (this *html_cleaner) try_update_article(candi *html.Node) {
	if candi == nil {
		return
	}
	sc := new_boilerpipe_score(candi)
	per := sc.words * 100 / (this.text_words + 1)
	if sc.words < wordwrap || per < w_current_line_l {
		return
	}
	this.article = candi
}

const (
	small_image_t = 190 // pixels
)

func trim_small_image(img *html.Node) (drop bool) {
	width, height := media_get_dim(img)

	if img.Parent == nil {
		return
	}
	if width > 0 && height > 0 && width*height < small_image_t*small_image_t && img.Parent.Data == "a" {
		img.Data = "input"
		drop = true
	} else if width == 1 && height == 1 {
		img.Data = "input"
		drop = true
	}
	return
}

func trim_invisible_image(img *html.Node) (drop bool) {
	width, werr := strconv.ParseInt(node_get_attribute(img, "width"), 0, 32)
	height, herr := strconv.ParseInt(node_get_attribute(img, "height"), 0, 32)

	if werr != nil || herr != nil || img.Parent == nil {
		return
	}
	// set width height explicit zero
	if width == 0 || height == 0 {
		img.Data = "input"
		drop = true
	}
	return
}
func remove_children(a *html.Node) {
	for a.FirstChild != nil {
		a.RemoveChild(a.FirstChild)
	}
}

func trim_display_none(n *html.Node) {
	st := node_get_attribute(n, "style")
	if strings.Contains(st, "display") && (strings.Contains(st, "none")) {
		//		log.Println("hide-node display:none", n.Data)
		n.Data = "input"
	}
}

// reserve id, class, href, src, width, height, alt
func (this *html_cleaner) clean_attributes(n *html.Node) {
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		this.clean_attributes(child)
	}
	attrs := []html.Attribute{}
	for _, attr := range n.Attr {
		if attr.Key == "id" || attr.Key == "class" || attr.Key == "href" || attr.Key == "src" || attr.Key == "width" || attr.Key == "height" || attr.Key == "alt" {
			attrs = append(attrs, attr)
		}
	}
	if len(attrs) != len(n.Attr) {
		n.Attr = attrs
	}
}

// clean-body wraps text-node with p
func (this *html_cleaner) clean_body() {
	this.clean_block_node(this.article)
}

//整理html文档，将block-level/inline-level混合的节点改成只有block-level的节点
//对已只有inline-level的节点，删除行前后的空白符
//将包含inline-level的节点展开成更为简单的形式，去掉想<font><span><strong>等等格式节点
func (this *html_cleaner) clean_block_node(n *html.Node) {
	blks := node_has_block_children(n)
	inlines := node_has_inline_children(n)

	// has bocks and inlines
	if blks && inlines {
		child := n.FirstChild
		for child != nil {
			if node_is_inline(child) {
				p := child.PrevSibling
				if p == nil || p.Data != "p" {
					p = create_element("p")
					n.InsertBefore(p, child)
				}
				n.RemoveChild(child)
				p.AppendChild(child)
				child = p.NextSibling
			} else {
				child = child.NextSibling
			}
		}
		inlines = false
	}

	// only inlines
	if blks == false && inlines {
		this.clean_inline_node(n)
		this.trim_empty_spaces(n)
	}

	// only blocks
	if blks && !inlines {
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			this.clean_block_node(child)
		}
	}
}

// flatten inlines text image a object video audio seq
// n is element-node
// inline node may have div element
func (this *html_cleaner) clean_inline_node(n *html.Node) {
	inlines := this.flatten_inline_node(n)

	for child := n.FirstChild; child != nil; child = n.FirstChild {
		n.RemoveChild(child)
	}
	for _, inline := range inlines {
		p := inline.Parent
		if p != nil {
			p.RemoveChild(inline) //			this.article.RemoveChild(child)

		}
		n.AppendChild(inline)
	}
}

//img video audio object embed保留原内容
//text-node保持原内容
//如果inline-level节点包含table/div/ul/ol等等block-level的节点，将这些节点保留
//其他inline-level的节点都直接使用text-node代替
func (this *html_cleaner) flatten_inline_node(n *html.Node) []*html.Node {
	inlines := []*html.Node{}
	for i := n.FirstChild; i != nil; i = i.NextSibling {
		switch {
		case i.Type == html.TextNode:
			fallthrough
		case i.Data == "img":
			fallthrough
		case i.Data == "object" || i.Data == "video" || i.Data == "audio" || i.Data == "embed":
			inlines = append(inlines, i)
		case node_is_block(i) == true:
			fallthrough
		case i.Type == html.ElementNode && i.Data == "a":
			this.clean_inline_node(i)
			inlines = append(inlines, i)
		case i.Type == html.ElementNode:
			x := this.flatten_inline_node(i)
			t := make([]*html.Node, len(inlines)+len(x))
			copy(t, inlines)
			copy(t[len(inlines):], x)
			inlines = t
		}
	}
	return inlines
}

/*
func (this *html_cleaner) clean_form() {
	if this.forms == nil || len(this.forms) == 0 {
		return
	}
	for _, form := range this.forms {
		form.Parent.RemoveChild(form)
	}
}
*/

//节点中没有可显示内容，也没有form等等后续需要处理的节点就是空节点
func (this *html_cleaner) clean_empty_nodes(n *html.Node) {
	child := n.FirstChild
	for child != nil {
		next := child.NextSibling
		this.clean_empty_nodes(child)
		child = next
	}

	if !node_is_not_empty(n) && n.Parent != nil {
		parent := n.Parent
		parent.RemoveChild(n)
	}
}

//删除行前后空白
func (this *html_cleaner) trim_empty_spaces_func(n *html.Node, trim func(string) string) {
	child := n.FirstChild
	for child != nil {
		if child.Type == html.TextNode {
			child.Data = trim(child.Data)
		} else {
			this.trim_empty_spaces_func(child, trim)
		}
		if node_is_not_empty(child) {
			break
		}
		next := child.NextSibling
		n.RemoveChild(child)
		child = next
	}
}

func (this *html_cleaner) trim_empty_spaces(n *html.Node) {
	this.trim_empty_spaces_func(n, func(o string) string {
		return strings.TrimLeftFunc(o, unicode.IsSpace)
	})

	this.trim_empty_spaces_func(n, func(o string) string {
		return strings.TrimRightFunc(o, unicode.IsSpace)
	})

}

const (
	link_img_as_words_c = 4
)

func (this *html_cleaner) link_density() int {
	switch {
	case this.text_words == 0 && this.links == 0:
		return 0
	case this.text_words == 0 && this.links > 0:
		return 100
	default:
		return (this.anchor_words + this.link_imgs*link_img_as_words_c) * 100 / (this.text_words + this.link_imgs*link_img_as_words_c)
	}
}

func (this *html_cleaner) String() string {
	return fmt.Sprint("cleaner links:", this.links,
		", texts:", this.text_words,
		", article:", this.article.Data,
		", linkd:", this.link_density(),
		", tables:", len(this.tables),
		", imgs:", this.imgs,
		", linkimgs:", this.link_imgs,
		", uls:", len(this.uls),
		", ols:", len(this.ols),
		", lis:", this.lis,
		", forms:", len(this.forms),
		", h1:", len(this.header1s),
		", h2:", len(this.header2s),
		", h3:", len(this.header3s))
}

/*
func new_html_cleaner(u string) *html_cleaner {
	rtn := &html_cleaner{}
	rtn.current_url, _ = url.Parse(u)
	return rtn
}
*/
func (cleaner *html_cleaner) grab_keywords(meta *html.Node) {
}

func (cleaner *html_cleaner) grab_description(meta *html.Node) {
}

func (this *html_cleaner) fix_forms() {
	if len(this.forms) == 0 {
		return
	}
	for _, form := range this.forms {
		score := new_boilerpipe_score_omit_table(form, false, false)
		pcnt := score.words * 100 / (1 + this.text_words)
		if pcnt > 33 {
			form.Data = "div"
		}
		//		log.Println("fix form", pcnt, form)
	}
}

func (this *html_cleaner) fix_a_href(a *html.Node) {
	href := node_get_attribute(a, "href")
	uri, err := url.Parse(href)
	if err != nil {
		return
	}
	if this.current_url == nil {
		return
	}
	abs := this.current_url.ResolveReference(uri)
	node_update_attribute(a, "href", abs.String())
}

//return local_filepath, words, images
func clean_fragment(cont, uri string) (string, *DocumentSummary) {
	doc, err := html.Parse(strings.NewReader(cont))
	if err != nil {
		return cont, &DocumentSummary{}
	}

	article := html_clean_root(doc, uri)
	_, body := flat_html(article)
	body.Data = "div" // remvoe body

	var buf bytes.Buffer
	err = html.Render(&buf, body)
	return buf.String(), new_docsummary(body)
}
