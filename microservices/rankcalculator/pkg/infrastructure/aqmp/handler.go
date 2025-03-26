package aqmp

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"

	appevent "rankcalculator/pkg/app/event"
)

type BaseEvent struct {
	Type string `json:"type"`
}

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
	_, err := h.amqpChannel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete
		// when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	msgs, err := h.amqpChannel.Consume(
		queueName, "", true, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	go func() {
		for d := range msgs {
			eventData := d.Body
			evt, err := h.createEvent(eventData)
			if err != nil {
				log.Printf("Failed to create event: %v", err)
				continue
			}
			if evt == nil {
				log.Printf("Unknown event type: %s", eventData)
				continue
			}

			h.eventHandler.Handle(evt)
		}
	}()

	log.Printf("Listening for messages on queue: %s", queueName)
}

func (h *AMQPHandler) createEvent(data []byte) (appevent.Event, error) {
	event := BaseEvent{}
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	switch event.Type {
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
