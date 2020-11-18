package module

import (
	"fmt"
	"log"
	"github.com/dlintw/goconf"
	mgo "gopkg.in/mgo.v2"
)

var (
	MONGO_DB_OTT_NAME 	string
	MONGO_DB_OTT_MODE 	string
	mongoSession 		*mgo.Session

	CommonConfig *goconf.ConfigFile
)

func init() {

	var err error
	CommonConfig, err = goconf.ReadConfigFile("config/common.config")
	if err != nil {
		fmt.Println("ReadConfigFile", err)
		Sentry_log(err)
	}
	mongoURI, _ := CommonConfig.GetString("MONGODB_OTT", "uri")
	mongoMode, _ := CommonConfig.GetString("MONGODB_OTT", "mode")

	mongoDBDialInfo, err := mgo.ParseURL(mongoURI)
	if err != nil {
		log.Fatalf("ParseURL Mongo: %s\n", err)
		Sentry_log(err)
	}

	MONGO_DB_OTT_NAME = mongoDBDialInfo.Database
	MONGO_DB_OTT_MODE = mongoMode

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	mongoSession, err = mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
		Sentry_log(err)
	}

	// set Mode
	switch MONGO_DB_OTT_MODE {
	case "Nearest":
		mongoSession.SetMode(mgo.Nearest, true)
	case "SecondaryPreferred":
		mongoSession.SetMode(mgo.SecondaryPreferred, true)
	default:
		mongoSession.SetMode(mgo.Primary, true)
	}
}



func GetCollection() (*mgo.Session, *mgo.Database, error) {
	// Request a socket connection from the session to process our query.
	// Close the session when the goroutine exits and put the connection back
	// into the pool.
	sessionCopy := mongoSession.Copy()
	
	// set Mode
	switch MONGO_DB_OTT_MODE {
	case "Nearest":
		sessionCopy.SetMode(mgo.Nearest, true)
	case "SecondaryPreferred":
		sessionCopy.SetMode(mgo.SecondaryPreferred, true)
	default:
		sessionCopy.SetMode(mgo.Primary, true)
	}

	db := sessionCopy.DB(MONGO_DB_OTT_NAME)
	return sessionCopy, db, nil
}

// func GetCollection2222() (*mgo.Session /**MongoSession*/, *mgo.Database, error) {
// 	// return Mongo.getCollection(config)

// 	session, err := mgo.Dial(MongoDBHosts)
// 	if err != nil {
// 		return &mgo.Session{}, &mgo.Database{}, ErrConnectDatabase
// 	}

// 	session.SetMode(mgo.Monotonic, true)

// 	db := session.DB(DB_NAME)

// 	if AuthUserName != "" || AuthPassword != "" {
// 		err = db.Login(AuthUserName, AuthPassword)
// 	}

// 	if err != nil {
// 		Sentry_log(err)
// 		return &mgo.Session{}, &mgo.Database{}, ErrLogin
// 	}

// 	return session, db, err
// }
