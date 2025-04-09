package amqp

import (
	"encoding/json"
	"fmt"
	"log"
	"rankcalculator/pkg/app/dispatcher"

	amqp "github.com/rabbitmq/amqp091-go"

	appevent "rankcalculator/pkg/app/event"
)

func NewAMQPDispatcher(amqpChannel *amqp.Channel, queueName string) dispatcher.EventDispatcher {
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

	// Объявляем fanout exchange
	err = a.amqpChannel.ExchangeDeclare("valuator", "fanout", true, false, false, false, nil)
	if err != nil {
		return err
	}

	// Публикуем В EXCHANGE, routing key игнорируется для fanout
	err = a.amqpChannel.Publish(
		"valuator", "", false, false, // <-- пустой routing key
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
		return err
	}

	log.Printf("Published event: %s", event.GetType())
	return nil
}
