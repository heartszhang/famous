package backend

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type feedcategory_op coll_op
type feedtag_op coll_op

func new_feedtag_operator() feedtag_operator {
	return feedtag_op{coll: "feed_sources"}
}

func (this feedtag_op) all() ([]string, error) {
	var v []struct {
		Tags []string `bson:"tags"`
	}
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		return coll.Find(nil).All(&v)
	})
	var r []string
	for _, i := range v {
		r = append(r, i.Tags...)
	}
	return r, err
}
func (this feedtag_op) drop(tag string) error {
	return nil
}
func (this feedtag_op) save(tag string) (interface{}, error) {
	return nil, nil
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
