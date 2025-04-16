package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"html/template"
	"log"
	"net/http"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"

	amqpadapter "valuator/pkg/Infrastructure/amqp"
	"valuator/pkg/Infrastructure/repository"
	"valuator/pkg/app/service"
)

type connectionContainer struct {
	RedisMain     *redis.Client
	RegionClients *map[string]*redis.Client
	AMQPConn      *amqp.Connection
	AMQPChannel   *amqp.Channel
}

func newConnectionContainer() *connectionContainer {
	container := &connectionContainer{}

	container.RedisMain = newRedisClient(getEnv("DB_MAIN", "redis-main:6379"))
	container.RegionClients = &map[string]*redis.Client{
		"RU":   newRedisClient(getEnv("DB_RU", "redis-ru:6379")),
		"EU":   newRedisClient(getEnv("DB_EU", "redis-eu:6379")),
		"ASIA": newRedisClient(getEnv("DB_ASIA", "redis-asia:6379")),
	}

	var err error
	container.AMQPConn, err = amqp.Dial(getEnv("AMQP_URL", "amqp://guest:guest@rabbitmq:5672/"))
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	container.AMQPChannel, err = container.AMQPConn.Channel()
	if err != nil {
		log.Fatal("Failed to open a AMQP channel:", err)
	}

	return container
}

func newRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: addr})
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

	port := getEnv("PORT", "8082")
	fmt.Println("Listening on port " + port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeHandler(w, r, connections)
	})
	http.HandleFunc("/summary", func(w http.ResponseWriter, r *http.Request) {
		summaryHandler(w, r, connections)
	})
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/health", healthCheck)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request, conn *connectionContainer) {
	if r.Method == http.MethodPost {
		_ = r.ParseForm()
		region := r.FormValue("region")
		text := r.FormValue("text")

		amqpDispatcher := amqpadapter.NewAMQPDispatcher(conn.AMQPChannel, "text")
		shardManager := repository.NewShardManager(conn.RedisMain, conn.RegionClients, region)
		textRepo := repository.NewTextRepository(shardManager)
		textService := service.NewTextService(textRepo, amqpDispatcher)

		hash, err := textService.EvaluateText(text)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/summary?id=%s", hash), http.StatusPermanentRedirect)
		return
	}

	tmpl, err := template.ParseFiles("pages/index.html")
	if err != nil {
		http.Error(w, "Template Error", 500)
		log.Println(err)
		return
	}
	_ = tmpl.Execute(w, nil)
}

func summaryHandler(w http.ResponseWriter, r *http.Request, conn *connectionContainer) {
	hash := r.URL.Query().Get("id")

	shardManager := repository.NewShardManager(conn.RedisMain, conn.RegionClients, "")
	textRepo := repository.NewTextRepository(shardManager)
	text, err := textRepo.FindByHash(hash)
	if err != nil {
		http.Error(w, "Not Found", 404)
		log.Println(err)
		return
	}

	tmpl, err := template.ParseFiles("pages/summary.html")
	if err != nil {
		http.Error(w, "Template Error", 500)
		log.Println(err)
		return
	}

	_ = tmpl.Execute(w, map[string]interface{}{
		"Rank":       text.GetRank(),
		"Similarity": boolToInt(text.GetSimilarity()),
	})
}

func aboutHandler(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("pages/about.html")
	if err != nil {
		http.Error(w, "Template Error", 500)
		log.Println(err)
		return
	}
	_ = tmpl.Execute(w, map[string]interface{}{
		"Name":  "Ilya Lezhnin",
		"Email": "Flipdoge87@gmail.com",
	})
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
