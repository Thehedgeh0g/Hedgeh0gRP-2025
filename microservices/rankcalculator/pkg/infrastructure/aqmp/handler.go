package aqmp

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"

	appevent "rankcalculator/pkg/app/event"
)

type AMQPHandler struct {
	eventHandler appevent.Handler
	amqpChannel  *amqp.Channel
}

func NewAMQPHandler(eventHandler appevent.Handler, connection *amqp.Channel) *AMQPHandler {
	return &AMQPHandler{
		eventHandler: eventHandler,
		amqpChannel:  connection,
	}
}

func (h *AMQPHandler) Listen(queueName string) {
	msgs, err := h.amqpChannel.Consume(
		queueName, "", true, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	go func() {
		for d := range msgs {
			eventType := d.Type
			eventData := d.Body

			evt, err := h.createEvent(eventType, eventData)
			if err != nil {
				log.Printf("Failed to create event: %v", err)
				continue
			}
			if evt == nil {
				log.Printf("Unknown event type: %s", eventType)
				continue
			}

			h.eventHandler.Handle(evt)
		}
	}()

	log.Printf("Listening for messages on queue: %s", queueName)
}

func (h *AMQPHandler) createEvent(eventType string, data []byte) (appevent.Event, error) {
	switch eventType {
	case "valuator.TextAdded":
		var event appevent.TextAddedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		return &event, nil
	default:
		return nil, nil
	}
}
