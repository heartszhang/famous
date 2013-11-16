package backend

import (
	"github.com/heartszhang/curl"
	"github.com/heartszhang/feedfeed"
)

// /feed/entry/image.json/{url}/{entry_id}
func image_description(uri string) (feedfeed.FeedImage, error) {
	imgo := new_imagecache_operator()
	v, err := imgo.find(uri)

	if err == nil {
		return v, nil
	}
	c := curl.NewCurl(config.ImageDir)
	cache, err := c.Get(uri)

	v.Mime = cache.Mime
	v.Code = cache.StatusCode
	v.OriginLocal = cache.Local
	if err != nil {
		return v, err
	}
	v.ThumbnailLocal, v.Mime, v.Width, v.Height, err = curl.NewThumbnail(cache.Local, config.ThumbnailDir, config.ThumbnailWidth, 0)
	go imgo.save(uri, v)
	return v, err
}
