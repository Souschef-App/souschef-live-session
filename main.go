package main

import (
	"fmt"
	"souschef/internal/server"
)

func main() {
	fmt.Println("Initializing...")
	server.StartWebSocket(":8080")
}
