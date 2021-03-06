package backend

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	"github.com/heartszhang/gfwlist"
	"github.com/qiniu/log"
)

type FeedTick struct {
	Tick  int64        `json:"tick"`
	Feeds []ReadEntity `json:"feeds,omitempty"`
}
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
	EntryDefaultPageCount int64  `json:"entry_defaultpagecount" bson:"entry_defaultpagecount"`
}

func pwd() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}
func init() {
	config := &backend_context.config
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
	config.SummaryDuplicateCount = 5

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

	config.Language = "zh-CN"
	config.EntryDefaultPageCount = 10
	backend_context.startup = time.Now()

	load_gfwrules()
}
func load_gfwrules() {
	fp := filepath.Join(backend_config().DataFolder, "gfwlist.txt")
	reader, err := os.Open(fp)
	if err == nil {
		defer reader.Close()
		backend_context.ruler, err = gfwlist.NewGfwRuler(reader)
	}
	log.Println("gfwlist-load", err, fp)
}
func (this FeedsBackendConfig) Address() string {
	return fmt.Sprintf("%v:%d", this.BackendIp, this.BackendPort)
}

var backend_context struct {
	sync.Mutex
	config       FeedsBackendConfig
	startup      time.Time
	working      int64
	feed_updates []ReadEntity
	ruler        gfwlist.GfwRuler
}

func BackendConfig() FeedsBackendConfig {
	return backend_config()
}

func backend_config() FeedsBackendConfig {
	//	backend_context.Lock()
	//	defer backend_context.Unlock()
	return backend_context.config
}

const (
	update_period = 5 * 60 //seconds
)

func backend_tick() FeedTick {
	backend_context.Lock()
	defer backend_context.Unlock()
	u := backend_context.feed_updates
	backend_context.feed_updates = nil
	r := time.Since(backend_context.startup).Nanoseconds() / int64(time.Second)
	if r > update_period {
		go update_work()
	}
	return FeedTick{Tick: r, Feeds: u}
}
