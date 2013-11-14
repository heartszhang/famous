package googlefeedservice

import (
	"encoding/json"
	"fmt"
	"github.com/heartszhang/curl"
	"net/http"
	"net/url"
	"os"
)

type GoogleFeedApiFindResult struct {
	//	Error *GoogleFeedApiError
	ResponseDetails string                 `json:"responseDetails,omitempty"`
	ResponseStatus  int                    `json:"responseStatus"`
	ResponseData    *GoogleFeedApiFindData `json:"responseData,omitempty"`
}
type GoogleFeedApiFindData struct {
	Query   string                   `json:"query,omitempty"`
	Entries []GoogleFeedApiFindEntry `json:"entries,omitempty"`
}

type GoogleFeedApiFindEntry struct {
	Url            string `json:"url,omitempty"`
	Title          string `json:"title,omitempty"`
	ContentSnippet string `json:"contentSnippet"`
	Link           string `json:"link,omitempty"`
}

type GoogleFeedApiLoadResult struct {
	ResponseDetails string                 `json:"responseDetails,omitempty"`
	ResponseStatus  int                    `json:"responseStatus"`
	ResponseData    *GoogleFeedApiLoadData `json:"responseData,omitempty"`
}

type GoogleFeedApiLoadData struct {
	Feed *GoogleFeedApiFeed `json:"feed,omitempty"`
}

type GoogleFeedApiFeed struct {
	FeedUrl     string                   `json:"feedUrl,omitempty"`
	Title       string                   `json:"title,omitempty"`
	Link        string                   `json:"link,omitempty"`
	Author      string                   `json:"author,omitempty"`
	Description string                   `json:"description,omitempty"`
	Type        string                   `json:"type,omitempty"`
	Entries     []GoogleFeedApiFeedEntry `json:"entries,omitempty"`
}

type GoogleFeedApiFeedEntry struct {
	MediaGroups    []GoogleFeedApiMediaGroup `json:"mediaGroups,omitempty"`
	Title          string                    `json:"title,omitempty"`
	Link           string                    `json:"link,omitempty"`
	Author         string                    `json:"author,omitempty"`
	PublishedDate  string                    `json:"publishedDate"`
	ContentSnippet string                    `json:"contentSnippet,omitempty"`
	Content        string                    `json:"content,omitempty"`
	Categories     []string                  `json:"categories,omitempty"`
}

type GoogleFeedApiMediaGroup struct {
	Contents []GoogleFeedApiMediaContent `json:"contents,omitempty"`
}

type GoogleFeedApiMediaContent struct {
	Url         string                        `json:"url,omitempty"`
	Type        string                        `json:"type,omitempty"`
	Medium      string                        `json:"medium,omitempty"`
	Height      uint                          `json:"height"`
	Width       uint                          `json:"width"`
	IsDefault   *bool                         `json:"isDefault,omitempty"`
	Title       string                        `json:"title,omitempty"`
	Description string                        `json:"description,omitempty"`
	Keywords    string                        `json:"keywords,omitempty"`
	Thumbnails  []GoogleFeedApiMediaThumbnail `json:"thumbnail,omitempty"`
	Categories  []string                      `json:"category,omitempty"`
	//	Player      *GoogleFeedApiMediaPlayer     `json:"player,omitempty"`
	//	Embed       *GoogleFeedApiMediaEmbed      `json:"embed,omitempty"`
}

type GoogleFeedApiMediaThumbnail struct {
	Url    string `json:"url,omitempty"`
	Height uint   `json:"height"`
	Width  uint   `json:"width"`
	Time   string `json:"time,omitempty"`
}

type GoogleFeedApiMediaPlayer struct {
	Url    string `json:"url,omitempty"`
	Height uint   `json:"height"`
	Width  uint   `json:"width"`
}

type GoogleFeedApiService interface {
	Find(q, hl string) (GoogleFeedApiFindResult, error)
	Load(uri, hl string, num int, scoring bool) (GoogleFeedApiLoadResult, error)
}

type google_feedapi struct {
	temp_folder string
	interceptor func(*http.Request)
}

func NewGoogleFeedApi(refer, tmp string) GoogleFeedApiService {
	return &google_feedapi{
		temp_folder: tmp,
		interceptor: func(r *http.Request) {
			r.Header.Add("refer", refer)
		}}
}

func (this google_feedapi) Find(q, hl string) (GoogleFeedApiFindResult, error) {
	p := googlefeedapi_find_param{q: q, hl: hl}
	c := curl.NewInterceptCurler(this.temp_folder, this.interceptor)
	cache, err := c.Get(p.build(), curl.CurlProxyPolicyAlwayseProxy)
	v := GoogleFeedApiFindResult{}
	if err != nil {
		return v, err
	}
	f, err := os.Open(cache.Local)
	if err != nil {
		return v, err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&v)
	return v, err
}

func (this google_feedapi) Load(uri, hl string, num int, scoring bool) (GoogleFeedApiLoadResult, error) {
	p := googlefeedapi_load_param{q: uri, hl: hl, num: num, scoring: scoring}
	c := curl.NewInterceptCurler(this.temp_folder, this.interceptor)
	cache, err := c.Get(p.build(), curl.CurlProxyPolicyAlwayseProxy)
	v := GoogleFeedApiLoadResult{}
	if err != nil {
		return v, err
	}
	f, err := os.Open(cache.Local)
	if err != nil {
		return v, err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&v)
	return v, err
}

const (
	feed_service = `https://ajax.googleapis.com/ajax/services/feed/`
	find_url     = feed_service + `find`
	load_url     = feed_service + `load`
)

type googlefeedapi_find_param struct {
	q  string
	hl string // default en, nil means en
}

func (this googlefeedapi_find_param) build() string {
	var (
		q, v, hl string
	)
	if this.hl != "" {
		hl = fmt.Sprintf("hl=%v", url.QueryEscape(this.hl))
	}
	v = `v=1.0`
	q = fmt.Sprintf("q=%v", url.QueryEscape(this.q))
	x := find_url + "?" + v
	if hl != "" {
		x += "&" + hl
	}
	x += "&" + q
	return x
}

type googlefeedapi_load_param struct {
	q       string
	hl      string // default en
	num     int    // default 4
	scoring bool   // nil or h
}

func (this googlefeedapi_load_param) build() string {
	var (
		q, hl, num, scoring string
	)
	v := `v=1.0`
	if this.hl != "" {
		hl = fmt.Sprintf("hl=%v", url.QueryEscape(this.hl))
	}
	q = fmt.Sprintf("q=%v", url.QueryEscape(this.q))
	if this.num > 0 {
		num = fmt.Sprintf("num=%v", this.num)
	}
	if this.scoring {
		scoring = fmt.Sprintf("scoring=%v")
	}
	x := load_url + "?" + v
	if hl != "" {
		x += "&" + hl
	}
	if num != "" {
		x += "&" + num
	}
	if scoring != "" {
		x += "&" + scoring
	}
	x += "&" + q
	return x
}
