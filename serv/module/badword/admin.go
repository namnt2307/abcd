package badword

import (
	"log"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"gopkg.in/mgo.v2/bson"
)

func GetBadwords(status, textSearch, sort string, page, limit int) []BadwordObjectStruct {

	var listBadwords = make([]BadwordObjectStruct, 0)
	// Connect DB
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return listBadwords
	}
	defer session.Close()
	Where,_ := BuildQuerySelect(status, textSearch)
	if sort != "asc" {
		sort = "created_at"
	} else {
		sort = "-created_at"
	}
	err = db.C(COLLECTION_BLACKLIST).Find(Where).Sort(sort).Skip(page * limit).Limit(limit).All(&listBadwords)
	if err != nil && err.Error() != "not found" {
		log.Println("GetBadwords err ", err)
		return listBadwords
	}

	return listBadwords

}

func GetTotalBadword(status, textSearch string) int {
	session, db, err := GetCollection()
	if err != nil {
		return 0
	}
	defer session.Close()

	where, _ := BuildQuerySelect(status, textSearch)

	total, _ := db.C(COLLECTION_BLACKLIST).Find(where).Count()
	return total
}
func BuildQuerySelect(status, textSearch string) (bson.M, bson.M) {
	
	Where := bson.M{}
	Select := bson.M{
		"_id":        0,
		"updated_at": 0,
		"reply":      0,
	}

	if status != "" {
		int_status,_ := StringToInt(status)
		Where["status"] = int_status
	}

	if textSearch != "" {
		// Where["content"] = bson.M{"$search": textSearch}
		Where["content"] = bson.RegEx{
			Pattern: "^" + textSearch,
			Options: "i",
		}


	}
	return Where, Select
}
