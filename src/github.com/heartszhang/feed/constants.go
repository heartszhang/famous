package feed

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
	FeedEntry_flag_none   uint = 0
	FeedEntry_flag_readed      = 1 << iota
	FeedEntry_flag_star
	FeedEntry_flag_save
)

const (
	Feed_status_text_empty uint64 = 1 << iota
	Feed_status_text_little
	Feed_status_text_many
	Feed_status_image_empty
	Feed_status_image_one
	Feed_status_image_many
	Feed_status_media_empty // image, audio , video
	Feed_status_media_one
	Feed_status_media_many
	Feed_status_linkdensity_low
	Feed_status_linkdensity_high
	Feed_status_format_flowdocument
	Feed_status_format_text
	Feed_status_mp4
	Feed_status_flv
	Feed_status_content_ready
	Feed_status_content_empty
	Feed_status_content_inline
	Feed_status_content_external_ready
	Feed_status_content_external_empty
	Feed_status_content_unresolved
	Feed_status_content_unavail
	Feed_status_content_duplicated
	Feed_status_content_mediainline
	Feed_status_summary_ready
	Feed_status_summary_empty
	Feed_status_summary_inline
	Feed_status_summary_external_ready
	Feed_status_summary_external_empty
	Feed_status_summary_unresolved
	Feed_status_summary_unavail
	Feed_status_summary_duplicated
	Feed_status_summary_mediainline
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
	Uri       string `json:"uri,omitempty" bson:"uri,omitempty"`     // url
	Alias     string `json:"alias,omitempty" bson:"alias,omitempty"` // title may be
	Local     string `json:"local,omitempty" bson:"local,omitempty"` // downloaded origin html
	Words     int    `json:"words" bson:"words"`                     // words after cleaned
	Sentences int    `json:"sentences" bson:"sentences"`             // sentences after cleaned
	Links     int    `json:"links" bson:"links"`                     // links after cleaned
	Density   int    `json:"density" bson:"density"`                 // density of original doc
	Length    int64  `json:"length" bson:"length"`
}

type FeedMedia struct {
	Title       string `json:"title,omitempty" bson:"title,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	Uri         string `json:"uri,omitempty" bson:"uri,omitempty"` // original url
	Width       int    `json:"width" bson:"width"`                 // -1 :unknown
	Height      int    `json:"height" bson:"height"`               // -1 : unknown
	Length      int64  `json:"length" bson:"length"`               // bytes
	Duration    int    `json:"duration" bson:"duration"`           // seconds, only for vidoe/audio
	Thumbnail   string `json:"thumbnail,omitempty" bson:"thumbnail,omitempty"`
	Mime        string `json:"mime,omitempty" bson:"mime,omitempty"`
}

type FeedContent struct {
	Uri     string      `json:"uri" bson:"uri"`
	FlowDoc string      `json:"doc,omitempty" bson:"doc,omitempty"`
	Local   string      `json:"local,omitempty" bson:"local,omitempty"`
	Words   uint        `json:"words" bson:"words"`
	Density uint        `json:"density" bson:"density"`
	Links   uint        `json:"links" bson:"links"`
	Status  uint64      `json:"status" bson:"status"`
	Images  []FeedMedia `json:"images,omitempty" bson:"images,omitempty"`
	Medias  []FeedMedia `json:"media,omitempty" bson:"media,omitempty"` // videos or audios
}

type FeedTextStatus struct {
	WordCount     int    `json:"wordcount" bson:"wordcount"`
	LinkWordCount int    `json:"link_wordcount" bson:"link_wordcount"`
	LinkCount     int    `json:"linkcount" bson:"linkcount"`
	Status        uint64 `json:"status" bson:"status"`
}

var (
	feedsource_types = map[string]uint{
		"":         Feed_type_unknown,
		"rss":      Feed_type_rss,
		"rss20":    Feed_type_rss,
		"atom":     Feed_type_atom,
		"atom10":   Feed_type_atom,
		"rss+xml":  Feed_type_rss,
		"atom+xml": Feed_type_atom,
		"feed":     Feed_type_atom,
		"blog":     Feed_type_blog,
		"tweet":    Feed_type_tweet,
		"weibo":    Feed_type_sina_weibo,
		"qqweibo":  Feed_type_qq_weibo,
	}
)

func FeedSourceType(typ string) uint {
	if v, ok := feedsource_types[typ]; ok {
		return v
	}
	return Feed_type_unknown
}

type FeedImage struct {
	Mime      string `json:"mime,omitempty" bson:"mime,omitempty"`
	Uri       string `json:"uri,omitempty" bson:"uri,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty" bson:"thumbnail,omitempty"` // local cache
	Origin    string `json:"origin,omitempty" bson:"origin,omitempty"`       // local cache
	Width     int    `json:"width" bson:"width"`
	Height    int    `json:"height" bson:"height"`
}

const (
	link_rel_self      = "self"
	link_rel_related   = "related"
	link_rel_alternate = "alternate"
	link_rel_enclosure = "enclosure"
	link_rel_via       = "via"
	link_rel_hub       = "hub"
	link_rel_icon      = "icon"
)
