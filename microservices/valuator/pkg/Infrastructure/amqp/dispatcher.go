package amqp

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	appevent "valuator/pkg/app/event"
)

func NewAMQPDispatcher(amqpChannel *amqp.Channel, queueName string) appevent.EventDispatcher {
	return &amqpDispatcher{
		amqpChannel: amqpChannel,
		queueName:   queueName,
	}
}

type amqpDispatcher struct {
	amqpChannel *amqp.Channel
	queueName   string
}

func (a *amqpDispatcher) Dispatch(event appevent.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = a.amqpChannel.Publish(
		"", a.queueName, false, false,
		amqp.Publishing{ContentType: "application/json", Body: body},
	)
	if err != nil {
		return err
	}

	log.Printf("Published event: %s", event.GetType())
	return nil
}
