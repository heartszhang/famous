package backend

import (
	"strconv"
)

import ()

const (
	feed_media_type_url = 1 << iota
	feed_media_type_video
	feed_media_type_audio
	feed_media_type_image
)

type FeedLink struct {
	media_type   int    `json:"-"`             // feed_media_type_url
	Uri          string `json:"uri"`           // url
	Alias        string `json:"alias"`         // title may be
	Local        string `json:"local"`         // downloaded origin html
	CleanedLocal string `json:"cleaned_local"` // cleaned-doc local rel path
	Words        int    `json:"words"`         // words after cleaned
	Sentences    int    `json:"sentences"`     // sentences after cleaned
	Links        int    `json:"links"`         // links after cleaned
	Density      int    `json:"density"`       // density of original doc
	Readable     bool   `json:"readable"`      // cleaned doc has perfect content

	Images []FeedImage `json:"images"`
	Video  []FeedMedia `json:"videos"`
	Audio  []FeedMedia `json:"audios"`
}

type FeedMedia struct {
	media_type  int    `json:"-"` // must be feed_media_type_video or _audio
	Title       string `json:"title"`
	Description string `json:"desc"`
	Uri         string `json:"uri"`      // original url
	Local       string `json:"local"`    // image : download rel path, video : extraced flv/mp4 url
	Width       int    `json:"width"`    // -1 :unknown
	Height      int    `json:"height"`   // -1 : unknown
	Size        int64  `json:""`         // -1 : unknown
	Duration    int    `json:"duration"` // seconds, only for vidoe/audio
}

type FeedImage FeedMedia

type FeedAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Id    uint64 `json:"id"` // for tweet etc
}

type FeedTitle struct {
	Main   string   `json:"main"`   // primary title
	Others []string `json:"second"` // secondary or alternative titles, not including main
}

type FeedEntry struct {
	Id     string `json:"_id"`
	Flags  int    `json:"flags"`
	Source string `json:"src"`  // source's _id
	Type   int    `json:"type"` // feed type, defined by feed_type

	Title   FeedTitle   `json:"title"`
	Author  FeedAuthor  `json:"author"`
	PubDate uint64      `json:"pubdate"` // unix time
	Summary string      `json:"summary"`
	Content string      `json:"content"`
	Tags    []string    `json:"tags"`
	Images  []FeedImage `json:"images"`
	Video   []FeedMedia `json:"videos"`
	Audio   []FeedMedia `json:"audios"`
	Links   []FeedLink  `json:"links"`
	Words   int         `json:"words"`
	Density int         `json:"density"` // percent
	Statis  uint64      `json:"statis"`
}

const (
	uint64_bits = 64
)

// /feeds/entries_since.json/{since_unixtime:[0-9]+}/{category:[0-9]+}/{count:[0-9]+}/{page:[0-9]+}
func web_feeds_entries_since(w http.ResponseWriter, r *http.Request) {
	var (
		since                = unixtime_now()
		category             = feed_category_root
		count    uint        = 20
		page     uint        = 0
		err      error       = nil
		fe       []FeedEntry = make([]FeedEntry, 0)
	)

	if err == nil {
		since, err = strconv.ParseUint(r.URL.Query().Get(":since_unixtime"), 0, uint64_bits)
	}
	if err == nil {
		category, err = strconv.ParseUint(r.URL.Query().Get(":category"), 0, uint64_bits)
	}
	if err == nil {
		count, err = strconv.ParseUint(r.URL.Query().Get(":count"), 0, 0)
	}
	if err != nil {
		page = strconv.ParseUint(r.URL.Query().Get(":page"), 0, 0)
	}
	if err != nil {
		fe, err = feeds_entries_since(since, category, count, page)
	}
	if err != nil {
		write_error(w, err)
	} else {
		webapi_write_as_json(w, fe)
	}
}

