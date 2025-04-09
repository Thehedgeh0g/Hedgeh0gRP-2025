package cli

import (
	"encoding/json"
	"eventslogger/pkg/app/service"
	"fmt"
)

const logFormat = "level: %s; message: %s"

func NewCliLoggerService() service.LoggerService {
	return &loggerService{}
}

type loggerService struct{}

func (l *loggerService) Log(level string, args ...[]any) error {
	message, err := json.Marshal(args)
	if err != nil {
		return err
	}
	fmt.Printf(logFormat, level, message)
	return nil
}
