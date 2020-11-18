package module

import (
	"strings"
	"time"
	"github.com/go-redis/redis"
)

var (
	// redisdbKV     *redis.Client
	redisDbKVRead  *redis.ClusterClient
	redisDbKVWrite *redis.ClusterClient
)

type RedisKVModelStruct int

func init() {
	hostsRead, _ := CommonConfig.GetString("REDIS_KV", "hostRead")
	addrHostsRead := strings.Split(hostsRead, ",")
	redisDbKVRead = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:         addrHostsRead,
		DialTimeout:   10 * time.Second,
		ReadTimeout:   30 * time.Second,
		WriteTimeout:  30 * time.Second,
		RouteRandomly: true,
	})

	hostsWrite, _ := CommonConfig.GetString("REDIS_KV", "hostWrite")
	addrHostsWrite := strings.Split(hostsWrite, ",")
	redisDbKVWrite = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        addrHostsWrite,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
	
	// check Redis
	if redisDbKVRead.Incr("REDIS_PING").Val() <= 0 {
		Sentry_log_with_msg("APIv3 redisDbKVRead not valid")
	}

	if redisDbKVWrite.Incr("REDIS_PING").Val() <= 0 {
		Sentry_log_with_msg("APIv3 redisDbKVWrite not valid")
	}

}

func (this *RedisKVModelStruct) GetString(key string) (string, error) {
	// check key SUB
	if strings.Contains(key, "USC_LIST_SUBCRIPTION_TEMP_") != true {
		key = KEY_VERSION + key
	}

	val, err := redisDbKVRead.Get(key).Result()
	return val, err
}

func (this *RedisKVModelStruct) SetString(key string, value string, ttl time.Duration) error {
	// check key SUB
	if strings.Contains(key, "USC_LIST_SUBCRIPTION_TEMP_") != true {
		key = KEY_VERSION + key
	}
	
	ttl = ttl * time.Second
	err := redisDbKVWrite.Set(key, value, ttl).Err()
	return err
}

func (this *RedisKVModelStruct) Del(key string) int64 {
	key = KEY_VERSION + key
	val := redisDbKVWrite.Del(key).Val()
	return val
}

func (this *RedisKVModelStruct) Exists(key string) int64 {
	key = KEY_VERSION + key
	val := redisDbKVRead.Exists(key).Val()
	return val
}

func (this *RedisKVModelStruct) Incr(key string) int64 {
	// check key VIEW
	if strings.Contains(key, "VIEW_") != true {
		key = KEY_VERSION + key
	}

	val := redisDbKVWrite.Incr(key).Val()
	return val
}

func (this *RedisKVModelStruct) IncrBy(key string, value int64) int64 {
	// check key VIEW
	if strings.Contains(key, "VIEW_") != true {
		key = KEY_VERSION + key
	}
	
	val := redisDbKVWrite.IncrBy(key, value).Val()
	return val
}

func (this *RedisKVModelStruct) Expire(key string, ttl time.Duration) error {
	key = KEY_VERSION + key
	ttl = ttl * time.Second
	err := redisDbKVWrite.Expire(key, ttl).Err()
	return err
}
