package backend

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

type feedcontent_op coll_op

func new_feedcontent_operator() feedcontent_operator {
	return feedcontent_op{coll: "feedcontent_hashs"}
}

func (this feedcontent_op) touch(hash uint64) (v uint, err error) {
	err = do_in_session(this.coll, func(coll *mgo.Collection) error {
		change := mgo.Change{
			ReturnNew: true,
			Upsert:    true,
			Update:    bson.M{"$set": bson.M{"ttl": time.Now()}, "$inc": bson.M{"count": 1}},
		}
		doc := struct {
			Count uint `bson:"count"`
		}{}
		_, err := coll.Find(bson.M{"hash": hash}).Apply(change, &doc)
		log.Println(doc)
		v = doc.Count
		return err
	})
	return
}
