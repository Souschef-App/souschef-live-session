package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"souschef/internal/message"

	"github.com/gorilla/websocket"
)

func handleMessage(conn *websocket.Conn) {
	for {
		// Read a message from the WebSocket client
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print the received message
		fmt.Printf("Received message: %s\n", p)

		if messageType == websocket.TextMessage {
			if err := routeMessage(conn, p); err != nil {
				fmt.Println("An error occured:", err)
				transmit(conn, message.ServerError, err.Error())
			}
		}

	}
}

func routeMessage(conn *websocket.Conn, data []byte) error {
	var msg message.ClientMessage

	// Decode JSON
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&msg); err != nil {
		fmt.Println("Failed to decode message.")
		return fmt.Errorf("invalid message")
	}

	// Delegate to appropriate handler
	handler, exists := commandHandlers[msg.Type]
	if !exists {
		return fmt.Errorf("unknown command")
	}

	return handler(conn, msg.Payload)
}
