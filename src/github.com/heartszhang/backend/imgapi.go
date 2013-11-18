package backend

import (
	"github.com/heartszhang/curl"
	"github.com/heartszhang/feedfeed"
	vt "github.com/heartszhang/videothumbnail"
)

// /api/image/video.thumbnail?uri=
func image_videothumbnail(uri string) (vt.VideoDescription, error) {
	return vt.DescribeVideo(uri)
}

// /feed/entry/image.json/{url}/{entry_id}
func image_description(uri string) (feedfeed.FeedImage, error) {
	imgo := new_imagecache_operator()
	v, err := imgo.find(uri)

	if err == nil {
		return v, nil
	}
	c := curl.NewCurl(config.ImageFolder)
	cache, err := c.Get(uri)

	v.Mime = cache.Mime
	v.Code = cache.StatusCode
	v.OriginLocal = cache.Local
	if err != nil {
		return v, err
	}
	v.ThumbnailLocal, v.Mime, v.Width, v.Height, err = curl.NewThumbnail(cache.Local, config.ThumbnailFolder, config.ThumbnailWidth, 0)
	go imgo.save(uri, v)
	return v, err
}

func image_description_cached(uri string) (feedfeed.FeedImage, error) {
	imgo := new_imagecache_operator()
	v, err := imgo.find(uri)
	return v, err
}

func image_dimension(uri string) (feedfeed.FeedImage, error) {
	imgo := new_imagecache_operator()
	v, err := imgo.find(uri)
	if err == nil {
		return v, err
	}
	v.Mime, v.Width, v.Height, _, err = curl.DescribeImage(uri)
	return v, err
}
