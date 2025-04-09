package event

type Event interface {
	GetType() string
}

type BaseEvent struct {
	Type string `json:"type"`
}

type TextAddedEvent struct {
	TextHash string
}

func (e *TextAddedEvent) GetType() string {
	return "valuator.TextAdded"
}

func NewRankCalculatedEvent(hash string, rank float64) Event {
	return RankCalculatedEvent{
		Entity:    "RankCalculator",
		Hash:      hash,
		Rank:      rank,
		BaseEvent: BaseEvent{Type: "log.rankCalculated"},
	}
}

type RankCalculatedEvent struct {
	Entity string
	Hash   string
	Rank   float64
	BaseEvent
}

func (e RankCalculatedEvent) GetType() string {
	return "log.rankCalculated"
}
