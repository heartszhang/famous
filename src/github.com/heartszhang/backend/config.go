package backend

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type FeedsBackendConfig struct {
	Ip            string         `json:"web_ip"`
	Port          uint           `json:"port"`
	DbAddress     string         `json:"db_address"` // ip:port or ip
	DbName        string         `json:"db_name"`
	Categories    []FeedCategory `json:"categories,omitempty"`
	DataDir       string         `json:"data_dir,omitempty"` //absolute
	Usage         uint64         `json:"usage"`              //bytes
	ImageDir      string         `json:"image,omitempty"`    //absolute
	DocumentDir   string         `json:"document,omitempty"` //absolute
	FeedSourceDir string         `json:"feed_source,omitmepty"`
	FeedEntryDir  string         `json:"feed_entry,omitempty"`
	Proxy         string         `json:"proxy, omitempty"` // "127.0.0.1:8087"
	CategoryMask  uint64         `json:"category_mask"`    // masked all categories}
}

func init() {
	config.Ip = "127.0.0.1"
	config.Port = 8002
	config.DbAddress = "127.0.0.1"
	config.DbName = "backend"
	config.DataDir = "data/"
	config.ImageDir = config.DataDir + "image/"
	config.DocumentDir = config.DataDir + "fulltext/"
	config.FeedSourceDir = config.DataDir + "sources/"
	config.FeedEntryDir = config.DataDir + "entries/"
	config.Categories = make([]FeedCategory, 0)
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

func BackendConfig() FeedsBackendConfig {
	locker.Lock()
	defer locker.Unlock()
	return config
}

func BackendStatus() FeedsStatus {
	locker.Lock()
	defer locker.Unlock()
	return FeedsStatus{Runned: status.runned_nano()}
}
