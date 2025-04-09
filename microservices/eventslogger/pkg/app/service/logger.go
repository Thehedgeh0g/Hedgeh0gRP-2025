package service

type LoggerService interface {
	Log(level string, args ...[]any) error
}
