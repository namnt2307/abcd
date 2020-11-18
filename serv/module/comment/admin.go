package comment

import (
	"log"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"gopkg.in/mgo.v2/bson"
)

func GetComments(contentID, parentId, userID, dateFilter, status, textSearch, sort string, page, limit int) []CommentObjectStruct {

	var listComments = make([]CommentObjectStruct, 0)
	// Connect DB
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return listComments
	}
	defer session.Close()
	Where, Select := BuildQuerySelect(contentID, parentId, userID, dateFilter, status, textSearch)
	if sort != "asc" {
		sort = "created_at"
	} else {
		sort = "-created_at"
	}

	//get comment reply
	if parentId != "" {
		Select["reply"] = bson.M{"$slice": []int{page * limit, limit}}
		err = db.C(COLLECTION_COMMENT).Find(Where).Select(Select).Sort(sort).All(&listComments)
		if err != nil && err.Error() != "not found" {
			log.Println("GetComments err ", err)
			return listComments
		}
		return listComments
	}

	err = db.C(COLLECTION_COMMENT).Find(Where).Select(Select).Sort(sort).Skip(page * limit).Limit(limit).All(&listComments)
	if err != nil && err.Error() != "not found" {
		log.Println("GetComments err ", err)
		return listComments
	}

	return listComments

}

func GetTotalComment(contentId, parentId, userId, dateFilter, status, textSearch string) int {
	session, db, err := GetCollection()
	if err != nil {
		return 0
	}
	defer session.Close()

	where, _ := BuildQuerySelect(contentId, parentId, userId, dateFilter, status, textSearch)

	total, _ := db.C(COLLECTION_COMMENT).Find(where).Count()
	return total
}
func BuildQuerySelect(contentID, parentId, userID, dateFilter, status, textSearch string) (bson.M, bson.M) {

	Where := bson.M{}
	Select := bson.M{
		"_id":        0,
		"updated_at": 0,
		"reply":      0,
	}

	if contentID != "" {
		Where["content_id"] = contentID
	}
	if userID != "" {
		Where["user.id"] = userID
	}
	if status != "" {
		int_status, _ := StringToInt(status)
		Where["status"] = int_status
	}
	if parentId != "" {
		Where["id"] = parentId
		Select = bson.M{
			"content_id": 0,
			"message":    0,
			"status":     0,
			"created_at": 0,
			"updated_at": 0,
		}
	}

	if dateFilter != "" {
		start, end := GetTimeStampStartEndFromString(dateFilter)
		if start > 0 && end > 0 && start < end {
			Where["$and"] = []bson.M{
				bson.M{"created_at": bson.M{"$gte": start}},
				bson.M{"created_at": bson.M{"$lte": end}},
			}
		}
	}
	if textSearch != "" {
		Where["$text"] = bson.M{"$search": textSearch}
		Select["score"] = bson.M{"$meta": "textScore"}
	}
	return Where, Select
}
