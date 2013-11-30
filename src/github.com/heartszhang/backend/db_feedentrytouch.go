package backend

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type feedentrytouch_op coll_op

func new_feedentrytouch_operator() feedentrytouch_operator {
	return feedentrytouch_op{"feedentry_hashs"}
}

func (this feedentrytouch_op) touch(hashes []string) ([]string, error) {
	var rtn []string
	err := do_in_session(this.coll, func(coll *mgo.Collection) error {
		for _, hash := range hashes {
			ci, err := coll.Upsert(bson.M{"hash": hash}, bson.M{"$setOnInsert": bson.M{"hash": hash, "ttl": time.Now()}})
			if ci != nil {
				rtn = append(rtn, hash)
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
	return rtn, err
}
