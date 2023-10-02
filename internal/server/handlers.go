package server

import (
	"souschef/api"
	"souschef/internal/utils"
)

var messageHandlers = map[string]func([]byte) error{
	"start_session": handleStartSession,
	"stop_session":  handleStopSession,
}

func handleStartSession(data []byte) error {
	var cmd api.SessionCommand
	if err := utils.DecodeJSON(data, &cmd); err != nil {
		return err
	}

	// HANDLE STARTING LOGIC

	return nil
}

func handleStopSession(data []byte) error {
	var cmd api.SessionCommand
	if err := utils.DecodeJSON(data, &cmd); err != nil {
		return err
	}

	// HANDLE STOPPING LOGIC

	return nil
}
