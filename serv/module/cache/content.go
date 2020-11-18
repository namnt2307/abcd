package cache

import (
	"net/http"
	"strings"
	"time"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	content "cm-v5/serv/module/content"
	vod "cm-v5/serv/module/vod"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func UpdateCacheContentById(c *gin.Context) {
	contentId := c.Param("content_id")
	var statusRibItem = 0
	//detail
	for _, val := range platforms {
		_, err := content.GetContent(contentId, val, false)
		if err == nil && statusRibItem != 1 {
			statusRibItem = 1
		}

		//Khanh DT-11009
		content.GetInfoVod(contentId, "", val, false)
	}

	content.ContentProfileMySQl(contentId, false)

	//Clear cache related video
	content.GetListIDRelatedVideoByZrange(contentId, 0, 30, false)

	//ribbon item
	lstRibitem, _ := GetRibItemV3(contentId)
	var lstRibIds []string
	for _, val := range lstRibitem {
		for _, ribID := range val.Rib_ids {

			var isNew = true
			for _, ribIDInLst := range lstRibIds {
				if ribIDInLst == ribID {
					isNew = false
					break
				}
			}
			if isNew == true {
				lstRibIds = append(lstRibIds, ribID)
			}
		}

	}

	go func(lstRibitem []RibbonItemV3ObjectStruct, lstRibIds []string) {
		for _, val := range lstRibitem {
			ClearCacheRibbonItemV3(val.Rib_item_id, val.Rib_ids, statusRibItem)
		}
		// update cache ribbon
		if len(lstRibIds) > 0 {
			for _, ribID := range lstRibIds {
				ClearCacheRibbonV3(ribID)
			}

		}
	}(lstRibitem, lstRibIds)

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}

func UpdateCacheContentEpisodeById(c *gin.Context) {
	var group_id string = c.Param("group_id")
	go func(group_id string) {
		time.Sleep(10 * time.Second)
		vod.GetVodByListGroup(group_id, 0, 0, 100, false)
	}(group_id)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}

func GetRibItemV3(ref_id string) ([]RibbonItemV3ObjectStruct, error) {
	var where = bson.M{
		"status": bson.M{"$ne": -1},
	}
	if ref_id != "" {
		where["ref_id"] = bson.M{"$in": strings.Split(ref_id, ",")}
	}

	var RibbonItemV3 = make([]RibbonItemV3ObjectStruct, 0)
	// Connect MongoDB
	session, db, err := GetCollection()
	if err != nil {
		return RibbonItemV3, err
	}
	defer session.Close()

	err = db.C(COLLECTION_RIB_ITEM_V3).Find(where).All(&RibbonItemV3)
	if err != nil {
		return RibbonItemV3, err
	}
	return RibbonItemV3, nil
}
