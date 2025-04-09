package cli

import (
	"encoding/json"
	"fmt"
	"log"

	"eventslogger/pkg/app/service"
)

const logFormat = "level: %s; message: %s"

func NewCliLoggerService() service.LoggerService {
	return &loggerService{}
}

type loggerService struct{}

func (l *loggerService) Log(level string, args map[string]any) error {
	message, err := json.Marshal(args)
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf(logFormat, level, message))
	return nil
}
