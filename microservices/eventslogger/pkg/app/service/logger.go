package service

type LoggerService interface {
	Log(level string, args map[string]any) error
}
