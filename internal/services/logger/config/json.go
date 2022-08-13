package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type LoggerServiceConfig struct {
	Port          string   `json:"port"`
	KnownServices []string `json:"knownServices"`
	TLS           bool     `json:"tls"`
}

func ReadConfig(pathKey string) (LoggerServiceConfig, error) {
	var conf LoggerServiceConfig
	path := os.Getenv(pathKey)
	if path == "" {
		return LoggerServiceConfig{}, errors.New("cannot find service config path or it's empty")
	}
	file, openErr := os.Open(pathKey)
	if openErr != nil {
		return LoggerServiceConfig{}, openErr
	}
	data, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		return LoggerServiceConfig{}, readErr
	}
	unmarshErr := json.Unmarshal(data, &conf)
	if unmarshErr != nil {
		return LoggerServiceConfig{}, unmarshErr
	}
	return conf, nil
}
