package utils

import (
	"github.com/gorilla/websocket"
)

func Broadcast(channelName string, messege []byte) {
	DeleteGroups[channelName].Wait()
	BroadcastGroups[channelName].Add(1)
	for i := 0; i < len(Channels[channelName]); i++ {
		Channels[channelName][i].WriteMessage(websocket.TextMessage, messege)
	}
	BroadcastGroups[channelName].Done()
}
