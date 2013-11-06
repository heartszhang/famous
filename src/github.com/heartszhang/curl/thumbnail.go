package curl

import (
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
)

func NewThumbnail(localfile, tmp string, width, height uint) (filepath, mediatype string, w, h int, err error) {
	f, err := os.Open(localfile)
	if err != nil {
		return
	}
	defer f.Close()
	img, mediatype, err := image.Decode(f)
	if err != nil {
		return
	}
	w = img.Bounds().Max.X
	h = img.Bounds().Max.Y

	imgnew := resize.Resize(width, height, img, resize.MitchellNetravali)
	mediatype = "image/jpeg"
	of, err := ioutil.TempFile(tmp, "jpeg.")
	if err != nil {
		return
	}
	defer of.Close()
	err = jpeg.Encode(of, imgnew, &jpeg.Options{90})
	if err != nil {
		return
	}
	filepath = of.Name()
	return
}

func NewJpegThumbnail(localfile, tmp string, width, height uint) (filepath string, w, h int, err error) {
	fp, mt, w, h, err := NewThumbnail(localfile, tmp, width, height)
	if err != nil {
		return
	}
	filepath = fp + "." + mt
	os.Rename(fp, filepath)
	return
}
