package module

import (
	"errors"
	"strconv"
	"time"
	"strings"
	"github.com/go-redis/redis"
)

var (
	// redisdbUSC     *redis.Client
	redisDbUscRead     *redis.ClusterClient
	redisDbUscWrite    *redis.ClusterClient
	// RedisServerUSC string
)

type RedisUSCModelStruct int

func init() {
	hostsRead, _ := CommonConfig.GetString("REDIS_USC", "hostRead")
	addrHostsRead := strings.Split(hostsRead, ",")
	redisDbUscRead = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:         addrHostsRead,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		RouteRandomly: true,
	})

	hostsWrite, _ := CommonConfig.GetString("REDIS_USC", "hostWrite")
	addrHostsWrite := strings.Split(hostsWrite, ",")
	redisDbUscWrite = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:         addrHostsWrite,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	// check Redis 
	if redisDbUscRead.Incr("REDIS_PING").Val() <= 0 {
		Sentry_log_with_msg("APIv3 redisDbUscRead not valid")
	}

	if redisDbUscWrite.Incr("REDIS_PING").Val() <= 0 {
		Sentry_log_with_msg("APIv3 redisDbUscWrite not valid")
	}

}

func (this *RedisUSCModelStruct) GetString(key string) (string, error) {
	val, err := redisDbUscRead.Get(key).Result()
	return val, err
}

func (this *RedisUSCModelStruct) SetString(key string, value string, ttl time.Duration) error {
	ttl = ttl * time.Second
	err := redisDbUscWrite.Set(key, value, ttl).Err()
	return err
}

func (this *RedisUSCModelStruct) Incr(key string) int64 {
	val := redisDbUscWrite.Incr(key).Val()
	return val
}

func (this *RedisUSCModelStruct) IncrBy(key string, value int64) int64 {
	val := redisDbUscWrite.IncrBy(key, value).Val()
	return val
}

func (this *RedisUSCModelStruct) GetInt(key string) (int, error) {
	val, err := redisDbUscRead.Get(key).Result()
	if err != nil {
		return 0, err
	}

	if i, err := strconv.Atoi(val); err == nil {
		return i, nil
	}
	return 0, errors.New("GetInt error: " + key)
}

func (this *RedisUSCModelStruct) SetInt(key string, value int, ttl time.Duration) error {
	ttl = ttl * time.Second
	err := redisDbUscWrite.Set(key, value, ttl).Err()
	return err
}

func (this *RedisUSCModelStruct) HDel(key string, value string) int64 {
	val := redisDbUscWrite.HDel(key, value).Val()
	return val
}

func (this *RedisUSCModelStruct) HSet(key string, field string, value interface{}) error {
	err := redisDbUscWrite.HSet(key, field, value).Err()
	return err
}

func (this *RedisUSCModelStruct) HGet(key string, field string) (string, error) {
	val, err := redisDbUscRead.HGet(key, field).Result()
	return val, err
}

func (this *RedisUSCModelStruct) HMGet(key string, fields []string) ([]interface{}, error) {
	val, err := redisDbUscRead.HMGet(key, fields...).Result()
	return val, err
}

func (this *RedisUSCModelStruct) HGetAll(key string) (map[string]string, error) {
	val, err := redisDbUscRead.HGetAll(key).Result()
	return val, err
}

func (this *RedisUSCModelStruct) ZRangeByScore(key string, offset int64, limit int64) ([]string, error) {
	vals, err := redisDbUscRead.ZRangeByScore(key, redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: offset,
		Count:  limit,
	}).Result()
	return vals, err
}

func (this *RedisUSCModelStruct) ZRevRange(key string, start int64, stop int64) []string {
	vals := redisDbUscRead.ZRevRange(key, start, stop).Val()
	return vals
}

func (this *RedisUSCModelStruct) ZRange(key string, start, stop int64) ([]string, error) {
	vals, err := redisDbUscRead.ZRange(key, start, stop).Result()
	return vals, err
}

func (this *RedisUSCModelStruct) ZAdd(key string, score float64, value string) error {
	err := redisDbUscWrite.ZAdd(key, redis.Z{
		Score:  score,
		Member: value,
	}).Err()
	return err
}

func (this *RedisUSCModelStruct) ZRem(key string, value string) int64 {
	val := redisDbUscWrite.ZRem(key, value).Val()
	return val
}

func (this *RedisUSCModelStruct) ZCountAll(key string) int64 {
	val := redisDbUscRead.ZCount(key, "-inf", "+inf").Val()
	return val
}

func (this *RedisUSCModelStruct) Exists(key string) int64 {
	val := redisDbUscRead.Exists(key).Val()
	return val
}

func (this *RedisUSCModelStruct) ZRank(key string, value string) int64 {
	val := redisDbUscRead.ZRank(key, value).Val()
	return val
}

func (this *RedisUSCModelStruct) ZScore(key string, value string) float64 {
	val := redisDbUscRead.ZScore(key, value).Val()
	return val
}

func (this *RedisUSCModelStruct) Del(key string) int64 {
	val := redisDbUscWrite.Del(key).Val()
	return val
}
