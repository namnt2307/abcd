package page

import (
	// "errors"
	"fmt"

	. "cm-v5/schema"
	. "cm-v5/serv/module"
	seo "cm-v5/serv/module/seo"

	"gopkg.in/mgo.v2/bson"
)

func GetInfoPageRibbonsV3(page_id, platform string, limit int, cacheActive bool) ([]PageRibbonsOutputObjectStruct, error) {
	var PageRibbonsV3Output = make([]PageRibbonsOutputObjectStruct, 0)

	var keyCache = KV_DETAIL_RIBBONS_V3 + "_" + platform + "_" + page_id + fmt.Sprint(limit)
	if cacheActive {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &PageRibbonsV3Output)
			if err == nil {
				return PageRibbonsV3Output, nil
			}
		}
	}

	listRibbonV3, err := GetListRibbonV3(page_id, platform, cacheActive)
	if err != nil {
		Sentry_log(err)
		return PageRibbonsV3Output, err
	}

	dataByte, _ := json.Marshal(listRibbonV3)
	err = json.Unmarshal(dataByte, &PageRibbonsV3Output)
	if err != nil {
		Sentry_log(err)
		return PageRibbonsV3Output, err
	}

	var platformInfo = Platform(platform)
	for k, val := range listRibbonV3 {
		PageRibbonsV3Output[k].Ribbon_items = make([]RibbonItemOutputObjectStruct, 0)

		// Mapping properties
		var PropertiesMapping PropertiesStruct
		switch platformInfo.Type {
		case "web":
			PropertiesMapping.Line = val.Properties.Web.Line
			PropertiesMapping.Is_title = val.Properties.Web.Is_title
			PropertiesMapping.Is_slide = val.Properties.Web.Is_slide
			PropertiesMapping.Is_refresh = val.Properties.Web.Is_refresh
			PropertiesMapping.Is_view_all = val.Properties.Web.Is_view_all
			PropertiesMapping.Thumb = val.Properties.Web.Thumb
		case "smarttv":
			PropertiesMapping.Line = val.Properties.Smarttv.Line
			PropertiesMapping.Is_title = val.Properties.Smarttv.Is_title
			PropertiesMapping.Is_slide = val.Properties.Smarttv.Is_slide
			PropertiesMapping.Is_refresh = val.Properties.Smarttv.Is_refresh
			PropertiesMapping.Is_view_all = val.Properties.Smarttv.Is_view_all
			PropertiesMapping.Thumb = val.Properties.Smarttv.Thumb
		case "app":
			PropertiesMapping.Line = val.Properties.App.Line
			PropertiesMapping.Is_title = val.Properties.App.Is_title
			PropertiesMapping.Is_slide = val.Properties.App.Is_slide
			PropertiesMapping.Is_refresh = val.Properties.App.Is_refresh
			PropertiesMapping.Is_view_all = val.Properties.App.Is_view_all
			PropertiesMapping.Thumb = val.Properties.App.Thumb
		}

		PageRibbonsV3Output[k].Properties = PropertiesMapping
	}

	// Write Redis
	dataByte, _ = json.Marshal(PageRibbonsV3Output)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE)

	return PageRibbonsV3Output, nil
}

func GetListRibbonV3(page_id, platform string, cacheActive bool) ([]RibbonsV3ObjectStruct, error) {
	var RibbonsV3 = make([]RibbonsV3ObjectStruct, 0)
	var keyCache = LIST_RIBBONS_V3 + "_" + platform + "_" + page_id
	if cacheActive {
		value, err := mRedis.GetString(keyCache)
		if err == nil && value != "" {
			err = json.Unmarshal([]byte(value), &RibbonsV3)
			if err == nil {
				return RibbonsV3, nil
			}
		}
	}

	var platformInfo = Platform(platform)

	// Connect MongoDB
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return RibbonsV3, err
	}
	defer session.Close()

	var where = bson.M{
		"menus.id":        page_id,
		"menus.platforms": platformInfo.Id,
		"platforms":       platformInfo.Id,
		"status":          1,
		"type":            bson.M{"$in": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}},
	}

	err = db.C(COLLECTION_RIB_V3).Find(where).Sort("odr").All(&RibbonsV3)
	if err != nil {
		mRedis.Del(keyCache)
		Sentry_log(err)
		return RibbonsV3, err
	}

	if len(RibbonsV3) <= 0 {
		mRedis.Del(keyCache)
		return RibbonsV3, nil
	}

	for k, val := range RibbonsV3 {
		RibbonsV3[k].Seo = seo.FormatSeoRibbon(val.Id, val.Slug, val.Name, 10, "", cacheActive)

	}

	// Write cache
	dataByte, _ := json.Marshal(RibbonsV3)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_LV1)

	return RibbonsV3, nil
}

func GetPageIdBySlug(pageSlug string) (string, error) {

	var keyCache = "DETAIL_PAGE_SLUG_" + pageSlug
	value, err := mRedis.GetString(keyCache)
	if err == nil && value != "" {
		return value, nil
	}

	var pageId string = ""

	// // Check slug validate
	// slugSplits := strings.Split(pageSlug, fmt.Sprintf(SEO_MENU_URL, ""))
	// if len(slugSplits) != 2 {
	// 	return pageId, errors.New("Slug not validated")
	// }
	// pageSlug = slugSplits[1]

	//Connect mysql
	db_mysql, err := ConnectMySQL()

	if err != nil {
		return pageId, err
	}
	defer db_mysql.Close()

	// sqlRaw := fmt.Sprintf(`SELECT a.id FROM menu_item as a WHERE a.slug = '%s' AND a.status = 1  ORDER BY a.odr ASC`, pageSlug)
	menuItemIdObj, err := db_mysql.Query(`
		SELECT a.id 
		FROM menu_item as a
		WHERE a.slug = ? AND a.status = 1  ORDER BY a.odr ASC`, pageSlug)

	if err != nil {
		return pageId, err
	}

	//fomart result in db
	for menuItemIdObj.Next() {
		err = menuItemIdObj.Scan(&pageId)
		if err != nil {
			return pageId, err
		}
	}

	mRedis.SetString(keyCache, pageId, TTL_REDIS_LV1)
	return pageId, nil
}
