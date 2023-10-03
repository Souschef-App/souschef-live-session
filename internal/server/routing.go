package server

import (
	"fmt"
	"souschef/client"
	"souschef/internal/utils"
)

func RouteMessage(data []byte) error {
	var msg client.Message

	// Decode JSON
	if err := utils.DecodeJSON(data, &msg); err != nil {
		return err
	}

	// Delegate to appropriate handler
	handler, exists := messageHandlers[msg.Type]
	if !exists {
		return fmt.Errorf("unknown command type: %s", msg.Type)
	}

	return handler(msg.Payload)
}
