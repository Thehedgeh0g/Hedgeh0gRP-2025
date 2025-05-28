package main

import (
	"auth/pkg/app/service"
	"auth/pkg/infrastucture/jwt"
	"auth/pkg/infrastucture/redis/provider"
	"auth/pkg/infrastucture/redis/repository"
	"auth/pkg/infrastucture/transport"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
)

type connectionContainer struct {
	RedisMain *redis.Client
}

func newConnectionContainer() *connectionContainer {
	container := &connectionContainer{}

	container.RedisMain = newRedisClient(
		getEnv("DB_MAIN", "redis-main:6379"),
		getEnv("REDIS_PASSWORD", "pass"),
	)

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
	userRepo := repository.NewUserRepository(connections.RedisMain)
	userService := service.NewUserService(userRepo)
	tokenRepo := repository.NewTokenRepository(connections.RedisMain)
	tokenProvider := provider.NewTokenProvider(connections.RedisMain)
	tokenService := jwt.NewTokenService(getEnv("JWT_KEY", "secret"), tokenRepo)
	handler := transport.NewHandler(tokenService, userService, tokenProvider)

	r := mux.NewRouter()

	r.HandleFunc("/auth/login", handler.Login).Methods("POST")
	r.HandleFunc("/auth/update-token", handler.UpdateToken).Methods("POST")
	r.HandleFunc("/auth/register", handler.Register).Methods("POST")
	r.HandleFunc("/auth/logout", handler.Logout).Methods("POST")
	port := getEnv("PORT", "8080")
	log.Println(fmt.Sprintf("Starting server on %s", port))
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		fmt.Println(err)
		return
	}
}
