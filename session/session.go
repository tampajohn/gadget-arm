package session

import (
	"bitbucket.org/tampajohn/gadget-arm/errors"
	"gopkg.in/mgo.v2"
	"os"
)

var (
	mgoSession *mgo.Session
)

func Get(connectionVariable string) *mgo.Session {
	if mgoSession == nil {
		cs := os.Getenv(connectionVariable)
		var err error
		mgoSession, err = mgo.Dial(cs)
		errors.Check(err)
	}
	return mgoSession.Clone()
}
