package session

import (
	"os"
	"strings"
	"sync"

	"crypto/tls"
	"crypto/x509"
	"net"
	"time"

	"github.com/jqatampa/gadget-arm/errors"
)

var sessions = make(map[string]*mgo.Session)
var mutex = &sync.Mutex{}

func Get(connectionVariable string, cert ...string) *mgo.Session {
	if sessions[connectionVariable] == nil {
		mutex.Lock()
		defer mutex.Unlock()
		if sessions[connectionVariable] == nil {
			var cs string
			var ssl bool

			if strings.HasPrefix(connectionVariable, "mongodb://") {
				cs = connectionVariable
			} else {
				cs = os.Getenv(connectionVariable)
			}

			if strings.Contains(connectionVariable, "ssl=true") {
				connectionVariable = strings.Replace(connectionVariable, "ssl=true", "", -1)
				connectionVariable = strings.Replace(connectionVariable, "?&", "?", -1)
				connectionVariable = strings.Replace(connectionVariable, "&&", "&", -1)
				ssl = true
			}

			var err error

			var session *mgo.Session
			if cert != nil || ssl {
				session, err = dialWithSSL(cs, cert)
			} else {
				session, err = mgo.Dial(cs)
			}

			errors.Check(err)

			session.SetSocketTimeout(10 * time.Second)
			session.SetSyncTimeout(10 * time.Second)

			// http://godoc.org/labix.org/v2/mgo#Session.SetMode
			session.SetMode(mgo.Monotonic, true)
			sessions[connectionVariable] = session
		}
	} else {
		sessions[connectionVariable].Refresh()
	}

	return sessions[connectionVariable].Copy()
}

func dialWithSSL(cs, certs []string) (session *mgo.Session, err error) {
	tlsConfig := &tls.Config{}

	if certs != nil {
		roots := x509.NewCertPool()
		roots.AppendCertsFromPEM([]byte(certs[0]))
		tlsConfig.RootCAs = roots
	} else {
		tlsConfig.InsecureSkipVerify = true
	}

	dialInfo, err := mgo.ParseURL(cs)
	if err != nil {
		return
	}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	session, err = mgo.DialWithInfo(dialInfo)
	return
}
