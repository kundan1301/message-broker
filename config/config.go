package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"runtime"
)

/*Config struct stores configs defined in config.json file
NodeIp will be used to communicate with other node in cluster/autoscaling group
ConfigUrl will be used to load config from that url, if its specified file config will be overridden
*/
type Config struct {
	Host         string `json:"host"`
	MqttPort     string `json:"mqttPort"`
	HttpPort     string `json:"httpPort"`
	MqttTlsPort  string `json:"mqttTlsPort"`
	HttpTlsPort  string `json:"httpTlsPort"`
	NodeIP       string `json:"nodeIP"`
	ConfigUrl    string `json:"configUrl"`
	AuthUrl      string `json:"authUrl"`
	SubscribeUrl string `json:"subscribeUrl"`
	publishUrl   string `json:"publishUrl"`
}

func loadConfigFromHTTP(configURL string) (*Config, error) {
	resp, err := http.Get(configURL)
	if err != nil {
		log.Println("error in loading config from config url: ", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
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
