package cleaner

import (
	"code.google.com/p/go.net/html"
	"io/ioutil"
	"log"
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
func (this html_extractor) MakeFragmentReadable(doc *html.Node) (*html.Node, *DocumentSummary, error) {
	article := html_clean_fragment(doc)
	of, err := write_file(doc, this.temp_dir)
	log.Println("clean-fragment", of, err)

	doc1, article := readabilitier_make_readable(article)
	of, err = write_file(doc1, this.temp_dir)
	log.Println("make-readable", of, err)

	article, images := boiler_clean_by_link_density(article)
	//of, err = write_file(doc1, this.temp_dir)
	//log.Println("clean-by-density", of, err)

	article = boiler_clean_form_prefix(article)
	//	of, err = write_file(doc1, this.temp_dir)
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

	//	s2, _ := WriteHtmlFile2(doc)
	doc1, article := readabilitier_make_readable(article)

	//	s2, _ = WriteHtmlFile2(doc1)
	article, images := boiler_clean_by_link_density(article)

	//	h4ml, _ := WriteHtmlFile2(doc1)
	article = boiler_clean_form_prefix(article)
	//	h5ml, err := WriteHtmlFile2(doc1)
	return article, new_docsummary(doc1, images), nil
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
