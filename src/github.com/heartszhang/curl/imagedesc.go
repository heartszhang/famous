package curl

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime"
	"net/http"
	"strings"
)

func DescribeImage(uri string) (mediatype string, width, height int, filelength int64, err error) {
	c := curler{
		dial_timeo: connection_speedup_timeout,
	}
	resp, err := c.do_get(uri, nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		filelength = int64(-resp.StatusCode)
		err = fmt.Errorf("%v: %v", resp.StatusCode, http.StatusText(resp.StatusCode))
		return
	}
	filelength = resp.ContentLength // -1 means unknown
	if filelength < 16 && filelength >= 0 {
		err = fmt.Errorf("%v: %v", filelength, "content-length insufficient info")
		return
	}
	ct := resp.Header.Get("Content-Type")
	mediatype, _, _ = mime.ParseMediaType(ct)
	types := strings.Split(mediatype, "/")
	if types[0] != "image" {
		err = fmt.Errorf("%v: unknown mime %v", uri, mediatype)
		return
	}
	ic, mediatype, err := image.DecodeConfig(resp.Body)
	width = ic.Width
	height = ic.Height
	return
}
