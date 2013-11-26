package baidu

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/heartszhang/curl"
	"github.com/heartszhang/oauth2"
	"io"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"
)

const (
	bcms_urlbase = "http://bcms.api.duapp.com/rest/2.0/bcms/"
)

func sortedkeys_from_map(vals map[string][]string) []string {
	var v []string
	for k, _ := range vals {
		v = append(v, k)
	}
	sort.Strings(v)
	return v
}

func bcms_sign(method, uri string, params interface{}, client_secret string) string {
	vals := oauth2.HttpQueryValues(params)
	keys := sortedkeys_from_map(vals)
	var t []string
	for _, i := range keys {
		val := vals[i]
		for _, j := range val {
			t = append(t, fmt.Sprintf("%v=%v", i, j))
		}
	}
	q := fmt.Sprintf("%v %v%v%v", method, uri, strings.Join(t, "&"), client_secret)
	q = url.QueryEscape(q)
	h := md5.New()
	io.WriteString(h, q)
	sign := fmt.Sprintf("%x", h.Sum(nil))
	return sign
}

type BaiduMessageQueue interface {
	FetchOne() (string, error)
	FetchAny(msgid, fetch_num *uint) ([]string, error)
	FetchAnyAsJson(msgid, fetchnum *uint, v interface{}) error
	FetchOneAsJson(v interface{}) error
}

// via bcms rest api
func NewBcms(qname, accesstoken, clientid, clientsecret string) BaiduMessageQueue {
	return bcms_queue{
		queue_name:    qname,
		access_token:  accesstoken,
		client_id:     clientid,
		client_secret: clientsecret,
	}
}

// via iweizhi2.duapp.com/popup.json
// only support popup_one
func NewBcmsProxy(popaddress string) BaiduMessageQueue {
	return bcms_proxy{popaddress}
}

type bcms_proxy struct {
	uri string
}
type bcms_queue struct {
	queue_name    string
	access_token  string
	client_id     string
	client_secret string
}

// bcms rest api common parameters
type bcms_common struct {
	method       string `param:"method"`
	timestamp    uint   `param:"timestamp"`
	access_token string `param:"access_token"`
	client_id    string `param:"client_id"`
	sign         string `param:"sign"`
	expires      *uint  `param:"expires"`
	apiver       *uint  `param:"v"`
}

type bcms_fetch struct {
	bcms_common
	msg_id    *uint `param:"msg_id"`
	fetch_num *uint `param:"fetch_num"`
}

func (this bcms_queue) FetchAny(msgid, fetch_num *uint) ([]string, error) {
	v := make([]string, 0)
	var err error
	q := bcms_fetch{
		bcms_common: bcms_common{
			method:       "fetch",
			timestamp:    unix_current(),
			access_token: this.access_token,
			client_id:    this.client_id,
		},
		msg_id:    msgid,
		fetch_num: fetch_num,
	}

	var bcmsr struct {
		RequestId int64 `json:"request_id"`
		Response  struct {
			MessageNum int `json:"message_num"`
			Messages   []struct {
				MsgId   string `json:"msg_id"`
				Message string `json:"message"`
			} `json:"messages,omitempty"`
		} `json:"response_params"`
	}
	uri := bcms_urlbase + this.queue_name
	q.sign = bcms_sign("POST", uri, q, this.client_secret) // use client_secret to sign post-params
	vals := oauth2.HttpQueryValues(q)
	c := curl.NewCurlerDetail("", curl.CurlProxyPolicyNoProxy, 0, nil)
	err = c.PostFormAsJson(uri, vals, &bcmsr)

	if err == nil {
		for _, i := range bcmsr.Response.Messages {
			v = append(v, i.Message)
		}
	}
	return v, err
}

func (this bcms_queue) FetchAnyAsJson(msgid, fetchnum *uint, v interface{}) error {
	resultv := reflect.ValueOf(v)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		panic("result argument must be slice address")
	}
	slicev := resultv.Elem()
	jsons, err := this.FetchAny(msgid, fetchnum)
	if err != nil {
		return err
	}
	elemt := slicev.Type().Elem()
	for i, m := range jsons {
		if i < slicev.Len() {
			x := slicev.Index(i).Addr().Interface()
			err = json.Unmarshal([]byte(m), x)
		} else {
			elemp := reflect.New(elemt)
			err = json.Unmarshal([]byte(m), elemp.Interface())
			slicev = reflect.Append(slicev, elemp.Elem())
		}
		if err != nil {
			break
		}
	}
	resultv.Elem().Set(slicev)
	return err
}

func (this bcms_queue) FetchOne() (string, error) {
	v, err := this.FetchAny(nil, nil)
	if len(v) > 0 {
		return v[0], err
	}
	return "", err
}

func (this bcms_queue) FetchOneAsJson(v interface{}) error {
	m, err := this.FetchOne()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(m), v)
	return err
}

func (this bcms_proxy) FetchOneAsJson(v interface{}) error {
	client := curl.NewCurlerDetail("", curl.CurlProxyPolicyNoProxy, 0, nil)
	return client.GetAsJson(this.uri, v)
}

func (this bcms_proxy) FetchOne() (string, error) {
	panic("not implemented yet")
}
func (this bcms_proxy) FetchAnyAsJson(msgid, fetchnum *uint, v interface{}) error {
	panic("not implemented yet")
}

func unix_current() uint {
	return uint(time.Now().Unix())
}

func (this bcms_proxy) FetchAny(msgid, fetch_num *uint) ([]string, error) {
	panic("not implemented yet")
}
