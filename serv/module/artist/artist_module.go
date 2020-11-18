package artist

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	seo "cm-v5/serv/module/seo"
)

func GetArtistByListID(listPeopleID []string, cacheActive bool) ([]ArtistObjectStruct, error) {
	var ArtistObjects = make([]ArtistObjectStruct, 0)
	if len(listPeopleID) <= 0 {
		return ArtistObjects, errors.New("GetArtistByListID: Empty data 1")
	}

	if cacheActive {
		dataRedis, err := mRedis.HMGet(REFIX_REDIS_ARTIST_HASH, listPeopleID)
		if err == nil && len(dataRedis) > 0 {
			for _, val := range dataRedis {
				if str, ok := val.(string); ok {
					var ArtistObject ArtistObjectStruct
					err = json.Unmarshal([]byte(str), &ArtistObject)
					if err == nil {
						continue
					}
					ArtistObjects = append(ArtistObjects, ArtistObject)
				}
			}
		}
	}

	//reset ArtistObjects = []
	// ArtistObjects = make([]ArtistObjectStruct, 0)
	if len(ArtistObjects) == len(listPeopleID) {
		return ArtistObjects, nil
	}

	//Connect mysql
	db_mysql, err := ConnectMySQL()
	if err != nil {
		return ArtistObjects, err
	}
	defer db_mysql.Close()

	queryPeopleID := strings.Join(listPeopleID, "','")
	sqlRaw := fmt.Sprintf(`
	SELECT p.id, p.name, p.slug, COALESCE(p.gender, ''), COALESCE(p.birthday, ''), p.status, 
	COALESCE(p.info, ''), COALESCE(p.country_id, ''), COALESCE(p.job, ''), COALESCE(t.name, '') 
	FROM people as p LEFT JOIN tag as t ON p.country_id = t.id WHERE p.id IN ('%s') `, queryPeopleID)
	listPeople, err := db_mysql.Query(sqlRaw)
	if err != nil {
		return ArtistObjects, err
	}

	for listPeople.Next() {
		var ArtistObject ArtistObjectStruct
		err = listPeople.Scan(&ArtistObject.Id, &ArtistObject.Name, &ArtistObject.Slug, &ArtistObject.Gender, &ArtistObject.Birthday, &ArtistObject.Status, &ArtistObject.Info, &ArtistObject.Country.Id, &ArtistObject.Job, &ArtistObject.Country.Name)
		if err != nil {
			continue
		}

		// Lay avatar nghe si lien quan
		imageArtist, _ := GetImageArtist(ArtistObject.Id, cacheActive)
		ArtistObject.Images.Avatar = BuildImage(imageArtist.Url)

		//Seo
		seoArtistRelated := seo.FormatSeoArtist(ArtistObject , 0 , "")
		ArtistObject.Seo = seoArtistRelated
		ArtistObjects = append(ArtistObjects, ArtistObject)
	}
	if len(ArtistObjects) <= 0 {
		return ArtistObjects, errors.New("GetArtistByListID: Empty data")
	}

	// Write cache
	for _, val := range ArtistObjects {
		dataByte, _ := json.Marshal(val)
		// Write cache follow ID
		mRedis.HSet(REFIX_REDIS_ARTIST_HASH, val.Id, string(dataByte))
	}

	return ArtistObjects, nil
}

func GetVodIdsArtist(peopleId, platform string, page, limit, sort int, cacheActive bool) ([]string, error) {
	var listVodID []string
	if peopleId == "" {
		return listVodID, errors.New("GetVodIdsArtist: Empty data")
	}
	var keyCache = ARTIST_VOD_ID + "_" + peopleId + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit) + "_" + fmt.Sprint(sort) + "_" + platform
	var mRedis RedisModelStruct
	if cacheActive {
		dataCache, err := mRedis.GetString(keyCache)
		if err == nil && dataCache != "" {
			err = json.Unmarshal([]byte(dataCache), &listVodID)
			if err == nil {
				return listVodID, nil
			}
		}
	}

	//Connect mysql
	db_mysql, err := ConnectMySQL()
	if err != nil {
		return listVodID, err
	}
	defer db_mysql.Close()

	var orderBy = "b.release_date ASC"
	switch sort {
	case 1:
		orderBy = "b.views DESC"
	case 2:
		orderBy = "b.release_date ASC"
	case 3:
		orderBy = "b.release_date DESC"
	}

	valid := regexp.MustCompile("^[A-Za-z0-9_-]+$")
	if !valid.MatchString(peopleId) {
		return listVodID, errors.New("peopleId not valid")
	}

	platformInfo := Platform(platform)
	sqlRaw := fmt.Sprintf(`
		SELECT DISTINCT (a.entity_id) FROM entity_people as a
		LEFT JOIN entity_vod as b ON b.id = a.entity_id
		LEFT JOIN entity_vod_platform as c ON c.entity_id = b.id
		WHERE a.people_id = '%s' AND b.status IN (3, 5) AND b.type IN (1,3,4,5)
		AND c.platform_id = %s ORDER BY %s LIMIT %d , %d`, peopleId, fmt.Sprint(platformInfo.Id), orderBy, page*limit, limit)

	dataVods, err := db_mysql.Query(sqlRaw)
	if err != nil {
		return listVodID, err
	}

	for dataVods.Next() {
		var vodId string
		err := dataVods.Scan(&vodId)
		if err != nil {
			return listVodID, err
		}
		listVodID = append(listVodID, vodId)
	}

	// Write cache
	dataByte, _ := json.Marshal(listVodID)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_LV1)

	return listVodID, nil
}

func GetImageArtist(peopleID string, cacheActive bool) (ImageArttistObjectStruct, error) {
	var ImageArttist ImageArttistObjectStruct
	if peopleID == "" {
		return ImageArttist, errors.New("GetImageArtist: Empty data")
	}
	var keyCache = ARTIST_IMAGES + "_" + peopleID
	if cacheActive {
		valueCache, err := mRedis.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &ImageArttist)
			if err == nil {
				return ImageArttist, nil
			}
		}
	}

	//Connect mysql
	db_mysql, err := ConnectMySQL()
	if err != nil {
		return ImageArttist, err
	}
	defer db_mysql.Close()

	err = db_mysql.QueryRow(`
		SELECT a.id, a.image_name, a.image_type, a.resolution_type, a.url 
		FROM image as a JOIN people_image as b ON b.image_id = a.id 
		WHERE b.people_id = ?`, peopleID).Scan(&ImageArttist.Id, &ImageArttist.Image_name, &ImageArttist.Image_type, &ImageArttist.Resolution_type, &ImageArttist.Url)
	if err != nil {
		return ImageArttist, err
	}

	// Write Redis
	dataByte, _ := json.Marshal(ImageArttist)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_LV1)

	return ImageArttist, nil
}
