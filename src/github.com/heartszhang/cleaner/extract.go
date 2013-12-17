// cleaner export functions
package cleaner

import (
	"code.google.com/p/go.net/html"
	"io/ioutil"
	//	"log"
)

type Extractor interface {
	MakeFragmentReadable(node *html.Node) (*html.Node, *DocumentSummary, error)
	CleanFragment(doc *html.Node) (*html.Node, *DocumentSummary, error)
	MakeHtmlReadable(doc *html.Node, url string) (*html.Node, *DocumentSummary, error)
}

type html_extractor struct {
	temp_dir string
}

func NewExtractor(tmpdir string) Extractor {
	return html_extractor{temp_dir: tmpdir}
}

//doc可能没有html/body父节点
func (this html_extractor) MakeFragmentReadable(doc *html.Node) (*html.Node, *DocumentSummary, error) {
	//清理确定无疑的非正文内容
	article := html_clean_fragment(doc)
	return this.make_article_readable(article)
}
func (this html_extractor) make_article_readable(article *html.Node) (*html.Node, *DocumentSummary, error) {
	//查找文档正文节点，并将其平面化
	doc1, article := readabilitier_make_readable(article)
	write_file(doc1, this.temp_dir)
	//	log.Println("make-readable", of, err)

	// 去除正文中的广告群
	article, images := boiler_clean_by_link_density(article)
	write_file(doc1, this.temp_dir)
	//log.Println("clean-by-density", of, err)

	// 对于以table为主的论坛页面，取出其中的正文table节点
	article = boiler_clean_form_prefix(article)
	write_file(doc1, this.temp_dir)
	//	log.Println("clean-form", of, err)
	return article, new_docsummary(doc1, images), nil

}
func (this html_extractor) CleanFragment(doc *html.Node) (*html.Node, *DocumentSummary, error) {
	article := html_clean_fragment(doc)
	doc1, article := readabilitier_make_readable(article)
	return article, new_docsummary(doc1, nil), nil
}

// cleaned html
// return filepath, *SummaryScore, error
func (this html_extractor) MakeHtmlReadable(doc *html.Node, url string) (*html.Node, *DocumentSummary, error) {
	article := html_clean_root(doc, url)
	write_file(doc, this.temp_dir)
	//	log.Println("1-step", n)
	return this.make_article_readable(article)
}

func write_file(doc *html.Node, temp string) (string, error) {
	of, err := ioutil.TempFile(temp, "html.")
	if err != nil {
		return "", err
	}
	defer of.Close()

	html.Render(of, doc)
	return of.Name(), nil
}
