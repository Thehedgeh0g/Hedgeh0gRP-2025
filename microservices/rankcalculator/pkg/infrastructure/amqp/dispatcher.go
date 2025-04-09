package amqp

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	appevent "rankcalculator/pkg/app/event"
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
		return fmt.Errorf("could not marshal event: %w", err)
	}

	_, err = a.amqpChannel.QueueDeclare(
		a.queueName, // name
		false,       // durable
		false,       // delete
		// when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
		return err
	}

	err = a.amqpChannel.Publish(
		"", a.queueName, false, false,
		amqp.Publishing{ContentType: "application/json", Body: body},
	)
	if err != nil {
		log.Fatalf("Failed to publish messages: %v", err)
		return err
	}

	log.Printf("Published event: %s", event.GetType())
	return nil
}
