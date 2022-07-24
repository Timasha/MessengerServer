package main

import (
	"log"
	"logger/internal/authService/api"
	"logger/internal/authService/config"
	"logger/internal/authService/database"
)

func main() {
	conf, readConfigErr := config.ReadConfig("AUTHSERVICE_PATH_CONFIG")
	if readConfigErr != nil {
		log.Fatalf("Read config error: %v", readConfigErr)
	}
	db, connErr := database.ConnectDB(conf)
	if connErr != nil {
		log.Fatalf("Cant connect db error: %v", connErr)
	}
	database.SetRepository(db)
	migrErr := database.Migrate("AUTHSERVICE_PATH_MIGRATIONS")
	if migrErr != nil {
		log.Fatalf("Migration error: %v", migrErr)
	}
	api.WebStart(conf)
}
