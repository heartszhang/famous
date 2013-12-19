package google

import (
	"fmt"
	"net/http"
	"strings"

	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"github.com/heartszhang/curl"
	"github.com/heartszhang/feed"
	"github.com/heartszhang/oauth2"
	"github.com/heartszhang/unixtime"
)

type GoogleFeedApiService interface {
	Find(q, hl string) ([]feed.FeedEntity, error)
	Load(uri, hl string, num int, scoring bool) (feed.FeedEntity, error)
}

type find_result struct {
	ResponseDetails string `json:"responseDetails,omitempty"`
	ResponseStatus  int    `json:"responseStatus"`
	ResponseData    *struct {
		Query   string `json:"query,omitempty"`
		Entries []struct {
			Url            string `json:"url,omitempty"`
			Title          string `json:"title,omitempty"`
			ContentSnippet string `json:"contentSnippet"`
			Website        string `json:"link,omitempty"`
		} `json:"entries,omitempty"`
	} `json:"responseData,omitempty"`
}

const (
	default_update_minutes int64 = 120
)

func (x find_result) to_feedentities() ([]feed.FeedEntity, error) {
	if x.ResponseStatus != http.StatusOK {
		return nil, gf_error{x.ResponseStatus, x.ResponseDetails}
	}
	var v []feed.FeedEntity
	if x.ResponseData == nil {
		return v, nil
	}
	for _, e := range x.ResponseData.Entries {
		v = append(v, feed.FeedEntity{
			FeedSource: feed.FeedSource{
				Uri:         e.Url,
				Period:      default_update_minutes, // minutes
				Name:        strip_html_tags(e.Title),
				Description: strip_html_tags(e.ContentSnippet),
				WebSite:     e.Website,
			},
		})
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
				Title          string        `json:"title,omitempty"`
				Link           string        `json:"link,omitempty"`
				Author         string        `json:"author,omitempty"`
				PublishedDate  unixtime.Time `json:"publishedDate"`
				ContentSnippet string        `json:"contentSnippet,omitempty"`
				Content        string        `json:"content,omitempty"`
				Categories     []string      `json:"categories,omitempty"`
			} `json:"entries,omitempty"`
		} `json:"feed,omitempty"`
	} `json:"responseData,omitempty"`
}

func media_type(mime string) string {
	s := strings.Split(mime, "/")
	return s[0]
}
func (x load_result) to_feedentity() (feed.FeedEntity, error) {
	if x.ResponseStatus != http.StatusOK || x.ResponseData == nil || x.ResponseData.Feed == nil {
		return feed.FeedEntity{}, gf_error{x.ResponseStatus, x.ResponseDetails}
	}

	var v []feed.FeedEntry
	f := x.ResponseData.Feed
	s := feed.FeedSource{
		Name:        strip_html_tags(f.Title),
		Uri:         f.FeedUrl,
		WebSite:     f.Website,
		Description: f.Description,
		Type:        feed.FeedSourceType(f.Type),
		Period:      default_update_minutes, // minutes
	}

	for _, e := range f.Entries {
		ne := feed.FeedEntry{
			Parent:  s.Uri,
			Title:   e.Title,
			Uri:     e.Link,
			Summary: e.ContentSnippet,
			Content: e.Content,
			Tags:    e.Categories,
			PubDate: int64(e.PublishedDate),
			Author:  e.Author,
		}
		v = append(v, ne)
	}
	return feed.FeedEntity{s, v}, nil
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

func (this google_feedapi) Find(q, hl string) ([]feed.FeedEntity, error) {
	p := struct {
		q  string `param:"q"`
		hl string `param:"hl"` // default en, nil means en
		v  string `param:"v"`
	}{q: q, hl: hl, v: "1.0"}
	v := find_result{}
	uri := find_url + "?" + oauth2.HttpQueryEncode(p)
	c := curl.NewCurlerDetail(this.temp_folder, curl.CurlProxyPolicyAlwayseProxy, 0, this, nil)
	err := c.GetAsJson(uri, &v)
	if err != nil {
		return nil, err
	}
	x, err := v.to_feedentities()
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
func (this google_feedapi) Load(uri, hl string, num int, scoring bool) (feed.FeedEntity, error) {
	p := struct {
		q       string `param:"q"`
		hl      string `param:"hl"`      // default en
		num     *int   `param:"num"`     // default 4
		scoring string `param:"scoring"` // nil or h
		v       string `param:"v"`
	}{q: uri, hl: hl, num: make_num(num), scoring: make_scoring(scoring), v: "1.0"}
	c := curl.NewCurlerDetail(this.temp_folder, curl.CurlProxyPolicyAlwayseProxy, 0, this, nil)
	v := load_result{}
	url := load_url + "?" + oauth2.HttpQueryEncode(p)
	err := c.GetAsJson(url, &v)

	if err != nil {
		return feed.FeedEntity{}, err
	}
	s, err := v.to_feedentity()
	return s, err
}

const (
	feed_service = `https://ajax.googleapis.com/ajax/services/feed/`

	find_url = feed_service + `find`
	load_url = feed_service + `load`
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

type gf_error struct {
	code   int
	reason string
}

func (this gf_error) Error() string {
	return fmt.Sprintf("%v: %v", this.code, this.reason)
}
