package topviews

import (
	// . "cm-v5/schema"
	. "cm-v5/serv/module"

	jsoniter "github.com/json-iterator/go"
)

var mRedis RedisModelStruct
var mRedisKV RedisKVModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary
