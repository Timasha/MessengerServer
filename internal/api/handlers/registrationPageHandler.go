package handlers

import (
	"MessengerServer/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func registrationFormHandler(config utils.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "registrationForm.html", map[string]string{
			"domain": config.Domain,
		})
	}
}
