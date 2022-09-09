package logger

type LogLevel string

const (
	Info  LogLevel = "info"
	Warn           = "warn"
	Error          = "error"
	Fatal          = "fatal"
)
