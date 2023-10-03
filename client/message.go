package client

import (
	"encoding/json"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type SessionMessage struct {
	HostID string         `json:"host_id"`
	Cmd    SessionCommand `json:"command"`
}

type TaskMessage struct {
	UserID string      `json:"user_id"`
	Cmd    TaskCommand `json:"command"`
}
