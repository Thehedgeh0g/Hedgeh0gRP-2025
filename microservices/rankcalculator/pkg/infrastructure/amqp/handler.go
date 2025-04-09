package amqp

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"rankcalculator/pkg/app/handler"

	appevent "rankcalculator/pkg/app/event"
)

type BaseEvent struct {
	Type string `json:"type"`
}

type AMQPConsumer struct {
	eventHandler handler.Handler
	amqpChannel  *amqp.Channel
}

func NewAMQPConsumer(eventHandler handler.Handler, connection *amqp.Channel) *AMQPConsumer {
	return &AMQPConsumer{
		eventHandler: eventHandler,
		amqpChannel:  connection,
	}
}

func (h *AMQPConsumer) Consume(queueName string) {
	err := h.amqpChannel.ExchangeDeclare("valuator", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}
	_, err = h.amqpChannel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete
		// when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	err = h.amqpChannel.QueueBind(queueName, "", "valuator", false, nil)
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

func (h *AMQPConsumer) createEvent(data []byte) (appevent.Event, error) {
	baseEvent := BaseEvent{}
	err := json.Unmarshal(data, &baseEvent)
	if err != nil {
		return nil, err
	}
	switch baseEvent.Type {
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
