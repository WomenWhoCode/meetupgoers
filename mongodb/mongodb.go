package mongodb

import (
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"
)

const (
	MongoDBHosts = os.Getenv("DBHost")
	AuthDatabase = os.Getenv("DBName")
	AuthUserName = os.Getenv("DBUserName")
	AuthPassword = os.Getenv("DBPassword")
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
