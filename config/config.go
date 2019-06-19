package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"runtime"

	"github.com/go-redis/redis"

	customHttp "github.com/kundan1301/message-broker/http"
)

/*Config struct stores configs defined in config.json file
NodeIp will be used to communicate with other node in cluster/autoscaling group
ConfigUrl will be used to load config from that url, if its specified file config will be overridden
*/

type Config struct {
	Host                string               `json:"host"`
	MqttPort            string               `json:"mqttPort"`
	HttpPort            string               `json:"httpPort"`
	MqttTlsPort         string               `json:"mqttTlsPort"`
	HttpTlsPort         string               `json:"httpTlsPort"`
	NodeIP              string               `json:"nodeIP"`
	ConfigUrl           string               `json:"configUrl"`
	AuthUrl             string               `json:"authUrl"`
	SubscribeUrl        string               `json:"subscribeUrl"`
	PublishUrl          string               `json:"publishUrl"`
	RedisOptions        redis.Options        `json:"redisOptions"`
	RedisClusterOptions redis.ClusterOptions `json:"redisClusterOptions"`
	UseRedisCluster     bool                 `json:useRedisCluster`
}

func loadConfigFromHTTP(configURL string) (*Config, error) {
	body, err := customHttp.Get(configURL)
	if err != nil {
		log.Println("error in reading config resp: ", err)
		return nil, err
	}
	var config Config
	err = json.Unmarshal(body, &config)
	if err != nil {
		log.Println("Unmarshal http config  error: ", err)
		return nil, err
	}
	return &config, nil

}

func LoadConfig() (*Config, error) {
	_, dirname, _, _ := runtime.Caller(0)
	filePath := path.Join(filepath.Dir(dirname), "config.json")
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("Read config file error: ", err)
		return nil, err
	}
	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Println("Unmarshal config file error: ", err)
		return nil, err
	}
	if config.ConfigUrl != "" {
		config2, err := loadConfigFromHTTP(config.ConfigUrl)
		if err == nil {
			overrideConfig(&config, config2)
		}
	}
	return &config, nil
}

// merge & override value of first from second
func overrideConfig(config1 *Config, config2 *Config) {

}
