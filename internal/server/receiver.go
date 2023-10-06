package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	message "souschef/data"

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
				// TODO: Improve error handling to gracefully handle errors
				// without prematurely closing the connection.
				// -------------------------------------------
				// Consider using a custom error structure for non-urgent errors,
				// so that when a custom error is encountered, it can be reported
				// without severing the connection.
				fmt.Println("An error occured:", err)
				return
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
		return err
	}

	// Delegate to appropriate handler
	handler, exists := commandHandlers[msg.Type]
	if !exists {
		return fmt.Errorf("unknown command")
	}

	return handler(conn)
}
