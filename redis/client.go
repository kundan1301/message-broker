package redis

import (
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/kundan1301/message-broker/config"
)

const (
	connectedBrokerIPKey = "connectedBrokerIP"
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

// return previous broker if connected
func CheckPrevConn(clientID string) string {
	key := connectedBrokerIPKey + "/" + clientID
	return GetRedisClient().Get(key).Val()
}

func SetNewConnInfo(clientID, node string) {
	key := connectedBrokerIPKey + "/" + clientID
	log.Println(clientID, node)
	err := GetRedisClient().Set(key, node, 0).Err()
	if err != nil {
		log.Println("Error in setting node", err)
	}
}
