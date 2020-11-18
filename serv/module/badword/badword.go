package badword

import (
	"log"
	"strings"
	"time"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"gopkg.in/mgo.v2/bson"
)

func UpdateWord(id string, content string) error {
	objBadword := ProcessWord(content)
	//connect db
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return err
	}
	defer session.Close()

	var where = bson.M{"_id": bson.ObjectIdHex(id)}

	update := bson.M{
		"$set": bson.M{
			"content":    objBadword.Content,
			"key":        objBadword.Key,
			"type":       objBadword.Type,
			"updated_at": objBadword.Updated_at,
		},
	}
	err = db.C(COLLECTION_BLACKLIST).Update(where, update)
	if err != nil {
		return err
	} else {
		// Update cache
		UpdateRedisBadWord("bad-words", "normal")
		UpdateRedisBadWord("bad-words", "like")

		return nil
	}
}

func AddWord(content string) error {
	//connect db
	objBadword := ProcessWord(content)
	// log.Println(objBadword)
	session, db, err := GetCollection()
	if err != nil {
		return err
	}

	defer session.Close()

	err = db.C(COLLECTION_BLACKLIST).Insert(objBadword)
	if err != nil {
		log.Println(err)
		return err
	} else {
		UpdateRedisBadWord("bad-words", "normal")
		UpdateRedisBadWord("bad-words", "like")

		return nil
	}
}

func UpdateRedisBadWord(key string, typeword string) {
	key = key + "_" + typeword
	mRedis.Del(key)
	mRedis.Del(key + "-count")
}

// *key* : (type =2)
// key*: (type = 3)
// *key: (type = 1)

func ProcessWord(badword string) BadwordObjectStruct {

	var objBadword BadwordObjectStruct
	objBadword.Type = 0
	objBadword.Content = badword

	// Get time to day
	t := time.Now()
	CurentTime := t.Format("2006-01-02 15:04:05")
	objBadword.Created_at = CurentTime
	objBadword.Updated_at = CurentTime

	arr := strings.Split(badword, "")

	if string(arr[0]) == "*" && arr[len(arr)-1] == "*" {
		objBadword.Type = 2
		arr[0] = ""
		arr[len(arr)-1] = ""
		objBadword.Content = strings.Join(arr, "")
	} else if arr[0] == "*" {
		objBadword.Type = 1
		arr[0] = ""
		objBadword.Content = strings.Join(arr, "")
	} else if arr[len(arr)-1] == "*" {
		objBadword.Type = 3
		arr[len(arr)-1] = ""
		objBadword.Content = strings.Join(arr, "")
	}

	objBadword.Key = badword
	objBadword.Status = 1

	return objBadword
}

func DelBadword(id string, status string) error {
	//connect db
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return err
	}
	defer session.Close()

	var where = bson.M{"_id": bson.ObjectIdHex(id)}
	// log.Println("============")
	update := bson.M{
		"$set": bson.M{
			"status": 0,
		},
	}
	err = db.C(COLLECTION_BLACKLIST).Update(where, update)
	if err != nil {
		return err
	} else {
		// Update cache
		UpdateRedisBadWord("bad-words", "normal")
		UpdateRedisBadWord("bad-words", "like")

		return nil
	}
}
