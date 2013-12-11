package pubsub

import (
	"fmt"
	"github.com/heartszhang/curl"
	"github.com/heartszhang/oauth2"
	"github.com/heartszhang/unixtime"
)

type PubSubscriber interface {
	Subscribe(uri string) (int, error)
	Unsubscribe(uri string) (int, error)
	Retrieve(uri string, count int) (string, error)
}

type pubsuber struct {
	verify           string
	service_provider string
	proxy_policy     int
}

func NewGooglePubSubscriber() PubSubscriber {
	return pubsuber{"async", hub_google, curl.CurlProxyPolicyAlwayseProxy}
}

func NewSuperFeedrPubSubscriber(verify_mode, user, password string) PubSubscriber {
	if verify_mode != "sync" && verify_mode != "async" {
		verify_mode = "async"
	}
	svc := fmt.Sprintf("https://%v:%v@push.superfeedr.com", user, password)
	return pubsuber{verify_mode, svc, 0}
}

const (
	hub_callback = "http://iweizhi2.duapp.com/hub_callback"
	//	hub_superfeed = "https://Hearts:Refresh@push.superfeedr.com"
	hub_google = "https://pubsubhubbub.appspot.com/subscribe"
)

func strptr(s string) *string {
	return &s
}

type pubsub_param struct {
	mode     string  `param:"hub.mode"`
	topic    string  `param:"hub.topic"`
	callback string  `param:"hub.callback"`
	secret   *string `param:"hub.secret"`
	verify   *string `param:"hub.verify"`
	format   *string `param:"format"`
	retrive  *bool   `param:"retrieve"`
}

func (this pubsuber) Subscribe(uri string) (int, error) {
	p := pubsub_param{
		mode:     "subscribe",
		topic:    uri,
		callback: hub_callback,
		format:   strptr("json"),
		verify:   &this.verify,
	}
	resp, err := curl.NewCurlerDetail("", this.proxy_policy, 0, nil, nil).PostForm(this.service_provider, oauth2.HttpQueryValues(p))
	return resp, err
}

type pubunsub_param struct {
	mode     string  `param:"hub.mode"`
	topic    string  `param:"hub.topic"`
	callback string  `param:"hub.callback"`
	verify   *string `param:"hub.verify"`
}

func (this pubsuber) Unsubscribe(uri string) (int, error) {
	p := pubunsub_param{
		mode:     "unsubscribe",
		topic:    uri,
		callback: hub_callback,
		verify:   &this.verify,
	}
	resp, err := curl.NewCurlerDetail("", this.proxy_policy, 0, nil, nil).PostForm(this.service_provider, oauth2.HttpQueryValues(p))
	return resp, err
}

type pubretrive_param struct {
	mode     string  `param:"hub.mode"`
	topic    string  `param:"hub.topic"`
	count    *int    `param:"count"`
	format   *string `param:"format"`
	callback *string `param:"callback"`
}

func (this pubsuber) Retrieve(uri string, count int) (string, error) {
	p := struct {
		mode   string `param:"hub.mode"`
		topic  string `param:"hub.topic"`
		count  int    `param:"count"`
		format string `param:"format"`
	}{"retrieve", uri, count, "json"}
	u := this.service_provider + "?" + oauth2.HttpQueryEncode(p)
	return curl.NewCurlerDetail("", this.proxy_policy, 0, nil, nil).GetAsString(u)
}

type PubsubMessage struct { // same as FeedSource
	//	XMLName xml.Name `json:"-" xml:"feed"`
	Status struct {
		StatusCode        int               `json:"code" xml:"code,attr"`
		StatusReason      string            `json:"http,omitempty" xml:"http,omitempty"`
		Feed              string            `json:"feed" xml:"feed,attr,omitempty"`
		Period            int64             `json:"period" xml:"period"`
		LastParse         unixtime.UnixTime `json:"lastParse" xml:"last_parse"`
		LastMaintenanceAt unixtime.UnixTime `json:"lastMaintenanceAt" xml:"last_maintenance_at"`
		NextFetch         unixtime.UnixTime `json:"nextFetch" xml:"next_fetch"`
		LastFetch         unixtime.UnixTime `json:"lastFetch" xml:"last_fetch"`
		EntriesCount      int               `json:"entriesCountSinceLastMaintenance" xml:"entries_count_since_last_maintenance"`
	} `json:"status" xml:"status"`
	Title         string               `json:"title,omitempty" xml:"title"`
	Subtitle      string               `json:"subtitle,omitempty" xml:"subtitle,omitempty"`
	StandardLinks PubsubStandardLink   `json:"standardLinks,omitempty" xml:"-"`
	Updated       unixtime.UnixTime    `json:"updated" xml:"updated"`
	Items         []PubsubMessageEntry `json:"items,omitempty" xml:"entry,omitempty"`
	Links         []PubsubLink         `json:"-" xml:"link,omitempty"`
}
type PubsubLink struct {
	Rel   string `json:"rel,omitempty" xml:"rel,omitempty"`
	Mime  string `json:"type,omitempty" xml:"type,omitempty"`
	Href  string `json:"href,omitempty" xml:"href,omitempty"`
	Title string `json:"title,omitempty" xml:"title,omitempty"`
}
type PubsubStandardLink struct {
	Self    []PubsubLink `json:"self,omitempty" xml:"-"`
	Picture []PubsubLink `json:"picture,omitempty" xml:"-"`
}
type PubsubMessageEntry struct {
	StandardLinks *PubsubStandardLink `json:"standardLinks,omitempty" xml:"-"`
	Uri           string              `json:"permalinkUrl,omitempty" xml:"id,omitempty"`
	Verb          string              `json:"verb,omitempty" xml:"-"`
	Content       string              `json:"content,omitempty" xml:"content"`
	Summary       string              `json:"summary,omitempty" xml:"summary"`
	Published     unixtime.UnixTime   `json:"published" xml:"published"`
	Updated       unixtime.UnixTime   `json:"updated" xml:"updated"`
	Title         string              `json:"title,omitempty" xml:"title"`
	Links         []PubsubLink        `json:"-" xml:"link,omitempty"`
	Categories    []string            `json:"categories,omitempty" xml:""`
}
