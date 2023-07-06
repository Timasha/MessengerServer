package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Domain         string `json:"domain"`
	Ip             string `json:"ip"`
	ServerPort     string `json:"serverPort"`
	JWTKey         string `json:"jwtKey"`
	DBUser         string `json:"dbUser"`
	DBUserPassword string `json:"dbUserPassword"`
}

func ReadConfig(path string) (config Config, err error) {
	file, openErr := os.Open(path)
	if openErr != nil {
		return Config{}, openErr
	}
	data, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		return Config{}, readErr
	}
	unmarshalErr := json.Unmarshal(data, &config)
	if unmarshalErr != nil {
		return Config{}, unmarshalErr
	}
	return config, nil
}
