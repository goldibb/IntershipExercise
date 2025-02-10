package main

import (
	"IntershipExercise/cmd/api"
	"log"
)

func main() {
	server := api.NewAPIServer("localhost:8080", db)
	if err := server.Run(); err != nil {
		log.Fatalf("Error: Unable to start server: %v", err)
	}

}
