package handlers

import (
	"MessengerServer/internal/utils"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

func handleMessege(db *gorm.DB, config utils.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		conn, connErr := utils.Upgrader.Upgrade(c.Writer, c.Request, c.Request.Header)
		if connErr != nil {
			log.Printf("Websocket connection error: %v", connErr)
			c.Writer.WriteString(connErr.Error())
		}
		channelName := c.Param("channel")
		utils.Channels[channelName] = append(utils.Channels[channelName], conn)
		for {
			msgType, msg, readMsgErr := conn.ReadMessage()
			if msgType == websocket.CloseAbnormalClosure || msgType == websocket.CloseMessage || msgType == websocket.CloseNormalClosure || msgType == -1 {
				utils.BroadcastGroups[channelName].Wait()
				utils.DeleteGroups[channelName].Wait()
				utils.DeleteGroups[channelName].Add(1)
				for i := 0; i < len(utils.Channels[channelName]); i++ {
					if utils.Channels[channelName][i] == conn {
						utils.Channels[channelName] = append(utils.Channels[channelName][:i], utils.Channels[channelName][i+1:]...)
						conn.Close()
						utils.DeleteGroups[channelName].Done()
						return
					}
				}
			}
			if readMsgErr != nil {
				utils.BroadcastGroups[channelName].Wait()
				utils.DeleteGroups[channelName].Wait()
				utils.DeleteGroups[channelName].Add(1)
				for i := 0; i < len(utils.Channels[channelName]); i++ {
					if utils.Channels[channelName][i] == conn {
						utils.Channels[channelName] = append(utils.Channels[channelName][:i], utils.Channels[channelName][i+1:]...)
						conn.Close()
						utils.DeleteGroups[channelName].Done()
						return
					}
				}
			}
			message, unmarshalErr := utils.ReadMessege(msg)
			if unmarshalErr != nil {
				utils.BroadcastGroups[channelName].Wait()
				utils.DeleteGroups[channelName].Wait()
				utils.DeleteGroups[channelName].Add(1)
				for i := 0; i < len(utils.Channels[channelName]); i++ {
					if utils.Channels[channelName][i] == conn {
						utils.Channels[channelName] = append(utils.Channels[channelName][:i], utils.Channels[channelName][i+1:]...)
						conn.Close()
						utils.DeleteGroups[channelName].Done()
						return
					}
				}
			}
			tokenLogin, parseTokenErr := utils.ParseToken(message.JWT, config.JWTKey)
			if parseTokenErr != nil {
				if parseTokenErr.Error() == "Token is expired" || tokenLogin != message.Login {
					conn.WriteMessage(websocket.TextMessage, utils.SystemMessegeMarshal("", 3, message.Channel)) // MsgType 3 is server message. It means login token is expired.
				} else {
					conn.WriteMessage(websocket.TextMessage, utils.SystemMessegeMarshal(parseTokenErr.Error(), 1, message.Channel))
					utils.BroadcastGroups[channelName].Wait()
					utils.DeleteGroups[channelName].Wait()
					utils.DeleteGroups[channelName].Add(1)
					for i := 0; i < len(utils.Channels[channelName]); i++ {
						if utils.Channels[channelName][i] == conn {
							utils.Channels[channelName] = append(utils.Channels[channelName][:i], utils.Channels[channelName][i+1:]...)
							conn.Close()
							utils.DeleteGroups[channelName].Done()
							return
						}
					}
				}
			} else if message.MsgType == 0 { // MsgType zero means user want to delete some message
				deleteStatus := utils.DeleteMessege(message, db)
				conn.WriteMessage(websocket.TextMessage, utils.SystemMessegeMarshal(deleteStatus, 1, message.Channel))
			} else if message.MsgType == 1 { // MsgType 1 is basic message. It can also be used to print system messages
				msgDB := utils.Message{
					Login:   message.Login,
					Time:    message.Time,
					Msg:     message.Msg,
					Channel: message.Channel,
				}
				db.Create(&msgDB)
				utils.LastMsgID += 1
				var messageToSent utils.MessageWS = utils.MessageWS{
					ID:      utils.LastMsgID,
					MsgType: 1,
					Login:   message.Login,
					Time:    message.Time,
					Msg:     message.Msg,
					Channel: message.Channel,
				}
				data, marshErr := json.Marshal(messageToSent)
				utils.CheckErr("Marshal messege error: %v", marshErr)
				utils.Broadcast(message.Channel, data)
			} else if message.MsgType == 2 { // MsgType two means user is connected to channel
				data := utils.SystemMessegeMarshal(("User " + message.Login + " is connected"), 2, message.Channel)
				utils.Broadcast(message.Channel, data)
			}
		}
	}
}
