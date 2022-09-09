package database

import (
	"context"
	"errors"
	"io/ioutil"
	"messengerServer/internal/services/authorization/config"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
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
		("postgres://" + conf.DbLogin + ":" + conf.DbPassword + "@" + conf.DbIp + ":" + conf.DbPort + "/messenger"))
	if dbConnErr != nil {
		time.Sleep(time.Second * 5)
	} else {
		return db, dbConnErr
	}
	for dbConnErr != nil {
		db.Conn, dbConnErr = pgxpool.Connect(context.Background(),
			("postgres://" + conf.DbLogin + ":" + conf.DbPassword + "@" + conf.DbIp + ":" + conf.DbPort + "/messenger"))
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
	files, readDirErr := os.ReadDir(path)
	if readDirErr != nil {
		return readDirErr
	}
	for _, fileStat := range files {
		if fileStat.IsDir() {
			continue
		}
		fileName := fileStat.Name()
		if fileName[len(fileStat.Name())-4:] != ".sql" {
			continue
		}
		file, openErr := os.Open(path + fileName)
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
	}

	return nil
}
func (db DB) GetUser(login string) (User, error) {
	var user User
	err := db.Conn.QueryRow(context.Background(), "SELECT * FROM users WHERE login=$1", login).Scan(&user.Login, &user.Password, &user.RefreshBodies)
	if err.Error() == pgx.ErrNoRows.Error() {
		return user, ErrNoRows
	}
	return user, err
}

func (db DB) CreateUser(user User) (int64, error) {
	cmdTag, err := db.Conn.Exec(context.Background(), "INSERT INTO users VALUES ($1,$2,$3)", user.Login, user.Password, user.RefreshBodies)
	return cmdTag.RowsAffected(), err
}

func (db DB) UpdateUser(login string, user User) (int64, error) {

	var sqlChanges []string
	if user.Login != "" {
		sqlChanges = append(sqlChanges, "SET login = $1")
	}
	if user.Password != "" {
		sqlChanges = append(sqlChanges, "SET password = $2")
	}
	if user.RefreshBodies != nil {
		sqlChanges = append(sqlChanges, "SET refreshBodies = $3")
	}
	sql := "UPDATE users " + strings.Join(sqlChanges, ",") + " WHERE login = $4"

	cmdTag, err := db.Conn.Exec(context.Background(), sql, user.Login, user.Password, user.RefreshBodies, login)
	return cmdTag.RowsAffected(), err
}
func (db DB) AddRefreshBody(login, refreshBody string) (int64, error) {
	cmdTag, err := db.Conn.Exec(context.Background(), "UPDATE users SET refreshBodies = array_append(refreshBodies, $1) WHERE login = $2", refreshBody, login)
	return cmdTag.RowsAffected(), err
}
