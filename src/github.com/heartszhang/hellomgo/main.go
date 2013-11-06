package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	session *mgo.Session
)

func clone_session() (*mgo.Session, error) {
	if session == nil {
		//		var err error
		if s, err := mgo.Dial("127.0.0.1"); err != nil {
			return nil, err
		} else {
			session = s
		}
	}
	return session.Clone(), nil
}

func close_session() {
	if session != nil {
		sess := session
		session = nil
		sess.Close()
	}
}

func do_in_session(collection string, act func(*mgo.Collection) error) error {
	sess, err := clone_session()
	if err != nil {
		return err
	}
	defer sess.Close()

	c := sess.DB("test").C(collection)
	return act(c)
}

type ImageCache struct {
	Width  int `bson:"width"`
	Height int `bson:"height"`
}
type ImageCacheW struct {
	uri        string `bson:"uri"`
	ImageCache `bson:",inline"`
}

func main() {
	v := ImageCache{Width: 1, Height: 2}
	wv := ImageCacheW{"what", v}
	err := do_in_session("image_cache", func(coll *mgo.Collection) error {
		u, err := coll.Upsert(bson.M{"uri": "what"}, bson.M{"$setOnInsert": wv})
		fmt.Println(u)
		return err
	})
	fmt.Println(err)
}
