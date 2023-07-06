package handlers

import (
	"MessengerServer/internal/utils"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handleConnectionChannel(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		readedBody, readBodyErr := ioutil.ReadAll(c.Request.Body)
		if readBodyErr != nil {
			log.Printf("Read http body error: %v", readBodyErr)
		}
		channelStatus := utils.CheckChannel(utils.Channel{ChannelName: string(readedBody)}, db)
		c.Writer.Write(utils.WriteJSONConnect(db, string(readedBody), channelStatus))
	}
}
