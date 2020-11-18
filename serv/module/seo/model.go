package seo

import (
	. "cm-v5/schema"
	. "cm-v5/serv/module"
	"encoding/json"
	"errors"
	"fmt"
)

var mRedis RedisModelStruct

func GetSlugDefault(t int) (string, error) {
	var slug string = ""
	db_mysql, err := ConnectMySQL()
	if err != nil {
		return slug, err
	}
	defer db_mysql.Close()

	sqlRaw := `
		SELECT s.slug FROM seo_config AS sc
		LEFT JOIN seo AS s ON s.id = sc.seo_id
		WHERE sc.type = ? AND s.status = 1`

	err = db_mysql.QueryRow(sqlRaw, t).Scan(&slug)
	if err != nil {
		return slug, err
	}

	return slug, nil
}

func GetSEO(ref_id string, cache bool) (SEOObjStruct, error) {
	var SEO SEOObjStruct
	var key string = fmt.Sprintf("seo_info:%s", ref_id)

	if cache {
		value, err := mRedis.GetString(key)
		if err == nil && value != "" {
			err = json.Unmarshal([]byte(value), &SEO)
			if err != nil {
				return SEO, errors.New("empty")
			}
			return SEO, nil
		}
	}

	db_mysql, err := ConnectMySQL()
	if err != nil {
		return SEO, err
	}
	defer db_mysql.Close()

	sqlRaw := `
		SELECT slug, canonical_tag, meta_robots, title,
		title_seo_tag, meta_description, seo_text, alternate, deeplink
		FROM seo
		WHERE ref_id = ? AND status = 1`

	err = db_mysql.QueryRow(sqlRaw, ref_id).Scan(
		&SEO.Slug,
		&SEO.Canonical_tag,
		&SEO.Meta_robots,
		&SEO.Title,
		&SEO.Title_seo_tag,
		&SEO.Meta_description,
		&SEO.Seo_text,
		&SEO.Alternate,
		&SEO.Deeplink,
	)
	if err != nil {
		mRedis.SetString(key, "", TTL_REDIS_30_MINUTE)
		return SEO, err
	}

	dataByte, _ := json.Marshal(SEO)
	mRedis.SetString(key, string(dataByte), TTL_REDIS_30_MINUTE)
	return SEO, nil
}

func CheckExistsSEOBySlug(slug string) string {
	var ref_id string
	key := fmt.Sprintf("exists_seo:%s", slug)
	value, err := mRedis.GetString(key)
	if err == nil && value != "" {
		return value
	}

	db_mysql, err := ConnectMySQL()
	if err != nil {
		return ref_id
	}
	defer db_mysql.Close()

	sqlRaw := `SELECT ref_id FROM seo WHERE slug = ?`
	err = db_mysql.QueryRow(sqlRaw, slug).Scan(&ref_id)
	if err != nil {
		return ref_id
	}

	mRedis.SetString(key, ref_id, TTL_REDIS_2_HOURS)
	return ref_id
}
