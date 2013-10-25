package backend

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type FeedSourceOperator interface {
	Save(feeds []FeedSource) ([]FeedSource, error)
	Upsert(f *FeedSource) error
	Find(uri string) (*FeedSource, error)
	TimeoutSources() ([]FeedSource, error)
	AllSources() ([]FeedSource, error)
	Touch(uri string, ttl int) error
	Drop(uri string) error
	Disable(uri string, dis bool) error
	Update(f *FeedSource) error
}

func NewFeedSourceOperator() FeedSourceOperator {
	return &feedsource_op{}
}

type feedsource_op struct{}

func (feedsource_op) Save(feeds []FeedSource) (inserted []FeedSource, err error) {
	inserted = make([]FeedSource, 0)
	err = do_in_session("feedsources", func(coll *mgo.Collection) error {
		for _, f := range feeds {
			ci, err := coll.Upsert(bson.M{"uri": f.Uri}, bson.M{"$setOnInsert": f})
			if err != nil {
				return err
			}
			if ci.UpsertedId != nil {
				inserted = append(inserted, f)
			}
		}
		return nil
	})
	return
}

func (feedsource_op) Upsert(f *FeedSource) error {
	return do_in_session("feedsources", func(coll *mgo.Collection) error {
		_, err := coll.Upsert(bson.M{"uri": f.Uri}, bson.M{"$setOnInsert": f})
		return err
	})
}

func (feedsource_op) Drop(uri string) error {
	return do_in_session("feedsources", func(coll *mgo.Collection) error {
		return coll.Remove(bson.M{"uri": uri})
	})
}

func (feedsource_op) Disable(uri string, dis bool) error {
	return do_in_session("feedsources", func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"uri": uri}, bson.M{"$set": bson.M{"disabled": dis}})
	})
}

func (feedsource_op) Update(f *FeedSource) error {
	return do_in_session("feedsources", func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"uri": f.Uri},
			bson.M{
				"$set":      bson.M{"period": f.Period},
				"$addToSet": bson.M{"tags": bson.M{"$each": f.Tags}}})
	})
}

func (feedsource_op) Find(uri string) (*FeedSource, error) {
	rtn := new(FeedSource)
	err := do_in_session("feedsources", func(coll *mgo.Collection) error {
		err := coll.Find(bson.M{"uri": uri}).One(rtn)
		return err
	})
	return rtn, err
}

func (feedsource_op) AllSources() (feds []FeedSource, err error) {
	feds = make([]FeedSource, 0)
	err = do_in_session("feedsources", func(coll *mgo.Collection) error {
		return coll.Find(bson.M{"disabled": false}).All(&feds)
	})
	return
}
func (feedsource_op) TimeoutSources() ([]FeedSource, error) {
	rtn := make([]FeedSource, 0)
	err := do_in_session("feedsources", func(coll *mgo.Collection) error {
		return coll.Find(bson.M{"disabled": false, "due_at": bson.M{"$lt": unixtime_now()}}).All(&rtn)
	})
	return rtn, err
}

func (feedsource_op) Touch(uri string, ttl int) error {
	dl := time.Now().Add(time.Duration(ttl) * time.Minute)
	return do_in_session("feedsources", func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"uri": uri}, bson.M{"$set": bson.M{"due_at": dl}})
	})
}

type FeedEntryOperator interface {
	Save([]FeedEntry) ([]FeedEntry, error)
	SaveOne(FeedEntry) (interface{}, error)
	TopN(skip, limit int) ([]FeedEntry, error)
	TopNByCategory(skip, limit int, category string) ([]FeedEntry, error)
	TopNByFeed(skip, limit int, feed string) ([]FeedEntry, error)
	MarkRead(link string, readed bool) error
	SetContent(link string, filepath string, words int, imgs []FeedImage) error
}

func NewFeedEntryOperator() FeedEntryOperator {
	return new(feedentry_op)
}

type feedentry_op struct {
}

func (feedentry_op) MarkRead(uri string, readed bool) error {
	return do_in_session("entries", func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"uri": uri}, bson.M{"$set": bson.M{"readed": readed}})
	})
}

func (feedentry_op) Save(entries []FeedEntry) ([]FeedEntry, error) {
	inserted := make([]FeedEntry, 0)
	err := do_in_session("entries", func(coll *mgo.Collection) error {
		for _, entry := range entries {
			iid, err := insert_entry(coll, entry)
			if err != nil {
				return err
			}
			if iid != nil {
				inserted = append(inserted, entry)
			}
		}
		return nil
	})
	return inserted, err
}

func (feedentry_op) SaveOne(entry FeedEntry) (uid interface{}, err error) {
	do_in_session("entries", func(coll *mgo.Collection) error {
		uid, err = insert_entry(coll, entry)
		return err
	})
	return uid, err
}

func (feedentry_op) TopN(skip, limit int) ([]FeedEntry, error) {
	rtn := make([]FeedEntry, 0)
	err := do_in_session("entries", func(coll *mgo.Collection) error {
		return coll.Find(bson.M{"readed": false}).Sort("-created").Skip(skip).Limit(limit).All(&rtn)
	})
	return rtn, err
}

func (feedentry_op) TopNByFeed(skip, limit int, feed string) ([]FeedEntry, error) {
	rtn := make([]FeedEntry, 0)
	err := do_in_session("entries", func(coll *mgo.Collection) error {
		return coll.Find(bson.M{"readed": false, "feed": feed}).Sort("-created").Skip(skip).Limit(limit).All(&rtn)
	})
	return rtn, err
}

func (feedentry_op) TopNByCategory(skip, limit int, tag string) ([]FeedEntry, error) {
	rtn := make([]FeedEntry, 0)
	err := do_in_session("entries", func(coll *mgo.Collection) error {
		return coll.Find(bson.M{"readed": false, "tags": tag}).
			Sort("-created").
			Skip(skip).
			Limit(limit).
			All(rtn)
	})
	return rtn, err
}

func insert_entry(coll *mgo.Collection, entry FeedEntry) (interface{}, error) {
	ci, err := coll.Upsert(bson.M{"uri": entry.Uri}, bson.M{"$setOnInsert": &entry})
	return ci.UpsertedId, err
}

func (feedentry_op) SetContent(uri, filepath string, words int, imgs []FeedImage) error {
	status := feed_content_failed
	imgc := len(imgs)
	if len(filepath) > 0 && (words+imgc*128) > 192 {
		status = feed_content_ready
	}
	cs := FeedContent{Uri: uri, Local: filepath, Words: uint(words), Status: status, Images: imgs}

	return do_in_session("entries", func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"uri": uri}, bson.M{"$set": bson.M{
			"status":  status,
			"content": cs,
		}, "$push": bson.M{"images": bson.M{"$each": imgs}}})
	})
}
