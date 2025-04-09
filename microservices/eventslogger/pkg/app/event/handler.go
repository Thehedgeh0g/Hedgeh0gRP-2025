package event

import (
	"errors"
	"fmt"

	"eventslogger/pkg/app/service"
)

var (
	ErrUnknownEventType = errors.New("unknown event type")
)

type Handler interface {
	Handle(event Event)
}

func NewHandler(logger service.LoggerService) Handler {
	return &handler{
		logger: logger,
	}
}

type handler struct {
	logger service.LoggerService
}

func (h *handler) Handle(event Event) {
	var err error
	switch e := event.(type) {
	case *RankCalculatedEvent:
		err = h.handleRankCalculatedEvent(*e)
	case *SimilarityCalculatedEvent:
		err = h.handleSimilarityCalculatedEvent(*e)
	default:
		err = ErrUnknownEventType
	}
	if err != nil {
		h.handleError(err)
		return
	}
	fmt.Printf("Событие: %s обработано успешно", event.GetType())
}

func (h *handler) handleRankCalculatedEvent(event RankCalculatedEvent) error {
	logData := []any{
		event.Entity,
		event.Rank,
		event.Hash,
	}
	return h.logger.Log("info", logData)
}

func (h *handler) handleSimilarityCalculatedEvent(event SimilarityCalculatedEvent) error {
	logData := []any{
		event.Entity,
		event.Similarity,
		event.Hash,
	}
	return h.logger.Log("info", logData)
}

func (h *handler) handleError(err error) {
	fmt.Printf("не удалось обработать событие: %s", err.Error())
}
