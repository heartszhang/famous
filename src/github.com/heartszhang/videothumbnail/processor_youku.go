package videothumbnail

import (
	"math/rand"
	"regexp"
	"fmt"
	"net/url"
	"github.com/heartszhang/curl"
)

var (
	ykrs1 = regexp.MustCompile(`id_([^/._])+\.html`)  // http://v.youku.com/v_show/id_XMjI4MDM4NDc2.html
	ykrs2 = regexp.MustCompile(`sid/([^/._])+/`)  // http://player.youku.com/player.php/sid/XMjU0NjI2Njg4/v.swf
)

func (this processor_youku) process(uri *url.URL) (VideoDescription, error) {
	ids := ykrs1.FindStringSubmatch(uri.Path)
	if len(ids) != 2 {
		ids = ykrs2.FindStringSubmatch(uri.Path)
	}
	var vid string
	if len(ids) == 2 {
		vid = ids[1]
	}
	if vid == "" {
		return VideoDescription{}, vd_error{"no video id", uri.String()}
	}
	var v struct {
		Data []struct {
			Logo    string   `json:"logo,omitempty"`
			Tags    []string `json:"tags,omitempty"`
			Title   string   `json:"title,omitempty"`
			Seconds float64  `json:"seconds"`
		} `json:"data,omitempty"`
	}
	u := fmt.Sprintf("http://v.youku.com/player/getPlayList/VideoIDS/%v/timezone/+08/version/5/source/out?password=&ran=%v&n=%v", vid, rand.Int(), rand.Int())
	err := curl.NewCurl("").GetAsJson(u, &v)
	if len(v.Data) == 0 {
		return VideoDescription{}, err
	}
	return VideoDescription{Image: v.Data[0].Logo, Thumbnail: v.Data[0].Logo, Title: v.Data[0].Title}, nil
}