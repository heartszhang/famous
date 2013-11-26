package pubsub

import (
	"github.com/heartszhang/curl"
	"github.com/heartszhang/oauth2"
	"time"
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

func NewSuperFeedrPubSubscriber(verify_mode string) PubSubscriber {
	if verify_mode != "sync" && verify_mode != "async" {
		verify_mode = "async"
	}
	return pubsuber{verify_mode, hub_superfeed, 0}
}

const (
	//hub_callback = "http://famousfront.appspot.com/hub/callback"
	hub_callback = "http://iweizhi2.duapp.com/hub_callback"
	//	hub_superfeed = "https://pubsubhubbub.superfeedr.com"
	hub_superfeed  = "https://Hearts:Refresh@push.superfeedr.com"
	hub_google     = "https://pubsubhubbub.appspot.com/subscribe"
	superfeed_user = "Hearts"
	password       = "Refresh"
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
	resp, err := curl.NewCurlerDetail("", this.proxy_policy, 0, nil).PostForm(this.service_provider, oauth2.HttpQueryValues(p))
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
	resp, err := curl.NewCurlerDetail("", this.proxy_policy, 0, nil).PostForm(this.service_provider, oauth2.HttpQueryValues(p))
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
	return curl.NewCurlerDetail("", this.proxy_policy, 0, nil).GetAsString(u)
}

type unixtime int64

func (this unixtime) Time() time.Time {
	return time.Unix(int64(this), 0)
}

type PubsubMessage struct {
	Status struct {
		StatusCode        int      `json:"code"`
		StatusReason      string   `json:"http,omitempty"`
		Feed              string   `json:"feed"`
		LastParse         unixtime `json:"lastParse"`
		Period            unixtime `json:"period"`
		LastMaintenanceAt unixtime `json:"lastMaintenanceAt"`
		NextFetch         unixtime `json:"nextFetch"`
		LastFetch         unixtime `json:"lastFetch"`
	}
	Title         string               `json:"title,omitempty"`
	Subtitle      string               `json:"subtitle,omitempty"`
	StandardLinks PubsubStandardLink   `json:"standardLinks,omitempty"`
	Updated       unixtime             `json:"updated"`
	Items         []PubsubMessageEntry `json:"items,omitempty"`
}
type PubsubLink struct {
	Mime  string `json:"type,omitempty"`
	Href  string `json:"href,omitempty"`
	Title string `json:"title,omitempty"`
}
type PubsubStandardLink struct {
	Self    []PubsubLink `json:"self,omitempty"`
	Picture []PubsubLink `json:"picture,omitempty"`
}
type PubsubMessageEntry struct {
	StandardLinks *PubsubStandardLink `json:"standardLinks,omitempty"`
	PermalinkUrl  string              `json:"permalinkUrl,omitempty"`
	Verb          string              `json:"verb,omitempty"`
	Content       string              `json:"content,omitempty"`
	Summary       string              `json:"summary,omitempty"`
	Published     unixtime            `json:"published"`
	Updated       unixtime            `json:"updated"`
	Title         string              `json:"title,omitempty"`
	Categories    []string            `json:"categories,omitempty"`
}
