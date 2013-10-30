package cleaner

import (
	"code.google.com/p/go.net/html"
)

func MakeFragmentReadable(doc *html.Node) (*html.Node, *DocSummary, error) {
	article := html_clean_fragment(doc)
	doc1, article := readabilitier_make_readable(article)
	article = boiler_clean_by_link_density(article)
	article = boiler_clean_form_prefix(article)
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
