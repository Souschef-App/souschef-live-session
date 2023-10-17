package server

import (
	"context"
	"fmt"
	"net/http"
	"souschef/internal/message"
	"souschef/internal/session"
	"sync"

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
	connections   = make(map[*websocket.Conn]string)
)

func StartWebSocket(addr string) {
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

	// userID := r.Header.Get("UserID")
	userID := r.URL.Query().Get("UserID")
	if userID == "" {
		reportError(conn, fmt.Errorf("userID missing in header or invalid"))
		return
	}

	registerConnection(conn, userID)
	handleMessage(conn) // Blocking call
	unregisterConnection(conn)
}

func registerConnection(conn *websocket.Conn, userID string) {
	connectionsMu.Lock()

	fmt.Println("New User Joined with ID:", userID)

	// Enforce single connection per user.
	// Replace existing connection with new if user is already connected.
	for oldConn, oldUserID := range connections {
		if userID == oldUserID {
			oldConn.Close()
			delete(connections, oldConn)
			break
		}
	}

	// Associate connection with userID
	connections[conn] = userID
	fmt.Printf("Client connected. Total connections: %d\n", len(connections))
	connectionsMu.Unlock()

	// Associate userID with Helper
	session.Live.AddHelper(userID)

	// Update user immediately if session is running
	if session.Live.IsRunning {
		task, err := session.Live.AssignTask(userID)
		if err != nil {
			transmit(conn, message.ServerError, err.Error())
			return
		}

		transmit(conn, message.ServerTaskNew, task)
	}
}

func unregisterConnection(conn *websocket.Conn) {
	connectionsMu.Lock()

	// Remove userID to Helper association
	userID, exist := connections[conn]
	if exist {
		session.Live.RemoveHelper(userID)
	}

	// Remove connection to userID association
	delete(connections, conn)

	conn.Close()
	fmt.Printf("Client disconnected. Total connections: %d\n", len(connections))
	connectionsMu.Unlock()

	if len(connections) == 0 {
		if err := server.Shutdown(context.TODO()); err != nil {
			fmt.Println("HTTP server shutdown error:", err)
		}
	}
}
