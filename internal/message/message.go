package message

// Receive
type ClientMessage struct {
	Type ClientMessageEnum `json:"type"`
}

type ClientMessageEnum string

const (
	ClientStartSession  ClientMessageEnum = "session_start"
	ClientStopSession   ClientMessageEnum = "session_stop"
	ClientCompletedTask ClientMessageEnum = "task_completed"
	ClientRerollTask    ClientMessageEnum = "task_reroll"
)

// Send
type ServerMessage struct {
	Type    ServerMessageEnum `json:"type"`
	Payload interface{}       `json:"payload"`
}

// Send
type ServerMessageEnum string

const (
	ServerError         ServerMessageEnum = "error"
	ServerTaskNew       ServerMessageEnum = "task_new"
	ServerTaskCompleted ServerMessageEnum = "task_completed"
	ServerMealCompleted ServerMessageEnum = "meal_completed"
)
