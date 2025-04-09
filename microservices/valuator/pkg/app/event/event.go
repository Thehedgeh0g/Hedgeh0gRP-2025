package event

type Event interface {
	GetType() string
}

type BaseEvent struct {
	Type string `json:"type"`
}

func NewTextAddedEvent(hash string) Event {
	return TextAddedEvent{
		TextHash: hash,
		BaseEvent: BaseEvent{
			Type: "valuator.TextAdded",
		},
	}
}

type TextAddedEvent struct {
	TextHash string
	BaseEvent
}

func (e TextAddedEvent) GetType() string {
	return "valuator.TextAdded"
}

func NewSimilarityCalculatedEvent(hash string, similarity bool) Event {
	return SimilarityCalculatedEvent{
		Entity:     "Valuator",
		Hash:       hash,
		Similarity: similarity,
		BaseEvent:  BaseEvent{},
	}
}

type SimilarityCalculatedEvent struct {
	Entity     string
	Hash       string
	Similarity bool
	BaseEvent
}

func (e SimilarityCalculatedEvent) GetType() string {
	return "log.similarityCalculated"
}
