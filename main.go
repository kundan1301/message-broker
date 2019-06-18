package main

import (
	"fmt"
	"log"

	"github.com/kundan1301/message-broker/config"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("configure broker config error: ", err)
	}
	//log.Println(config.HttpPort)
	fmt.Printf("%+v\n", *config)

}
