package mongodb

import (
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"
)

func ConnDB() *mgo.Session {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{os.Getenv("DBHost")},
		Timeout:  600 * time.Second,
		Database: os.Getenv("DBName"),
		Username: os.Getenv("DBUserName"),
		Password: os.Getenv("DBPassword"),
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	session, err := mgo.DialWithInfo(mongoDBDialInfo)

	if err != nil {
		panic(err)
	}
	return session
}
