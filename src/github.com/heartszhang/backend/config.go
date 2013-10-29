package backend

import (
	"fmt"
	feed "github.com/heartszhang/feedfeed"
	"os"
	"sync"
	"time"
)

type FeedsBackendConfig struct {
	Ip            string              `json:"web_ip"`
	Port          uint                `json:"port"`
	DbAddress     string              `json:"db_address"` // ip:port or ip
	DbName        string              `json:"db_name"`
	Categories    []feed.FeedCategory `json:"categories,omitempty"`
	DataDir       string              `json:"data_dir,omitempty"` //absolute
	Usage         uint64              `json:"usage"`              //bytes
	ImageDir      string              `json:"image,omitempty"`    //absolute
	ThumbnailDir  string              `json:"thumbnail,omitempty"`
	DocumentDir   string              `json:"document,omitempty"` //absolute
	FeedSourceDir string              `json:"feed_source,omitmepty"`
	FeedEntryDir  string              `json:"feed_entry,omitempty"`
	Proxy         string              `json:"proxy, omitempty"` // "127.0.0.1:8087"

	SummaryThreshold uint `json:"summary_threshold" bson:"summary_threshuld"`
	ThumbnailWidth   uint `json:"thumbnail_width"`
}

func init() {
	config.Ip = "127.0.0.1"
	config.Port = 8002
	config.DbAddress = "127.0.0.1"
	config.DbName = "backend"
	config.DataDir = "data/"
	config.ImageDir = config.DataDir + "images/"
	config.ThumbnailDir = config.DataDir + "thumbnails/"
	config.DocumentDir = config.DataDir + "fulltext/"
	config.FeedSourceDir = config.DataDir + "sources/"
	config.FeedEntryDir = config.DataDir + "entries/"
	config.SummaryThreshold = 250
	//	config.Categories = make([]feed.FeedCategory, 0)
	config.ThumbnailWidth = 320
	os.MkdirAll(config.ImageDir, 0644)
	os.MkdirAll(config.DocumentDir, 0644)
	os.MkdirAll(config.FeedSourceDir, 0644)
	os.MkdirAll(config.FeedEntryDir, 0644)
	status.startat = time.Now()
}

func (this FeedsBackendConfig) Address() string {
	return fmt.Sprintf("%v:%d", this.Ip, this.Port)
}

type FeedsStatus struct {
	startat time.Time `json:"-"`
	Runned  int64     `json:"runned"` // seconds
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
