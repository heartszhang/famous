package videothumbnail

import (
	"fmt"
	"github.com/heartszhang/curl"
	"net/url"
	"strconv"
	"strings"
)

func DescribeVideo(uri string) (v VideoDescription, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}
	domain := domain_from_host(u.Host)
	processor := processor_from_domain(domain)
	v, err = processor.process(u)
	return
}

type VideoDescription struct {
	Image     string
	Thumbnail string
	Title     string
	Tag       string
	Duration  int
}

func domain_from_host(host string) string {
	secs := strings.Split(host, ".")
	secc := len(secs)
	if secc < 3 {
		return host
	}
	begin := secc - 2
	return strings.Join(secs[begin:], ".")
}

type thumbnail_processor interface {
	process(uri *url.URL) (VideoDescription, error)
}
type processor_dummy struct{}

func (this processor_dummy) process(uri *url.URL) (VideoDescription, error) {
	return VideoDescription{}, vd_error{"unknown-processor", uri.String()}
}

type processor_56 struct{}

func (this processor_56) process(uri *url.URL) (VideoDescription, error) {
	uid := this.uid_from_uri(uri)
	if uid == "" {
		return VideoDescription{}, vd_error{"56 unrecognize uid", uri.String()}
	}
	vxmlurl := this.metaurl_from_uid(uid)
	var vxml struct {
		Info struct {
			Image     string `json:"bimg,omitempty"`
			Vid       string `json:"vid"`
			Uid       string `json:"textid"`
			Thumbnail string `json:"img"`
			Key       string `json:"key"`
			Tag       string `json:"tag"`
			Title     string `json:"Subject"`
			Duration  string `json:"duration"`
			Files     []struct {
				Length   string `json:"filesize"`
				Duration string `json:"totaltime"`
				Url      string `json:"url"`
				Type     string `json:"type"`
			} `json:"rfiles,omitempty"`
		} `json:"info"`
		Reason      string `json:"msg"`
		SegmentSize int    `json:"segsize"`
		Status      string `json:"status"`
		P2p         int    `json:"p2p"`
	}
	err := curl.NewCurl("").GetAsJson(vxmlurl, &vxml)
	v := VideoDescription{
		Image:     vxml.Info.Image,
		Thumbnail: vxml.Info.Thumbnail,
		Title:     vxml.Info.Title,
		Tag:       vxml.Info.Tag,
		Duration:  atoi(vxml.Info.Duration),
	}
	return v, err
}

func atoi(x string) int {
	v, _ := strconv.Atoi(x)
	return v
}

func (this processor_56) uid_from_uri(uri *url.URL) string {
	secs := strings.Split(uri.Path, "/")
	var x string
	for _, sec := range secs {
		if strings.HasPrefix(sec, "v_") {
			x = sec
			break
		}
	}
	repl := strings.NewReplacer("v_", "", ".swf", "", ".html", "")
	uid := repl.Replace(x)
	return uid
}
func (this processor_56) metaurl_from_uid(uid string) string {
	// http://vxml.56.com/json/{$matches[1]}/?src=out
	return fmt.Sprintf("http://vxml.56.com/json/%v/?src=out", uid)
}

func processor_from_domain(domain string) thumbnail_processor {

	switch domain {
	case "56.com":
		return processor_56{}
	case "youku.com":
		fallthrough
	case "tudou.com":
		fallthrough
	case "ku6.com":
		fallthrough
	case "letv.com":
		fallthrough
	case "sina.com":
		fallthrough
	case "iqiyi.com":
		fallthrough
	case "sohu.com":
		fallthrough
	case "qq.com":
		fallthrough
	case "xunlei.com":
		fallthrough
	default:
		return processor_dummy{}
	}
}

type vd_error struct {
	code   string
	reason string
}

func (this vd_error) Error() string {
	return fmt.Sprintf("%v: %v", this.code, this.reason)
}
