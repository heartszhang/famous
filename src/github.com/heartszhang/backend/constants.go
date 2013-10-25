package backend

const (
	feed_media_type_none    uint = 0
	feed_media_type_unknown      = 1 << iota
	feed_media_type_url
	feed_media_type_video
	feed_media_type_audio
	feed_media_type_image

	feed_media_type_media = feed_media_type_video | feed_media_type_audio
)

const (
	feed_flag_none   uint = 0
	feed_flag_readed      = 1 << iota
	feed_flag_star
	feed_flag_save
)

const (
	feed_status_fulltext_ready uint64 = 1 << iota
	feed_status_thumbnail_ready
	feed_status_has_image
	feed_status_has_audio
	feed_status_has_video
	feed_status_has_url
	feed_status_invisible
	feed_status_text_only
	feed_status_image_only
	feed_status_image_gallery_only
	feed_status_video_only
	feed_status_audio_only
	feed_status_mp4
	feed_status_flv
	feed_content_unresolved
	feed_content_ready
	feed_content_failed
	feed_content_unavail
	feed_content_summary
)

const (
	feed_category_root uint64 = 0
	feed_category_none        = feed_category_root
	feed_category_all         = 1<<64 - 1 //math.MaxUint64
)

const (
	feed_type_unknown uint = 0
	feed_type_rss          = 1 << iota // ignore rss version
	feed_type_atom                     // ignore atom version
	feed_type_sina_weibo
	feed_type_qq_weibo
	feed_type_blog
	feed_type_tweet

	feed_type_feed = feed_type_rss | feed_type_atom
)

type FeedLink struct {
	media_type   uint   // feed_media_type...
	Uri          string `uri,omitempty`           // url
	Alias        string `alias,omitempty`         // title may be
	Local        string `local,omitempty`         // downloaded origin html
	CleanedLocal string `cleaned_local,omitempty` // cleaned-doc local rel path
	Words        int    `words`                   // words after cleaned
	Sentences    int    `sentences`               // sentences after cleaned
	Links        int    `links`                   // links after cleaned
	Density      int    `density`                 // density of original doc
	Length       int64  `length`
	Readable     bool   `readable` // cleaned doc has perfect content

	Images []FeedImage `images,omitempty`
	Videos []FeedMedia `videos,omitempty`
	Audios []FeedMedia `audios,omitempty`
}

type FeedMedia struct {
	media_type  uint
	Title       string `title,omitempty`
	Description string `desc,omitempty`
	Uri         string `uri,omitempty`   // original url
	Local       string `local,omitempty` // image : download rel path, video : extraced flv/mp4 url
	Width       int    `width`           // -1 :unknown
	Height      int    `height`          // -1 : unknown
	Length      int64  `length`
	Duration    int    `duration` // seconds, only for vidoe/audio
	Mime        string `mime,omitempty`
}

type FeedImage FeedMedia

type FeedAuthor struct {
	Name  string `name,omitempty`
	Email string `email,omitempty`
	Id    uint64 `id` // for tweet, weibo etc
}

type FeedTitle struct {
	Main   string   `main,omitempty`   // primary title
	Others []string `second,omitempty` // secondary or alternative titles, not including main
}

type FeedContent struct {
	Uri      string      `uri`
	Local    string      `local`
	FullText string      `-`
	Words    uint        `words`
	Density  uint        `density`
	Links    uint        `links`
	Status   uint64      `status`
	Images   []FeedImage `images`
	Medias   []FeedMedia `media`
}

type FeedEntry struct {
	Id       string      `bson:"_id" json:"id"`
	Flags    uint        `flags`
	Source   string      `src,omitempty` // source's uri
	Type     uint        `type`          // feed_type...
	Uri      string      `uri, omitempty`
	Title    FeedTitle   `title,omitempty`
	Author   FeedAuthor  `author,omitempty`
	PubDate  int64       `pubdate` // unix time
	Summary  string      `summary,omitempty`
	Content  FeedContent `content,omitempty`
	Tags     []string    `tags,omitempty`
	Images   []FeedImage `images,omitempty`
	Videos   []FeedMedia `videos,omitempty`
	Audios   []FeedMedia `audios,omitempty`
	Links    []FeedLink  `links,omitempty`
	Readed   bool        `readed`
	Words    uint        `words`   // of sumary
	Density  uint        `density` // percent
	Status   uint64      `status`
	Category uint64      `category`
}

const (
	uint64_bits = 64
)

type FeedSource struct {
	Name        string    `name, omitempty`
	Uri         string    `uri,omitempty` // rss/atom url
	Local       string    `local`
	Period      uint      `period`   // minutes
	DueAt       int64     `due_at`   // unixtime_nano
	Category    uint64    `category` //categories
	Type        uint      `type`     // feed_type...
	Disabled    bool      `disabled` //auto refresh enabled
	EnableProxy bool      `enable_proxy`
	Update      int64     `update`            // the last time, we refreshed, unix-time
	WebSite     string    `website,omitempty` // home
	Tags        []string  `tags,omitempty`
	Media       FeedMedia `media,omitmepty`
	Description string    `description, omitempty`
}

type FeedCategory struct {
	Mask uint64 `mask`
	Name string `name`
}
