package data

type ClientMessage struct {
	Type string `json:"type"`
}

type ServerMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
