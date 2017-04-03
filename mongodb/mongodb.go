package mongodb

import (
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"os"
	"os"
)

var (
	dbHost     string
	dbName     string
	dbUser     string
	dbPassword string
)

func setUp() {
	dbHost = os.Getenv("DBHost")
	dbName = os.Getenv("DBName")
	dbUser = os.Getenv("DBUser")
	dbPassword = os.Getenv("DBPassword")
	if dbHost == "" {
		log.Fatal("$DB Host related info must be set")
	}
	if dbName == "" {
		log.Fatal("$DB Name related info must be set")
	}
}

func ConnDB() *mgo.Session {
	setUp()
	dBDialInfo := &mgo.DialInfo{
		Addrs:    []string{dbHost},
		Timeout:  600 * time.Second,
		Database: dbName,
		Username: dbUser,
		Password: dbPassword,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	session, err := mgo.DialWithInfo(dBDialInfo)

	if err != nil {
		panic(err)
	}
	return session
}

func RenameCollection(src string, dst string, session *mgo.Session) {
	var dat interface{}
	err := session.DB(dbName).C(dst).DropCollection()
	if err != nil {
		println(err)
	}
	err = session.Run(bson.D{{"renameCollection", dbName+"."+src}, {"to", dbName+"."+dst}}, dat)
	if err != nil {
		panic(err)
	}
}
