package event

type Event interface {
	GetType() string
}

type TextAddedEvent struct {
	TextHash string
}

func (e *TextAddedEvent) GetType() string {
	return "valuator.TextAdded"
}

func NewRankCalculatedEvent(hash string, rank float64) Event {
	return RankCalculatedEvent{
		Entity: "RankCalculator",
		Hash:   hash,
		Rank:   rank,
	}
}

type RankCalculatedEvent struct {
	Entity string
	Hash   string
	Rank   float64
}

func (e RankCalculatedEvent) GetType() string {
	return "log.rankCalculated"
}
