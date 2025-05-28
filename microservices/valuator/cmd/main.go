package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"

	publicapi "valuator/api"
	"valuator/pkg/Infrastructure/transport"
)

func newConnectionContainer() *transport.ConnectionContainer {
	container := &transport.ConnectionContainer{}

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

	webAPI := transport.NewPublicWeb(connections, "secret")
	handler := publicapi.NewStrictHandler(webAPI, []publicapi.StrictMiddlewareFunc{})

	port := getEnv("PORT", "8082")
	router := mux.NewRouter()

	router.PathPrefix("/api/").Handler(jwtMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("API request %s %s", r.Method, r.URL.Path)
		publicapi.Handler(handler).ServeHTTP(w, r)
	})))

	staticDir := "./static"
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))),
	)

	router.HandleFunc("/", serveHTML("index.html"))
	router.HandleFunc("/summary", serveHTML("summary.html"))
	router.HandleFunc("/login", serveHTML("login.html"))
	router.HandleFunc("/register", serveHTML("register.html"))
	router.HandleFunc("/about", serveHTML("about.html"))

	server := http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	fmt.Println("Listening on port " + port)
	log.Fatal(server.ListenAndServe())
}

func serveHTML(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pages/"+filename)
	}
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем токен из куки
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "No token cookie found", http.StatusUnauthorized)
			return
		}

		tokenString := tokenCookie.Value

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(getEnv("JWT_KEY", "secret")), nil
		})

		if err != nil || !token.Valid {
			updateReq, err := http.NewRequest("POST", "http://localhost:8080/auth/update-token", nil)
			if err != nil {
				http.Error(w, "Failed to create update request", http.StatusUnauthorized)
				return
			}

			updateReq.AddCookie(tokenCookie)

			client := &http.Client{}
			resp, err := client.Do(updateReq)
			if err != nil || resp.StatusCode != http.StatusOK {
				http.Error(w, "Failed to refresh token", http.StatusUnauthorized)
				return
			}

			for _, c := range resp.Cookies() {
				if c.Name == "token" {
					http.SetCookie(w, c)
					tokenString = c.Value

					token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
						if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
							return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
						}
						return []byte(getEnv("JWT_KEY", "secret")), nil
					})
					break
				}
			}

			if err != nil || !token.Valid {
				http.Error(w, "Token refresh failed", http.StatusUnauthorized)
				return
			}
		}

		ctx := context.WithValue(r.Context(), "user", token.Claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
