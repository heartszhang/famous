package backend

import (
	"time"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type feedentry_op coll_op

func new_feedentry_operator() feedentry_operator {
	return feedentry_op{coll: "feed_entries"}
}

func (this feedentry_op) mark(uri string, newmark uint) error {
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"uri": uri}, bson.M{"$bit": bson.M{"flags": bson.M{"$or": newmark}}})
	})
}

func (this feedentry_op) umark(uri string, markbit uint) error {
	mask := ^markbit
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"uri": uri}, bson.M{"$bit": bson.M{"flags": bson.M{"$and": mask}}})
	})
}

func (this feedentry_op) umark_category(category string, markbit uint) error {
	mask := ^markbit
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"categories": category}, bson.M{"$bit": bson.M{"flags": bson.M{"$and": mask}}})
	})
}

func (this feedentry_op) mark_category(category string, newmark uint) error {
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"categories": category}, bson.M{"$bit": bson.M{"flags": bson.M{"$or": newmark}}})
	})
}

func (this feedentry_op) umark_source(source string, markbit uint) error {
	mask := ^markbit
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"src": source}, bson.M{"$bit": bson.M{"flags": bson.M{"$and": mask}}})
	})
}

func (this feedentry_op) mark_source(source string, newmark uint) error {
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"src": source}, bson.M{"$bit": bson.M{"flags": bson.M{"$or": newmark}}})
	})
}

func (this feedentry_op) save(entries []ReadEntry) ([]ReadEntry, error) {
	var inserted []ReadEntry
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		for _, entry := range entries {
			iid, err := insert_entry(coll, entry)
			if iid != nil {
				inserted = append(inserted, entry)
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
	return inserted, err
}

func (this feedentry_op) save_one(entry ReadEntry) (uid interface{}, err error) {
	do_in_session(this.coll, func(coll *mgo.Collection) error {
		uid, err = insert_entry(coll, entry)
		return err
	})
	return uid, err
}

func (this feedentry_op) unread_count(source string) (int, error) {
	var rtn int
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		var err error
		rtn, err = coll.Find(bson.M{"src": source, "readed": false}).Count()
		return err
	})
	return rtn, err
}

func (this feedentry_op) unread_count_category(cate string) (int, error) {
	var rtn int
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		var err error
		rtn, err = coll.Find(bson.M{"categories": cate, "readed": false}).Count()
		return err
	})
	return rtn, err
}

func (this feedentry_op) unread_count_sources() (v []feedentry_unreadcount, err error) {
	err = do_in_session(this.coll, func(coll *mgo.Collection) error {
		job := &mgo.MapReduce{
			Map:    "function() { emit(this.src, 1) }",
			Reduce: "function(key, values) { return Array.sum(values) }",
		}
		_, err := coll.Find(bson.M{"readed": false}).MapReduce(job, &v)
		return err
	})
	return
}

func (this feedentry_op) unread_count_categories() (v []feedentry_unreadcount, err error) {
	err = do_in_session(this.coll, func(coll *mgo.Collection) error {
		job := &mgo.MapReduce{
			Map:    `function() { for(var c in this.categories){emit(c,1);} }`,
			Reduce: "function(key, values) { return Array.sum(values) }",
		}
		_, err := coll.Find(bson.M{"readed": false}).MapReduce(job, &v)
		return err
	})
	return
}
func (this feedentry_op) topn(skip, limit int) ([]ReadEntry, error) {
	var rtn []ReadEntry
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Find(bson.M{"readed": false}).Sort("-pubdate").Skip(skip).Limit(limit).All(&rtn)
	})
	return rtn, err
}

func (this feedentry_op) topn_by_feedsource(skip, limit int, source string) ([]ReadEntry, error) {
	var rtn []ReadEntry
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Find(bson.M{"readed": false, "src": source}).Sort("-pubdate").Skip(skip).Limit(limit).All(&rtn)
	})
	return rtn, err
}

func (this feedentry_op) topn_by_category(skip, limit int, tag string) ([]ReadEntry, error) {
	var rtn []ReadEntry
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Find(bson.M{"readed": false, "tags": tag}).
			Sort("-pubdate").
			Skip(skip).
			Limit(limit).
			All(rtn)
	})
	return rtn, err
}

func insert_entry(coll *mgo.Collection, entry ReadEntry) (interface{}, error) {
	xe := struct {
		ReadEntry `bson:",inline"`
		TTL       time.Time `bson:"ttl"`
	}{entry, time.Now()}
	ci, err := coll.Upsert(bson.M{"uri": entry.Uri}, bson.M{"$setOnInsert": &xe})
	if ci == nil {
		return nil, err
	}
	return ci.UpsertedId, err
}

/*
func (this feedentry_op) setcontent(uri, filepath string, words int, imgs []feed.FeedMedia) error {
	status := feed.Feed_status_content_unresolved
	imgc := len(imgs)
	if len(filepath) > 0 && (words+imgc*128) > 192 {
		status = feed.Feed_status_content_ready
	}
	cs := feed.FeedContent{Uri: uri, Local: filepath, Words: uint(words), Status: status, Images: imgs}

	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"uri": uri}, bson.M{"$set": bson.M{
			"status":  status,
			"content": cs,
		}, "$push": bson.M{"images": bson.M{"$each": imgs}}})
	})
}
*/
