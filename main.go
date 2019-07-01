package main

import (
	"log"

	"github.com/kundan1301/message-broker/broker"
	"github.com/kundan1301/message-broker/config"
	customHttp "github.com/kundan1301/message-broker/http"
	customRedis "github.com/kundan1301/message-broker/redis"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("configure broker config error: ", err)
	}
	customHttp.InitHttpClient()
	customRedis.InitClient(config)
	broker.RunningBroker = nil
	var err1 error
	broker.RunningBroker, err1 = broker.NewBroker(config)
	if err1 != nil {
		log.Fatal("error in intializing new broker")
	}
	broker.RunningBroker.Start()

}
