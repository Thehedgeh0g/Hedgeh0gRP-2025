package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"

	"rankcalculator/pkg/app/handler"
	"rankcalculator/pkg/app/service"
	amqpClient "rankcalculator/pkg/infrastructure/amqp"
	"rankcalculator/pkg/infrastructure/centrifugo"
	"rankcalculator/pkg/infrastructure/repository"
)

type connectionContainer struct {
	RedisMain     *redis.Client
	RegionClients *map[string]*redis.Client
	AMQPConn      *amqp.Connection
	AMQPChannel   *amqp.Channel
}

func newConnectionContainer() *connectionContainer {
	container := &connectionContainer{}

	container.RedisMain = newRedisClient(
		getEnv("DB_MAIN", "redis-main:6379"),
		getEnv("REDIS_PASSWORD", "pass"),
	)
	container.RegionClients = &map[string]*redis.Client{
		"RU": newRedisClient(
			getEnv("DB_RU", "redis-ru:6379"),
			getEnv("REDIS_PASSWORD", "pass"),
		),
		"EU": newRedisClient(
			getEnv("DB_EU", "redis-eu:6379"),
			getEnv("REDIS_PASSWORD", "pass"),
		),
		"ASIA": newRedisClient(
			getEnv("DB_ASIA", "redis-asia:6379"),
			getEnv("REDIS_PASSWORD", "pass"),
		),
	}

	var err error
	amqpUser := getEnv("AMQP_USER", "guest")
	amqpPassword := getEnv("AMQP_PASS", "guest")
	container.AMQPConn, err = amqp.Dial("amqp://" + amqpUser + ":" + amqpPassword + "@rabbitmq:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	container.AMQPChannel, err = container.AMQPConn.Channel()
	if err != nil {
		log.Fatal("Failed to open a AMQP channel:", err)
	}

	return container
}

func newRedisClient(
	addr, pass string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
	})
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	connections := newConnectionContainer()
	defer connections.AMQPConn.Close()
	defer connections.AMQPChannel.Close()

	centrifugoClient := centrifugo.NewCentrifugoClient()
	shardManager := repository.NewShardManager(connections.RedisMain, connections.RegionClients, "")
	textRepository := repository.NewTextRepository(shardManager)
	amqpDispatcher := amqpClient.NewAMQPDispatcher(connections.AMQPChannel, "text")
	textService := service.NewStatisticsService(textRepository, amqpDispatcher, centrifugoClient)
	eventHandler := handler.NewHandler(textService)
	amqpHandler := amqpClient.NewAMQPConsumer(eventHandler, connections.AMQPChannel)
	var forever chan struct{}
	amqpHandler.Consume("text")
	<-forever
}
