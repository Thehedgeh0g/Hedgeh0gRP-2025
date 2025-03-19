package event

import (
	"errors"
	"fmt"

	"rankcalculator/pkg/app/service"
)

var (
	ErrUnknownEventType = errors.New("unknown event type")
)

type Handler interface {
	Handle(event Event)
}

func NewHandler(statisticsService service.StatisticsService) Handler {
	return &handler{statisticsService: statisticsService}
}

type handler struct {
	statisticsService service.StatisticsService
}

func (h *handler) Handle(event Event) {
	var err error
	switch e := event.(type) {
	case *TextAddedEvent:
		err = h.handleTextAddedEvent(*e)
	default:
		err = ErrUnknownEventType
	}
	if err != nil {
		h.handleError(err)
		return
	}
	fmt.Printf("Событие: %s обработано успешно", event.GetType())
}

func (h *handler) handleTextAddedEvent(event TextAddedEvent) error {
	return h.statisticsService.CalculateRank(event.TextHash)
}

func (h *handler) handleError(err error) {
	fmt.Printf("не удалось обработать событие: %s", err.Error())
}
