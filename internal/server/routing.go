package server

import (
	"encoding/json"
	"errors"
	"fmt"
)

func RouteMessage(data []byte) error {
	var message map[string]interface{}

	// Decode JSON
	if err := json.Unmarshal(data, &message); err != nil {
		return err
	}

	// Check if JSON contains a 'type' field of type string
	messageType, ok := message["type"].(string)
	if !ok {
		errMsg := "message is missing a 'type' field or it's not a string"
		return errors.New(errMsg)
	}

	// Delegate to appropriate handler
	handler, exists := messageHandlers[messageType]
	if !exists {
		return fmt.Errorf("unknown command type: %s", messageType)
	}

	return handler(data)
}
