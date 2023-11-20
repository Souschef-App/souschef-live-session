package message

import "encoding/json"

// Receive
type ClientMessage struct {
	Type    ClientMessageEnum `json:"type"`
	Payload json.RawMessage   `json:"payload"`
}

// Receive
type ClientMessageEnum string

const (
	ClientHandshake               ClientMessageEnum = "handshake"
	ClientGuestHandshake          ClientMessageEnum = "guest_handshake"
	ClientStartSession            ClientMessageEnum = "session_start"
	ClientStopSession             ClientMessageEnum = "session_stop"
	ClientCompletedTask           ClientMessageEnum = "task_completed"
	ClientRerolledTask            ClientMessageEnum = "task_rerolled"
	ClientCompletedBackgroundTask ClientMessageEnum = "task_background_completed"
)

// Send
type ServerMessage struct {
	Type    ServerMessageEnum `json:"type"`
	Payload any               `json:"payload"`
}

// Send
type ServerMessageEnum string

const (
	ServerError              ServerMessageEnum = "error"
	ServerHandshake          ServerMessageEnum = "server_handshake"
	ServerClientConnected    ServerMessageEnum = "client_connected"
	ServerClientDisconnected ServerMessageEnum = "client_disconnected"
	ServerMealCompleted      ServerMessageEnum = "meal_completed"
	ServerTaskNew            ServerMessageEnum = "task_new"
	ServerFeedSnapshot       ServerMessageEnum = "feed_snapshot"
	ServerTimestampUpdate    ServerMessageEnum = "timestamp_update"
)
