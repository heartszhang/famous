package backend

import (
	"github.com/heartszhang/curl"
)

type ImageCache struct {
	Mime           string `json:"mime,omitempty" bson:"mime,omitempty"`
	ThumbnailLocal string `json:"thumbnail,omitempty" bson:"thumbnail,omitempty"`
	OriginLocal    string `json:"origin,omitempty" bson:"origin,omitempty"`
	Code           int    `json:"code" bson:"code"`
	Width          int    `json:"width" bson:"width"`
	Height         int    `json:"height" bson:"height"`
}

func image_get_or_cache(uri string) (ImageCache, error) {
	imgo := new_imagecache_operator()
	v, err := imgo.find(uri)

	if err == nil {
		return v, nil
	}
	c := curl.NewCurl(config.ImageDir)
	cache, err := c.Get(uri, 0)

	v.Mime = cache.Mime
	v.Code = cache.StatusCode
	v.OriginLocal = cache.Local
	if err != nil {
		return v, err
	}
	v.ThumbnailLocal, v.Mime, v.Width, v.Height, err = curl.NewThumbnail(cache.Local, config.ThumbnailDir, config.ThumbnailWidth, 0)
	imgo.save(uri, v)
	return v, err
}
