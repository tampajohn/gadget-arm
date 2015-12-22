package session

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/jqatampa/gadget-arm/errors"
	"gopkg.in/mgo.v2"
)

var sessions = make(map[string]*mgo.Session)
var mutex = &sync.Mutex{}

func Get(connectionVariable string) *mgo.Session {
	if sessions[connectionVariable] == nil {
		mutex.Lock()
		if sessions[connectionVariable] == nil {
			var cs string

			if strings.HasPrefix(connectionVariable, "mongodb://") {
				cs = connectionVariable
			} else {
				cs = os.Getenv(connectionVariable)
			}

			var err error

			session, err := mgo.Dial(cs)

			errors.Check(err)

			// http://godoc.org/labix.org/v2/mgo#Session.SetMode
			session.SetMode(mgo.Monotonic, true)
			sessions[connectionVariable] = session
		}
		mutex.Unlock()
	}

	if err := sessions[connectionVariable].Ping(); err != nil {
		sessions[connectionVariable].Refresh()
		log.Printf("Refreshing session: %v", sessions[connectionVariable])
	}

	return sessions[connectionVariable].Copy()
}
