package redis

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/kundan1301/message-broker/config"
)

var redisClient *redis.Client
var redisClusterClient *redis.ClusterClient

type CustomRedisClient interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

func InitClient(config *config.Config) {
	if config.UseRedisCluster {
		redisClusterClient = redis.NewClusterClient(&config.RedisClusterOptions)
	} else {
		redisClient = redis.NewClient(&config.RedisOptions)
	}
}

func GetRedisClient() CustomRedisClient {
	if redisClient == nil {
		return redisClusterClient
	}
	return redisClient
}
