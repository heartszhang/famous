package backend

import (
	"labix.org/v2/mgo"
)

var (
	session *mgo.Session
)

func clone_session() (*mgo.Session, error) {
	if session == nil {
		//		var err error
		if s, err := mgo.Dial(config.DbAddress); err != nil {
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

	c := sess.DB(config.DbName).C(collection)
	return act(c)
}
