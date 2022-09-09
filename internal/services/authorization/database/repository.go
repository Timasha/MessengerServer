package database

import "errors"

type Repository interface {
	Close()
	Migrate(string) error
	GetUser(string) (User, error)
	CreateUser(User) (int64, error)
	UpdateUser(string, User) (int64, error)
	AddRefreshBody(string, string) (int64, error)
}

var ErrNoRows error = errors.New("no object or row found")

var impl Repository

func SetRepository(repo Repository) {
	impl = repo
}

func Close() {
	impl.Close()
}
func Migrate(pathKey string) error {
	return impl.Migrate(pathKey)
}
func GetUser(login string) (User, error) {
	return impl.GetUser(login)
}

func CreateUser(user User) (int64, error) {
	return impl.CreateUser(user)
}
func UpdateUser(login string, user User) (int64, error) {
	return impl.UpdateUser(login, user)
}
func AddRefreshBody(login, refreshBody string) (int64, error) {
	return impl.AddRefreshBody(login, refreshBody)
}
