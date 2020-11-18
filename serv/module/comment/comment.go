package comment

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"html"

	"unicode/utf8"

	recommendation "cm-v5/serv/module/recommendation"
	"github.com/gin-gonic/gin"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	vod "cm-v5/serv/module/vod"
	"gopkg.in/mgo.v2/bson"

	// Thu DT-9922
	"bytes"
	"strings"

	. "github.com/ahmetb/go-linq"
	"github.com/garyburd/redigo/redis"
)

var LISTSTATUSCOMMENT = []int{1, 2, -1}

func PostComment(userId, platform, contentId, parentId, message string, old_msg string, c *gin.Context) (CommentOutputStruct, error) {
	message = strings.TrimSpace(message)
	// cắt khoảng trắng cuối sentence
	var keyCache = KV_USER_POST_COMMENT + "_" + userId
	result := mRedis.Incr(keyCache)
	if result <= 1 {
		mRedis.Expire(keyCache, TTL_KVCACHE)
	}

	var commentOutput CommentOutputStruct
	if result < 6 {

		platformInfo := Platform(platform)

		//check content
		dataVod, err := vod.GetVodDetail(contentId, platformInfo.Id, true)
		if err != nil {
			// fmt.Println("PostComment err", err)
			return commentOutput, errors.New("Nội dung không tồn tại")
		}

		// check type vod only single or season or trailer is allow comment
		if dataVod.Type != VOD_TYPE_SEASON && dataVod.Type != VOD_TYPE_MOVIE && dataVod.Type != VOD_TYPE_TRAILER {
			return commentOutput, errors.New("Nội dung không cho phép bình luận")
		}

		infoUser, err := GetInfoUserById(userId)
		if err != nil {
			return commentOutput, err
		}
		t := time.Now()
		var commentObj CommentObjectStruct
		commentObj.Id = UUIDV4Generate()
		commentObj.Content_id = contentId
		commentObj.Message = message
		commentObj.Status = 1
		commentObj.Created_at = t.Unix()
		commentObj.Updated_at = t.Unix()
		commentObj.User.Id = userId
		commentObj.User.Avatar = infoUser.Avatar
		commentObj.User.Name = infoUser.Given_name
		commentObj.User.Gender = infoUser.Gender
		// Thu add DT-9997
		if commentObj.User.Name == "" {
			if infoUser.Mobile != "" {
				value := infoUser.Mobile
				runes := []rune(value)
				commentObj.User.Name = string(runes[0:3]) + "***" + string(runes[len(value)-3:])
			} else if infoUser.Email != "" {
				arr_string := strings.Split(infoUser.Email, "@")
				runes := []rune(arr_string[0])
				if len(arr_string[0])-4 > 0 {
					commentObj.User.Name = "***" + string(runes[len(arr_string[0])-4:]) + "@" + arr_string[1]
				} else {
					commentObj.User.Name = "***" + "@" + arr_string[1]
				}
				// commentObj.User.Name = "***" + string(runes[len(arr_string[0])-4:]) + "@" + arr_string[1]
			}

		} else {
			// Thu add DT-14239
			runes_name_phone := []rune(commentObj.User.Name)
			if len(commentObj.User.Name) == 10 && string(runes_name_phone[0]) == "0" {
				commentObj.User.Name = string(runes_name_phone[0:3]) + "***" + string(runes_name_phone[len(commentObj.User.Name)-3:])
			} else {
				re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
				if re.MatchString(commentObj.User.Name) == true {
					arr_string := strings.Split(commentObj.User.Name, "@")
					runes_email := []rune(arr_string[0])
					if len(arr_string[0])-4 > 0 {
						commentObj.User.Name = "***" + string(runes_email[len(arr_string[0])-4:]) + "@" + arr_string[1]
					} else {
						commentObj.User.Name = "***" + "@" + arr_string[1]
					}
				}
			}
			// Thu add DT-14239
		}

		commentObj.Oldmessage = ""

		if message != strings.TrimSpace(old_msg) {
			commentObj.Oldmessage = old_msg
			commentObj.Status = 2
			commentObj.Message = MakeFirstUpperCase(message)
		}

		platforms := Platform(platform)
		commentObj.Platforms = platforms
		if parentId != "" {
			err := AddCommentReply(commentObj, parentId)
			if err != nil {
				return commentOutput, err
			}
		} else { //insert new comment
			err := AddCommentMain(commentObj)
			if err != nil {
				return commentOutput, err
			}
		}

		//push kafka recommandation
		var TrackingData TrackingDataStruct
		TrackingUserData, _ := recommendation.GetTracking(c, contentId, TrackingData)
		recommendation.PushDataToKafka(COMMENT, TrackingUserData)

		//convert to comment ouput
		dataByte, _ := json.Marshal(commentObj)
		err = json.Unmarshal(dataByte, &commentOutput)
		return commentOutput, err
	}

	return commentOutput, errors.New("Vui lòng không spam bình luận")

}

