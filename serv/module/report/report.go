package report

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
)

func GetReportType(platform string, page, limit int, cacheActive bool) ([]ReportTypeStruct, error) {
	platformInfo := Platform(platform)

	var ReportTypeOutputs = make([]ReportTypeStruct, 0)
	var keyCache = KV_REDIS_REPORT_TYPE + "_" + strconv.Itoa(platformInfo.Id) + "_" + strconv.Itoa(page) + "_" + strconv.Itoa(limit)

	if cacheActive {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &ReportTypeOutputs)
			if err == nil {
				return ReportTypeOutputs, nil
			}
		}
	}

	db_mysql, err := ConnectMySQL()
	if err != nil {
		return ReportTypeOutputs, err
	}
	defer db_mysql.Close()

	sqlRaw := fmt.Sprintf(`
		SELECT id,title,description FROM report_type 
		WHERE status = 1
		LIMIT %d , %d`, page*limit, limit)

	listReportType, err := db_mysql.Query(sqlRaw)
	if err != nil {
		return ReportTypeOutputs, err
	}

	for listReportType.Next() {
		var ReportTypeObj ReportTypeStruct
		err = listReportType.Scan(&ReportTypeObj.Id, &ReportTypeObj.Title, &ReportTypeObj.Description)
		if err != nil {
			continue
		}

		// Lay avatar nghe si lien quan
		ReportTypeOutputs = append(ReportTypeOutputs, ReportTypeObj)
	}

	// Write Redis
	dataByte, _ := json.Marshal(ReportTypeOutputs)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_REDIS_1_HOURS)

	return ReportTypeOutputs, nil
}

func UserReportContent(ReportContent UserReportContentStruct) error {

	var keyCache = PREFIX_REDIS_USER_REPORT + "_" + ReportContent.User_id
	result := mRedis.Incr(keyCache)
	mRedis.Expire(keyCache, TTL_REDIS_1_HOURS)

	//limit 5 time in 1h
	if result < 6 {
		t := time.Now()
		CurentTime := t.Format("2006-01-02 15:04:05")
		//insert
		db_mysql, err := ConnectMySQL()
		if err != nil {
			return err
		}
		defer db_mysql.Close()

		StmtIns, err := db_mysql.Prepare(`
			INSERT report
			SET id=?, message=?, created_at=?, updated_at=?, status=?
			, entity_id=?, report_id=?, user_id=?, user_agent=?, access_token=?
			, video_profile=?, audio_profile=?, subtitle=?`)
		if err != nil {
			return err
		}
		_, err = StmtIns.Exec(string(UUIDV4()), ReportContent.Message, CurentTime, CurentTime, 0, ReportContent.Entity_id,
			ReportContent.Report_id, ReportContent.User_id, ReportContent.User_agent, ReportContent.Access_token, ReportContent.Video_profile, ReportContent.Audio_profile, ReportContent.Subtitle)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Please report slow down")
	}

	return nil
}
