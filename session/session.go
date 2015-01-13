package session

import (
	"os"
	"strings"

	"bitbucket.org/tampajohn/gadget-arm/errors"
	"gopkg.in/mgo.v2"
)

var (
	mgoSession *mgo.Session
)

func Get(connectionVariable string) *mgo.Session {
	if mgoSession == nil {
		var cs string

		if strings.HasPrefix(connectionVariable, "mongodb://") {
			cs = connectionVariable
		} else {
			cs = os.Getenv(connectionVariable)
		}

		var err error
		mgoSession, err = mgo.Dial(cs)

		// http://godoc.org/labix.org/v2/mgo#Session.SetMode
		mgoSession.SetMode(mgo.Monotonic, true)

		errors.Check(err)
	}
	return mgoSession.Clone()
}
