package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"

	"rankcalculator/pkg/app/handler"
	"rankcalculator/pkg/app/service"
	amqpClient "rankcalculator/pkg/infrastructure/amqp"
	"rankcalculator/pkg/infrastructure/centrifugo"
	"rankcalculator/pkg/infrastructure/repository"
)

var redisClient *redis.Client
var amqpConn *amqp.Connection
var amqpChannel *amqp.Channel

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Адрес Redis
	})
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

	centrifugoClient := centrifugo.NewCentrifugoClient()
	textRepository := repository.NewTextRepository(redisClient)
	amqpDispatcher := amqpClient.NewAMQPDispatcher(amqpChannel, "text")
	textService := service.NewStatisticsService(textRepository, amqpDispatcher, centrifugoClient)
	eventHandler := handler.NewHandler(textService)
	amqpHandler := amqpClient.NewAMQPConsumer(eventHandler, amqpChannel)
	var forever chan struct{}
	amqpHandler.Consume("text")
	<-forever
}
