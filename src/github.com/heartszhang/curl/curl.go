package curl

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"github.com/qiniu/iconv"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Cache struct {
	Length      int64
	Mime        string
	Charset     string
	Local       string
	Disposition string
	LocalUtf8   string
	LengthUtf8  int64
	StatusCode  int
}

const (
	CurlProxyPolicyUseProxy = iota
	CurlProxyPolicyNoProxy
	CurlProxyPolicyAlwayseProxy
)

type Curler interface {
	Get(uri string, proxypolicy int) (Cache, error)
	GetUtf8(uri string, proxypolicy int) (Cache, error)
}

type curler struct {
	data_dir string
}

func NewCurl(datadir string) Curler {
	return &curler{data_dir: datadir}
}

const (
	connection_speedup_timeout = 4
	connection_timeout         = 10 //seconds
	response_timeout           = 20 // seconds
)

func new_timeout_dialer(timeo int) func(string, string) (net.Conn, error) {
	return func(network, addr string) (net.Conn, error) {
		return net.DialTimeout(network, addr, time.Duration(timeo)*time.Second)
	}
}

/*
func timeout_dialer(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, time.Duration(connection_timeout*time.Second))
}
*/
func client_do_get(client *http.Client, uri string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err == nil {
		req.Header.Add("accept-encoding", "gzip,deflate")
		resp, err = client.Do(req)
	}
	return
}

func do_get(uri string, useproxy int) (*http.Response, error) {
	return do_get_timeo(uri, useproxy, connection_timeout)
}

func do_get_timeo(uri string, useproxy int, timeo int) (*http.Response, error) {
	trans := &http.Transport{
		Dial: new_timeout_dialer(timeo),
		ResponseHeaderTimeout: time.Duration(response_timeout * time.Second),
	}

	noretry := true
	switch useproxy {
	case CurlProxyPolicyNoProxy:
		trans.Proxy = nil
	case CurlProxyPolicyAlwayseProxy:
		trans.Proxy = http.ProxyFromEnvironment
	case CurlProxyPolicyUseProxy:
		fallthrough
	default:
		noretry = false
	}
	client := &http.Client{Transport: trans}
	resp, err := client_do_get(client, uri)
	if err != nil && noretry == false {
		fmt.Println("try again with proxy", uri, err)
		trans.Proxy = http.ProxyFromEnvironment
		resp, err = client_do_get(client, uri)
	}
	return resp, err
}
func extract_charset(mime_header string) (mediatype, charset string, err error) {
	mediatype, params, err := mime.ParseMediaType(mime_header)
	charset = params["charset"]
	return
}

func mime_should_convert(mime, charset string, ignore_empty_mime bool) bool {
	if charset != "" {
		return true
	}
	types := strings.Split(mime, "/")
	if len(types) != 2 {
		return ignore_empty_mime
	}
	switch types[0] {
	case "text":
		return true
	case "application":
		is_xml := strings.Contains(types[1], "+xml")
		return is_xml
	default:
		return false
	}
}

// <meta http-equiv="" content=xxx/>...
// <meta charset=''/>
// return content-type
func detect_charset_by_token(attrs []html.Attribute) (string, bool) {
	var http_equiv, content, charset string
	for _, attr := range attrs {
		switch attr.Key {
		case "http-equiv":
			http_equiv = attr.Val
		case "content":
			content = attr.Val
		case "charset":
			charset = attr.Val
		}
	}
	switch {
	case strings.ToLower(http_equiv) == "content-type":
		return content, true
	case len(charset) > 0:
		return "text/html; charset=" + charset, true
	}
	return "", false
}

func html_detect_content_type(head []byte) string {
	reader := bytes.NewReader(head)
	z := html.NewTokenizer(reader)
	expect_html_root := true
	for tt := z.Next(); tt != html.ErrorToken; tt = z.Next() {
		t := z.Token()
		switch {
		case t.Data == "meta" && (tt == html.StartTagToken || tt == html.SelfClosingTagToken):
			if ct, ok := detect_charset_by_token(t.Attr); ok == true {
				return ct
			}
		case t.Data == "head" && tt == html.EndTagToken:
			break
			// un-html file
		case expect_html_root && (tt == html.StartTagToken || tt == html.SelfClosingTagToken):
			if t.Data == "html" {
				expect_html_root = false
			} else {
				break
			}
		}
	}
	return ""
}

// like <?xml version="1.0" encoding="us-ascii"?>
func xml_detect_content_type(head string) string {
	dft := "text/xml; charset=gbk"
	x := `encoding="`
	i := strings.Index(head, x)
	if i == -1 {
		return dft
	}

	x2 := head[i+len(x):]
	i = strings.Index(x2, `"`)
	if i == -1 {
		return dft
	}
	return x2[:i]
}

