package backend

import (
	"github.com/heartszhang/cleaner"
	"github.com/heartszhang/curl"
	feed "github.com/heartszhang/feedfeed"
)

// since_unixtime , 0: from now
// categories, categories mask, every bit represent a category
// count: entries per page
// page: page no, start at 0
func feeds_entries_since(since_unixtime int64, category string, count uint, page uint) ([]feed.FeedEntry, error) {

	return []feed.FeedEntry{}, nil
}

func feedentry_unread(source string, count uint, page uint) ([]feed.FeedEntry, error) {
	c := curl.NewCurl(backend_config().FeedEntryDir)
	cache, err := c.GetUtf8(source, curl.CurlProxyPolicyUseProxy)
	if err != nil || cache.LocalUtf8 == "" {
		return nil, err
	}
	v, err := feed.CreateFeedEntriesRss2(cache.LocalUtf8)
	v = feed_entries_unreaded(v)
	v = feed_entries_clean(v)
	v = feed_entries_statis(v)
	v = feed_entries_clean_summary(v)
	v = feed_entries_clean_fulltext(v)
	v = feed_entries_autotag(v)
	v = feed_entries_backup(v)
	return v, err
}

func feedsource_all() ([]feed.FeedSource, error) {
	dbo := new_feedsource_operator()
	fs, err := dbo.all()
	return fs, err
}

func feedentry_mark(uri string, flags uint) (uint, error) {
	dbo := new_feedentry_operator()
	err := dbo.mark(uri, flags)
	return flags, err
}

// /feed/entry/umark.json/{id}/{flags}
func feedentry_umark(uri string, flags uint) (uint, error) {
	dbo := new_feedentry_operator()
	err := dbo.mark(uri, flags)
	return flags, err
}

/*
	media_type   uint   // feed_media_type...
	Uri          string `json:"uri,omitempty" bson:"uri,omitempty"`                     // url
	Alias        string `json:"alias,omitempty" bson:"alias,omitempty"`                 // title may be
	Local        string `json:"local,omitempty" bson:"local,omitempty"`                 // downloaded origin html
	CleanedLocal string `json:"cleaned_local,omitempty" bson:"cleaned_local,omitempty"` // cleaned-doc local rel path
	Words        int    `json:"words" bson:"words"`                                     // words after cleaned
	Sentences    int    `json:"sentences" bson:"sentences"`                             // sentences after cleaned
	Links        int    `json:"links" bson:"links"`                                     // links after cleaned
	Density      int    `json:"density" bson:"density"`                                 // density of original doc
	Length       int64  `json:"length" bson:"length"`
	Readable     bool   `json:"readable" bson:"readable"` // cleaned doc has perfect content

	Images []FeedMedia `json:"images,omitempty" bson:"images,omitempty"`
	Videos []FeedMedia `json:"videos,omitempty" bson:"videos,omitempty"`
	Audios []FeedMedia `json:"audios,omitempty" bson:"audios,omitempty"`
*/
// /feed/entry/full_text.json/{url}/{entry_id}
func feedentry_fulltext(uri string, entry_uri string) (v feed.FeedLink, err error) {
	c := curl.NewCurl(config.DocumentDir)
	cache, err := c.GetUtf8(uri, 0)
	v.Uri = uri
	v.Local = cache.LocalUtf8
	v.Length = cache.LengthUtf8
	if err != nil {
		return v, err
	}
	doc, err := html_create_from_file(cache.LocalUtf8)
	if err != nil {
		return v, err
	}
	article, sum, err := cleaner.MakeHtmlReadable(doc, uri)
	v.Images = make([]feed.FeedMedia, len(sum.Images))
	for idx, imguri := range sum.Images {
		v.Images[idx].Uri = imguri
	}
	v.Words = sum.WordCount
	v.Links = sum.LinkCount
	v.CleanedLocal, err = html_write_file(article, config.DocumentDir)
	return v, err
}

// /feed/entry/image.json/{url}/{entry_id}
func feedentry_image(url string, entry_id string) (feed.FeedMedia, error) {
	v := image_from_cache(url)
	if v.Local != "" {
		return v, nil
	}
	//	v := feed.FeedMedia{Uri: url}
	c := curl.NewCurl(config.ImageDir)
	cache, err := c.Get(url, 0)
	if err != nil {
		return v, err
	}
	v.Local = cache.Local
	v.Mime = cache.Mime
	v.Length = cache.Length
	v.Thumbnail, _, v.Width, v.Height, err = curl.NewThumbnail(cache.Local, config.ThumbnailDir, config.ThumbnailWidth, 0)
	go image_to_cache(v)
	return v, err
}

func image_from_cache(url string) feed.FeedMedia {
	return feed.FeedMedia{Uri: url}
}

func image_to_cache(img feed.FeedMedia) {
}

// /feed/entry/media.json/{url}/{entry_id}/{media_type:[0-9]+}

func feedentry_media(url string, entry_id string, media_type uint) (feed.FeedMedia, error) {
	return feed.FeedMedia{}, nil
}

// /feed/entry/drop.json/{id}

// id is mongo's _id
func feedentry_drop(id string) error {
	return nil
}

// select a idle category_id, assigned to category
func feedcategory_create(name string) (string, error) {
	fco := new_feedcategory_operator()
	uid, err := fco.save(name)

	if uid == nil {
		return "", err
	}
	return uid.(string), err
}

// id : isn't root or all, drop the category whoes name is name
// id : other, drop categories
// name : can be empty. if id is root or all, name cann't be empty
func feedcategory_drop(name string) error {
	return nil
}

// /tick.json

func tick() (FeedsStatus, error) {
	s := backend_status()
	return s, nil
}

func feedsource_subscribe(uri string, source_type uint) (v feed.FeedSource, err error) {
	fso := new_feedsource_operator()
	fs, err := fso.find(uri)
	if err == nil {
		return *fs, nil
	}
	curler := curl.NewCurl(backend_config().FeedSourceDir)
	cache, err := curler.GetUtf8(uri, curl.CurlProxyPolicyUseProxy)
	//	fmt.Println(cache, err)

	if cache.LocalUtf8 != "" {
		fstype := feed.DetectFeedSourceType(cache.LocalUtf8)
		switch fstype {
		case feed.Feed_type_rss:
			v, err = feed.CreateFeedSourceRss2(cache.LocalUtf8)
		case feed.Feed_type_atom:
			v, err = feed.CreateFeedSourceAtom(cache.LocalUtf8)
		}
	}
	if err == nil {
		err = fso.upsert(&v)
	}
	return v, err
}

func feedsource_unsubscribe(url string) error {
	fso := new_feedsource_operator()
	err := fso.drop(url)
	return err
}

func feedcategory_all() ([]string, error) {
	fco := new_feedcategory_operator()
	return fco.all()
}

func meta() (FeedsBackendConfig, error) {
	return backend_config(), nil
}

func meta_cleanup() error {
	// clean temp files
	// entries
	// thumbnails
	// images
	return nil
}

func source_type_map(sourcetype string) uint {
	v, ok := feed.FeedSourceTypes[sourcetype]
	if !ok {
		v = feed.Feed_type_unknown
	}
	return v
}

func feedtag_all() ([]string, error) {
	fto := new_feedtag_operator()
	return fto.all()
}
