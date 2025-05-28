package main

import (
	"fmt"
	"net/http"
	"os"
	"protokey/pkg/app/model"
	"protokey/pkg/infrastructure/transport"
)

func main() {
	storage := model.NewStorage()
	handler := transport.NewHandler(storage)

	http.HandleFunc("/set", handler.Set)
	http.HandleFunc("/get", handler.Get)
	http.HandleFunc("/keys", handler.Keys)

	fmt.Println("ProtoKey server listening on :" + getEnv("PORT", "6370"))
	err := http.ListenAndServe(":"+getEnv("PORT", "6370"), nil)
	if err != nil {
		return
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
