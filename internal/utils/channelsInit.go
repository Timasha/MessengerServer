package utils

import (
	"sync"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

func InitInformation(db *gorm.DB) {
	var dbChan []Channel
	var count int64
	db.Find(&dbChan, Channel{}).Count(&count)
	Channels = make(map[string][]*websocket.Conn, count)
	for i := 0; i < len(dbChan); i++ {
		Channels[dbChan[i].ChannelName] = make([]*websocket.Conn, 0)
		BroadcastGroups[dbChan[i].ChannelName] = &sync.WaitGroup{}
		DeleteGroups[dbChan[i].ChannelName] = &sync.WaitGroup{}
	}
	var lastMsg Message
	db.Last(&lastMsg, Message{})
	LastMsgID = lastMsg.ID
}
