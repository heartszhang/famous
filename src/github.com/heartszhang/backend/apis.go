package backend

import (
	"github.com/heartszhang/feed"
	"github.com/heartszhang/pubsub"
	"github.com/heartszhang/unixtime"
	"log"
	"math/rand"
	"sync/atomic"
	"time"
)

// /tick.json
func tick() (FeedTick, error) {
	s := backend_tick()
	return s, nil
}

func meta() (FeedsBackendConfig, error) {
	return backend_config(), nil
}

func update_popup() (*FeedEntity, error) {
	fs, fes, err := feedentries_updated()
	if err == nil {
		return &FeedEntity{FeedSource: *fs, Entries: fes}, nil
	}
	return nil, err
}

func meta_cleanup() error {
	// clean temp files
	// entries
	// thumbnails
	// images
	return nil
}

func source_type_map(sourcetype string) uint {
	v, ok := feed.FeedSourceTypes[sourcetype]
	if !ok {
		v = feed.Feed_type_unknown
	}
	return v
}

func feedtag_all() ([]string, error) {
	fto := new_feedtag_operator()
	return fto.all()
}

func backend_push_update(fs feed.FeedSource, fes []feed.FeedEntry, err error) {
	if err != nil {
		return
	}
	backend_context.feed_updates = append(backend_context.feed_updates, FeedEntity{fs, fes})
}

func update_work() {
	backend_context.Lock()
	defer backend_context.Unlock()
	w := atomic.AddInt64(&backend_context.working, 1)
	defer atomic.AddInt64(&backend_context.working, -1)
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
	newfs.LastTouch = unixtime.UnixTimeNow()
	newfs.LastUpdate = newfs.LastTouch
	newfs.NextTouch = unixtime.UnixTime(newfs.Period) + newfs.LastTouch
	err = feedsource_save(newfs)
	fes = feedentry_filter(fes)
	backend_push_update(newfs, fes, err)
	ps := pubsub.NewSuperFeedrPubSubscriber("async")
	sc, err := ps.Subscribe(fs.Uri)
	if err != nil {
		log.Println("pubsub-google", sc, err)
	}
	ps = pubsub.NewGooglePubSubscriber()
	sc, err = ps.Subscribe(fs.Uri)
	if err != nil {
		log.Println("pubsub-superfeedr", sc, err)
	}
	log.Println("update-tick", fs.Name, err)
}