func CheckRuleMessage(message string) (string, error) {

	//không được trống và nhiều hơn 300 kí tự
	if message == "" || utf8.RuneCountInString(message) > 300 {
		return "", errors.New("Bình luận không hợp lệ hoặc chứa hơn 300 kí tự")
	}

	//không được tồn tại thẻ html
	if m, _ := regexp.MatchString("\\<[\\S\\s]+?\\>", message); m {
		return "", errors.New("Bình luận chứa kí tự không hợp lệ")
	}
	//không được tồn tại link
	if m, _ := regexp.MatchString("(?:(?:https?|ftp):\\/\\/)?[\\w/\\-?=%.]+\\.[\\w/\\-?=%.]+", message); m {
		return "", errors.New("Bình luận chứa kí tự không hợp lệ")
	}

	//EscapeString
	return html.EscapeString(message), nil
}
func AddCommentReply(commentObj CommentObjectStruct, parentId string) error {
	var keyCacheZrange = PREFIX_REDIS_ZRANGE_COMMENT + "_" + commentObj.Content_id
	var keyCacheHash = PREFIX_REDIS_HASH_COMMENT + "_" + commentObj.Content_id
	var listIDCmt = []string{parentId}

	listComments, err := GetCommentByListID(commentObj.Content_id, listIDCmt, 0, 10, true)
	if err != nil {
		return err
	}

	if len(listComments) == 0 {
		return errors.New("Bình luận không tồn tại")
	}
	go func(parentId string, commentObj CommentObjectStruct) {

		//connect db
		session, db, err := GetCollection()
		if err != nil {
			Sentry_log(err)
			return
		}
		defer session.Close()
		where := bson.M{"id": parentId}
		update := bson.M{
			"$push": bson.M{
				"reply": commentObj,
			},
		}
		err = db.C(COLLECTION_COMMENT).Update(where, update)
		if err != nil {
			log.Println("AddComment err", err)
			return
		}
	}(parentId, commentObj)

	//write cache zrange
	dataByte, _ := json.Marshal(commentObj)
	var tempComment CommentOutputStruct
	err = json.Unmarshal(dataByte, &tempComment)
	if err != nil {
		return err
	}
	listComments[0].Reply = append(listComments[0].Reply, tempComment)

	// clear cache zrange
	mRedis.HDel(keyCacheHash, parentId)
	mRedis.ZRem(keyCacheZrange, parentId)

	dataByte, _ = json.Marshal(listComments[0])
	mRedis.HSet(keyCacheHash, parentId, string(dataByte))
	mRedis.ZAdd(keyCacheZrange, float64(listComments[0].Created_at), parentId)

	return nil
}
func AddCommentMain(commentObj CommentObjectStruct) error {
	var keyCacheZrange = PREFIX_REDIS_ZRANGE_COMMENT + "_" + commentObj.Content_id
	var keyCacheHash = PREFIX_REDIS_HASH_COMMENT + "_" + commentObj.Content_id

	go func(commentObj CommentObjectStruct) {

		//connect db
		session, db, err := GetCollection()
		if err != nil {
			Sentry_log(err)
			log.Println("AddComment err", err)
			return
		}
		defer session.Close()

		err = db.C(COLLECTION_COMMENT).Insert(commentObj)
		if err != nil {
			log.Println("AddComment err", err)
			return
		}

	}(commentObj)

	//write cache zrange
	dataByte, _ := json.Marshal(commentObj)
	mRedis.HSet(keyCacheHash, commentObj.Id, string(dataByte))
	mRedis.ZAdd(keyCacheZrange, float64(commentObj.Created_at), commentObj.Id)
	return nil
}
func GetCommentsByContentId(contentID string, page, limit int, cacheActive bool) CommentOutputPageStruct {
	var dataPage CommentOutputPageStruct

	listCommentID := GetListIDCommentByContentId(contentID, page, limit, cacheActive)
	dataComment, err := GetCommentByListID(contentID, listCommentID, page, limit, cacheActive)
	if err != nil {
		return dataPage
	}
	dataPage.Items = dataComment
	dataPage.Metadata.Page = page
	dataPage.Metadata.Limit = limit
	dataPage.Metadata.Total = GetTotalCommentByContentId(contentID)
	return dataPage

}

