package main

import (
	"os"

	"protocli/pkg/client"
	"protocli/pkg/commands"
)

func main() {
	baseURL := os.Getenv("PROTOKEY_URL")
	if baseURL == "" {
		baseURL = "http://127.0.0.1:6370"
	}

	c := client.NewProtoKeyClient(baseURL)
	handler := commands.NewCommandHandler(c)
	handler.Handle(os.Args[1:])
}
