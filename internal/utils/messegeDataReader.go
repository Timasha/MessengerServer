package utils

import "encoding/json"

func ReadMessege(data []byte) (result MessageWS, err error) {
	err = json.Unmarshal(data, &result)
	return
}
