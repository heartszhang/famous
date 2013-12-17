package videothumbnail

import (
	"regexp"
	"fmt"
	"io/ioutil"
	"net/url"
	"github.com/heartszhang/curl"
	"strconv"
)

func (this processor_qq) process(uri *url.URL) (VideoDescription, error) {
	cid := uri.Query().Get("cid")
	vid := uri.Query().Get("vid")
	if cid != "" && vid != "" {
		return this.image_from_cidvid(cid, vid), nil
	}
	r1s := r1.FindStringSubmatch(uri.Path)
	if len(r1s) != 0 && vid != "" {
		return this.image_from_cidvid(r1s[1], vid), nil
	}
	r3s := r3.FindStringSubmatch(uri.Path)
	if len(r3s) == 3 {
		return this.image_from_cidvid(r3s[1], r3s[2]), nil
	}
	c := curl.NewCurlerDetail("", 0, 0, nil, nil)
	cache, err := c.GetUtf8(uri.String())
	if err != nil {
		return VideoDescription{}, err
	}
	return this.image_from_htmlbody(cache.LocalUtf8)
}

var (
	r1 = regexp.MustCompile(`/cover/o/([^/.]+)\.html`)
	r2 = regexp.MustCompile(`/cover/d/([^/.]+)\.html`)
	r3 = regexp.MustCompile(`/cover/d/([^/._]+)/([^./_]+)\.html`)
)

func (this processor_qq) image_from_cidvid(cid, vid string) VideoDescription {
	x := fmt.Sprintf("http://vpic.video.qq.com/%v/%v_1.jpg", cid, vid)
	return VideoDescription{Image: x, Thumbnail: x}
}

/*
 * http://v.qq.com/cover/o/o9tab7nuu0q3esh.html?vid=97abu74o4w3_0
 * http://v.qq.com/play/97abu74o4w3.html
 * http://v.qq.com/cover/d/dtdqyd8g7xvoj0o/9SfqULsrtSb.html
 */
var (
	rc = regexp.MustCompile(`var COVER_INFO\s*=\s*{([^}]*)};`)
	rv = regexp.MustCompile(`var VIDEO_INFO\s*=\s*{([^}]*)};`)
)

func (this processor_qq) image_from_htmlbody(fp string) (VideoDescription, error) {
	vd := VideoDescription{}
	c, err := ioutil.ReadFile(fp)
	if err != nil {
		return vd, err
	}
	content := string(c)
	coverinfo := rc.FindStringSubmatch(content)
	videoinfo := rv.FindStringSubmatch(content)
	if len(coverinfo) != 2 || len(videoinfo) != 2 {
		return vd, vd_error{"no videoinfo", this.origin_uri}
	}
	var vid, cid string
	jsjson_foreach_field(coverinfo[1], func(n, v string) {
		switch n {
		case "pic":
			vd.Image = v
			vd.Thumbnail = v
		case "id":
			cid = v
		case "title":
			vd.Title = v
		}
	})
	jsjson_foreach_field(videoinfo[1], func(n, v string) {
		switch n {
		case "vid":
			vid = v
		case "duration":
			vd.Seconds, _ = strconv.Atoi(v)
		}
	})
	if vd.Image != "" {
		return vd, nil
	}
	if cid == "" || vid == "" {
		return vd, vd_error{"invalid vid", this.origin_uri}
	}
	return this.image_from_cidvid(cid, vid), nil
}