func GetListIDCommentByContentId(contentID string, page, limit int, cacheActive bool) []string {
	var listCommentID []string
	var keyCache = PREFIX_REDIS_ZRANGE_COMMENT + "_" + contentID
	var start = page * limit
	var stop = (start + limit) - 1

	if cacheActive {
		vals := mRedis.ZRevRange(keyCache, int64(start), int64(stop))
		return vals
	}

	dataComments, err := GetDataInMongo(contentID, listCommentID, page, limit)
	if err != nil {
		return listCommentID
	}

	for _, val := range dataComments {
		listCommentID = append(listCommentID, val.Id)
	}

	return listCommentID
}

func GetCommentByListID(contentID string, listID []string, page, limit int, cacheActive bool) ([]CommentOutputStruct, error) {
	var keyCacheHash = PREFIX_REDIS_HASH_COMMENT + "_" + contentID
	var commentsOutput = make([]CommentOutputStruct, 0)

	if len(listID) == 0 {
		return commentsOutput, errors.New("Không có bình luận nào")
	}

	if cacheActive {
		dataRedis, err := mRedis.HMGet(keyCacheHash, listID)
		if err == nil && len(dataRedis) > 0 {
			for _, val := range dataRedis {
				if str, ok := val.(string); ok {
					var commentObj CommentOutputStruct
					err = json.Unmarshal([]byte(str), &commentObj)
					if err != nil {
						continue
					}
					commentsOutput = append(commentsOutput, commentObj)
				}
			}
		}
		return commentsOutput, err

	}

	// if data in cache != total listID then get data from database
	if len(commentsOutput) != len(listID) {
		listObj, err := GetDataInMongo(contentID, listID, page, limit)
		if err != nil && err.Error() != "not found" {
			return commentsOutput, err
		}
		commentsOutput = listObj
	}

	return commentsOutput, nil
}

