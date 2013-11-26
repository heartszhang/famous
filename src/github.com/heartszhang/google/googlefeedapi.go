package google

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"fmt"
	"github.com/heartszhang/curl"
	"github.com/heartszhang/feedfeed"
	"github.com/heartszhang/oauth2"
	"net/http"
	"strings"
)

type googlefeed_error struct {
	code   int
	reason string
}

func (this googlefeed_error) Error() string {
	return fmt.Sprintf("%v: %v", this.code, this.reason)
}

type find_result struct {
	ResponseDetails string `json:"responseDetails,omitempty"`
	ResponseStatus  int    `json:"responseStatus"`
	ResponseData *struct {
		Query   string `json:"query,omitempty"`
		Entries []struct {
			Url            string `json:"url,omitempty"`
			Title          string `json:"title,omitempty"`
			ContentSnippet string `json:"contentSnippet"`
			Website        string `json:"link,omitempty"`
		} `json:"entries,omitempty"`
	} `json:"responseData,omitempty"`
}

func findresult_to_findentry(x find_result) ([]feedfeed.FeedSourceFindEntry, error) {
	if x.ResponseStatus != 200 {
		return nil, googlefeed_error{x.ResponseStatus, x.ResponseDetails}
	}
	v := make([]feedfeed.FeedSourceFindEntry, 0)
	if x.ResponseData == nil {
		return v, nil
	}
	for _, e := range x.ResponseData.Entries {
		v = append(v, feedfeed.FeedSourceFindEntry{e.Url,
			strip_html_tags(e.Title),
			strip_html_tags(e.ContentSnippet),
			e.Website, false})
	}
	return v, nil
}

type load_result struct {
	ResponseDetails string `json:"responseDetails,omitempty"`
	ResponseStatus  int    `json:"responseStatus"`

	ResponseData *struct {
		Feed *struct {
			FeedUrl     string `json:"feedUrl,omitempty"`
			Title       string `json:"title,omitempty"`
			Website     string `json:"link,omitempty"`
			Author      string `json:"author,omitempty"`
			Description string `json:"description,omitempty"`
			Type        string `json:"type,omitempty"`
			Entries     []struct {
				MediaGroups []struct {
					Contents []struct {
						Url         string   `json:"url,omitempty"`
						Type        string   `json:"type,omitempty"`
						Medium      string   `json:"medium,omitempty"`
						Height      int      `json:"height"`
						Width       int      `json:"width"`
						IsDefault   *bool    `json:"isDefault,omitempty"`
						Title       string   `json:"title,omitempty"`
						Description string   `json:"description,omitempty"`
						Keywords    string   `json:"keywords,omitempty"`
						Categories  []string `json:"category,omitempty"`
						Thumbnails  []struct {
							Url    string `json:"url,omitempty"`
							Height int    `json:"height"`
							Width  int    `json:"width"`
							Time   string `json:"time,omitempty"`
						} `json:"thumbnail,omitempty"`
					} `json:"contents,omitempty"`
				} `json:"mediaGroups,omitempty"`
				Title          string   `json:"title,omitempty"`
				Link           string   `json:"link,omitempty"`
				Author         string   `json:"author,omitempty"`
				PublishedDate  string   `json:"publishedDate"`
				ContentSnippet string   `json:"contentSnippet,omitempty"`
				Content        string   `json:"content,omitempty"`
				Categories     []string `json:"categories,omitempty"`
			} `json:"entries,omitempty"`
		} `json:"feed,omitempty"`
	} `json:"responseData,omitempty"`
}

