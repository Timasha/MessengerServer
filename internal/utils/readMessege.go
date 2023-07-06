package utils

import (
	"encoding/json"
	"log"
	"time"
)

func SystemMessegeMarshal(msg string, msgType int, currentChan string) (data []byte) {
	var marshErr error

	data, marshErr = json.Marshal(MessageWS{
		MsgType: msgType,
		Login:   "System",
		Time:    time.Now().Unix(),
		Msg:     msg,
		Channel: currentChan,
	})
	if marshErr != nil {
		log.Fatalf("Marshal messege json error: %v", marshErr)
	}
	return
}
