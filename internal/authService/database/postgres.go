package database

import (
	"context"
	"errors"
	"io/ioutil"
	"logger/internal/authService/config"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	Conn *pgxpool.Pool
}

func ConnectDB(conf config.AuthServiceConfig) (DB, error) {
	db := DB{}
	var dbConnErr error
	i := 0

	db.Conn, dbConnErr = pgxpool.Connect(context.Background(),
		("postgres://" + conf.DbLogin + ":" + conf.DbPassword + "@" + conf.DbIp + ":" + conf.DbPort + "/logger"))
	if dbConnErr != nil {
		time.Sleep(time.Second * 5)
	} else {
		return db, dbConnErr
	}
	for dbConnErr != nil {
		db.Conn, dbConnErr = pgxpool.Connect(context.Background(),
			("postgres://" + conf.DbLogin + ":" + conf.DbPassword + "@" + conf.DbIp + ":" + conf.DbPort + "/logger"))
		time.Sleep(time.Second * 5)
		i++
		if i == 5 {
			return DB{}, dbConnErr
		}
	}
	return db, dbConnErr
}
func (db DB) Close() {
	db.Conn.Close()
}
func (db DB) Migrate(pathKey string) error {
	path := os.Getenv(pathKey)
	if path == "" {
		return errors.New("empty migrations path")
	}
	file, openErr := os.Open((path + "users.sql"))
	if openErr != nil {
		return openErr
	}
	data, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		return readErr
	}
	_, queryErr := db.Conn.Query(context.Background(), string(data))
	if queryErr != nil {
		return queryErr
	}
	return nil
}
func (db DB) GetUser(login string) (User, error) {
	var user User
	err := db.Conn.QueryRow(context.Background(), "SELECT * FROM users WHERE login=$1", login).Scan(&user.Login, &user.Password)
	return user, err
}
