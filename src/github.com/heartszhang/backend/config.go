package backend

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
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

type FeedsStatus struct {
	startat time.Time `json:"-"`
	Runned  int64     `json:"runned"` // seconds
	Error   string    `json:"error,omitempty"`
}

func (this FeedsStatus) runned_nano() int64 {
	return int64(time.Since(this.startat).Seconds())
}

var (
	locker sync.Mutex
	config FeedsBackendConfig
	status FeedsStatus
)

func backend_config() FeedsBackendConfig {
	locker.Lock()
	defer locker.Unlock()
	return config
}

func BackendConfig() FeedsBackendConfig {
	return backend_config()
}

func backend_status() FeedsStatus {
	locker.Lock()
	defer locker.Unlock()
	return FeedsStatus{Runned: status.runned_nano()}
}
