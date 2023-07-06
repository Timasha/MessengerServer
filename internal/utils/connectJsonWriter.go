package utils

import (
	"encoding/json"

	"gorm.io/gorm"
)

type connectJson struct {
	ChannelStatus string    `json:"channelStatus"`
	MsgHistory    []Message `json:"msgHistory"`
}

func WriteJSONConnect(db *gorm.DB, channel, channelStatus string) []byte {
	var messeges []Message
	var data []byte
	var marshalErr error
	if channelStatus == "channel_exist" {
		db.Find(&messeges, Message{Channel: channel})
		var ConnJson connectJson = connectJson{
			ChannelStatus: channelStatus,
			MsgHistory:    messeges,
		}
		data, marshalErr = json.Marshal(ConnJson)
		CheckErr("Marshal msg history error: %v", marshalErr)
	} else if channelStatus == "channel_not_exist" {
		data, marshalErr = json.Marshal(connectJson{
			ChannelStatus: channelStatus,
			MsgHistory:    messeges,
		})
		CheckErr("Marshal msg history error: %v", marshalErr)
	}
	return data
}
