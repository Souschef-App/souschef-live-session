package server

import (
	"encoding/json"
	"fmt"
	"souschef/data"
	"souschef/internal/message"
	"souschef/internal/session"

	"github.com/gorilla/websocket"
)

var (
	waitingQueue    = map[*websocket.Conn]*data.User{}
	commandHandlers = map[message.ClientMessageEnum]func(*websocket.Conn, json.RawMessage) error{
		message.ClientHandshake:     handleClientHandshake,
		message.ClientStartSession:  handleSessionStart,
		message.ClientStopSession:   handleSessionStop,
		message.ClientCompletedTask: handleTaskComplete,
		message.ClientRerollTask:    handleTaskReroll,
	}
)

func handleClientHandshake(conn *websocket.Conn, payload json.RawMessage) error {
	user := &data.User{}
	if err := json.Unmarshal(payload, user); err != nil {
		return err
	}

	registerConnection(conn, user)

	if session.Live.IsRunning {
		task, err := session.Live.AssignTask(user)
		if err != nil {
			return err
		}

		transmit(conn, message.ServerTaskNew, task)
	}

	return nil
}

func handleSessionStart(conn *websocket.Conn, _ json.RawMessage) error {
	user, exist := connections[conn]
	if !exist {
		return fmt.Errorf("user not found")
	}

	if err := session.Live.Start(user.ID); err != nil {
		return err
	}

	waitingQueue = map[*websocket.Conn]*data.User{}

	// Assign a task to all connected users
	for conn, user := range connections {
		task, err := session.Live.AssignTask(user)
		if err == nil {
			transmit(conn, message.ServerTaskNew, task)
		}
	}

	return nil
}

func handleSessionStop(conn *websocket.Conn, _ json.RawMessage) error {
	user, exist := connections[conn]
	if !exist {
		return fmt.Errorf("user not found")
	}

	waitingQueue = nil

	return session.Live.Stop(user.ID)
}

func handleTaskComplete(conn *websocket.Conn, _ json.RawMessage) error {
	user, exist := connections[conn]
	if !exist {
		return fmt.Errorf("user not found")
	}

	// 1. Complete the user's task
	if err := session.Live.CompleteTask(user); err != nil {
		return err
	}

	// 2. Check if session is still running
	if !session.Live.IsRunning {
		broadcast(message.ServerMealCompleted, nil)
		return nil // Session finished! :)
	}

	// 3. Attempt to assign task to helpers
	processWaitingQueue()

	// 4. Assign the user a new task
	task, err := session.Live.AssignTask(user)
	if err != nil {
		return err
	}

	// 4.1 If no task is available, add to waiting queue
	if task == nil {
		waitingQueue[conn] = user
	}

	// 5. Transmit the new task to the user
	transmit(conn, message.ServerTaskNew, task)

	return nil
}

func processWaitingQueue() {
	for conn, user := range waitingQueue {
		task, err := session.Live.AssignTask(user)
		if err == nil && task != nil {
			delete(waitingQueue, conn)
			transmit(conn, message.ServerTaskNew, task)
		}
	}
}

func handleTaskReroll(conn *websocket.Conn, _ json.RawMessage) error {
	user, exist := connections[conn]
	if !exist {
		return fmt.Errorf("user not found")
	}

	task, err := session.Live.RerollTask(user)
	if err != nil {
		return err
	}

	transmit(conn, message.ServerTaskNew, task)

	return nil
}
