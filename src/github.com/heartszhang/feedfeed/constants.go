package feedfeed

import ()

const (
	Feed_media_type_none    uint = 0
	Feed_media_type_unknown      = 1 << iota
	Feed_media_type_url
	Feed_media_type_video
	Feed_media_type_audio
	Feed_media_type_image

	Feed_media_type_media = Feed_media_type_video | Feed_media_type_audio
)

const (
	Feed_flag_none   uint = 0
	Feed_flag_readed      = 1 << iota
	Feed_flag_star
	Feed_flag_save
)

const (
	Feed_content_ready uint64 = 1 << iota
	Feed_content_empty
	Feed_content_inline
	Feed_content_external_ready
	Feed_content_external_empty
	Feed_status_has_audio
	Feed_status_has_video
	Feed_status_has_url
	Feed_status_has_image
	Feed_status_invisible
	Feed_status_text_empty
	Feed_status_text_little
	Feed_status_text_many
	Feed_status_image_empty
	Feed_status_image_one
	Feed_status_image_many
	Feed_status_media_empty // image, audio , video
	Feed_status_media_one
	Feed_status_media_many
	Feed_status_media_inline
	Feed_status_linkdensity_low
	Feed_status_linkdensity_high
	Feed_status_format_flowdocument
	Feed_status_format_text
	Feed_status_mp4
	Feed_status_flv
	Feed_content_unresolved
	Feed_summary_ready
	Feed_summary_empty
	Feed_content_unavail
)

const (
	Feed_category_root uint64 = 0
	Feed_category_none        = Feed_category_root
	Feed_category_all         = 1<<64 - 1 //math.MaxUint64
)

const (
	Feed_type_unknown uint = 0
	Feed_type_rss          = 1 << iota // ignore rss version
	Feed_type_atom                     // ignore atom version
	Feed_type_sina_weibo
	Feed_type_qq_weibo
	Feed_type_blog
	Feed_type_tweet

	Feed_type_feed = Feed_type_rss | Feed_type_atom
)

type FeedLink struct {
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
	//	TTL    time.Time   `json:"ttl" bson:"ttl"`
}

type FeedMedia struct {
	media_type  uint
	Title       string `json:"title,omitempty" bson:"title,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	Uri         string `json:"uri,omitempty" bson:"uri,omitempty"`     // original url
	Local       string `json:"local,omitempty" bson:"local,omitempty"` // image : download rel path, video : extraced flv/mp4 url
	Width       int    `json:"width" bson:"width"`                     // -1 :unknown
	Height      int    `json:"height" bson:"height"`                   // -1 : unknown
	Length      int64  `json:"length" bson:"length"`
	Duration    int    `json:"duration" bson:"duration"` // seconds, only for vidoe/audio
	Mime        string `json:"mime,omitempty" bson:"mime,omitempty"`
	Thumbnail   string `json:"thumbnail,omitempty" bson:"thumbnail,omitempty"`
}

type FeedAuthor struct {
	Name  string `json:"name,omitempty" bson:"name,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Id    uint64 `json:"id" bson:"id"` // for tweet, weibo etc
}

type FeedTitle struct {
	Main   string   `json:"main,omitempty" bson:"main,omitempty"`     // primary title
	Others []string `json:"second,omitempty" bson:"second,omitempty"` // secondary or alternative titles, not including main
}

type FeedContent struct {
	Uri      string      `json:"uri" bson:"uri"`
	Local    string      `json:"local" bson:"local"`
	FullText string      `json:"-" bson:"-"`
	Words    uint        `json:"words" bson:"words"`
	Density  uint        `json:"density" bson:"density"`
	Links    uint        `json:"links" bson:"links"`
	Status   uint64      `json:"status" bson:"status"`
	Images   []FeedMedia `json:"images" bson:"images"`
	Medias   []FeedMedia `json:"media" bson:"media"`
}

