package packages_premium

import (
	. "cm-v5/serv/module"

	jsoniter "github.com/json-iterator/go"
)

var mRedisKV RedisKVModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary
