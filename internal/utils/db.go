package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbLogger logger.Interface = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	logger.Config{
		SlowThreshold:             time.Second,   // Slow SQL threshold
		LogLevel:                  logger.Silent, // Log level
		IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
		Colorful:                  false,         // Disable color
	},
)

func ConnectDb(dbUser, dbUserPassword string) *gorm.DB {
	var db *gorm.DB
	var dbErr error
	for db == nil {
		db, dbErr = gorm.Open(postgres.Open("host=localhost user="+dbUser+" password="+dbUserPassword+" dbname=messenger port=5432 sslmode=disable"), &gorm.Config{
			Logger: dbLogger,
		})
		if dbErr != nil {
			log.Fatalf("Database connection error: %v", dbErr)
		}
		if db == nil {
			fmt.Printf("Cannot connect to database. Trying again.")
			time.Sleep(time.Second)
		}
	}
	db.AutoMigrate(&User{}, &Message{}, &Channel{})
	return db
}
func CheckUser(reqUser User, db *gorm.DB) string {
	var countLogin, countPassword int64
	db.First(&User{}, User{Login: reqUser.Login, Password: Hash([]byte(reqUser.Password))}).Count(&countPassword)
	if countPassword == 1 {
		return "user_exist"
	} else if db.First(&User{}, User{Login: reqUser.Login}).Count(&countLogin); countLogin == 1 {
		return "invalid_password"
	} else {
		return "user_not_exist"
	}
}
func CheckChannel(reqChan Channel, db *gorm.DB) string {
	var count int64
	db.First(&Channel{}, Channel{reqChan.ChannelName}).Count(&count)
	if count == 1 {
		return "channel_exist"
	} else if count == 0 {
		return "channel_not_exist"
	}
	return ""
}
func DeleteMessege(message MessageWS, db *gorm.DB) string {
	var msgCount int64
	id, parseErr := strconv.ParseUint(message.Msg, 10, 64)
	if parseErr != nil {
		log.Printf("Parse to uint64 error: %v", parseErr)
	}
	db.First(&Message{}, Message{ID: id, Login: message.Login}).Count(&msgCount)
	if msgCount == 1 {
		go db.Delete(&Message{}, Message{ID: id})
		return "successful\n"
	} else if db.First(&Message{}, Message{ID: id}).Count(&msgCount); msgCount == 1 {
		return "wrong_user\n"
	} else {
		return "messege_not_exist\n"
	}
}
