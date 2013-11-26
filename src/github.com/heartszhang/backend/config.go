package backend

import (
	"fmt"
	feed "github.com/heartszhang/feedfeed"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

type FeedsBackendConfig struct {
	BackendIp             string `json:"web_ip"`
	BackendPort           uint   `json:"port"`
	DbAddress             string `json:"db_address"` // ip:port or ip
	DbName                string `json:"db_name"`
	CacheUsage            uint64 `json:"usage"`              //bytes
	DataFolder            string `json:"data_dir,omitempty"` //absolute
	ImageFolder           string `json:"image,omitempty"`    //absolute
	ThumbnailFolder       string `json:"thumbnail,omitempty"`
	DocumentFolder        string `json:"document,omitempty"` //absolute
	FeedSourceFolder      string `json:"feed_source,omitmepty"`
	FeedEntryFolder       string `json:"feed_entry,omitempty"`
	CleanFolder           string `json:"clean,omitempty"`
	FailedFolder          string `json:"fails,omitempty"`
	FlowDocumentFolder    string `json:"flowdocs,omitempty"`
	ProxyAddress          string `json:"proxy, omitempty"` // "127.0.0.1:8087"
	SummaryThreshold      uint   `json:"summary_threshold" bson:"summary_threshuld"`
	SummaryMinWords       int    `json:"summary_minwords" bson:"summary_minwords"`
	ThumbnailWidth        uint   `json:"thumbnail_width" bson:"thumbnail_width"`
	SummaryDuplicateCount uint   `json:"summary_duplicatecount" bson:"summary_duplicatecount"`
	LinkDensityThreshuld  uint   `json:"linkdensity_threshold" bson:"linkdensity_threshuld"`
	Language              string `json:"language,omitempty" bson:"language,omitempty"`
}

func pwd() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}
func init() {
	config.BackendIp = "127.0.0.1"
	config.BackendPort = 8002
	config.DbAddress = "127.0.0.1"
	config.DbName = "backend"
	config.DataFolder = filepath.Join(pwd(), "data/")
	config.ImageFolder = filepath.Join(config.DataFolder, "images/")
	config.ThumbnailFolder = filepath.Join(config.DataFolder, "thumbnails/")
	config.FeedEntryFolder = filepath.Join(config.DataFolder, "entries/")
	config.DocumentFolder = filepath.Join(config.DataFolder, "fulltexts/")
	config.FeedSourceFolder = filepath.Join(config.DataFolder, "sources/")
	config.CleanFolder = filepath.Join(config.DataFolder, "cleans/")
	config.FailedFolder = filepath.Join(config.DataFolder, "clean_fails/")
	config.FlowDocumentFolder = filepath.Join(config.DataFolder, "flowdocs/")
	config.SummaryThreshold = 250
	config.SummaryMinWords = 25
	config.SummaryDuplicateCount = 3

	config.ThumbnailWidth = 400
	os.Mkdir(config.DataFolder, 0644)
	os.Mkdir(config.ImageFolder, 0644)
	os.Mkdir(config.ThumbnailFolder, 0644)
	os.Mkdir(config.FeedEntryFolder, 0644)
	os.Mkdir(config.DocumentFolder, 0644)
	os.Mkdir(config.FeedSourceFolder, 0644)
	os.Mkdir(config.CleanFolder, 0644)
	os.Mkdir(config.FailedFolder, 0644)
	os.Mkdir(config.FlowDocumentFolder, 0644)
	status.startat = time.Now()

	config.Language = "zh-CN"
}

func (this FeedsBackendConfig) Address() string {
	return fmt.Sprintf("%v:%d", this.BackendIp, this.BackendPort)
}

type feed_status struct {
	startat time.Time `json:"-"`
	Error   string    `json:"error,omitempty"`
}
type feed_tick struct {
	Tick  int64         `json:"tick"`
	Feeds []feed_update `json:"feeds,omitempty"`
}

func (this feed_status) runned_nano() int64 {
	return int64(time.Since(this.startat).Seconds())
}

var (
	locker       sync.Mutex
	config       FeedsBackendConfig
	status       feed_status
	working      int64
	feed_updates = make([]feed_update, 0)
)

func backend_config() FeedsBackendConfig {
	locker.Lock()
	defer locker.Unlock()
	return config
}

const (
	update_period = 5 * 60 //seconds
)

func backend_tick() feed_tick {
	locker.Lock()
	defer locker.Unlock()
	u := feed_updates
	feed_updates = nil
	r := status.runned_nano()
	if r > update_period {
		go update_work()
	}
	return feed_tick{Tick: r, Feeds: u}
}

func backend_push_update(fs feed.FeedSource, fes []feed.FeedEntry, err error) {
	if err != nil {
		return
	}
	// locker has been locked
	feed_updates = append(feed_updates, feed_update{fs, fes})
}

func update_work() {
	locker.Lock()
	defer locker.Unlock()
	w := atomic.AddInt64(&working, 1)
	defer atomic.AddInt64(&working, -1)
	if w != 1 {
		return
	}
	fss, err := feedsource_expired(time.Now().Unix())
	if err != nil || len(fss) == 0 {
		return
	}
	idx := rand.Intn(len(fss))
	fs := fss[idx]
	newfs, fes, err := feed_fetch(fs.Uri)
	newfs.Type = fs.Type
	newfs.EnableProxy = fs.EnableProxy
	newfs.Categories = append(newfs.Categories, fs.Categories...)
	if newfs.Logo == "" {
		newfs.Logo = fs.Logo
	}

	newfs.Disabled = fs.Disabled
	newfs.LastTouch = time.Now().Unix()
	newfs.LastUpdate = newfs.LastTouch
	newfs.NextTouch = int64(newfs.Period) + newfs.LastTouch
	err = feedsource_save(newfs)
	fes = feedentry_filter(fes)
	backend_push_update(newfs, fes, err)
}
