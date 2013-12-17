package curl

import (
	"encoding/json"
	"fmt"
	"github.com/heartszhang/gfwlist"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	CurlProxyPolicyUseProxy = iota
	CurlProxyPolicyNoProxy
	CurlProxyPolicyAlwayseProxy
)
const (
	gfw_url = `http://autoproxy-gfwlist.googlecode.com/svn/trunk/gfwlist.txt`
)

type curler struct {
	data_dir     string
	proxy_policy int
	dial_timeo   int
	interceptor  CurlerRoundTrip
	ruler        gfwlist.GfwRuler
}

func NewCurl(datadir string) Curler {
	return &curler{data_dir: datadir, proxy_policy: 0, dial_timeo: connection_timeout, interceptor: nil}
}

func NewInterceptCurler(datadir string, intercepter func(*http.Request)) Curler {
	return &curler{data_dir: datadir,
		interceptor: roundtrip_wrapper(intercepter),
		dial_timeo:  connection_timeout}
}

func NewCurlerDetail(datadir string, proxypolicy,
	dialtimeo int, intercepter CurlerRoundTrip, ruler gfwlist.GfwRuler) Curler {
	if dialtimeo == 0 {
		dialtimeo = connection_timeout
	}
	return &curler{data_dir: datadir,
		proxy_policy: proxypolicy,
		interceptor:  intercepter,
		dial_timeo:   dialtimeo,
		ruler:        ruler,
	}
}

const (
	connection_speedup_timeout = 4
	connection_timeout         = 14 // seconds
	response_timeout           = 20 // seconds
)

func (this *curler) new_proxy() proxy_func {
	if this.proxy_policy == CurlProxyPolicyAlwayseProxy {
		return http.ProxyFromEnvironment
	}
	if this.proxy_policy == CurlProxyPolicyNoProxy {
		return nil
	}
	if this.ruler != nil {
		return func(req *http.Request) (*url.URL, error) {
			uri := req.URL.String()
			blocked := this.ruler.IsBlocked(uri)
			if blocked {
				return http.ProxyFromEnvironment(req)
			}
			return nil, nil
		}
	}
	return nil
}

// ProxyFromEnvironment(req *Request) (*url.URL, error)
func (this *curler) do_post(uri string, form url.Values) (*http.Response, error) {
	retry := this.proxy_policy == CurlProxyPolicyUseProxy && this.ruler == nil
	proxy := this.new_proxy()
	if this.proxy_policy == CurlProxyPolicyAlwayseProxy {
		proxy = http.ProxyFromEnvironment
	}
	resp, err := this.new_client(proxy, nil).PostForm(uri, form)
	if retry && err != nil {
		resp, err = this.new_client(http.ProxyFromEnvironment, nil).PostForm(uri, form)
	}
	return resp, err
}

type curler_transport struct {
	transport   http.RoundTripper
	interceptor CurlerRoundTrip
	cache       *Cache
}

//RoundTrip(req *Request) (resp *Response, err error)
func (this *curler_transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if this.cache != nil && this.cache.StatusCode == http.StatusOK {
		//		req.Header.Set("Cache-Control", "max-age=0")
		if this.cache.ETag != "" {
			req.Header.Set("if-none-match", this.cache.ETag)
		}
		if this.cache.LastModified != "" {
			req.Header.Set("if-modified-since", this.cache.LastModified)
		}
	}

	if this.interceptor != nil {
		this.interceptor.RoundTrip(req)
	}
	resp, err = this.transport.RoundTrip(req)
	// we dont process not-modified here
	return
}

type proxy_func func(*http.Request) (*url.URL, error)

//Proxy func(*Request) (*url.URL, error)
func (this *curler) new_curler_transport(proxy proxy_func, cache *Cache) http.RoundTripper {
	trans := &http.Transport{
		Dial: new_timeout_dialer(this.dial_timeo),
		ResponseHeaderTimeout: time.Duration(response_timeout) * time.Second,
		Proxy: proxy,
	}
	return &curler_transport{
		transport:   trans,
		interceptor: this.interceptor,
		cache:       cache,
	}
}
func (this *curler) new_client(proxy proxy_func, cache *Cache) *http.Client {
	trans := this.new_curler_transport(proxy, cache)
	return &http.Client{Transport: trans}
}

func (this *curler) do_get(uri string, cache *Cache) (*http.Response, error) {
	retry := this.proxy_policy == CurlProxyPolicyUseProxy
	var proxy = this.new_proxy()
	resp, err := this.new_client(proxy, cache).Get(uri)
	if err != nil && retry {
		resp, err = this.new_client(http.ProxyFromEnvironment, cache).Get(uri)
	}
	return resp, err
}

