package videothumbnail

import (
	"net/url"
	"github.com/heartszhang/curl"
	"strings"
	"fmt"
)


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
		Tags:       []string{vxml.Info.Tag},
		Seconds:  atoi(vxml.Info.Duration),
	}
	return v, err
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