func media_type(mime string) string {
	s := strings.Split(mime, "/")
	return s[0]
}
func loadresult_to_feedsource(x load_result) (*feedfeed.FeedSource, []feedfeed.FeedEntry, error) {
	if x.ResponseStatus != 200 || x.ResponseData == nil || x.ResponseData.Feed == nil {
		return nil, nil, googlefeed_error{x.ResponseStatus, x.ResponseDetails}
	}

	v := make([]feedfeed.FeedEntry, 0)
	f := x.ResponseData.Feed
	s := feedfeed.FeedSource{
		Name:        f.Title,
		Uri:         f.FeedUrl,
		WebSite:     f.Type + f.Website,
		Description: f.Description,
	}
	for _, e := range f.Entries {
		ne := feedfeed.FeedEntry{
			Title:   feedfeed.FeedTitle{Main: e.Title},
			Uri:     e.Link,
			Summary: e.ContentSnippet,
			Content: e.Content,
			Tags:    e.Categories,
		}
		for _, m := range e.MediaGroups {
			for _, c := range m.Contents {
				fm := feedfeed.FeedMedia{
					Title:       c.Title,
					Description: c.Description,
					Width:       c.Width, Height: c.Height,
					Mime: c.Type,
					Uri:  c.Url,
				}
				switch media_type(c.Type) {
				case "image":
					ne.Images = append(ne.Images, fm)
				case "vidoe":
					ne.Videos = append(ne.Videos, fm)
				case "audio":
					ne.Audios = append(ne.Audios, fm)
				}
				for _, t := range c.Thumbnails {
					fmt := feedfeed.FeedMedia{
						Width: t.Width, Height: t.Height, Uri: t.Url,
					}
					ne.Images = append(ne.Images, fmt)
				}
			}
		}
		v = append(v, ne)
	}
	return &s, v, nil
}
type GoogleFeedApiService interface {
	Find(q, hl string) ([]feedfeed.FeedSourceFindEntry, error)
	Load(uri, hl string, num int, scoring bool) (*feedfeed.FeedSource, []feedfeed.FeedEntry, error)
}

type google_feedapi struct {
	temp_folder string
	refer       string
}

func NewGoogleFeedApi(refer, tmp string) GoogleFeedApiService {
	if refer == "" {
		refer = "https://heartszhang.github.com/google"
	}
	return &google_feedapi{temp_folder: tmp, refer: refer}
}

func (this google_feedapi) RoundTrip(r *http.Request) {
	r.Header.Set("refer", this.refer)
}

func (this google_feedapi) Find(q, hl string) ([]feedfeed.FeedSourceFindEntry, error) {
	p := struct {
		q  string `param:"q"`
		hl string `param:"hl"` // default en, nil means en
		v  string `param:"v"`
	}{q: q, hl: hl, v: "1.0"}
	c := curl.NewCurlerDetail(this.temp_folder, curl.CurlProxyPolicyAlwayseProxy, 0, this)
	v := find_result{}
	uri := find_url + "?" + oauth2.HttpQueryEncode(p)
	err := c.GetAsJson(uri, &v)
	if err != nil {
		return nil, err
	}
	x, err := findresult_to_findentry(v)
	return x, err
}
func make_num(num int) *int {
	if num <= 0 {
		return nil
	}
	return &num
}
func make_scoring(scoring bool) string {
	if scoring {
		return "h"
	}
	return ""
}
func (this google_feedapi) Load(uri, hl string, num int, scoring bool) (*feedfeed.FeedSource, []feedfeed.FeedEntry, error) {
	p := struct {
		q       string `param:"q"`
		hl      string `param:"hl"`      // default en
		num     *int   `param:"num"`     // default 4
		scoring string `param:"scoring"` // nil or h
		v       string `param:"v"`
	}{q: uri, hl: hl, num: make_num(num), scoring: make_scoring(scoring), v: "1.0"}
	c := curl.NewCurlerDetail(this.temp_folder, curl.CurlProxyPolicyAlwayseProxy, 0, this)
	v := load_result{}
	url := load_url + "?" + oauth2.HttpQueryEncode(p)
	err := c.GetAsJson(url, &v)

	if err != nil {
		return nil, nil, err
	}
	s, e, err := loadresult_to_feedsource(v)
	return s, e, err
}

const (
	feed_service = `https://ajax.googleapis.com/ajax/services/feed/`
	find_url     = feed_service + `find`
	load_url     = feed_service + `load`
)

type googlefeedapi_find_param struct {
	q  string `param:"q"`
	hl string `param:"hl"` // default en, nil means en
	v  string `param:"v"`
}

func (this googlefeedapi_find_param) build() string {
	return oauth2.HttpQueryEncode(this)
}

type googlefeedapi_load_param struct {
	q       string
	hl      string // default en
	num     *int   // default 4
	scoring *bool  // nil or h
	v       string
}

func (this googlefeedapi_load_param) build() string {
	return oauth2.HttpQueryEncode(this)
}

func strip_html_tags(htm string) string {
	reader := strings.NewReader(htm)
	root := &html.Node{Type: html.ElementNode, Data: "article", DataAtom: atom.Article}
	frags, _ := html.ParseFragment(reader, root)
	var txt string
	for _, f := range frags {
		txt += extract_html_text(f)
	}
	return txt
}

func extract_html_text(node *html.Node) string {
	if node.Type == html.TextNode {
		return strings.Replace(node.Data, "\n", "", -1)
	}
	if node.Type != html.ElementNode {
		return ""
	}
	var v string
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		v += extract_html_text(child)
	}
	return v
}
