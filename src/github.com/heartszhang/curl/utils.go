package curl

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"mime"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func new_timeout_dialer(timeo int) func(string, string) (net.Conn, error) {
	return func(network, addr string) (net.Conn, error) {
		return net.DialTimeout(network, addr, time.Duration(timeo)*time.Second)
	}
}

func extract_charset(mime_header string) (mediatype, charset string, err error) {
	mediatype, params, err := mime.ParseMediaType(mime_header)
	charset = params["charset"]
	return
}

// text/*
// application/*+xml
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
		is_xml := strings.Contains(types[1], "xml")
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
FORBEGIN:
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
				break FORBEGIN
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
	return "text/xml; charset=" + x2[:i]
}

// DetectContentType will treat all xml as utf-8 encoded. so some extrac work should be done
func file_detect_content_type(local, mime string) string {
	f, err := os.Open(local)
	if err != nil {
		return "application/octec-stream" // treat as binary
	}
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

// only process type/subtype, without parameters
func MimeToExt(typesubtype string) string {
	types := strings.Split(typesubtype, "/")
	switch len(types) > 1 {
	case true:
		return types[1]
	}
	return typesubtype
}

func resolve_dir(dir string) string {
	if dir == "" {
		return os.TempDir()
	}
	return dir
}
