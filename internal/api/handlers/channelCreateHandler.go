package handlers

import (
	"MessengerServer/internal/utils"
	"io/ioutil"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

func handleCreateChannel(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		readedBody, readBodyErr := ioutil.ReadAll(c.Request.Body)
		if readBodyErr != nil {
			c.Writer.WriteString(readBodyErr.Error())
		}
		channelStatus := utils.CheckChannel(utils.Channel{ChannelName: string(readedBody)}, db)
		if channelStatus == "channel_exist" {
			c.Writer.WriteString(channelStatus)
		} else if channelStatus == "channel_not_exist" {
			go db.Create(utils.Channel{ChannelName: string(readedBody)})
			utils.Channels[string(readedBody)] = make([]*websocket.Conn, 0)
			utils.BroadcastGroups[string(readedBody)] = &sync.WaitGroup{}
			utils.DeleteGroups[string(readedBody)] = &sync.WaitGroup{}
			c.Writer.WriteString(channelStatus)
		}
	}
}
