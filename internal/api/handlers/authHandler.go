package handlers

import (
	"MessengerServer/internal/utils"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handleAuthData(db *gorm.DB, config utils.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		body, bodyReadErr := ioutil.ReadAll(c.Request.Body)
		if bodyReadErr != nil {
			log.Printf("Read body error: %v \n", bodyReadErr)
		}
		authData, readAuthDataErr := utils.ReadAuthData(body)
		utils.CheckErr("Read auth data error: %v", readAuthDataErr)
		userStatus := utils.CheckUser(authData, db)
		if userStatus == "user_not_exist" {
			c.Writer.WriteString(userStatus)
		} else if userStatus == "invalid_password" {
			c.Writer.WriteString(userStatus)
		} else if userStatus == "user_exist" {
			token, tokenErr := utils.CreateToken([]byte(config.JWTKey), authData.Login)
			utils.CheckErr("Token create error: %v", tokenErr)
			c.Writer.WriteString(token)
		}
	}
}
