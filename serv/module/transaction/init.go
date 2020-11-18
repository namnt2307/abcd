package transaction

import (
	. "cm-v5/serv/module"
	jsoniter "github.com/json-iterator/go"
)

var mRedisKV RedisKVModelStruct
var mRedis RedisModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary
