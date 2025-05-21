package main

import (
	"eventslogger/pkg/infrastructure/cli"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"

	"eventslogger/pkg/app/event"
	amqpClient "eventslogger/pkg/infrastructure/amqp"
)

var amqpConn *amqp.Connection
var amqpChannel *amqp.Channel

func init() {
	var err error
	amqpUser := getEnv("AMQP_USER", "guest")
	amqpPassword := getEnv("AMQP_PASS", "guest")
	amqpConn, err = amqp.Dial("amqp://" + amqpUser + ":" + amqpPassword + "@rabbitmq:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
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
	amqpHandler.Listen("log")
	<-forever
}
