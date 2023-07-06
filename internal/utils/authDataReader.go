package utils

import (
	"encoding/json"
)

func ReadAuthData(data []byte) (authData User, err error) {
	unmarshalErr := json.Unmarshal(data, &authData)
	if unmarshalErr != nil {
		return User{}, unmarshalErr
	}
	return authData, nil
}
