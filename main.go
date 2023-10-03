package main

import (
	"fmt"
	"souschef/internal/server"
)

func main() {
	fmt.Println("Initializing...")

	// TODO: (Optional) Kafka message:
	// 1. Request valid websocket IP

	server.StartWebSocket()

	// TODO: Kafka message:
	// 1. (Optional) Send websocket IP
	// 2. Request mealplan

	// TODO: Kafka message:
	// 1. Notify "ready for client connections"
}
