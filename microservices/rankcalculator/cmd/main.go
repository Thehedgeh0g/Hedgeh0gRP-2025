package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"log"

	"rankcalculator/pkg/app/event"
	"rankcalculator/pkg/app/service"
	amqpClient "rankcalculator/pkg/infrastructure/aqmp"
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
	amqpConn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer func(amqpConn *amqp.Connection) {
		err2 := amqpConn.Close()
		if err2 != nil {
			panic(err2)
		}
	}(amqpConn)
}

func main() {
	var err error

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

	textRepository := repository.NewTextRepository(redisClient)
	textService := service.NewStatisticsService(textRepository)
	eventHandler := event.NewHandler(textService)
	amqpHandler := amqpClient.NewAMQPHandler(eventHandler, amqpChannel)
	amqpHandler.Listen("texts")
}
