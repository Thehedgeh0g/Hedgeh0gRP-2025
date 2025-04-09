package event

type Event interface {
	GetType() string
}

func NewSimilarityCalculatedEvent(entity, hash string, similarity bool) Event {
	return SimilarityCalculatedEvent{
		Entity:     entity,
		Hash:       hash,
		Similarity: similarity,
		BaseEvent:  BaseEvent{},
	}
}

func NewRankCalculatedEvent(entity, hash string, rank float64) Event {
	return RankCalculatedEvent{
		Entity:    entity,
		Hash:      hash,
		Rank:      rank,
		BaseEvent: BaseEvent{},
	}
}

type BaseEvent struct {
	Type string `json:"type"`
}

type RankCalculatedEvent struct {
	Entity string
	Hash   string
	Rank   float64
	BaseEvent
}

func (r RankCalculatedEvent) GetType() string {
	return "log.rankCalculated"
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
