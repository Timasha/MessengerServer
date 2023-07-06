package handlers

import (
	"MessengerServer/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Middleware(config utils.Config, db *gorm.DB, flags utils.Flags) {
	r := gin.Default()
	utils.InitInformation(db)
	r.LoadHTMLFiles(flags.TemplatePath)
	r.GET("/", registrationFormHandler(config))
	r.POST("/connectChannel", handleConnectionChannel(db))
	r.GET("/ws/:channel", handleMessege(db, config))
	r.POST("/createChan", handleCreateChannel(db))
	r.POST("/register", handleRegister(db))
	r.POST("/authData", handleAuthData(db, config))
	r.Run(config.Ip + ":" + "8080")
}
