package event

import (
	"errors"

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
}

func (h *handler) handleRankCalculatedEvent(event RankCalculatedEvent) error {
	logData := map[string]any{
		"Entity": event.Entity,
		"Rank":   event.Rank,
		"Hash":   event.Hash,
	}
	return h.logger.Log("info", logData)
}

func (h *handler) handleSimilarityCalculatedEvent(event SimilarityCalculatedEvent) error {
	logData := map[string]any{
		"Entity":     event.Entity,
		"Similarity": event.Similarity,
		"Hash":       event.Hash,
	}
	return h.logger.Log("info", logData)
}

func (h *handler) handleError(err error) {
}
