package server

import (
	"fmt"
	"souschef/internal/message"
	"souschef/internal/session"

	"github.com/gorilla/websocket"
)

var commandHandlers = map[message.ClientMessageEnum]func(*websocket.Conn) error{
	message.ClientStartSession:  handleSessionStart,
	message.ClientStopSession:   handleSessionStop,
	message.ClientCompletedTask: handleTaskComplete,
	message.ClientRerollTask:    handleTaskReroll,
}

var waitingQueue = map[string]*websocket.Conn{}

func handleSessionStart(conn *websocket.Conn) error {
	userID, exist := connections[conn]
	if !exist {
		return fmt.Errorf("user not found")
	}

	err := session.Live.Start(userID)
	if err != nil {
		return err
	}

	waitingQueue = map[string]*websocket.Conn{}

	// Assign a task to all connected users
	for conn := range connections {
		userID := connections[conn]
		task, err := session.Live.AssignTask(userID)
		if err == nil {
			transmit(conn, message.ServerTaskNew, task)
		}
	}

	return nil
}

func handleSessionStop(conn *websocket.Conn) error {
	userID, exist := connections[conn]
	if !exist {
		return fmt.Errorf("user not found")
	}

	waitingQueue = nil

	return session.Live.Stop(userID)
}

func handleTaskComplete(conn *websocket.Conn) error {
	userID, exist := connections[conn]
	if !exist {
		return fmt.Errorf("user not found")
	}

	// 1. Complete the user's task
	fmt.Println("Handeling client task completed")
	completedTask, err := session.Live.CompleteTask(userID)
	if err != nil {
		return err
	}

	// 2. Broadcast the task completion
	// TODO: Send task + user info
	broadcast(message.ServerTaskCompleted, completedTask)

	// 3. Check if session is still running
	fmt.Println("Session Running?:", session.Live.IsRunning)
	if !session.Live.IsRunning {
		broadcast(message.ServerMealCompleted, nil)
		return nil // Session finished! :)
	}

	// 4. Attempt to assign task to helpers
	processWaitingQueue()

	// 5. Assign the user a new task
	task, err := session.Live.AssignTask(userID)
	if err != nil {
		return err
	}

	// 5.1 If no task is available, add to waiting queue
	if task == nil {
		waitingQueue[userID] = conn
	}

	// 4. Transmit the new task to the user
	transmit(conn, message.ServerTaskNew, task)

	return nil
}

func processWaitingQueue() {
	for userID, conn := range waitingQueue {
		task, err := session.Live.AssignTask(userID)
		if err == nil && task != nil {
			delete(waitingQueue, userID)
			transmit(conn, message.ServerTaskNew, task)
		}
	}
}

func handleTaskReroll(conn *websocket.Conn) error {
	userID, exist := connections[conn]
	if !exist {
		return fmt.Errorf("user not found")
	}

	task, err := session.Live.RerollTask(userID)
	if err != nil {
		return err
	}

	transmit(conn, message.ServerTaskNew, task)

	return nil
}
