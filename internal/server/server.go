package server

import (
	"context"
	"fmt"
	"net/http"
	"souschef/data"
	"souschef/internal/message"
	"souschef/internal/session"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	server   *http.Server
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Origin validation logic here
			return true
		},
	}
	connectionsMu sync.Mutex
	connections   = make(map[*websocket.Conn]*data.User)
)

type LivefeedBroadcaster struct{}

func (l *LivefeedBroadcaster) Update(snapshot any) {
	broadcast(message.ServerFeedSnapshot, snapshot)
}

func StartWebSocket(addr string) {
	observer := &LivefeedBroadcaster{}
	ticker := time.NewTicker(60 * time.Second)
	defer func() {
		session.Live.Observable.UnregisterObserver(observer)
		ticker.Stop()
	}()

	session.Live.Observable.RegisterObserver(observer)

	go func() {
		for range ticker.C {
			broadcast(message.ServerTimestampUpdate, nil)
		}
	}()

	http.HandleFunc("/ws", handleWebSocket)
	fmt.Println("WebSocket server started on :8080/ws")

	server = &http.Server{Addr: addr}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println(err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: Close connection if no handshake occurs
	handleMessage(conn) // Blocking call
	unregisterConnection(conn)
}

// Gets called inside handlers.go
func registerConnection(conn *websocket.Conn, user *data.User) {
	connectionsMu.Lock()
	defer connectionsMu.Unlock()

	var users = []*data.User{}

	// Enforce single connection per user.
	for existConn, existUser := range connections {
		// Replace existing connection with new if user is already connected.
		if user.ID == existUser.ID {
			existConn.Close()
			delete(connections, existConn)
		}

		// Compile list of connected users for welcome snapshot
		users = append(users, existUser)
	}

	connections[conn] = user
	users = append(users, user)

	broadcast(message.ServerClientConnected, user)

	welcomeSnapshot := &data.WelcomeSnapshot{
		Users:    users,
		Tasks:    session.Live.RecipeManager.Registry,
		Livefeed: session.Live.Livefeed,
	}

	transmit(conn, message.ServerHandshake, welcomeSnapshot)

	fmt.Printf("Client connected. Total connections: %d\n", len(connections))
}

func unregisterConnection(conn *websocket.Conn) {
	connectionsMu.Lock()

	user, exist := connections[conn]
	if exist {
		// Gracefully handle uncompleted tasks
		session.Live.RecipeManager.UnassignTask(user.TaskID)
		processWaitingQueue()

		delete(connections, conn)

		// Notify other connected users
		broadcast(message.ServerClientDisconnected, user)
	}

	conn.Close()
	fmt.Printf("Client disconnected. Total connections: %d\n", len(connections))
	connectionsMu.Unlock()

	if len(connections) == 0 {
		if err := server.Shutdown(context.TODO()); err != nil {
			fmt.Println("HTTP server shutdown error:", err)
		}
	}
}
