package server

import (
	"fmt"
	"net/http"
	"souschef/internal/session"
	"sync"

	"github.com/gorilla/websocket"
)

var (
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
	err := http.ListenAndServe(addr, nil)
	if err != nil {
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

	userID := r.Header.Get("UserID")
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
}
