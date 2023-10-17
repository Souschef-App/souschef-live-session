package server

import (
	"encoding/json"
	"fmt"
	"souschef/internal/message"

	"github.com/gorilla/websocket"
)

// TODO: Implement retry mechanism
func transmit(conn *websocket.Conn, msgType message.ServerMessageEnum, payload any) {
	msg := message.ServerMessage{
		Type:    msgType,
		Payload: payload,
	}

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
func broadcast(msgType message.ServerMessageEnum, payload any) {
	msg := message.ServerMessage{
		Type:    msgType,
		Payload: payload,
	}

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
