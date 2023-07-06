package main

import (
	"MessengerServer/internal/api/handlers"
	"MessengerServer/internal/utils"
	"log"
)

func main() {
	flags := utils.ParceFlags()
	config, confReadErr := utils.ReadConfig(flags.ConfigPath)
	if confReadErr != nil {
		log.Fatalf("Config read error: %v", confReadErr)
	}
	db := utils.ConnectDb(config.DBUser, config.DBUserPassword)
	handlers.Middleware(config, db, flags)
}
