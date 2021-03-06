package backend

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func new_feedsource_operator() feedsource_operator {
	return feedsource_op{coll: "feed_sources"}
}

type feedsource_op coll_op

func (this feedsource_op) save(feeds []ReadSource) (inserted []ReadSource, err error) {
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
func (this feedsource_op) upsert(f ReadSource) error {
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		_, err := coll.Upsert(bson.M{"uri": f.Uri}, bson.M{"$setOnInsert": &f})
		return err
	})
}

func (this feedsource_op) drop(uri string) error {
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Remove(bson.M{"uri": uri})
	})
}

func (this feedsource_op) set_subscribe_state(uri string, s int) error {
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"uri": uri}, bson.M{"$set": bson.M{"subscribe_state": s}})
	})
}

func (this feedsource_op) save_one(f ReadSource) error {
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"uri": f.Uri}, bson.M{"$set": &f})
	})
}

func (this feedsource_op) update(f ReadSource) error {
	val := bson.M{
		"name":        f.Name,
		"period":      f.Period,
		"update":      f.Update,
		"last_touch":  f.LastTouch,
		"last_update": f.LastUpdate,
		"next_touch":  f.NextTouch,
	}

	if f.WebSite != "" {
		val["website"] = f.WebSite
	}
	if f.Hub != "" {
		val["hub"] = f.Hub
	}
	if f.Logo != "" {
		val["logo"] = f.Logo
	}
	return do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Update(bson.M{"uri": f.Uri}, bson.M{"$set": val})
	})
}

func (this feedsource_op) find(uri string) (ReadSource, error) {
	var rtn ReadSource
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		err := coll.Find(bson.M{"uri": uri}).One(&rtn)
		return err
	})
	return rtn, err
}

func (this feedsource_op) findbatch(uris []string) ([]ReadSource, error) {
	var rtn []ReadSource
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		err := coll.Find(bson.M{"uri": bson.M{"$in": uris}}).All(&rtn)
		return err
	})
	return rtn, err
}

func (this feedsource_op) expired(beforeunxtime int64) ([]ReadSource, error) {
	var rtn []ReadSource
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		err := coll.Find(bson.M{"next_touch": bson.M{"$lt": beforeunxtime}}).All(&rtn)
		return err
	})
	return rtn, err
}
func (this feedsource_op) all() (feds []ReadSource, err error) {
	err = do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Find(bson.M{"uri": bson.M{"$ne": ""}}).Sort("-last_touch").All(&feds)
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
