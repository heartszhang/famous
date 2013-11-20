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
	frag := `原文：<a href="http://www.coindesk.com/bitcoin-iranian-shoe-store-trade-sanctions/" target="_blank">Bitcoin helps Iranian shoe store overcome international trade sanctions</a><br />作者：Jon Southurst<br />时间：2013-11-06<br />本文由Fish翻译。<br /><div class="separator" style="clear: both; text-align: center;"><a href="http://persianshoes.com/wp-content/uploads/2013/10/S1-940x500.jpg" imageanchor="1" style="margin-left: 1em; margin-right: 1em;"><img border="0" src="http://persianshoes.com/wp-content/uploads/2013/10/S1-940x500.jpg" height="212" width="400" /></a></div><br />上周，有一个有趣的电子商务网站宣布上线，它只接受比特币（Bitcoin)，因为大多数海外客户无法支付别的货币。<br /><br />该公司名叫<a href="http://persianshoes.com/" target="_blank"> Persian Shoes</a> （波斯鞋业），有超过70年的历史， 主营手工制作的皮鞋，它坐落在伊朗第三大城市伊斯法罕。<br /><br />业主很高兴运货到世界各地，但如何支付是一个问题。由于变幻莫测的国际政治和过去几十年的历史，普通的电子商务渠道（在伊朗）被封锁。<br /><br />联合国，美国，欧盟和其它国家对伊朗全境实施贸易制裁，这意味着西联和主要的信用卡公司都不处理与伊朗相关的业务，即便仅仅是服装时尚业的买卖。<br /><br />在伊朗付款给某人的唯一方法是口袋里携带现金，或者使用一些容易转移和大多不受监管的数字货币。<br /><br />在第一天的销售中，这家仅接受比特币的商店卖出了四双皮鞋。 “以我的标准来说，这是一个很好的销售业绩！”首席执行官罗卡尼说。上线第一个星期，这家网上商店运作很成功。<br /><br />“使用比特币，到目前为止我们已经卖出10双皮鞋了，这远高于我们的预期。“他补充说。<br /><br />Persian Shoes 目前由三兄弟经营，他们准备大力扩张其父亲开创的这一生意。 罗卡尼的表弟居住在澳大利亚，向三兄弟介绍了比特币，并帮助他们建立了网站。<br /><br />该网站目前提供女士的皮革手袋，钱包和鞋子，还有七个品种的男士皮鞋。货品用美元标价，最低从80美元起。<br /><br />网站的FAQ页面上介绍了比特币，并指导新用户使用诸如Coinbase、Bitstamp，BitBargain，以及LocalBitcoins等比特币相关的线上服务。 它还有业务简介：<br /><br />“<i>我们的业务是生产和销售皮革制品。我们希望把我们的产品卖到全世界，客户越多越好。但问题是我们在伊朗经营，大多数的支付系统都不愿意为我们提供完整的服务，或者对我们的业务有巨大的风险。在推出这个网站之前，我们的国际销售一直局限于几个专门的客户，他们了解我们的产品质量，并为支付时可能会出现很多的麻烦做好了准备！当然，对于使用各种电子支付的人来说这听起来可能难以理解。然而，在出现比特币之前，如何收款已成为了我们扩大业务的头号障碍。</i>”<br /><br />除了国际贸易的限制，把比特币兑换为当地的货币（伊朗里亚尔）也面临着知识和技术的障碍。<br /><br />“兑换时有点棘手，因为比特币并不是常用的。我们保留一些比特币在手上，同时在localbitcoins上卖出一些”，罗卡尼说，缺乏网上购物的本地文化和互联网连接不畅阻碍了该公司与伊朗的其它地区的贸易。<br /><br />根据你的国籍或居住地，即使使用比特币，你购买这些鞋子也可能不被允许。正如美国人无法在哈瓦那度过一个周末或者享受真正的古巴雪茄，他们也被禁止与伊朗境内的企业或个人从事任何贸易活动。<br /><br />尽管这样， Persian Shoes 已经从美国客户那里接了几个订单，并且正等着看看是否顺利交付，以便进一步做广告推广业务。<br /><br />就个人以电子商务形式购买衣物问题，向美国海关和边境保护局（国土安全部的一部分）询问，得到的回应是：“不幸的是，对伊朗的制裁是禁止这类交易的。”<br /><br />其它国家对此类买卖的态度不太明朗，或者没有回复。<br /><br />大多数国际贸易制裁是关于技术或工业设备的，有关个人日常用品的交易（可能是由于缺乏支付选项）没有具体的说明，当然肯定不会提到皮鞋。<br /><br />虽然如此，你总还可以自由地浏览Persian Shoes的Facebook页面和网上商店。<br /><div><br /></div><div><div align="left" class="MsoNormal" style="background-color: white; color: #333333; font-family: Arial, serif; font-size: 14px; line-height: 22.875px; text-indent: 24pt;"><i style="line-height: 22px;">本文版权属于原出版公司及作者所有。</i></div><div align="left" class="MsoNormal" style="background-color: white; color: #333333; font-family: Arial, serif; font-size: 14px; line-height: 22.875px; text-indent: 24pt;"><i style="line-height: 22px;">©译者遵守知识共享署名-非商业性使用-相同方式共享 3.0许可协议。</i></div></div><div class="blogger-post-footer"><i>译文遵循<a href="http://creativecommons.org/licenses/by-nc-sa/3.0/deed.zh">CC3.0</a>版权标准。转载务必标明链接和“转自译者”。不得用于商业目的。点击<a href="http://us5.campaign-archive2.com/home?u=b4608e11fae8d90d93807e499&id=15355fc869">这里</a>查看和订阅《每日译者》手机报。<a href="http://www.my1510.cn/article.php?id=77401">穿墙查看</a>译者博客、书刊、音频和视频</i></div>`
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