// since_unixtime , 0: from now
// categories, categories mask, every bit represent a category
// count: entries per page
// page: page no, start at 0
func feeds_entries_since(since_unixtime int64, categories uint64, count uint, page uint) ([]FeedEntry, error) {
	return nil, nil
}

const (
	feed_flag_readed = 1 << iota
	feed_flag_star
	feed_flag_save
)

// /feed/entry/mark.json/{id}/{flags}

func feed_entry_mark(id string, flags int) (uint, error) {
	return 0, nil
}

// /feed/entry/umark.json/{id}/{flags}
func feed_entry_umark(id string, flags int) (uint, error) {
	return 0, nil
}

// /feed/entry/drop.json/{id}

// id is mongo's _id
func feed_entry_drop(id string) error {
	return nil
}

// select a idle category_id, assigned to category
func feed_category_create(name string) (uint64, error) {
	return 0, nil
}

// id : isn't root or all, drop the category whoes name is name
// id : other, drop categories
// name : can be empty. if id is root or all, name cann't be empty
func feed_category_drop(id uint64, name string) error {
	return nil
}

// /tick.json

func tick() (FeedsStatus, error) {
	s := BackendStatus()
	return s, nil
}

type FeedSource struct {
	Uri         string `json:"uri"`      // rss/atom url
	Period      uint   `json:"period"`   // hour
	Category    uint64 `json:"category"` //categories
	Type        string `json:"type"`     // source -type
	IsActive    bool   `json:"active"`   //auto refresh enabled
	EnableProxy bool   `json:"enable_proxy"`
	PubDate     uint64 `json:"pubdate"` // the last time, we refreshed, unix-time
	WebSite     string `json:"website"` // home
}

// /source/subscribe.json/{uri}/{source_type}/{category}
func web_source_subscribe(w http.ResponseWriter, r *http.Request) {
	var (
		url         string
		source_type       = 0
		category          = feed_category_root
		err         error = nil
	)
	url = r.URL.Query().Get(":uri")
	if err == nil {
		source_type, err = strconv.ParseInt(r.URL.Query().Get(":source_type"), 0, 0)
	}
	if err == nil {
		category, err = strconv.ParseUint(r.URL.Query().Get(":category"), 0, uint64_bits)
	}
	if err == nil {
		fs, err = source_subscribe(url, source_type, category)
	}
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, fs)
	}
}

func source_subscribe(url string, source_type int, category uint64) (FeedSource, error) {
	return FeedSource{}, nil
}

const (
	feed_type_unknown = 0
	feed_type_rss     = 1 << iota
	feed_type_atom
	feed_type_sina_weibo
	feed_type_qq_weibo
	feed_type_blog
	feed_type_tweet
)
const (
	feed_source_rss   = "rss"
	feed_source_atom  = "atom"
	feed_source_blog  = "blog"
	feed_source_tweet = "tweet"
)

// /source/unsubscribe.json/{uri}/{source_type}/{category}
func source_unsubscribe(url string, source_type string, category uint64) error {
	return nil
}

const (
	feed_category_root uint64 = 0
	feed_category_all  uint64 = 1<<63 - 1
)

type FeedCategory struct {
	Mask uint64 `json:"mask"`
	Name string `json:"name"`
}

// /meta/categories.json

func meta_categories() ([]FeedCategory, error) {
	return nil, nil
}

type FeedsMeta struct {
	Categories  []FeedCategory `json:"categories"`
	DataDir     string         `json:"data_dir"`     //absolute
	Usage       int64          `json:"usage"`        //bytes
	ImageDir    string         `json:"image_dir"`    //absolute
	DocumentDir string         `json:"document_dir"` //absolute
}

// /meta.json
func web_meta(w http.ResponseWriter, r *http.Request) {
}

func meta() (FeedsMeta, error) {
	return FeedsMeta{}, nil
}

// /meta/cleanup.json
func meta_cleanup() error {
	return nil
}
