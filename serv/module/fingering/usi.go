package fingering

import (
	"errors"
	"fmt"
	"time"

	. "cm-v5/serv/module"
)

var (
	mRedis           RedisModelStruct
	mRedisStandalone RedisStandaloneModelStruct
	MAX_INDEX int64 = 100000000
	MAX_CCU_STREAM int = 1
)

func GetUniqueStreamingID(user_id string) (string , error) {

	// check user is stream K+
	maxCcuStream, _ := mRedis.GetInt("MAX_ACTIVE_STREAMING_LIMIT")
	if maxCcuStream > 0 {
		MAX_CCU_STREAM = maxCcuStream
	}
	numStreaming, _ := mRedisStandalone.GetInt("ACTIVE_STREAMING_" + user_id)

	fmt.Println("ACTIVE_STREAMING_"+user_id, numStreaming)
	fmt.Println("MAX_CCU_STREAM", MAX_CCU_STREAM)

	
	if numStreaming >= MAX_CCU_STREAM {
		return "" , errors.New(fmt.Sprint(numStreaming) + " player Is Streaming")
	}

	USI_INDEX := mRedis.Incr("USI_INDEX")
	if USI_INDEX >= MAX_INDEX {
		mRedis.Del("USI_INDEX")
	}

	// code => User
	code := fmt.Sprintf("%X", USI_INDEX)
	mRedis.SetString("USI_USER_" + code , user_id , 3600)

	// User => code list
	var curTime = time.Now().Unix()
	mRedis.ZAdd("Z_USER_USI_"+user_id, float64(curTime), code)
	mRedis.Expire("Z_USER_USI_"+user_id, 12*3600)

	return code, nil
}

func KeepUniqueStreamingID(code, user_id string) error {
	valStr, _ := mRedis.GetString("USI_USER_" + code)
	if valStr != user_id {
		return errors.New("User not valid")
	}
	mRedis.Expire("USI_USER_" + code , 3600)
	return nil
}

func KeepStreamingByUserID(user_id string) error {
	mRedisStandalone.Expire("ACTIVE_STREAMING_" + user_id , 300)
	return nil
}
