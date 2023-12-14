package server

import (
	"context"
	"encoding/json"
	"fmt"
	"souschef/data"
	"souschef/internal/message"
	"souschef/internal/session"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	ticker          *time.Ticker = nil
	observer                     = &LivefeedBroadcaster{}
	waitingQueue                 = map[*websocket.Conn]*data.User{}
	commandHandlers              = map[message.ClientMessageEnum]func(*websocket.Conn, json.RawMessage) error{
		message.ClientHandshake:               handleClientHandshake,
		message.ClientGuestHandshake:          handleGuestHandshake,
		message.ClientCreateSession:           handleSessionCreate,
		message.ClientStartSession:            handleSessionStart,
		message.ClientStopSession:             handleSessionStop,
		message.ClientCompletedTask:           handleTaskCompleted,
		message.ClientRerolledTask:            handleTaskRerolled,
		message.ClientCompletedBackgroundTask: handleTaskBackgroundCompleted,
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
		} else if task == nil {
			transmit(conn, message.ServerTaskNew, nil)
		} else {
			transmit(conn, message.ServerTaskNew, task.ID)
		}
	}

	return nil
}

func handleGuestHandshake(conn *websocket.Conn, payload json.RawMessage) error {
	var guestName string
	if err := json.Unmarshal(payload, &guestName); err != nil {
		return err
	}

	guest := &data.User{
		ID:         uuid.NewString(),
		Name:       guestName,
		SkillLevel: data.Expert,
	}

	registerConnection(conn, guest)

	if session.Live.IsRunning {
		task, err := session.Live.AssignTask(guest)
		if err != nil {
			return err
		} else if task == nil {
			transmit(conn, message.ServerTaskNew, nil)
		} else {
			transmit(conn, message.ServerTaskNew, task.ID)
		}
	}

	return nil
}

type LivefeedBroadcaster struct{}

func (l *LivefeedBroadcaster) Update(snapshot any) {
	broadcast(message.ServerFeedSnapshot, snapshot)
}

func handleSessionCreate(conn *websocket.Conn, payload json.RawMessage) error {
	var mealplan data.MealPlan
	if err := json.Unmarshal(payload, &mealplan); err != nil {
		return err
	}

	// Add default recipe temporarily
	mealplan.Recipes = append(mealplan.Recipes, &data.DefaultRecipe)

	fmt.Println("Creating Session...")

	session.Live = session.CreateSession(mealplan)

	// Setup observer
	session.Live.Observable.RegisterObserver(observer)
	if ticker == nil {
		ticker = time.NewTicker(60 * time.Second)
	} else {
		ticker.Reset(60 * time.Second)
	}

	go func() {
		for range ticker.C {
			broadcast(message.ServerTimestampUpdate, nil)
		}
	}()

	session.Live.Start(mealplan.HostID)
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
		if err != nil {
			fmt.Println("Failed to assign task to user: ", user.ID)
		} else if task == nil {
			transmit(conn, message.ServerTaskNew, nil)
		} else {
			transmit(conn, message.ServerTaskNew, task.ID)
		}
	}

	return nil
}

func handleSessionStop(conn *websocket.Conn, _ json.RawMessage) error {
	user, exist := connections[conn]
	if !exist {
		return fmt.Errorf("user not found")
	}

	if err := session.Live.Stop(user.ID); err != nil {
		return err
	}

	waitingQueue = nil

	if err := server.Shutdown(context.TODO()); err != nil {
		fmt.Println("HTTP server shutdown error:", err)
	}

	return nil
}

func handleTaskCompleted(conn *websocket.Conn, _ json.RawMessage) error {
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
		// 5. Transmit no task available to the user
		transmit(conn, message.ServerTaskNew, nil)
	} else {
		// 5.1 Transmit the new task to the user
		transmit(conn, message.ServerTaskNew, task.ID)
	}

	return nil
}

func processWaitingQueue() {
	for conn, user := range waitingQueue {
		task, err := session.Live.AssignTask(user)
		if err == nil && task != nil {
			transmit(conn, message.ServerTaskNew, task.ID)
			delete(waitingQueue, conn)
		}
	}
}

func handleTaskRerolled(conn *websocket.Conn, _ json.RawMessage) error {
	user, exist := connections[conn]
	if !exist {
		return fmt.Errorf("user not found")
	}

	// Reroll guarantees non-nil task if no error
	task, err := session.Live.RerollTask(user)
	if err != nil {
		return err
	}

	transmit(conn, message.ServerTaskNew, task.ID)

	return nil
}

func handleTaskBackgroundCompleted(conn *websocket.Conn, payload json.RawMessage) error {
	user, exist := connections[conn]
	if !exist {
		return fmt.Errorf("user not found")
	}

	var taskID string
	if err := json.Unmarshal(payload, &taskID); err != nil {
		return err
	}

	if err := session.Live.CompleteBackgroundTask(user, taskID); err != nil {
		return err
	}

	// Check if session is still running
	if !session.Live.IsRunning {
		broadcast(message.ServerMealCompleted, nil)
	}

	processWaitingQueue()

	return nil
}
