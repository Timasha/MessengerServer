package logger

type Repository interface {
	Log()
	LocalLog()
}

var impl Repository

func SetRepository(repo Repository) {
	impl = repo
}