func GetDataInMongo(contentID string, listID []string, page, limit int) ([]CommentOutputStruct, error) {
	var keyCacheZrange = PREFIX_REDIS_ZRANGE_COMMENT + "_" + contentID
	var keyCacheHash = PREFIX_REDIS_HASH_COMMENT + "_" + contentID
	var listCommentsOutput = make([]CommentOutputStruct, 0)
	var listComments = make([]CommentObjectStruct, 0)
	// Connect DB
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return listCommentsOutput, err
	}
	defer session.Close()

	where := bson.M{
		"content_id": contentID,
		"status":     bson.M{"$ne": -1},
		// Thu DT-9922
		// "status":     1,
	}
	if len(listID) > 0 {
		where["id"] = bson.M{"$in": listID}

		//set page = 0 nếu truyền listid
		page = 0
	}
	err = db.C(COLLECTION_COMMENT).Find(where).Sort("-pin", "-created_at").Skip(page * limit).Limit(limit).All(&listComments)
	if err != nil && err.Error() != "not found" {
		Sentry_log(err)
		return listCommentsOutput, err
	}

	//remove cache if not exists
	if len(listID) > 0 && len(listComments) != len(listID) {
		for _, id := range listID {
			found := false
			for i := range listComments {
				if listComments[i].Id == id {
					found = true
					break
				}
			}
			if found == false {
				mRedis.HDel(keyCacheHash, id)
				mRedis.ZRem(keyCacheZrange, id)
			}
		}
	}

	//check status comment reply
	for i, comment := range listComments {
		var tempComment []CommentObjectStruct
		for _, reply := range comment.Reply {
			if reply.Status == 1 {
				tempComment = append(tempComment, reply)
			}
		}
		listComments[i].Reply = tempComment
	}

	dataByte, _ := json.Marshal(listComments)
	err = json.Unmarshal(dataByte, &listCommentsOutput)
	if err != nil {
		return listCommentsOutput, nil
	}
	//Write cache
	for _, comment := range listCommentsOutput {
		// clear cache zrange
		mRedis.HDel(keyCacheHash, comment.Id)
		mRedis.ZRem(keyCacheZrange, comment.Id)

		//set cache
		dataByte, _ := json.Marshal(comment)
		mRedis.HSet(keyCacheHash, comment.Id, string(dataByte))
		if comment.Pin == 1 {
			mRedis.ZAdd(keyCacheZrange, float64(comment.Created_at)*2, comment.Id)
		} else {
			mRedis.ZAdd(keyCacheZrange, float64(comment.Created_at), comment.Id)
		}

	}

	return listCommentsOutput, nil
}

func GetTotalCommentByContentId(contentID string) int {
	var keyCache = PREFIX_REDIS_ZRANGE_COMMENT + "_" + contentID
	total := mRedis.ZCountAll(keyCache)
	return int(total)
}

func DelComment(userId, contentId, commentId, parentId string, status int) error {

	// Thu DT-9922 thêm status động
	if ok, _ := In_array(status, LISTSTATUSCOMMENT); !ok {
		return errors.New("Status not correct")
	}

	var keyCacheZrange = PREFIX_REDIS_ZRANGE_COMMENT + "_" + contentId
	var keyCacheHash = PREFIX_REDIS_HASH_COMMENT + "_" + contentId

	var listID = []string{commentId}
	if parentId != "" {
		listID = []string{parentId}
	}

	dataCommentMain, err := GetCommentByListID(contentId, listID, 0, 10, true)

	if err != nil {
		return err
	}

	if len(dataCommentMain) == 0 {
		return errors.New("Bình luận không tồn tại")
	}

	dataComment := dataCommentMain[0]

	//Del comment in reply cmt
	if listID[0] == parentId {
		found := false
		//check comment user
		for _, val := range dataComment.Reply {
			if val.Id == commentId {
				if val.User.Id == userId {
					found = true
				}
				break
			}
		}

		if found == true {

			//async
			go func(commentId string) {
				//connect db
				session, db, err := GetCollection()
				if err != nil {
					Sentry_log(err)
					return
				}
				defer session.Close()

				where := bson.M{
					// "id":            parentId,
					"reply.id": commentId,
					// "reply.user.id": userId,
				}
				update := bson.M{
					"$set": bson.M{
						"reply.$.status": status,
					},
				}

				err = db.C(COLLECTION_COMMENT).Update(where, update)
				if err != nil {
					log.Println("DelComment err", err)
					return
				}
			}(commentId)

			//clear cache
			mRedis.ZRem(keyCacheZrange, listID[0])
			return nil
		}
	} else if dataComment.User.Id == userId { //check comment user

		//async
		go func(userId, commentId string) {
			//connect db
			session, db, err := GetCollection()
			if err != nil {
				Sentry_log(err)
				return
			}
			defer session.Close()

			where := bson.M{
				"id":      commentId,
				"user.id": userId,
			}
			update := bson.M{
				"$set": bson.M{
					"status": status,
				},
			}
			err = db.C(COLLECTION_COMMENT).Update(where, update)
			if err != nil {
				return
			}
		}(userId, commentId)

		//status = 1 => comment đã được duyệt
		//status = 2 => ẩn comment. Chỉ hiển thị với user này
		if status == 2 || status == 1 {
			dataComment.Status = status
			dataByte, _ := json.Marshal(dataComment)
			mRedis.HDel(keyCacheHash, commentId)
			mRedis.HSet(keyCacheHash, commentId, string(dataByte))

		} else { //ẩn comment với tất cả
			//clear cache
			mRedis.HDel(keyCacheHash, listID[0])
			mRedis.ZRem(keyCacheZrange, listID[0])
		}

		return nil

	}

	return errors.New("Bình luận không tồn tại")
}

