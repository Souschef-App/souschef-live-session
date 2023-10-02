package server

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var fakeToken = "1w2r2ol3XKhA98LYvHOmLggYcHtrVp2MH3VheZ4cdLA6VKmyzgYQXYtbyTfXqWux"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var connectionsMu sync.Mutex
var connectionsCount int

func StartWebSocket() {
	http.HandleFunc("/ws", handleWebSocket)
	fmt.Println("WebSocket server started on :8080/ws")
	http.ListenAndServe(":8080", nil)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		conn.Close()
		connectionsMu.Lock()
		connectionsCount--
		connectionsMu.Unlock()
		fmt.Printf("Client disconnected. Total connections: %d\n", connectionsCount)
	}()

	connectionsMu.Lock()
	connectionsCount++
	connectionsMu.Unlock()

	fmt.Printf("Client connected. Total connections: %d\n", connectionsCount)

	if !authenticateUser(r) {
		closeMessage := websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Authentication failed")
		if err := conn.WriteMessage(websocket.CloseMessage, closeMessage); err != nil {
			fmt.Println(err)
		}
		return
	}

	for {
		// Read a message from the WebSocket client
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print the received message
		fmt.Printf("Received message: %s\n", p)

		if messageType == websocket.TextMessage {
			if err := RouteMessage(p); err != nil {
				fmt.Println("Error handling message:", err)
				return
			}
		}
	}
}

func authenticateUser(r *http.Request) bool {
	fmt.Println("Authenticating user...")
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		fmt.Println("Invalid or missing Authorization header")
		return false
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	fmt.Println("Token:", token)
	return token == fakeToken
}
