package backend

import (
	"github.com/heartszhang/feed"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type image_op coll_op

func new_imagecache_operator() imagecache_operator {
	return image_op{coll: "image_cache"}
}

func (this image_op) find(uri string) (v feed.FeedImage, err error) {
	err = do_in_session(this.coll, func(coll *mgo.Collection) error {
		err := coll.Find(bson.M{"uri": uri}).One(&v)
		return err
	})
	return
}

func (this image_op) save(uri string, v feed.FeedImage) error {
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		_, err := coll.Upsert(bson.M{"uri": uri}, bson.M{"$setOnInsert": &v})
		return err
	})
	return err
}
