package main

import (
	"fmt"
	"log"

	"github.com/kundan1301/message-broker/config"
	customHttp "github.com/kundan1301/message-broker/http"
)

func main() {
	config, err := config.LoadConfig()
	customHttp.InitHttpClient()
	if err != nil {
		log.Fatal("configure broker config error: ", err)
	}
	fmt.Printf("%+v\n", *config)

}
