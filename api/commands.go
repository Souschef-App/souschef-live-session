package api

type Command struct {
	Type string `json:"type"`
}

type SessionCommand struct {
	Command
	HostUserID int `json:"hostUserID"`
}