func PinComment(contentId string, commentId string, status int) error {
	var listCommentID []string
	listCommentID = append(listCommentID, commentId)

	// get comment pin = 1
	var listComments = make([]CommentObjectStruct, 0)
	// Connect DB
	_, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
	}
	whereGet := bson.M{
		"content_id": contentId,
		"pin":        1,
	}
	err = db.C(COLLECTION_COMMENT).Find(whereGet).Skip(0).Limit(20).All(&listComments)
	if err != nil && err.Error() != "not found" {
		Sentry_log(err)
	}
	if len(listComments) > 0 {
		for i := range listComments {
			listCommentID = append(listCommentID, listComments[i].Id)
		}
	}

	//async update
	go func(contentId string) {
		//connect db
		session, db, err := GetCollection()
		if err != nil {
			Sentry_log(err)
			return
		}
		defer session.Close()

		where := bson.M{
			"content_id": contentId,
			"pin":        1,
		}
		update := bson.M{
			"$set": bson.M{
				"pin": 0,
			},
		}
		err = db.C(COLLECTION_COMMENT).Update(where, update)
		if err == nil && status == 0 {
			_, _ = GetDataInMongo(contentId, listCommentID, 0, 50)

		}
	}(contentId)

	if status == 1 {
		// Update
		go func(contentId string, commentId string, listCommentID []string) {
			//connect db
			session, db, err := GetCollection()
			if err != nil {
				Sentry_log(err)
				return
			}
			defer session.Close()

			where := bson.M{
				"content_id": contentId,
				"id":         commentId,
			}
			update := bson.M{
				"$set": bson.M{
					"pin": 1,
				},
			}
			err = db.C(COLLECTION_COMMENT).Update(where, update)
			if err == nil {
				_, _ = GetDataInMongo(contentId, listCommentID, 0, 50)

			}
		}(contentId, commentId, listCommentID)
	}
	return nil
}

func AddComment(commentObj CommentObjectStruct) error {
	var keyCacheZrange = PREFIX_REDIS_ZRANGE_COMMENT + "_" + commentObj.Content_id
	var keyCacheHash = PREFIX_REDIS_HASH_COMMENT + "_" + commentObj.Content_id

	go func(commentObj CommentObjectStruct) {

		//connect db
		session, db, err := GetCollection()
		if err != nil {
			Sentry_log(err)
			log.Println("AddComment err", err)
			return
		}
		defer session.Close()

		err = db.C(COLLECTION_COMMENT).Insert(commentObj)
		if err != nil {
			log.Println("AddComment err", err)
			return
		}

	}(commentObj)

	//write cache zrange
	dataByte, _ := json.Marshal(commentObj)
	mRedis.HSet(keyCacheHash, commentObj.Id, string(dataByte))
	mRedis.ZAdd(keyCacheZrange, float64(commentObj.Created_at), commentObj.Id)
	return nil
}

