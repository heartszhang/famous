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
		if s, err := mgo.Dial(backend_context.config.DbAddress); err != nil {
			return nil, err
		} else {
			session = s
		}
	}
	return session.Copy(), nil
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

	c := sess.DB(backend_context.config.DbName).C(collection)
	return act(c)
}

type coll_op struct {
	coll string
}
