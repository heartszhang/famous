package backend

import (
	"fmt"
	"github.com/heartszhang/cleaner"
	"github.com/heartszhang/curl"
	feed "github.com/heartszhang/feedfeed"
	google "github.com/heartszhang/googlefeedservice"
)

// since_unixtime , 0: from now
// categories, categories mask, every bit represent a category
// count: entries per page
// page: page no, start at 0
func feeds_entries_since(since_unixtime int64, category string, count uint, page uint) ([]feed.FeedEntry, error) {
	return []feed.FeedEntry{}, nil
}

func feedentry_unread(source string, count uint, page uint) ([]feed.FeedEntry, error, int) {
	c := curl.NewCurl(backend_config().FeedEntryFolder)
	cache, err := c.GetUtf8(source)

	if err != nil || cache.LocalUtf8 == "" {
		return nil, err, cache.StatusCode
	}
	//atom+xml;xml;html
	ext := curl.MimeToExt(cache.Mime)
	if ext != "xml" && ext != "atom+xml" && ext != "rss+xml" {
		return nil, fmt.Errorf("unsupported mime: %v", cache.Mime), 0
	}
	v, err := feed.MakeFeedEntries(cache.LocalUtf8)

	v = feed_entries_unreaded(v) // clean readed entries
	v = feed_entries_clean(v)
	v = feed_entries_clean_summary(v)
	v = feed_entries_clean_fulltext(v)
	v = feed_entries_autotag(v)
	v = feed_entries_statis(v)
	v = feed_entries_backup(v)
	return v, err, cache.StatusCode
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

// /feed/entry/full_text.json/{url}/{entry_id}
func feedentry_fulltext(uri string, entry_uri string) (v feed.FeedLink, err error) {
	c := curl.NewCurl(config.DocumentFolder)
	cache, err := c.GetUtf8(uri)
	v.Uri = uri
	v.Local = cache.LocalUtf8
	v.Length = cache.LengthUtf8
	if err != nil {
		return v, err
	}
	if curl.MimeToExt(cache.Mime) != "html" {
		return v, fmt.Errorf("unsupported mime %v", cache.Mime)
	}
	doc, err := html_create_from_file(cache.LocalUtf8)
	if err != nil {
		return v, err
	}
	article, sum, err := cleaner.NewExtractor(config.CleanFolder).MakeHtmlReadable(doc, uri)
	v.Images = make([]feed.FeedMedia, len(sum.Images))
	for idx, img := range sum.Images {
		v.Images[idx].Uri = img.Uri
		v.Images[idx].Width = int(img.Width)
		v.Images[idx].Height = int(img.Height)
		v.Images[idx].Description = img.Alt
	}
	v.Words = sum.WordCount
	v.Links = sum.LinkCount
	v.CleanedLocal, err = html_write_file(article, config.DocumentFolder)
	return v, err
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
	curler := curl.NewCurl(backend_config().FeedSourceFolder)
	cache, err := curler.GetUtf8(uri)
	ext := curl.MimeToExt(cache.Mime)
	if ext != "xml" && ext != "atom+xml" && ext != "rss+atom" {
		return v, fmt.Errorf("unsupported mime: %v", cache.Mime)
	}

	if cache.LocalUtf8 != "" {
		v, err = feed.MakeFeedSource(cache.LocalUtf8)
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

const (
	refer = "http://iweizhi2.duapp.com"
)

func feedsource_show(uri string) ([]feed.FeedSourceFindEntry, error) {
	fs, _, err := google.NewGoogleFeedApi(refer, config.FeedSourceFolder).Load(uri, config.Language, 4, false)
	if err != nil {
		return nil, err
	}
	v := make([]feed.FeedSourceFindEntry, 1)
	v[0].Summary = fs.Description
	v[0].Title = fs.Name
	v[0].Url = fs.Uri
	v[0].Website = fs.WebSite
	return v, err
}

func feedsource_find(q string) ([]feed.FeedSourceFindEntry, error) {
	svc := google.NewGoogleFeedApi(refer, config.FeedSourceFolder)
	v, err := svc.Find(q, config.Language)
	if err != nil {
		return v, err
	}
	uris := make([]string, 0)
	for _, ve := range v {
		uris = append(uris, ve.Url)
	}
	dbo := new_feedsource_operator()
	subed, err := dbo.findbatch(uris)
	if err != nil {
		return v, err
	}
	for _, fs := range subed {
		for i := 0; i < len(v); i++ {
			if v[i].Url == fs.Uri {
				v[i].Subscribed = true
			}
		}
	}
	return v, err
}
