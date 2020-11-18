package module

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"github.com/go-redis/redis"
)

var (
	KEY_VERSION string = "V5.1_"
	// redisdb     	*redis.Client
	redisDbRead  *redis.ClusterClient
	redisDbWrite *redis.ClusterClient
)

type RedisModelStruct int

func init() {
	hostsRead, _ := CommonConfig.GetString("REDIS", "hostRead")
	addrHostsRead := strings.Split(hostsRead, ",")
	redisDbRead = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:         addrHostsRead,
		DialTimeout:   10 * time.Second,
		ReadTimeout:   30 * time.Second,
		WriteTimeout:  30 * time.Second,
		RouteRandomly: true,
	})

	hostsWrite, _ := CommonConfig.GetString("REDIS", "hostWrite")
	addrHostsWrite := strings.Split(hostsWrite, ",")
	redisDbWrite = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        addrHostsWrite,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	// check Redis
	if redisDbRead.Incr("REDIS_PING").Val() <= 0 {
		Sentry_log_with_msg("APIv3 redisDbRead not valid")
	}

	if redisDbWrite.Incr("REDIS_PING").Val() <= 0 {
		Sentry_log_with_msg("APIv3 redisDbWrite not valid")
	}
}

func (this *RedisModelStruct) GetString(key string) (string, error) {
	// check key SUB
	if strings.Contains(key, "USC_LIST_SUBCRIPTION_TEMP_") != true {
		key = KEY_VERSION + key
	}
	val, err := redisDbRead.Get(key).Result()
	return val, err
}

func (this *RedisModelStruct) SetString(key string, value string, ttl time.Duration) error {
	// check key SUB
	if strings.Contains(key, "USC_LIST_SUBCRIPTION_TEMP_") != true {
		key = KEY_VERSION + key
	}
	ttl = ttl * time.Second
	err := redisDbWrite.Set(key, value, ttl).Err()
	return err
}

func (this *RedisModelStruct) Incr(key string) int64 {
	key = KEY_VERSION + key
	val := redisDbWrite.Incr(key).Val()
	return val
}

func (this *RedisModelStruct) IncrBy(key string, value int64) int64 {
	key = KEY_VERSION + key
	val := redisDbWrite.IncrBy(key, value).Val()
	return val
}

func (this *RedisModelStruct) GetInt(key string) (int, error) {
	key = KEY_VERSION + key
	val, err := redisDbRead.Get(key).Result()
	if err != nil {
		return 0, err
	}

	if i, err := strconv.Atoi(val); err == nil {
		return i, nil
	}
	return 0, errors.New("GetInt error: " + key)
}

func (this *RedisModelStruct) SetInt(key string, value int, ttl time.Duration) error {
	key = KEY_VERSION + key
	ttl = ttl * time.Second
	err := redisDbWrite.Set(key, value, ttl).Err()
	return err
}

func (this *RedisModelStruct) HSet(key string, field string, value interface{}) error {
	key = KEY_VERSION + key
	err := redisDbWrite.HSet(key, field, value).Err()
	return err
}

func (this *RedisModelStruct) HGet(key string, field string) (interface{}, error) {
	key = KEY_VERSION + key
	val, err := redisDbRead.HGet(key, field).Result()
	return val, err
}

func (this *RedisModelStruct) HDel(key string, field string) int64 {
	key = KEY_VERSION + key
	val := redisDbWrite.HDel(key, field).Val()
	return val
}

func (this *RedisModelStruct) HMGet(key string, fields []string) ([]interface{}, error) {
	key = KEY_VERSION + key
	val, _ := redisDbRead.HMGet(key, fields...).Result()
	return val, nil
}

func (this *RedisModelStruct) ZRangeByScore(key string, offset int64, limit int64) ([]string, error) {
	key = KEY_VERSION + key
	vals, err := redisDbRead.ZRangeByScore(key, redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: offset,
		Count:  limit,
	}).Result()
	return vals, err
}

func (this *RedisModelStruct) ZRevRange(key string, start int64, stop int64) []string {
	key = KEY_VERSION + key
	vals := redisDbRead.ZRevRange(key, start, stop).Val()
	return vals
}

func (this *RedisModelStruct) ZRange(key string, start, stop int64) ([]string, error) {
	key = KEY_VERSION + key
	vals, err := redisDbRead.ZRange(key, start, stop).Result()
	return vals, err
}

func (this *RedisModelStruct) ZRem(key string, value string) int64 {
	key = KEY_VERSION + key
	val := redisDbWrite.ZRem(key, value).Val()
	return val
}

func (this *RedisModelStruct) ZAdd(key string, score float64, value string) error {
	key = KEY_VERSION + key
	err := redisDbWrite.ZAdd(key, redis.Z{
		Score:  score,
		Member: value,
	}).Err()
	return err
}

func (this *RedisModelStruct) ZRank(key string, value string) int64 {
	key = KEY_VERSION + key
	val := redisDbRead.ZRank(key, value).Val()
	return val
}

func (this *RedisModelStruct) ZCountAll(key string) int64 {
	key = KEY_VERSION + key
	val := redisDbRead.ZCount(key, "-inf", "+inf").Val()
	return val
}

func (this *RedisModelStruct) Del(key string) int64 {
	key = KEY_VERSION + key
	val := redisDbWrite.Del(key).Val()
	return val
}

func (this *RedisModelStruct) Exists(key string) int64 {
	key = KEY_VERSION + key
	val := redisDbRead.Exists(key).Val()
	return val
}

func (this *RedisModelStruct) Expire(key string, ttl time.Duration) error {
	key = KEY_VERSION + key
	ttl = ttl * time.Second
	err := redisDbRead.Expire(key, ttl).Err()
	return err
}
