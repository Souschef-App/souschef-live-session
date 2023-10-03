package server

import (
	"fmt"
	"souschef/client"
	"souschef/internal/utils"
)

var messageHandlers = map[string]func([]byte) error{
	"session_msg": handleSessionMsg,
	"task_msg":    handleTaskMsg,
}

func handleTaskMsg(data []byte) error {
	var msg client.TaskMessage
	if err := utils.DecodeJSON(data, &msg); err != nil {
		return err
	}

	// HANDLE MEAL PLAN LOGIC
	fmt.Println(msg.Cmd)

	return nil
}

func handleSessionMsg(data []byte) error {
	var msg client.SessionMessage
	if err := utils.DecodeJSON(data, &msg); err != nil {
		return err
	}

	return nil
}