// Thu DT-9922
var (
	host, _         = CommonConfig.GetString("REDIS", "host")
	port, _         = CommonConfig.GetString("REDIS", "port")
	c, err          = redis.Dial("tcp", host+":"+port)
)


func FilterMessage(message string, id string) string {

	// Split String to array word
	terms := strings.Fields(strings.ToLower(message))

	// Check redis blackwords
	loop := CheckRedisBadWord("bad-words", "normal")
	// Filter msg
	for j := loop; j > 1; j-- {
		var input1 []string
		for i := 0; i < len(terms)-1; i++ {
			//
			if i+j < len(terms) {
				v := terms[i]
				for k := 1; k < j; k++ {
					v = v + " " + terms[i+k]
				}
				input1 = append(input1, v)
			}
		}
		message = SubFilterMessage(input1, " "+message+" ", id)
	}

	message = SubFilterMessage(terms, " "+message+" ", id)

	// Loai bo khoang trang dau cuoi
	return strings.TrimSpace(message)

}

func SubFilterMessage(input []string, msg string, id string) string {

	if err != nil {
		return msg
	}

	name_redis_user_word := "user-words-" + id

	c.Send("MULTI")
	c.Send("DEL", name_redis_user_word)
	c.Send("SADD", redis.Args{}.Add(name_redis_user_word).AddFlat(input)...)
	mRedis.Expire(name_redis_user_word, 120)
	c.Send("SINTER", name_redis_user_word, "bad-words_normal")

	reply, err := c.Do("EXEC")

	if err != nil {
		return msg
	}

	values, _ := redis.Values(reply, nil)
	curse_words, err := redis.Strings(values[2], nil)
	mRedis.Del(name_redis_user_word)

	if err != nil {
		return msg
	}
	if (len(curse_words)) > 0 {
		msg = strings.ToLower(msg)
		for _, v := range curse_words {
			msg = strings.Replace(msg, v, "***", 7)
		}
	}

	return msg
}

func CheckRedisBadWord(key string, typeword string) int {
	if err != nil {
		return 0
	}

	key = key + "_" + typeword
	loop := 0
	status := mRedis.Exists(key)
	statusCount := mRedis.Exists(key + "-count")

	if status <= 0 || statusCount <= 0 {
		mRedis.Del(key)
		mRedis.Del(key + "-count")
		// Get from Db add redis
		var listBadwords = make([]BadwordObjectStruct, 0)
		// Connect DB
		session, db, err := GetCollection()
		if err != nil {
			Sentry_log(err)
			return loop
		}
		defer session.Close()

		where := bson.M{
			"status": 1,
			"type":   0,
		}

		if typeword == "like" {
			where = bson.M{
				"status": 1,
				"type":   bson.M{"$ne": 0},
			}
		}

		err = db.C(COLLECTION_BLACKLIST).Find(where).All(&listBadwords)
		if err != nil && err.Error() != "not found" {
			Sentry_log(err)
			return loop
		}

		// TH : Ghi redis array cho key có *
		if typeword == "like" {
			if len(listBadwords) > 0 {
				var input []string
				for i := range listBadwords {
					// Count words
					word := strings.ToLower(listBadwords[i].Content)
					if len(strings.Fields(word)) > loop {
						loop = len(strings.Fields(word))
					}
					input = append(input, word+"[]"+fmt.Sprint(listBadwords[i].Type))
				}
				// Set redis

				c.Send("SADD", redis.Args{}.Add(key).AddFlat(input)...)
				mRedis.SetInt(key+"-count", loop, -1)

			}
		} else { // TH : Ghi redis array cho key chính xác
			if len(listBadwords) > 0 {
				var input []string
				for i := range listBadwords {
					// Count words
					word := strings.ToLower(listBadwords[i].Content)
					if len(strings.Fields(word)) > loop {
						loop = len(strings.Fields(word))
					}
					input = append(input, word)
				}
				// Set redis
				// c.Send("SADD", redis.Args{}.Add(key).AddFlat(input)...)
				c.Send("SADD", redis.Args{}.Add(key).AddFlat(input)...)
				mRedis.SetInt(key+"-count", loop, -1)

			}
		}

	} else {
		// Get value loop
		loop, _ = mRedis.GetInt(key + "-count")
	}

	return loop
}

