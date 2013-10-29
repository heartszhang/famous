package cleaner

import (
	"code.google.com/p/go.net/html"
)

func MakeFragmentReadable(doc *html.Node) (*html.Node, *DocSummary, error) {
	article := html_clean_fragment(doc)
	doc1, article := readabilitier_make_readable(article)
	return article, new_docsummary(doc1), nil
}

// cleaned html
// return filepath, *SummaryScore, error
func MakeHtmlReadable(doc *html.Node, url string) (*html.Node, *DocSummary, error) {
	/*	htmlfile, _, err := DownloadHtml(url)
		if err != nil {
			return "", &DocSummary{}, err
		}

		f, err := os.Open(htmlfile)
		if err != nil {
			return "", &DocSummary{}, err
		}
		defer f.Close()

		reader := bufio.NewReader(f)
		doc, err := html.Parse(reader)
		if err != nil {
			return "", &DocSummary{}, err
		}
	*/
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
