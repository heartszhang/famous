package backend

import (
	feed "github.com/heartszhang/feedfeed"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	//	"time"
)

func new_feedsource_operator() feedsource_operator {
	return feedsource_op{coll: "feed_sources"}
}

type feedsource_op coll_op

func (this feedsource_op) save(feeds []feed.FeedSource) (inserted []feed.FeedSource, err error) {
	inserted = make([]feed.FeedSource, 0)
	err = do_in_session(this.coll, func(coll *mgo.Collection) error {
		for _, f := range feeds {
			if f.Uri == "" {
				return db_error("feedsource.uri cannot be null")
			}
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

func (this feedsource_op) save_one(f feed.FeedSource) error {
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"uri": f.Uri}, bson.M{"$set": &f})
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

func (this feedsource_op) findbatch(uris []string) ([]feed.FeedSource, error) {
	rtn := make([]feed.FeedSource, 0)
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		err := coll.Find(bson.M{"uri": bson.M{"$in": uris}}).All(&rtn)
		return err
	})
	return rtn, err
}

func (this feedsource_op) expired(beforeunxtime int64) ([]feed.FeedSource, error) {
	var rtn []feed.FeedSource
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		err := coll.Find(bson.M{"next_touch": bson.M{"$lt": beforeunxtime}}).All(&rtn)
		return err
	})
	return rtn, err
}
func (this feedsource_op) all() (feds []feed.FeedSource, err error) {
	feds = make([]feed.FeedSource, 0)
	err = do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Find(bson.M{"disabled": false, "uri": bson.M{"$ne": ""}}).Sort("-last_touch").All(&feds)
	})
	return
}

func (this feedsource_op) touch(uri string, last, next, period int64) error {
	//	dl := time.Now().Add(time.Duration(ttl) * time.Minute)
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"uri": uri},
			bson.M{
				"$set": bson.M{
					"period":      period,
					"last_touch":  last,
					"last_update": last,
					"next_touch":  next,
				}})
		//		return coll.Update(bson.M{"uri": uri}, bson.M{"$set": bson.M{"due_at": dl}})
	})
}

type db_error string

func (this db_error) Error() string {
	return string(this)
}