type FeedEntry struct {
	Id         string       `bson:"_id,omitempty" json:"id,omitempty"`
	Flags      uint         `json:"flags" bson:"flags"`
	Source     string       `json:"src,omitempty" bson:"src,omitempty"` // source's uri
	Type       uint         `json:"type" bson:"type"`                   // feed_type...
	Uri        string       `json:"uri,omitempty" bson:"uri,omitempty"`
	Title      FeedTitle    `json:"title,omitempty" bson:"title,omitempty"`
	Author     *FeedAuthor  `json:"author,omitempty" bson:"author,omitempty"`
	PubDate    int64        `json:"pubdate" bson:"pubdate"` // unix time
	Summary    string       `json:"summary,omitempty" bson:"summary,omitempty"`
	Content    *FeedContent `json:"content,omitempty" bson:"content,omitempty"`
	Tags       []string     `json:"tags,omitempty" bson:"tags,omitempty"`
	Images     []FeedMedia  `json:"images,omitempty" bson:"images,omitempty"`
	Videos     []FeedMedia  `json:"videos,omitempty" bson:"videos,omitempty"`
	Audios     []FeedMedia  `json:"audios,omitempty" bson:"audios,omitempty"`
	Links      []FeedLink   `json:"links,omitempty" bson:"links,omitempty"`
	Words      uint         `json:"words" bson:"words"`     // of sumary
	Density    uint         `json:"density" bson:"density"` // percent
	Status     uint64       `json:"status" bson:"status"`
	Categories []string     `json:"categories,omitempty" bson:"category,omitempty"`
	//	TTL        time.Time    `json:"ttl" bson:"ttl"`
}

const (
	uint64_bits = 64
)

type FeedSource struct {
	Name        string     `json:"name,omitempty" bson:"name,omitempty"`
	Uri         string     `json:"uri,omitempty" bson:"uri,omitempty"` // rss/atom url
	Local       string     `json:"local" bson:"local"`
	Period      uint       `json:"period" bson:"period"`     // minutes
	Deadline    int64      `json:"deadline" bson:"deadline"` // unixtime_nano
	Type        uint       `json:"type" bson:"type"`         // feed_type...
	Disabled    bool       `json:"disabled" bson:"disabled"` //auto refresh enabled
	EnableProxy bool       `json:"enable_proxy" bson:"enable_proxy"`
	Update      int64      `json:"update" bson:"update"`                       // the last time, we refreshed, unix-time
	WebSite     string     `json:"website,omitempty" bson:"website,omitempty"` // home
	Tags        []string   `json:"tags,omitempty" bson:"tags,omitempty"`
	Categories  []string   `json:"categories,omitempty" bson:"categories,omitempty"`
	Unreaded    int        `json:"unreaded" bson:"unreaded"`
	Media       *FeedMedia `json:"media,omitempty" bson:"media,omitempty"`
	Description string     `json:"description,omitempty" bson:"description,omitempty"`
}

type FeedCategory struct {
	Mask uint64 `json:"mask" bson:"mask"`
	Name string `json:"name" bson:"name"`
}

type FeedTag struct {
	Value string `json:"value,omitempty" bson:"value,omitempty"`
	//	TTL   time.Time `json:"ttl" bson:"ttl"`
}

var (
	FeedSourceTypes = map[string]uint{
		"":         Feed_type_unknown,
		"rss":      Feed_type_rss,
		"atom":     Feed_type_atom,
		"atom10":   Feed_type_atom,
		"rss+xml":  Feed_type_rss,
		"atom+xml": Feed_type_atom,
		"feed":     Feed_type_atom,
		"blog":     Feed_type_blog,
		"tweet":    Feed_type_tweet,
		"weibo":    Feed_type_sina_weibo,
		"qqweibo":  Feed_type_qq_weibo}
)

type FeedImage struct {
	Mime           string `json:"mime,omitempty" bson:"mime,omitempty"`
	ThumbnailLocal string `json:"thumbnail,omitempty" bson:"thumbnail,omitempty"`
	OriginLocal    string `json:"origin,omitempty" bson:"origin,omitempty"`
	Code           int    `json:"code" bson:"code"`
	Width          int    `json:"width" bson:"width"`
	Height         int    `json:"height" bson:"height"`
}

type FeedSourceFindEntry struct {
	Url     string `json:"url,omitempty`
	Title   string `json:"title,omitempty`
	Summary string `json:"summary,omitempty`
	Website string `json:"website,omitempty"`
}