func FilterMessageLike(inputStr string) string {

	loop := CheckRedisBadWord("bad-words", "like")
	var notAllowed []string
	if loop > 0 {
		// Get redis
		c.Send("MULTI")
		c.Send("Smembers", "bad-words_like")
		// c.Send("GET", "bad-words_like")
		reply, _ := c.Do("EXEC")
		values, _ := redis.Values(reply, nil)
		if len(values) > 0 {
			notAllowed, _ = redis.Strings(values[0], nil)	
		}
	}

	// convert str into a slice
	strSlice := strings.Fields(inputStr)

	for j := loop; j > 1; j-- {
		var input1 []string
		for i := 0; i < len(strSlice)-1; i++ {
			//
			if i+j < len(strSlice) {
				v := strSlice[i]
				for k := 1; k < j; k++ {
					v = v + " " + strSlice[i+k]
				}
				input1 = append(input1, v)
			}
		}
		inputStr = CensorWord(inputStr, input1, notAllowed)
	}

	inputStr = CensorWord(inputStr, strSlice, notAllowed)
	return inputStr

}

func CensorWord(str string, strSlice []string, censored []string) string {
	// var newSlice []string
	str = " " + str + " "
	// check for empty slice
	if len(censored) <= 0 {
		return strings.TrimSpace(str)
	}

	//check each words in strSlice against censored slice
	for _, word := range strSlice {
		for _, forbiddenWord := range censored {
			// NOTE : change between Index and EqualFold to see the different result
			arr_forbiddenWord := strings.Split(forbiddenWord, "[]")
			if len(arr_forbiddenWord) == 2 {
				if index := strings.Index(strings.ToLower(word), arr_forbiddenWord[0]); index > -1 {
					if (arr_forbiddenWord[1] == "3" && index == 0) || (arr_forbiddenWord[1] == "1" && len(word) == len(arr_forbiddenWord[0])+index) || (arr_forbiddenWord[1] == "2") {
						str = strings.Replace(str, word, "***", 7)
					}

				}
			}

		}
	}
	return strings.TrimSpace(str)
}

func MakeFirstUpperCase(s string) string {
	if len(s) < 2 {
		return strings.ToLower(s)
	}

	bts := []byte(s)
	lc := bytes.ToUpper([]byte{bts[0]})
	rest := bts[1:]

	return string(bytes.Join([][]byte{lc, rest}, nil))
}

func FilterByUserID(dataComment CommentOutputPageStruct, page int, limit int, userId string) CommentOutputPageStruct {

	var comments = make([]CommentOutputStruct, 0)

	// stt := 0
	// start := page * limit
	// if page > 0 {
	// 	start++
	// }
	From(dataComment.Items).Where(func(c interface{}) bool {
		tam := 0
		if (c.(CommentOutputStruct).Status == 1) || (c.(CommentOutputStruct).Status == 2 && c.(CommentOutputStruct).User.Id == userId) {
			tam = 1
		}
		return tam == 1
	}).Select(func(c interface{}) interface{} {
		return c.(CommentOutputStruct)

	}).ToSlice(&comments)
	dataComment.Items = comments
	dataComment.Metadata.Page = page
	dataComment.Metadata.Limit = limit
	return dataComment

}