// DetectContentType will treat all xml as utf-8 encoded. so some extrac work should be done
func file_detect_content_type(local, mime string) string {
	f, err := os.Open(local)
	if err != nil {
		return "application/octec-stream"
	}
	defer f.Close()
	head := make([]byte, 512)
	_, err = f.Read(head)
	f.Close()

	ct := http.DetectContentType(head)

	if ct == "text/xml; charset=utf-8" { // this file may be encoded with other charset
		ct = xml_detect_content_type(string(head))
	} else if ct == "text/html; charset=utf-8" { // charset is hard coded in html.DetectContentType
		ct = html_detect_content_type(head)
	}

	return ct
}

// detect content-type via http://mimesniff.spec.whatwg.org/ if charset isn't declared
// treat gb2312 as gbk
// convert text/* and application/*+xml to utf-8
// only when converting sucessfully, Cache.LocalUtf8 would be set
// when convert non-utf8 encoded xml, file would be saved as utf-8, but xml declares another encoding
// xml-decoder should use a passthrough charset-reader
func (this *curler) GetUtf8(uri string, proxypolicy int) (Cache, error) {
	v, err := this.Get(uri, proxypolicy)
	if err != nil {
		return v, err
	}
	// text or application/*+xml
	if !mime_should_convert(v.Mime, v.Charset, true) {
		return v, err
	}
	if v.Charset == "" {
		ct := file_detect_content_type(v.Local, v.Mime)
		mime, cs, _ := extract_charset(ct)
		v.Charset = cs
		if v.Mime == "" {
			v.Mime = mime
			if !mime_should_convert(mime, cs, false) {
				return v, err
			}
		}
	}

	// some website use cht by declaring gb2312 encoded
	if v.Charset == "" || v.Charset == "gb2312" {
		v.Charset = "gbk"
	}
	if v.Charset == "utf-8" {
		v.LocalUtf8 = v.Local
		v.LengthUtf8 = v.Length
		return v, err
	}
	converter, err := iconv.Open("utf-8", v.Charset)
	if err != nil {
		return v, err
	}
	defer converter.Close()

	out, err := ioutil.TempFile(this.data_dir, "u-"+mime_to_ext(v.Mime, v.Disposition))
	if err != nil {
		return v, err
	}
	defer out.Close()

	in, err := os.Open(v.Local)
	if err != nil {
		return v, err
	}
	defer in.Close()

	log.Println(v.Local, " =>", out.Name())

	reader := iconv.NewReader(converter, in, 0)

	v.LengthUtf8, err = io.Copy(out, reader)

	if err == nil {
		v.LocalUtf8 = out.Name()
	}
	return v, err
}

// use env-proxy or goagent for all http sessions, if direct conn fail
// detect charset by mimetype
// use server-site filename as name-prefix
func (this *curler) Get(uri string, useproxy int) (Cache, error) {
	v := Cache{}
	resp, err := do_get(uri, useproxy)
	if err != nil {
		return v, err
	}
	v.StatusCode = resp.StatusCode
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return v, fmt.Errorf("%v", resp.Status)
	}
	mime, cs, err := extract_charset(resp.Header.Get("content-type"))
	if err == nil {
		v.Mime = strings.ToLower(mime)
		v.Charset = strings.ToLower(cs)
	}
	v.Disposition = resp.Header.Get("content-disposition")

	reader, err := uncompress(resp)
	if err != nil {
		return v, err
	}
	defer reader.Close()

	out, err := ioutil.TempFile(this.data_dir, mime_to_ext(v.Mime, v.Disposition))
	if err != nil {
		return v, err
	}
	defer out.Close()

	v.Length, err = io.Copy(out, reader)
	v.Local = out.Name()

	return v, nil
}

// Content-Disposition: attachment; filename=genome.jpeg;
// type/subtype; param=""
// use server-side filename first
// use subtype as ext, subtype may be rdf+xml etc.
func mime_to_ext(typesubtype, dispose string) string {
	_, params, _ := mime.ParseMediaType(dispose)
	filename := params["filename"]
	if filename != "" {
		return filename + "-"
	}

	types := strings.Split(typesubtype, "/")
	switch len(types) > 1 {
	case true:
		return types[1] + "."
	}
	return typesubtype
}
func MimeToExt(typesubtype string) string {
	types := strings.Split(typesubtype, "/")
	switch len(types) > 1 {
	case true:
		return types[1]
	}
	return typesubtype
}

// process gzip/deflate
func uncompress(resp *http.Response) (v io.ReadCloser, err error) {
	encoding := resp.Header.Get("content-encoding")
	switch encoding {
	default:
		v = resp.Body
	case "gzip":
		v, err = gzip.NewReader(resp.Body)
	case "deflate":
		v = flate.NewReader(resp.Body)
	}
	return
}

func init() {
	k := "HTTP_PROXY"
	pxy := os.Getenv(k)
	if pxy == "" {
		os.Setenv(k, "http://localhost:8087") // use goagent as default
	}
}
