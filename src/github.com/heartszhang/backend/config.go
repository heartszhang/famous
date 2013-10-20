package backend

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

type feeds_config struct {
	ip   string
	port uint
}

type FeedsBackendConfig struct {
	Ip   string `json:bind_ip`
	Port uint   `json:port`
}

func (this FeedsBackendConfig) Address() string {
	return fmt.Sprintf("%v:%d", this.Ip, this.Port)
}

type feeds_status struct {
	startat time.Time
}

type FeedsStatus struct {
	Runned int64 `json:"runned"` // seconds
}

func (this feeds_status) runned_nano() int64 {
	return int64(time.Since(this.startat).Seconds())
}

var (
	locker sync.Mutex
	config feeds_config
	status feeds_status

	bind_ip = flag.String("bind_ip", "127.0.0.1", "binding address should be localhost")
	port    = flag.Uint("port", 8002, "backend working port default 8002")
)

func init() {
	flag.Parse()
	config.ip = *bind_ip
	config.port = *port
	status.startat = time.Now()
}

func BackendConfig() FeedsBackendConfig {
	locker.Lock()
	defer locker.Unlock()
	return FeedsBackendConfig{Ip: config.ip, Port: config.port}
}

func BackendStatus() FeedsStatus {
	locker.Lock()
	defer locker.Unlock()
	return FeedsStatus{Runned: status.runned_nano()}
}
