package backend

import (
	feed "github.com/heartszhang/feedfeed"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

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
		feed.FeedEntry `bson:",inline"`
		TTL            time.Time `bson:"ttl"`
	}{entry, time.Now()}
	ci, err := coll.Upsert(bson.M{"uri": entry.Uri}, bson.M{"$setOnInsert": xe})
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
