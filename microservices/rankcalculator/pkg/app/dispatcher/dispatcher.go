package dispatcher

import "rankcalculator/pkg/app/event"

type EventDispatcher interface {
	Dispatch(event event.Event) error
}
