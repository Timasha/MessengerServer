package database

type Repository interface {
	Close()
	Migrate(string) error
	GetUser(string) (User, error)
}

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
