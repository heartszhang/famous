package backend

import (
	feed "github.com/heartszhang/feedfeed"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type feedsource_operator interface {
	save(feeds []feed.FeedSource) ([]feed.FeedSource, error)
	upsert(f *feed.FeedSource) error
	find(uri string) (*feed.FeedSource, error)
	expired() ([]feed.FeedSource, error)
	all() ([]feed.FeedSource, error)
	touch(uri string, ttl int) error
	drop(uri string) error
	disable(uri string, dis bool) error
	update(f *feed.FeedSource) error
}

func new_feedsource_operator() feedsource_operator {
	return &feedsource_op{coll: "feed_sources"}
}

type coll_op struct {
	coll string
}

type feedsource_op coll_op

func (this feedsource_op) save(feeds []feed.FeedSource) (inserted []feed.FeedSource, err error) {
	inserted = make([]feed.FeedSource, 0)
	err = do_in_session(this.coll, func(coll *mgo.Collection) error {
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

func (this feedsource_op) upsert(f *feed.FeedSource) error {
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		_, err := coll.Upsert(bson.M{"uri": f.Uri}, bson.M{"$setOnInsert": f})
		return err
	})
}

func (this feedsource_op) drop(uri string) error {
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Remove(bson.M{"uri": uri})
	})
}

func (this feedsource_op) disable(uri string, dis bool) error {
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"uri": uri}, bson.M{"$set": bson.M{"disabled": dis}})
	})
}

func (this feedsource_op) update(f *feed.FeedSource) error {
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"uri": f.Uri},
			bson.M{
				"$set":      bson.M{"period": f.Period},
				"$addToSet": bson.M{"tags": bson.M{"$each": f.Tags}}})
	})
}

func (this feedsource_op) find(uri string) (*feed.FeedSource, error) {
	rtn := new(feed.FeedSource)
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		err := coll.Find(bson.M{"uri": uri}).One(rtn)
		return err
	})
	return rtn, err
}

func (this feedsource_op) all() (feds []feed.FeedSource, err error) {
	feds = make([]feed.FeedSource, 0)
	err = do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Find(bson.M{"disabled": false}).All(&feds)
	})
	return
}
func (this feedsource_op) expired() ([]feed.FeedSource, error) {
	rtn := make([]feed.FeedSource, 0)
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Find(bson.M{"disabled": false, "due_at": bson.M{"$lt": feed.UnixTimeNow()}}).All(&rtn)
	})
	return rtn, err
}

func (this feedsource_op) touch(uri string, ttl int) error {
	dl := time.Now().Add(time.Duration(ttl) * time.Minute)
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"uri": uri}, bson.M{"$set": bson.M{"due_at": dl}})
	})
}

type feedcategory_operator interface {
	save(cate string) (interface{}, error)
	all() ([]string, error)
	drop(category string) error
}
type feedcategory_op coll_op
type feedtag_operator feedcategory_operator

func new_feedtag_operator() feedtag_operator {
	return feedcategory_op{coll: "feed_tags"}
}

func new_feedcategory_operator() feedcategory_operator {
	return feedcategory_op{coll: "feed_categories"}
}

func (this feedcategory_op) all() ([]string, error) {
	var v []struct {
		Name string `bson:"name"`
	}
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Find(nil).All(&v)
	})
	x := make([]string, len(v))
	for i, c := range v {
		x[i] = c.Name
	}
	return x, err
}

func (this feedcategory_op) drop(cate string) error {
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Remove(bson.M{"name": cate})
	})
}

func (this feedcategory_op) save(cate string) (uid interface{}, err error) {
	err = do_in_session(this.coll, func(coll *mgo.Collection) error {
		ci, err := coll.Upsert(bson.M{"name": cate}, bson.M{"$setOnInsert": bson.M{"name": cate}})
		uid = ci.UpsertedId
		return err
	})
	return
}

type feedentry_operator interface {
	save([]feed.FeedEntry) ([]feed.FeedEntry, error)
	save_one(feed.FeedEntry) (interface{}, error)
	topn(skip, limit int) ([]feed.FeedEntry, error)
	topn_by_category(skip, limit int, category string) ([]feed.FeedEntry, error)
	topn_by_feedsource(skip, limit int, feed string) ([]feed.FeedEntry, error)
	mark(link string, newmark uint) error
	umark(uri string, markbit uint) error
	setcontent(link string, filepath string, words int, imgs []feed.FeedMedia) error
}

type feedentry_op coll_op

func new_feedentry_operator() feedentry_operator {
	return feedentry_op{coll: "feed_entries"}
}

func (this feedentry_op) mark(uri string, newmark uint) error {
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		//	return coll.Update(bson.M{"uri": uri}, bson.M{"$set": bson.M{"flags": newmark}})
		return coll.Update(bson.M{"uri": uri}, bson.M{"$bit": bson.M{"flags": bson.M{"$or": newmark}}})
	})
}

func (this feedentry_op) umark(uri string, markbit uint) error {
	mask := ^markbit
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"uri": uri}, bson.M{"$bit": bson.M{"flags": bson.M{"$and": mask}}})
	})
}

func (this feedentry_op) save(entries []feed.FeedEntry) ([]feed.FeedEntry, error) {
	inserted := make([]feed.FeedEntry, 0)
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
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

func (this feedentry_op) save_one(entry feed.FeedEntry) (uid interface{}, err error) {
	do_in_session(this.coll, func(coll *mgo.Collection) error {
		uid, err = insert_entry(coll, entry)
		return err
	})
	return uid, err
}

func (this feedentry_op) topn(skip, limit int) ([]feed.FeedEntry, error) {
	rtn := make([]feed.FeedEntry, 0)
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Find(bson.M{"readed": false}).Sort("-created").Skip(skip).Limit(limit).All(&rtn)
	})
	return rtn, err
}

func (this feedentry_op) topn_by_feedsource(skip, limit int, source string) ([]feed.FeedEntry, error) {
	rtn := make([]feed.FeedEntry, 0)
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Find(bson.M{"readed": false, "source": source}).Sort("-created").Skip(skip).Limit(limit).All(&rtn)
	})
	return rtn, err
}

func (this feedentry_op) topn_by_category(skip, limit int, tag string) ([]feed.FeedEntry, error) {
	rtn := make([]feed.FeedEntry, 0)
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Find(bson.M{"readed": false, "tags": tag}).
			Sort("-created").
			Skip(skip).
			Limit(limit).
			All(rtn)
	})
	return rtn, err
}

func insert_entry(coll *mgo.Collection, entry feed.FeedEntry) (interface{}, error) {
	xe := struct {
		feed.FeedEntry
		TTL time.Time
	}{entry, time.Now()}
	ci, err := coll.Upsert(bson.M{"uri": entry.Uri}, bson.M{"$setOnInsert": &xe})
	return ci.UpsertedId, err
}

func (this feedentry_op) setcontent(uri, filepath string, words int, imgs []feed.FeedMedia) error {
	status := feed.Feed_content_failed
	imgc := len(imgs)
	if len(filepath) > 0 && (words+imgc*128) > 192 {
		status = feed.Feed_content_ready
	}
	cs := feed.FeedContent{Uri: uri, Local: filepath, Words: uint(words), Status: status, Images: imgs}

	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"uri": uri}, bson.M{"$set": bson.M{
			"status":  status,
			"content": cs,
		}, "$push": bson.M{"images": bson.M{"$each": imgs}}})
	})
}
