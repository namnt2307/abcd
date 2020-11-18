package artist

import (
	"errors"
	. "cm-v5/serv/module"
	. "cm-v5/schema"
	
)

func GetIdArtistBySlug(peopleSlug string, cacheActive bool) (string, error) {
	var peopleID string
	if peopleSlug == "" {
		return peopleID, errors.New("GetIdArtistBySlug: Empty data")
	}
	var keyCache = ARTIST_SLUG + "_" + peopleSlug
	if cacheActive {
		valueCache, err := mRedis.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &peopleID)
			if err == nil {
				return peopleID, nil

			}
		}
	}

	//Connect mysql
	db_mysql, err := ConnectMySQL()
	if err != nil {
		return peopleID, err
	}
	defer db_mysql.Close()

	err = db_mysql.QueryRow("SELECT id FROM people WHERE slug = ?", peopleSlug).Scan(&peopleID)
	if err != nil {
		return peopleID, err
	}

	// Write Redis
	dataByte, _ := json.Marshal(peopleID)
	mRedis.SetString(keyCache, string(dataByte), 3600*2)

	return peopleID, nil
}

func GetInfoArtist(peopleID string, cacheActive bool) (ArtistObjectStruct, error) {
	var ArtistOutputObject ArtistObjectStruct
	var listPeopleID []string
	var keyCache = PREFIX_REDIS_ARTIST + "_" + peopleID
	if cacheActive {
		valueCache, err := mRedis.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &ArtistOutputObject)
			if err == nil {
				return ArtistOutputObject, nil
			}
		}
	}

	listPeopleID = append(listPeopleID, peopleID)
	ArtistObjects, err := GetArtistByListID(listPeopleID, cacheActive)
	if len(ArtistObjects) <= 0 && err != nil {
		return ArtistOutputObject, err
	}

	ArtistOutputObject = ArtistObjects[0]

	// Write Redis
	dataByte, _ := json.Marshal(ArtistOutputObject)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_LV1)

	return ArtistOutputObject, nil
}
