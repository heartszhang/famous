package cleaner

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"log"
	"strings"
	"testing"
)

func TestCleaner(t *testing.T) {
	root := &html.Node{Type: html.ElementNode, Data: "article", DataAtom: atom.Article}
	frag := `<img src="http://pic.yupoo.com/fotomag/DgLAr6tN/qdAAr.jpg" class="insertimage" alt="" title="" border="0" align="right" />“<a href="http://www.leica.org.cn/Julia_Fullerton_Batton_Tokyo/" target="_blank">Tokyo 2013</a>”是女摄影师 Julia Fullerton-Batton 完成于2013年的一个拍摄项目，在这个项目中，摄影师通过一组拍摄于东京街头的女性肖像，来揭示当代日本社会中女性地位与角色的转变。<br/><br/>“在日本传统文化中，女性一直是在家照料孩子和丈夫的传统角色。而在现在，日本女性的角色已经重新定位，越来越多女性拒绝成为‘领工资丈夫’的附庸，选择不结婚，进入职场的生活方式。事实上，日本已经成为全球出生率最低的国家之一。<br/><br/>传统日本女性往往通过婚姻和家庭来传递自己的力量，现在她们有自己的职业和新潮喜好，并愈发改变着日本城市的景观，这组照片就是记录下其中最戏剧化的片段，将其呈现给人们。”<br/><br/><strong>关于摄影师</strong><br/>Julia Fullerton-Batton，英国女摄影师，擅长拍摄 Fine-Art 风格的环境肖像作品，我们曾经介绍过她的多个拍摄项目。<br/><br/><p align="center"><img src="http://pic.yupoo.com/fotomag/DgLAsmSv/q6uP2.jpg" class="insertimage" alt="" title="" border="0" /></p><br/><p align="center"><img src="http://pic.yupoo.com/fotomag/DgLAtf1n/srGWJ.jpg" class="insertimage" alt="" title="" border="0" /></p><br/><p align="center"><img src="http://pic.yupoo.com/fotomag/DgLAu4E7/GStkd.jpg" class="insertimage" alt="" title="" border="0" /></p><br/><p align="center"><img src="http://pic.yupoo.com/fotomag/DgLAvuFg/NxZZt.jpg" class="insertimage" alt="" title="" border="0" /></p><br/><p align="center"><img src="http://pic.yupoo.com/fotomag/DgLAvAxg/sd3S9.jpg" class="insertimage" alt="" title="" border="0" /></p><br/><p align="center"><img src="http://pic.yupoo.com/fotomag/DgLAwcUS/8GVNA.jpg" class="insertimage" alt="" title="" border="0" /></p><br/><p align="center"><img src="http://pic.yupoo.com/fotomag/DgLAvQZc/8OEkA.jpg" class="insertimage" alt="" title="" border="0" /></p><br/><p align="center"><img src="http://pic.yupoo.com/fotomag/DgLAyxqX/wmpNn.jpg" class="insertimage" alt="" title="" border="0" /></p><br/><br/>『<a href="http://www.leica.org.cn/" target="_blank">Leica中文摄影杂志</a>』推荐<a href="http://feedburner.google.com/fb/a/mailverify?uri=leicachina" target="_blank">使用Email的方式订阅</a>；在Apple Mac OS X下可获得最佳阅读体验<br/><img src="http://pic.yupoo.com/fotomag/Bl4vOcGk/F8RzA.jpg" class="insertimage" alt="" title="" border="0" align="right" /><br/>『<a href="http://iphoto.ly/" target="_blank">iPhoto.ly</a>』在苹果上阅读：<a href="http://itunes.apple.com/app/id373067909?mt=8" target="_blank">iPhone版</a>+<a href="http://itunes.apple.com/app/id460728960?mt=8" target="_blank">iPad版</a>，^_^<br/><br/><strong>Tips：</strong> <strong><span style="color: #4169E1;">关注我们： <a href="https://twitter.com/leica" target="_blank">Twitter</a>、<a href="http://fanfou.com/fotomag" target="_blank">饭否</a>、<a href="http://weibo.com/iphotoly" target="_blank">微博</a></span></strong><br/><br/><span style="font-size: 12px;"><span style="color: #FFA500;"><strong>『小建议』</strong></span>如果你在Email里看到这篇文章，可以<strong>转发</strong>给你的朋友；如果你在豆瓣里看到这篇文章，不妨<strong>推荐</strong>给更多人；或者干脆Copy下这篇文章的链接，发给你MSN上最喜欢的人；<strong>我们永远相信，分享是一种美德，Great People Share Knowledge</strong>...</span><br/>Tags - <a href="http://www.leica.org.cn/tags/%25E4%25BB%2596%25E4%25BB%25AC%25E5%259C%25A8%25E6%258B%258D%25E4%25BB%2580%25E4%25B9%2588/" rel="tag">他们在拍什么</a> , <a href="http://www.leica.org.cn/tags/%25E5%25A5%25B3%25E6%2591%2584%25E5%25BD%25B1%25E5%25B8%2588/" rel="tag">女摄影师</a>
<br />
<div>
<strong>推荐阅读</strong>
<ul>
	<li><a href="http://feedburner.google.com/fb/a/mailverify?uri=leicachina"><strong>『特别推荐』使用Email订阅Leica
中文摄影杂志</strong></a></li>
<li><a href="http://www.leica.org.cn/post/774/">『Black&White』12张值得推荐的黑白摄影照片</a></li>
</ul>
</div>
<div><img src="http://img.tongji.linezing.com/712108/tongji.gif" board="0" /></div><img src="http://www1.feedsky.com/t1/731643689/leica/feedsky/s.gif?r=http://www.leica.org.cn/Julia_Fullerton_Batton_Tokyo/" border="0" height="0" width="0" style="position:absolute" />`
	nodes, _ := html.ParseFragment(strings.NewReader(frag), root)
	for _, n := range nodes {
		root.AppendChild(n)
	}
	article, sum, _ := MakeFragmentReadable(root)
	log.Println(sum)
	print_html_doc(article)
	//	WriteHtmlFile(root, "../3.html.txt")
}

func print_html_doc(node *html.Node) {
	var buffer bytes.Buffer
	html.Render(&buffer, node) // ignore return error
	log.Println(buffer.String())
}

func trim(d string) string {
	if len(d) < 10 {
		return d
	}
	return d[:10]
}
