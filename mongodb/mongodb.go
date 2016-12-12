package mongodb

import (
	"time"

	mgo "gopkg.in/mgo.v2"
)

const (
	MongoDBHosts = "ds033036.mlab.com:33036"
	AuthDatabase = "heroku_vb4zpgmk"
	AuthUserName = ""
	AuthPassword = ""
)

func ConnDB() *mgo.Session {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  600 * time.Second,
		Database: AuthDatabase,
		Username: AuthUserName,
		Password: AuthPassword,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	session, err := mgo.DialWithInfo(mongoDBDialInfo)

	if err != nil {
		panic(err)
	}
	return session
}

//err = c.Find(bson.M{"name": "Ale"}).Select(bson.M{"phone": 0}).One(&result)
//if err != nil {
//	panic(err)
//}
