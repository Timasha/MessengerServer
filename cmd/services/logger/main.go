package main

import (
	"log"
	"messengerServer/internal/services/logger/api"
	"messengerServer/internal/services/logger/config"
)

func main() {
	conf, confErr := config.ReadConfig("LOGGERSERVICE_PATH_CONFIG")
	if confErr != nil {
		log.Fatalf("Read config error: %v", confErr)
	}
	api.WebStart(conf)
}
