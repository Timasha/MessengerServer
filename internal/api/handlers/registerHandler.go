package handlers

import (
	"MessengerServer/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handleRegister(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := utils.User{
			Login:    c.PostForm("login"),
			Password: utils.Hash([]byte(c.PostForm("password"))),
		}
		if strings.Trim(c.PostForm("login"), " ") == "" || strings.Trim(c.PostForm("password"), " ") == "" {
			c.String(http.StatusOK, "%v", "Empty login or password field. Try again")
		} else if utils.CheckUser(user, db) == "invalid_password" || utils.CheckUser(user, db) == "user_exist" {
			c.String(http.StatusOK, "%v", "User already exist. Try again")
		} else {
			db.Create(user)
			c.String(http.StatusOK, "%v", "User succesfuly created")
		}
	}
}
