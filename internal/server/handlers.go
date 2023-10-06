package server

import (
	"fmt"
	"souschef/data"
	"souschef/internal/session"

	"github.com/gorilla/websocket"
)

var commandHandlers = map[string]func(*websocket.Conn) error{
	"session_start": handleSessionStart,
	"session_stop":  handleSessionStop,
	"task_complete": handleTaskComplete,
}

func handleSessionStart(conn *websocket.Conn) error {
	userID, exist := connections[conn]
	if !exist {
		return fmt.Errorf("user not found")
	}

	err := session.Live.Start(userID)
	if err != nil {
		return err
	}

	// Assign a task to all connected users
	for conn := range connections {
		userID := connections[conn]
		task, err := session.Live.AssignTask(userID)
		if err == nil {
			msg := &data.ServerMessage{
				Type:    "task_new",
				Payload: task,
			}
			transmit(conn, msg)
		}
	}

	return nil
}

func handleSessionStop(conn *websocket.Conn) error {
	userID, exist := connections[conn]
	if !exist {
		return fmt.Errorf("user not found")
	}

	return session.Live.Stop(userID)
}

func handleTaskComplete(conn *websocket.Conn) error {
	userID, exist := connections[conn]
	if !exist {
		return fmt.Errorf("user not found")
	}

	// 1. Complete the user's task
	completedTask, err := session.Live.CompleteTask(userID)
	if err != nil {
		return err
	}

	// 2. Broadcast the task completion
	broadcastMsg := &data.ServerMessage{
		Type:    "task_completed",
		Payload: completedTask,
	}
	broadcast(broadcastMsg)

	// 3. Assign the user a new task
	newTask, err := session.Live.AssignTask(userID)
	if err != nil {
		return err
	}

	// 4. Transmit the new task to the user
	transmitMsg := &data.ServerMessage{
		Type:    "task_new",
		Payload: newTask,
	}
	transmit(conn, transmitMsg)

	return nil
}
