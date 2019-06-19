package main

import (
	"fmt"
	"log"

	"github.com/kundan1301/message-broker/broker"
	"github.com/kundan1301/message-broker/config"
	customHttp "github.com/kundan1301/message-broker/http"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("configure broker config error: ", err)
	}
	customHttp.InitHttpClient()
	b, err := broker.NewBroker(config)
	if err != nil {
		log.Fatal("error in intializing new broker")
	}

	fmt.Printf("%+v\n", b.Config)

}
