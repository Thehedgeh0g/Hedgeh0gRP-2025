package transport

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
	publicapi "valuator/api"
	amqpadapter "valuator/pkg/Infrastructure/amqp"
	repository "valuator/pkg/Infrastructure/redis"
	"valuator/pkg/app/service"
)

type PublicWeb interface {
	publicapi.StrictServerInterface
}

type ConnectionContainer struct {
	RedisMain     *redis.Client
	RegionClients *map[string]*redis.Client
	AMQPConn      *amqp.Connection
	AMQPChannel   *amqp.Channel
}

func newConnectionContainer() *ConnectionContainer {
	container := &ConnectionContainer{}

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
		log.Fatal("Failed to p.connectionContainerect to RabbitMQ:", err)
	}

	container.AMQPChannel, err = container.AMQPConn.Channel()
	if err != nil {
		log.Fatal("Failed to open a AMQP channel:", err)
	}

	return container
}

func jwtTokenFromContext(ctx context.Context) *Claims {
	val := ctx.Value("user")
	if val == nil {
		return nil
	}
	jsonData, _ := json.Marshal(val)
	token := Claims{}
	err := json.Unmarshal(jsonData, &token)
	if err != nil {
		return nil
	}
	return &token
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

type Claims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

func NewPublicWeb(connectionContainer *ConnectionContainer, secret string) PublicWeb {
	return &publicWeb{
		connectionContainer: connectionContainer,
		jwtSecret:           []byte(secret),
	}
}

type publicWeb struct {
	connectionContainer *ConnectionContainer
	jwtSecret           []byte
}

func (p *publicWeb) About(ctx context.Context, request publicapi.AboutRequestObject) (publicapi.AboutResponseObject, error) {
	return publicapi.About200Response{}, nil
}

func (p *publicWeb) Health(ctx context.Context, request publicapi.HealthRequestObject) (publicapi.HealthResponseObject, error) {
	return publicapi.Health200Response{}, nil
}

func (p *publicWeb) SendText(ctx context.Context, request publicapi.SendTextRequestObject) (publicapi.SendTextResponseObject, error) {
	token := jwtTokenFromContext(ctx)
	if token == nil {
		return publicapi.SendText401Response{}, nil
	}
	region := request.Body.Region
	text := request.Body.Text

	amqpDispatcher := amqpadapter.NewAMQPDispatcher(p.connectionContainer.AMQPChannel, "text")
	shardManager := repository.NewShardManager(p.connectionContainer.RedisMain, p.connectionContainer.RegionClients, region)
	textRepo := repository.NewTextRepository(shardManager)
	textService := service.NewTextService(textRepo, amqpDispatcher)

	hash, err := textService.EvaluateText(text, token.Login)
	return publicapi.SendText200JSONResponse(hash), err
}

func (p *publicWeb) Summary(ctx context.Context, request publicapi.SummaryRequestObject) (publicapi.SummaryResponseObject, error) {
	token := jwtTokenFromContext(ctx)
	if token == nil {
		return publicapi.Summary401Response{}, nil
	}

	hash := request.Params.Id

	shardManager := repository.NewShardManager(p.connectionContainer.RedisMain, p.connectionContainer.RegionClients, "")
	textRepo := repository.NewTextRepository(shardManager)
	text, err := textRepo.FindByHash(hash)
	if err != nil {
		return publicapi.Summary404Response{}, err
	}
	channel := "personal#" + text.GetHash()
	return publicapi.Summary200JSONResponse{
		CentrifugoToken: generateCentrifugoToken(token.Login, channel),
		Channel:         channel,
		Rank:            float32(text.GetRank()),
		Similarity:      text.GetSimilarity(),
	}, nil
}

func generateCentrifugoToken(identifier string, channel string) string {
	claims := jwt.MapClaims{
		"sub":      identifier,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"channels": []string{channel},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("my_secret"))
	if err != nil {
		log.Printf("Ошибка генерации токена: %v", err)
		return ""
	}

	return signedToken
}
