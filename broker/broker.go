package broker

import (
	"errors"
	"log"

	"github.com/kundan1301/message-broker/config"
)

// Broker store broker info
type Broker struct {
	config *config.Config
}

func checkConfig(config *config.Config) error {
	if config.Host == "" {
		return errors.New("Host is required")
	}
	if config.MqttPort == "" {
		return errors.New("MqttPort  is required")
	}
	if config.HttpPort == "" {
		return errors.New("HttpPort is required")
	}
	if config.UseRedisCluster {
		if config.RedisClusterOptions.Addrs == nil || len(config.RedisClusterOptions.Addrs) == 0 {
			return errors.New("Cluster host is required")
		}
	} else {
		if config.RedisOptions.Addr == "" {
			return errors.New("Redis host is required")
		}
	}
	return nil
}

// NewBroker intialize broker
func NewBroker(config *config.Config) (*Broker, error) {
	err := checkConfig(config)
	if err != nil {
		log.Println("error in creating new broker", err)
		return nil, err
	}
	broker := &Broker{
		config: config,
	}
	return broker, nil
}
