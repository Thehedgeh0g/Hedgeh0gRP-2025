package event

type EventDispatcher interface {
	Dispatch(event Event) error
}
