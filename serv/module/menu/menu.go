package menu

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	. "cm-v5/schema"
	. "cm-v5/serv/module"
	seo "cm-v5/serv/module/seo"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func GetListMenu(c *gin.Context) {
	platform := Platform(c.DefaultQuery("platform", "web"))
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	keyCache := LIST_MENU + "_" + platform.Type

	if cacheActive {
		val, err := LocalCache.GetValue(keyCache)
		if err == nil {
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", val))
			return
		}
	}

	dataContent, err := GetListMenuInfoMySQL(platform.Type, platform.Id, cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	// Write cache
	LocalCache.SetValue(keyCache, dataContent, TTL_LOCALCACHE)

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataContent))
}

func GetListMenuInfoMySQL(platform_type string, platform_id int, cacheActive bool) ([]MenuOutputObjectStruct, error) {
	var menuItems []MenuItemObjStruct
	var menuItemsNew []MenuItemObjStruct
	var menuOutputObject []MenuOutputObjectStruct
	var mRedis RedisModelStruct

	keyCache := LIST_MENU + "_" + fmt.Sprint(platform_id)
	if cacheActive {
		// Get data in cache
		dataCache, err := mRedis.GetString(keyCache)
		// fmt.Println("TIENNM GetListMenuInfoMySQL Read Redis:" , keyCache, err)
		// fmt.Println("TIENNM GetListMenuInfoMySQL Data Redis:" , dataCache)
		if err == nil && dataCache != "" {
			err = json.Unmarshal([]byte(dataCache), &menuOutputObject)
			if err == nil {
				return menuOutputObject, nil
			}
		}
	}

	//Connect mysql
	db_mysql, err := ConnectMySQL()
	if err != nil {
		return menuOutputObject, err
	}
	defer db_mysql.Close()

	sqlRaw := fmt.Sprintf(`
		SELECT a.id,a.name,a.slug,a.parent_id,COALESCE(a.icon, ''),
		COALESCE(a.tag_id,''),COALESCE(a.title_ribbon, ''),COALESCE(a.quantity_ribbon, 0),
		COALESCE(a.list_ribbon, ''),COALESCE(a.have_banner, 0) ,COALESCE(a.icon_text, '')
		FROM menu_item as a 
		INNER JOIN menu_item_platform as b ON a.id = b.menu_item_id 
		WHERE a.menu_id = '%s' AND a.status = 1 AND b.platform_id = %d
		ORDER BY a.odr ASC`, MENU_ID, platform_id)

	menuItemIdObj, err := db_mysql.Query(sqlRaw)

	if err != nil {
		return menuOutputObject, err
	}

	//fomart result in db
	for menuItemIdObj.Next() {
		var menuItem MenuItemObjStruct

		err = menuItemIdObj.Scan(&menuItem.Id, &menuItem.Name, &menuItem.Slug, &menuItem.Parent_id, &menuItem.Icon, &menuItem.Tag.Id, &menuItem.Title_ribbon, &menuItem.Quantity_ribbon, &menuItem.List_ribbon, &menuItem.Have_banner, &menuItem.Icon_text)
		if err != nil {
			return menuOutputObject, err
		}

		menuItem.Seo = seo.FormatSeoMenu(menuItem.Id, menuItem.Slug, menuItem.Name, cacheActive)
		menuItem.Sub_menu_ribbon, err = GetListRibbonFilterMenu(menuItem.List_ribbon, platform_id)
		menuItems = append(menuItems, menuItem)
	}

	var checkNewUIValid = false
	//get menu
	for _, val := range menuItems {
		if val.Parent_id == "0" {
			var menuItemNew MenuItemObjStruct
			dataByte, _ := json.Marshal(val)
			err = json.Unmarshal(dataByte, &menuItemNew)
			if err != nil {
				return menuOutputObject, err
			}
			menuItemsNew = append(menuItemsNew, menuItemNew)
		}

		if val.Slug == "trang-chu-new-ui" {
			checkNewUIValid = true
		}
	}

	dataByte, _ := json.Marshal(menuItemsNew)
	err = json.Unmarshal(dataByte, &menuOutputObject)
	if err != nil {
		return menuOutputObject, err
	}
	//get sub menu
	for k1, _ := range menuItemsNew {
		var subMenuItems []MenuItemObjStruct
		for k2, result := range menuItems {
			var subMenuItemsNew MenuItemObjStruct
			if menuItemsNew[k1].Id == menuItems[k2].Parent_id {
				dataByte, _ := json.Marshal(result)
				err = json.Unmarshal(dataByte, &subMenuItemsNew)
				if err != nil {
					// log.Println("GetListMenu in sub menu", err)
					continue
				}
				subMenuItems = append(subMenuItems, subMenuItemsNew)
			}
		}
		dataByte, _ = json.Marshal(subMenuItems)
		err = json.Unmarshal(dataByte, &menuOutputObject[k1].Sub_menu)
		if err != nil {
			// log.Println("Get sub menu getListMenu", err)
			continue
		}
	}

	if checkNewUIValid == false {
		// Push sentry
		Sentry_log_mysql_miss("Not New UI In Menu")
	}

	// Write cache
	dataByte, _ = json.Marshal(menuOutputObject)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_LV1)

	// fmt.Println("TIENNM GetListMenuInfoMySQL Write Redis:" , keyCache , errr , TTL_REDIS_LV1)
	return menuOutputObject, nil

}

