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
	feed_status_fulltext_ready uint = 1 << iota
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
	Uri          string `json:"uri,omitempty"`           // url
	Alias        string `json:"alias,omitempty"`         // title may be
	Local        string `json:"local,omitempty"`         // downloaded origin html
	CleanedLocal string `json:"cleaned_local,omitempty"` // cleaned-doc local rel path
	Words        int    `json:"words"`                   // words after cleaned
	Sentences    int    `json:"sentences"`               // sentences after cleaned
	Links        int    `json:"links"`                   // links after cleaned
	Density      int    `json:"density"`                 // density of original doc
	Length       int64  `json:"length"`
	Readable     bool   `json:"readable"` // cleaned doc has perfect content

	Images []FeedImage `json:"images,omitempty"`
	Video  []FeedMedia `json:"videos,omitempty"`
	Audio  []FeedMedia `json:"audios,omitempty"`
}

type FeedMedia struct {
	media_type  uint
	Title       string `json:"title,omitempty"`
	Description string `json:"desc,omitempty"`
	Uri         string `json:"uri,omitempty"`   // original url
	Local       string `json:"local,omitempty"` // image : download rel path, video : extraced flv/mp4 url
	Width       int    `json:"width"`           // -1 :unknown
	Height      int    `json:"height"`          // -1 : unknown
	Length      int64  `json:"length"`
	Duration    int    `json:"duration"` // seconds, only for vidoe/audio
	Mime        string `json:"mime,omitempty"`
}

type FeedImage FeedMedia

type FeedAuthor struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Id    uint64 `json:"id"` // for tweet, weibo etc
}

type FeedTitle struct {
	Main   string   `json:"main,omitempty"`   // primary title
	Others []string `json:"second,omitempty"` // secondary or alternative titles, not including main
}

type FeedContent struct {
	Uri     string      `json:"uri"`
	Local   string      `json:"local"`
	Words   uint        `json:"words"`
	Density uint        `json:"density"`
	Links   uint        `json:"links"`
	Status  uint64      `json:"status"`
	Images  []FeedImage `json:"images"`
	Media   []FeedMedia `json:"media"`
}

type FeedEntry struct {
	Id     string `json:"_id"`
	Flags  uint   `json:"flags"`
	Source string `json:"src,omitempty"` // source's uri
	Type   uint   `json:"type"`          // feed_type...

	Uri      string      `json:"uri, omitempty"`
	Title    FeedTitle   `json:"title,omitempty"`
	Author   FeedAuthor  `json:"author,omitempty"`
	PubDate  int64       `json:"pubdate"` // unix time
	Summary  string      `json:"summary,omitempty"`
	Content  FeedContent `json:"content,omitempty"`
	Tags     []string    `json:"tags,omitempty"`
	Images   []FeedImage `json:"images,omitempty"`
	Video    []FeedMedia `json:"videos,omitempty"`
	Audio    []FeedMedia `json:"audios,omitempty"`
	Links    []FeedLink  `json:"links,omitempty"`
	Words    uint        `json:"words"`
	Density  uint        `json:"density"` // percent
	Status   uint64      `json:"status"`
	Category uint64      `json:"category"`
}

const (
	uint64_bits = 64
)

type FeedSource struct {
	Name        string    `json:"name, omitempty"`
	Uri         string    `json:"uri,omitempty"` // rss/atom url
	Local       string    `json:"local"`
	Period      uint      `json:"period"`   // hour
	TouchAt     int64     `json:"touch_at"` // unixtime_nano
	Category    uint64    `json:"category"` //categories
	Type        uint      `json:"type"`     // feed_type...
	Disabled    bool      `json:"disabled"` //auto refresh enabled
	EnableProxy bool      `json:"enable_proxy"`
	PubDate     int64     `json:"pubdate"`           // the last time, we refreshed, unix-time
	WebSite     string    `json:"website,omitempty"` // home
	Tags        []string  `json:"tags,omitempty"`
	Media       FeedMedia `json:"media,omitmepty"`
	Description string    `json:"description, omitempty"`
}

type FeedCategory struct {
	Mask uint64 `json:"mask"`
	Name string `json:"name"`
}

type FeedsProfile struct {
	Categories   []FeedCategory `json:"categories,omitempty"`
	DataDir      string         `json:"data_dir,omitempty"`     //absolute
	Usage        uint64         `json:"usage"`                  //bytes
	ImageDir     string         `json:"image_dir,omitempty"`    //absolute
	DocumentDir  string         `json:"document_dir,omitempty"` //absolute
	Proxy        string         `json:"proxy, omitempty"`       // "127.0.0.1:8087"
	CategoryMask uint64         `json:"category_mask"`          // masked all categories
}

var ( // global status
	profile = FeedsProfile{}
)

func init() {
}

// need not lock
func feedsprofile() FeedsProfile {
	return profile
}

func (this FeedsProfile) content_dir() string {
	return this.DataDir + "content/"
}
