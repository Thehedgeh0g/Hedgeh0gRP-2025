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
