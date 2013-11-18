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
	frag := `<p>Beau每天和7周大的小狗一起睡午觉，他妈妈拍下了这些超可爱的照片。</p> <p><a href="http://jandan.net/2013/11/18/shyba-puppy.html"><img src="http://ww1.sinaimg.cn/mw600/66b3de17gw1eap0rxzl61j20dw0dw781.jpg" alt="太有爱：每天和小孩一起睡午觉的小狗" /></a></p> <p>一可爱幼童和一只7周大的的小狗被晒午觉照片，照片中的他们如“兄弟”。当妈妈Jessica Shyba 在网上发布了这两位的照片后，这一对已成为新进网络红人了。</p> <p>当这家人在当地动物收容所后院看到了这只杂种狗Theo 时，Jessica 就决定驯养它，让它睡在自家的狗窝里。但是她再也受不了它在夜里的叫声，最后让小狗和她一起睡在床上。</p> <p>一天下午，当Jessica 哄着Beau 睡午觉时，Theo跳上来，然后睡着了。Jessica 说：“我那时几乎是大叫起来，差不多要把他们两个都叫醒了。”</p> <p>从此以后，Theo每天和Beau一起睡午觉，彼此相互缠着至少要睡上两个钟头。 Jessica就开始用相机拍下这些可爱镜头，并把它们发布在网络的图片分享上。</p> <p><a href="http://jandan.net/2013/11/18/shyba-puppy.html"><img src="http://ww3.sinaimg.cn/mw600/66b3de17gw1eap0s079euj20dw0ch78h.jpg" alt="太有爱：每天和小孩一起睡午觉的小狗" /></a></p> <p><a href="http://jandan.net/2013/11/18/shyba-puppy.html"><img src="http://ww3.sinaimg.cn/mw600/66b3de17gw1eap0ryw25tj20dw0auq4n.jpg" alt="太有爱：每天和小孩一起睡午觉的小狗" /></a></p> <p><a href="http://jandan.net/2013/11/18/shyba-puppy.html"><img src="http://ww4.sinaimg.cn/mw600/66b3de17gw1eap0rv3t3oj20dw0dw78e.jpg" alt="太有爱：每天和小孩一起睡午觉的小狗" /></a></p> <p><a href="http://jandan.net/2013/11/18/shyba-puppy.html"><img src="http://ww4.sinaimg.cn/mw600/66b3de17gw1eap0rtxwmsj20dw0acwhi.jpg" alt="太有爱：每天和小孩一起睡午觉的小狗" /></a></p> <p><a href="http://jandan.net/2013/11/18/shyba-puppy.html"><img src="http://ww3.sinaimg.cn/mw600/66b3de17gw1eap0rsp2l1j20dw0dwgpr.jpg" alt="太有爱：每天和小孩一起睡午觉的小狗" /></a></p> <p><a href="http://jandan.net/2013/11/18/shyba-puppy.html"><img src="http://ww1.sinaimg.cn/mw600/66b3de17gw1eap0rr10p2j20dw0b00vd.jpg" alt="太有爱：每天和小孩一起睡午觉的小狗" /></a></p> <p><a href="http://jandan.net/2013/11/18/shyba-puppy.html"><img src="http://ww3.sinaimg.cn/mw600/66b3de17gw1eap0rpor1ij20ba0dwdij.jpg" alt="太有爱：每天和小孩一起睡午觉的小狗" /></a></p> <p>“每天，一到午睡时间Theo就与我们会和，耐心等待Beau睡着。”Jessica在她的博客上如此描述。<br /> “每当那时，它都非常困乏。所以当我把Beau轻轻放在床上时，它就跌跌撞撞地过来，扑通一下倒在Beau身上倒头就睡。”<br /> “这件事让我感到非常愉悦。”</p> <p>这家人认为Theo是只拳师犬、德国牧羊犬和拉布拉多犬的杂交后代。它一出生就连同它的10个兄弟姐妹一起被遗弃了。Jessica希望创作一本影集，把影集收入捐献给圣克鲁兹动物保护协会。</p> <p><a href="http://jandan.net/2013/11/18/shyba-puppy.html"><img src="http://ww3.sinaimg.cn/mw600/66b3de17gw1eap0rwl7cdj20dw0dwq72.jpg" alt="太有爱：每天和小孩一起睡午觉的小狗" /></a></p> <p><a href="http://jandan.net/2013/11/18/shyba-puppy.html"><img src="http://ww3.sinaimg.cn/mw600/66b3de17gw1eap0rohx69j20dw0dwtc8.jpg" alt="太有爱：每天和小孩一起睡午觉的小狗" /></a></p> <p><em>[<a href="http://jandan.net/2013/11/18/shyba-puppy.html">Vera</a> via <a target=_blank rel="external" href="http://www.mirror.co.uk/news/world-news/beau-shyba-puppy-theo-pictures-2802112">Mirror</a>]</em></p> <p><a href="http://jandan.net/2013/11/18/shyba-puppy.html">>>点这里浏览原文<<</a></p><div style="border-top: 1px solid #DCDCDC; padding: 5px 0;"></div><p>&#169; <a target="_blank" href="http://jandan.net/">煎蛋</a> / <a href="http://weibo.com/diggdigest/profile">超载鸡微博</a> / 图片托管于<a target="_blank" href="https://www.upyun.com/?utm_source=jandan2&utm_medium=ad&utm_campaign=upyun&md=jandan2">又拍云存储</a> / <a href="http://jandan.net/2013/08/15/chaozai-tee.html" target="_blank">煎蛋7周年纪念款TEE</a></p><img width='1' height='1' src='http://jandan.feedsportal.com/c/34036/f/617798/s/33cb4a03/sc/38/mf.gif' border='0'/><br clear='all'/><div class='mf-viral'><table border='0'><tr><td valign='middle'><a href="http://share.feedsportal.com/share/twitter/?u=http%3A%2F%2Fjandan.net%2F2013%2F11%2F18%2Fshyba-puppy.html&t=%E5%A4%AA%E6%9C%89%E7%88%B1%EF%BC%9A%E6%AF%8F%E5%A4%A9%E5%92%8C%E5%B0%8F%E5%AD%A9%E4%B8%80%E8%B5%B7%E7%9D%A1%E5%8D%88%E8%A7%89%E7%9A%84%E5%B0%8F%E7%8B%97" target="_blank"><img src="http://res3.feedsportal.com/social/twitter.png" border="0" /></a>&nbsp;<a href="http://share.feedsportal.com/share/facebook/?u=http%3A%2F%2Fjandan.net%2F2013%2F11%2F18%2Fshyba-puppy.html&t=%E5%A4%AA%E6%9C%89%E7%88%B1%EF%BC%9A%E6%AF%8F%E5%A4%A9%E5%92%8C%E5%B0%8F%E5%AD%A9%E4%B8%80%E8%B5%B7%E7%9D%A1%E5%8D%88%E8%A7%89%E7%9A%84%E5%B0%8F%E7%8B%97" target="_blank"><img src="http://res3.feedsportal.com/social/facebook.png" border="0" /></a>&nbsp;<a href="http://share.feedsportal.com/share/linkedin/?u=http%3A%2F%2Fjandan.net%2F2013%2F11%2F18%2Fshyba-puppy.html&t=%E5%A4%AA%E6%9C%89%E7%88%B1%EF%BC%9A%E6%AF%8F%E5%A4%A9%E5%92%8C%E5%B0%8F%E5%AD%A9%E4%B8%80%E8%B5%B7%E7%9D%A1%E5%8D%88%E8%A7%89%E7%9A%84%E5%B0%8F%E7%8B%97" target="_blank"><img src="http://res3.feedsportal.com/social/linkedin.png" border="0" /></a>&nbsp;<a href="http://share.feedsportal.com/share/gplus/?u=http%3A%2F%2Fjandan.net%2F2013%2F11%2F18%2Fshyba-puppy.html&t=%E5%A4%AA%E6%9C%89%E7%88%B1%EF%BC%9A%E6%AF%8F%E5%A4%A9%E5%92%8C%E5%B0%8F%E5%AD%A9%E4%B8%80%E8%B5%B7%E7%9D%A1%E5%8D%88%E8%A7%89%E7%9A%84%E5%B0%8F%E7%8B%97" target="_blank"><img src="http://res3.feedsportal.com/social/googleplus.png" border="0" /></a>&nbsp;<a href="http://share.feedsportal.com/share/email/?u=http%3A%2F%2Fjandan.net%2F2013%2F11%2F18%2Fshyba-puppy.html&t=%E5%A4%AA%E6%9C%89%E7%88%B1%EF%BC%9A%E6%AF%8F%E5%A4%A9%E5%92%8C%E5%B0%8F%E5%AD%A9%E4%B8%80%E8%B5%B7%E7%9D%A1%E5%8D%88%E8%A7%89%E7%9A%84%E5%B0%8F%E7%8B%97" target="_blank"><img src="http://res3.feedsportal.com/social/email.png" border="0" /></a></td><td valign='middle'></td></tr></table></div><br/><br/><a href="http://da.feedsportal.com/r/180264375713/u/49/f/617798/c/34036/s/33cb4a03/a2.htm"><img src="http://da.feedsportal.com/r/180264375713/u/49/f/617798/c/34036/s/33cb4a03/a2.img" border="0"/></a><img width="1" height="1" src="http://pi.feedsportal.com/r/180264375713/u/49/f/617798/c/34036/s/33cb4a03/a2t.img" border="0"/>`
	nodes, _ := html.ParseFragment(strings.NewReader(frag), root)
	for _, n := range nodes {
		root.AppendChild(n)
	}
	article, _, _ := NewExtractor("").MakeFragmentReadable(root)
	//	log.Println(sum)
	print_html_doc(article)
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
