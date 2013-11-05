package cleaner

import (
	"code.google.com/p/go.net/html"
	"io/ioutil"
)

func MakeFragmentReadable(doc *html.Node) (*html.Node, *DocSummary, error) {
	article := html_clean_fragment(doc)
	//	of, err := write_file(doc)

	doc1, article := readabilitier_make_readable(article)
	//	of, err = write_file(doc1)

	article = boiler_clean_by_link_density(article)
	//	of, err = write_file(doc1)

	article = boiler_clean_form_prefix(article)
	//	of, err = write_file(doc1)
	return article, new_docsummary(doc1), nil
}

func CleanFragment(doc *html.Node) (*html.Node, *DocSummary, error) {
	article := html_clean_fragment(doc)
	doc1, article := readabilitier_make_readable(article)
	return article, new_docsummary(doc1), nil
}

// cleaned html
// return filepath, *SummaryScore, error
func MakeHtmlReadable(doc *html.Node, url string) (*html.Node, *DocSummary, error) {
	article := html_clean_root(doc, url)

	//	s2, _ := WriteHtmlFile2(doc)
	doc1, article := readabilitier_make_readable(article)

	//	s2, _ = WriteHtmlFile2(doc1)
	article = boiler_clean_by_link_density(article)

	//	h4ml, _ := WriteHtmlFile2(doc1)
	article = boiler_clean_form_prefix(article)
	//	h5ml, err := WriteHtmlFile2(doc1)
	return article, new_docsummary(doc1), nil
}

func write_file(doc *html.Node) (string, error) {
	of, err := ioutil.TempFile("", "html.")
	if err != nil {
		return "", err
	}
	defer of.Close()

	html.Render(of, doc)
	return of.Name(), nil
}