// detect content-type via http://mimesniff.spec.whatwg.org/ if charset isn't declared
// treat gb2312 as gbk
// convert text/* and application/*+xml to utf-8
// only when converting sucessfully, Cache.LocalUtf8 would be set
// when convert non-utf8 encoded xml, file would be saved as utf-8, but xml declares another encoding
// xml-decoder should use a passthrough charset-reader
func (this *curler) GetUtf8(uri string) (Cache, error) {
	v, err := this.get(uri, false)
	if err != nil {
		return v, err
	}
	if v.LocalUtf8 != "" { // from cache, already converted
		return v, err
	}
	// text or application/*+xml
	if !mime_should_convert(v.Mime, v.Charset, true) {
		return v, fmt.Errorf("invalid-mime: %v, char: %v", v.Mime, v.Charset)
	}
	if v.Charset == "" {
		ct := file_detect_content_type(v.Local, v.Mime)
		mime, cs, _ := extract_charset(ct)
		v.Charset = cs
		if v.Mime == "" {
			v.Mime = mime
			if !mime_should_convert(mime, cs, false) {
				return v, fmt.Errorf("invalid-mime: %v, char: %v", mime, cs)
			}
		}
	}
	cacher := disk_cacher{data_folder: this.data_dir}

	// some website use cht by declaring gb2312 encoded
	if v.Charset == "" || v.Charset == "gb2312" {
		v.Charset = "gbk"
	}
	if v.Charset == "utf-8" {
		v.LocalUtf8 = v.Local
		v.LengthUtf8 = v.Length
		cacher.save_index(uri, v)
		return v, err
	}

	in, err := os.Open(v.Local)
	if err != nil {
		return v, err
	}
	defer in.Close()
	in2, err := new_utf8_reader(v.Charset, in)
	if err != nil {
		return v, err
	}
	defer in2.Close()

	out, err := ioutil.TempFile(this.data_dir, "u-"+mime_to_ext(v.Mime, v.Disposition))
	if err != nil {
		return v, err
	}
	defer out.Close()

	v.LengthUtf8, err = io.Copy(out, in2)

	if err == nil {
		v.LocalUtf8 = out.Name()
	}
	cacher.save_index(uri, v)
	return v, err
}

func (this *curler) PostForm(uri string, form url.Values) (int, error) {
	resp, err := this.do_post(uri, form)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, err
}
func (this *curler) PostFormAsJson(uri string, form url.Values, val interface{}) error {
	resp, err := this.do_post(uri, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return curler_error{resp.StatusCode, resp.Status} // content is ignored
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(val)
	return err
}

func (this *curler) PostFormAsString(uri string, form url.Values) (string, error) {
	resp, err := this.do_post(uri, form)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", curler_error{resp.StatusCode, resp.Status}
	}
	_, cs, _ := extract_charset(resp.Header.Get("content-type"))
	ireader, err := new_utf8_reader(cs, resp.Body)
	if err != nil {
		return "", err
	}
	d, err := ioutil.ReadAll(ireader)
	return string(d), err
}

func (this *curler) GetAsString(uri string) (rtn string, err error) {
	resp, err := this.do_get(uri, nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return rtn, curler_error{resp.StatusCode, resp.Status}
	}

	_, cs, _ := extract_charset(resp.Header.Get("content-type"))
	ireader, err := new_utf8_reader(cs, resp.Body)
	x, err := ioutil.ReadAll(ireader)
	if err != nil {
		return
	}
	return string(x), err
}

func (this *curler) GetAsJson(uri string, val interface{}) error {
	resp, err := this.do_get(uri, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return curler_error{resp.StatusCode, resp.Status}
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(val)
	return err
}

func (this *curler) get(uri string, savecache bool) (Cache, error) {
	cacher := disk_cacher{data_folder: this.data_dir}
	v, _ := cacher.load_index(uri)
	resp, err := this.do_get(uri, v)
	if err != nil {
		return Cache{Uri: uri, StatusCode: -1}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified {
		return *v, err
	}
	v = &Cache{Uri: uri}
	v.StatusCode = resp.StatusCode
	if resp.StatusCode != http.StatusOK {
		cacher.remove(uri)
		return *v, curler_error{resp.StatusCode, resp.Status}
	}
	v.ETag = resp.Header.Get("etag")
	v.LastModified = resp.Header.Get("last-modified")
	v.Expires = resp.Header.Get("expires")
	mime, cs, err := extract_charset(resp.Header.Get("content-type"))
	if err == nil {
		v.Mime = strings.ToLower(mime)
		v.Charset = strings.ToLower(cs)
	}
	v.Disposition = resp.Header.Get("content-disposition")

	out, err := ioutil.TempFile(this.data_dir, mime_to_ext(v.Mime, v.Disposition))
	if err != nil {
		return *v, err
	}
	defer out.Close()

	v.Length, err = io.Copy(out, resp.Body)
	v.Local = out.Name()
	if savecache {
		cacher.save_index(uri, *v)
	}
	return *v, nil
}

// use env-proxy or goagent for all http sessions, if direct conn fail
// detect charset by mimetype
// use server-site filename as name-prefix
func (this *curler) Get(uri string) (Cache, error) {
	return this.get(uri, true)
}

func init() {
	k := "HTTP_PROXY"
	pxy := os.Getenv(k)
	if pxy == "" {
		os.Setenv(k, "http://localhost:8087") // use goagent as default
	}
}

type CurlerRoundTrip interface {
	RoundTrip(*http.Request)
}
type roundtrip_wrapper func(*http.Request)

func (this roundtrip_wrapper) RoundTrip(r *http.Request) {
	this(r)
}
