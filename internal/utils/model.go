package utils

import (
	"sync"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{}

type User struct {
	Login    string `json:"login" gorm:"primary_key"`
	Password string `json:"password"`
}
type MessageWS struct {
	ID      uint64 `json:"id"`
	MsgType int    `json:"msgtype"`
	Login   string `json:"login"`
	Time    int64  `json:"time"`
	Msg     string `json:"msg"`
	JWT     string `json:"jwt"`
	Channel string `json:"channel"`
}
type Message struct {
	ID      uint64 `json:"id" gorm:"primaryKey"`
	Login   string `json:"login"`
	Time    int64  `json:"time"`
	Msg     string `json:"msg"`
	Channel string `json:"channel"`
}
type Channel struct {
	ChannelName string `gorm:"primary_key"`
}

var Channels map[string][]*websocket.Conn
var DeleteGroups map[string]*sync.WaitGroup = make(map[string]*sync.WaitGroup, 0)
var BroadcastGroups map[string]*sync.WaitGroup = make(map[string]*sync.WaitGroup, 0)
var LastMsgID uint64