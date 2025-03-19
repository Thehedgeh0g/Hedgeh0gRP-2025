package event

type Event interface {
	GetType() string
}

type TextAddedEvent struct {
	TextHash string
}

func (e TextAddedEvent) GetType() string {
	return "valuator.TextAdded"
}
