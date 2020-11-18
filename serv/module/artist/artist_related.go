package artist

import (
	"fmt"
	"strings"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
)

func GetArtistRelated(peopleID, platform string, page, limit int, cacheActive bool) (ArtistRelatedOutputObjectStruct, error) {
	var ArtistRelatedOutputObject ArtistRelatedOutputObjectStruct
	ArtistRelatedOutputObject.Items = make([]ItemArtistOutputObjectStruct, 0)

	var keyCache = KV_ARTIST_RELATED + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit) + "_" + peopleID + "_" + platform

	if cacheActive {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &ArtistRelatedOutputObject)
			// fmt.Println(len(ArtistRelatedOutputObject.Items))
			if err == nil {
				return ArtistRelatedOutputObject, nil
			}
		}
	}

	//Lay vodID cua nghe si da tham gia
	listVodID, err := GetVodIdsArtist(peopleID, platform, 0, 1000, 0, cacheActive)
	if err != nil || len(listVodID) <= 0 {
		return ArtistRelatedOutputObject, err
	}

	//Lay danh sach id nghe si lien quan
	listPeopleID, err := GetListPeopleIdRelated(listVodID, peopleID, page, limit, cacheActive)

	if err != nil || len(listPeopleID) <= 0 {
		return ArtistRelatedOutputObject, err
	}

	//Lay thong tin nghe si lien quan bang listPeopleID
	ArtistObjects, err := GetArtistByListID(listPeopleID, cacheActive)
	if len(ArtistObjects) <= 0 || err != nil {
		return ArtistRelatedOutputObject, err
	}

	dataByte, _ := json.Marshal(ArtistObjects)
	err = json.Unmarshal(dataByte, &ArtistRelatedOutputObject.Items)
	if err != nil {
		return ArtistRelatedOutputObject, err
	}

	//Kiem tra thong tin nghe si va lay thong tin SEO nghe si
	dataArtist, err := GetInfoArtist(peopleID, cacheActive)
	if err != nil || dataArtist.Id == "" {
		return ArtistRelatedOutputObject, err
	}
	ArtistRelatedOutputObject.Seo = dataArtist.Seo

	//total
	total, err := GetTotalArtitsRelated(listVodID, peopleID, cacheActive)

	//Pagination
	ArtistRelatedOutputObject.Metadata.Page = page
	ArtistRelatedOutputObject.Metadata.Limit = limit
	ArtistRelatedOutputObject.Metadata.Total = total

	// Write Redis
	dataByte, _ = json.Marshal(ArtistRelatedOutputObject)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE)

	return ArtistRelatedOutputObject, nil
}

func GetTotalArtitsRelated(listVodID []string, peopleID string, cacheActive bool) (int, error) {
	var keyCache = ARTIST_VOD_TOTAL + "_" + peopleID
	var total int

	if cacheActive {
		total, err := mRedis.GetInt(keyCache)
		if err == nil {
			return total, nil
		}
	}

	//Connect mysql
	db_mysql, err := ConnectMySQL()
	if err != nil {
		return total, err
	}
	defer db_mysql.Close()

	queryString := strings.Join(listVodID, "','")
	sqlRaw := fmt.Sprintf(`SELECT COUNT(DISTINCT people_id) FROM entity_people WHERE entity_id IN ('%s') AND status = 1 AND people_id != '%s'`, queryString, peopleID)
	totalPeople, err := db_mysql.Query(sqlRaw)
	if err != nil {
		return total, err
	}

	for totalPeople.Next() {
		err := totalPeople.Scan(&total)
		if err != nil {
			return total, err
		}
	}

	//Write cache
	mRedis.SetInt(keyCache, total, TTL_REDIS_LV1)
	return total, err
}

func GetListPeopleIdRelated(listVodID []string, peopleID string, page, limit int, cacheActive bool) ([]string, error) {
	var listPeopleID []string
	var keyCache = ARTIST_RELATED + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit) + "_" + peopleID
	if cacheActive {
		dataCache, err := mRedis.GetString(keyCache)
		if err == nil && dataCache != "" {
			err = json.Unmarshal([]byte(dataCache), &listPeopleID)
			if err == nil {
				return listPeopleID, nil
			}
		}
	}

	//Connect mysql
	db_mysql, err := ConnectMySQL()
	if err != nil {
		return listPeopleID, err
	}
	defer db_mysql.Close()

	queryString := strings.Join(listVodID, "','")
	sqlRaw := fmt.Sprintf(`SELECT DISTINCT people_id FROM entity_people WHERE entity_id IN ('%s') AND status = 1 AND people_id != '%s' ORDER BY created_at ASC LIMIT %d, %d `, queryString, peopleID, page*limit, limit)
	listPeople, err := db_mysql.Query(sqlRaw)
	if err != nil {
		return listPeopleID, err
	}

	for listPeople.Next() {
		var idPeople string
		err := listPeople.Scan(&idPeople)
		if err != nil {
			return listPeopleID, err
		}
		listPeopleID = append(listPeopleID, idPeople)
	}

	// Write cache
	dataByte, _ := json.Marshal(listPeopleID)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_LV1)

	return listPeopleID, nil
}
