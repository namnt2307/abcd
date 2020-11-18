package cache

import (
	"log"
	"net/http"
	"strconv"
	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)


func AddViewByContentId(c *gin.Context) {
	contentId := c.Param("content_id")
	views, err := strconv.ParseInt(c.Param("views"), 10, 64)
	// Add redis
	keyCache := "VIEW_" + contentId
	viewsAdd := mRedisKV.IncrBy(keyCache, views) 
	// Update DB
	var viewObj ViewObjectStruct
	viewObj.View = viewsAdd
	viewObj.ContentId = contentId
	//connect db
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		log.Println("AddViewByContentId err", err)
	}
	defer session.Close()
	var where = bson.M{
		"contentid": viewObj.ContentId,
	}
	_, err = db.C(COLLECTION_VIEW).Upsert(where, viewObj)

	if err != nil {
		log.Println("AddViewByContentId err", err)
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", viewObj))
}
