package curl

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime"
	"net/http"
)

func DescribeImage(uri string) (mediatype string, width, height int, filelength int64, err error) {
	resp, err := do_get_timeo(uri, 0, connection_speedup_timeout)
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

	ic, mediatype, err := image.DecodeConfig(resp.Body)
	width = ic.Width
	height = ic.Height
	return
}
