package server

import (
	"encoding/json"
	"fmt"
	message "souschef/data"

	"github.com/gorilla/websocket"
)

// TODO: Implement retry mechanism
func transmit(conn *websocket.Conn, msg *message.ServerMessage) {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Failed to serialize server message")
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		fmt.Println("Failed to send server message to:", conn.LocalAddr())
		return
	}
}

// TODO: Implement retry mechanism
func broadcast(msg *message.ServerMessage) {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Failed to serialize server message")
		return
	}

	for conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			fmt.Println("Failed to send server message to:", conn.LocalAddr())
			return
		}
	}
}

func reportError(conn *websocket.Conn, err error) {
	closeMsg := websocket.FormatCloseMessage(websocket.ClosePolicyViolation, err.Error())
	if err := conn.WriteMessage(websocket.CloseMessage, closeMsg); err != nil {
		fmt.Println(err)
	}

	fmt.Println("An error occured:", err)
}
