package main

import (
	"eventslogger/pkg/infrastructure/cli"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"eventslogger/pkg/app/event"
	amqpClient "eventslogger/pkg/infrastructure/amqp"
)

var amqpConn *amqp.Connection
var amqpChannel *amqp.Channel

func init() {
	var err error
	amqpConn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
}

func main() {
	var err error

	defer func(amqpConn *amqp.Connection) {
		err2 := amqpConn.Close()
		if err2 != nil {
			panic(err2)
		}
	}(amqpConn)

	amqpChannel, err = amqpConn.Channel()
	if err != nil {
		log.Fatal("Failed to open a amqpChannel:", err)
	}
	defer func(channel *amqp.Channel) {
		err2 := channel.Close()
		if err2 != nil {
			panic(err2)
		}
	}(amqpChannel)

	loggerService := cli.NewCliLoggerService()
	eventHandler := event.NewHandler(loggerService)
	amqpHandler := amqpClient.NewAMQPHandler(eventHandler, amqpChannel)
	var forever chan struct{}
	amqpHandler.Listen("text")
	<-forever
}