func GetListMenuInfo(platform_type string, platform_id int, cacheActive bool) ([]MenuOutputObjectStruct, error) {
	var menuOutputObject []MenuOutputObjectStruct

	var mRedis RedisModelStruct

	keyCache := LIST_MENU + "_" + platform_type

	if cacheActive {
		// Get data in cache
		dataCache, err := mRedis.GetString(keyCache)
		if err == nil && dataCache != "" {
			err = json.Unmarshal([]byte(dataCache), &menuOutputObject)
			if err == nil {
				return menuOutputObject, nil
			}
		}
	}

	// Connect MongoDB
	session, db, err := GetCollection()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var where = bson.M{
		"platform_id": platform_id,
	}

	err = db.C(COLLECTION_MENU).Find(where).All(&menuOutputObject)
	if err != nil {
		return menuOutputObject, err
	}

	//fomart result in data
	for index, element := range menuOutputObject {

		element.Seo = seo.FormatSeoMenu(element.Id, element.Slug, element.Name, cacheActive)
		element.Tag.Seo = seo.FormatSeoMenu(element.Id, element.Tag.Slug, element.Tag.Name, cacheActive)
		menuOutputObject[index] = element
	}

	//get menu
	for k1, _ := range menuOutputObject {
		for k2, result := range menuOutputObject {
			if menuOutputObject[k1].Id == menuOutputObject[k2].Parent_id {
				menuOutputObject[k1].Sub_menu = append(menuOutputObject[k1].Sub_menu, result)

			}
		}
	}

	// Write cache
	dataByte, _ := json.Marshal(menuOutputObject)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_LV1)

	return menuOutputObject, nil

}


func GetListRibbonFilterMenu(list_id_ribbon string, platform_id int) ([]MenuRibbonObjectStruct, error) {
	var menuRibbonObjects = make([]MenuRibbonObjectStruct, 0)

	if list_id_ribbon == "" {
		return menuRibbonObjects, nil
	}

	var listRibbonID = strings.Split(list_id_ribbon, ",")

	// Get data in DB
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return menuRibbonObjects, err
	}
	defer session.Close()
	var where = bson.M{
		"id":              bson.M{"$in": listRibbonID},
		"status":          1,
		"platforms":       platform_id,
		"menus.platforms": platform_id,
	}

	err = db.C(COLLECTION_RIB_V3).Find(where).All(&menuRibbonObjects)
	if err != nil && err.Error() != "not found" {
		return menuRibbonObjects, err
	}

	//create seo object
	for k, val := range menuRibbonObjects {
		menuRibbonObjects[k].Seo = seo.FormatSeoRibbon(val.Id, val.Slug_filter, val.Name_filter, 10, "", true)
	}

	return menuRibbonObjects, nil
}
