package module

import (
	"errors"
	"strconv"
	"time"
	"github.com/go-redis/redis"
)

var (
	redisDbStand     	*redis.Client
)

type RedisStandaloneModelStruct int

func init() {
	host, _ := CommonConfig.GetString("REDIS_STANDALONE", "host")
	port, _ := CommonConfig.GetString("REDIS_STANDALONE", "port")
	RedisServer := host + ":" + port

	redisDbStand = redis.NewClient(&redis.Options{
		Addr:         RedisServer,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	if redisDbStand.Incr("REDIS_PING").Val() <= 0 {
		Sentry_log_with_msg("APIv3 redisDbStand standalone not valid")
	}
}

func (this *RedisStandaloneModelStruct) GetString(key string) (string, error) {
	val, err := redisDbStand.Get(key).Result()
	return val, err
}

func (this *RedisStandaloneModelStruct) SetString(key string, value string, ttl time.Duration) error {
	ttl = ttl * time.Second
	err := redisDbStand.Set(key, value, ttl).Err()
	return err
}

func (this *RedisStandaloneModelStruct) Incr(key string) int64 {
	val := redisDbStand.Incr(key).Val()
	return val
}

func (this *RedisStandaloneModelStruct) IncrBy(key string, value int64) int64 {
	val := redisDbStand.IncrBy(key, value).Val()
	return val
}

func (this *RedisStandaloneModelStruct) GetInt(key string) (int, error) {
	val, err := redisDbStand.Get(key).Result()
	if err != nil {
		return 0, err
	}

	if i, err := strconv.Atoi(val); err == nil {
		return i, nil
	}
	return 0, errors.New("GetInt error: " + key)
}

func (this *RedisStandaloneModelStruct) SetInt(key string, value int, ttl time.Duration) error {
	ttl = ttl * time.Second
	err := redisDbStand.Set(key, value, ttl).Err()
	return err
}

func (this *RedisStandaloneModelStruct) HSet(key string, field string, value interface{}) error {
	err := redisDbStand.HSet(key, field, value).Err()
	return err
}

func (this *RedisStandaloneModelStruct) HGet(key string, field string) (interface{}, error) {
	val, err := redisDbStand.HGet(key, field).Result()
	return val, err
}

func (this *RedisStandaloneModelStruct) HDel(key string, field string) int64 {
	val := redisDbStand.HDel(key, field).Val()
	return val
}

func (this *RedisStandaloneModelStruct) HMGet(key string, fields []string) ([]interface{}, error) {
	val, err := redisDbStand.HMGet(key, fields...).Result()
	return val, err
}

func (this *RedisStandaloneModelStruct) ZRangeByScore(key string, offset int64, limit int64) ([]string, error) {
	vals, err := redisDbStand.ZRangeByScore(key, redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: offset,
		Count:  limit,
	}).Result()
	return vals, err
}

func (this *RedisStandaloneModelStruct) ZRevRange(key string, start int64, stop int64) []string {
	vals := redisDbStand.ZRevRange(key, start, stop).Val()
	return vals
}

func (this *RedisStandaloneModelStruct) ZRange(key string, start, stop int64) ([]string, error) {
	vals, err := redisDbStand.ZRange(key, start, stop).Result()
	return vals, err
}

func (this *RedisStandaloneModelStruct) ZRem(key string, value string) int64 {
	val := redisDbStand.ZRem(key, value).Val()
	return val
}

func (this *RedisStandaloneModelStruct) ZAdd(key string, score float64, value string) error {
	err := redisDbStand.ZAdd(key, redis.Z{
		Score:  score,
		Member: value,
	}).Err()
	return err
}

func (this *RedisStandaloneModelStruct) ZRank(key string, value string) int64 {
	val := redisDbStand.ZRank(key, value).Val()
	return val
}

func (this *RedisStandaloneModelStruct) ZCountAll(key string) int64 {
	val := redisDbStand.ZCount(key, "-inf", "+inf").Val()
	return val
}

func (this *RedisStandaloneModelStruct) Del(key string) int64 {
	val := redisDbStand.Del(key).Val()
	return val
}

func (this *RedisStandaloneModelStruct) Exists(key string) int64 {
	val := redisDbStand.Exists(key).Val()
	return val
}

func (this *RedisStandaloneModelStruct) Expire(key string, ttl time.Duration) error {
	ttl = ttl * time.Second
	err := redisDbStand.Expire(key, ttl).Err()
	return err
}
