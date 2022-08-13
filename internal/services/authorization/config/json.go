package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type AuthServiceConfig struct {
	Port string `json:"port"` // Порт для сервера

	JwtKey          string `json:"jwtKey"`      // Ключ для шифрования токенов авторизации
	JWTLifetime     int64  `json:"jwtLifetime"` // Время жизни токена авторизации
	RefreshLifetime int64  `json:"refreshLifetime"`

	DbLogin    string `json:"dbLogin"`    // Логин для авторизации в базе данных
	DbPassword string `json:"dbPassword"` // Пароль для авторизации в базе данных
	DbIp       string `json:"dbIp"`       // IP базы данных
	DbPort     string `json:"dbPort"`     // Порт базы данных

	TLS bool `json:"tls"` // Параметр отвечающий за включение и отключение TLS
}

func ReadConfig(pathKey string) (AuthServiceConfig, error) {
	var authServiceConfig AuthServiceConfig
	path := os.Getenv(pathKey)
	if path == "" {
		return AuthServiceConfig{}, errors.New("cannot find service config path or it's empty")
	}
	file, openErr := os.Open(path)
	if openErr != nil {
		return AuthServiceConfig{}, openErr
	}
	data, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		return AuthServiceConfig{}, readErr
	}
	unmarshErr := json.Unmarshal(data, &authServiceConfig)
	if unmarshErr != nil {
		return AuthServiceConfig{}, unmarshErr
	}
	return authServiceConfig, nil
}
