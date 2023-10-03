package client

type SessionCommand int

const (
	StartSession SessionCommand = iota
	StopSession
)

type TaskCommand int

const (
	CompleteTask TaskCommand = iota
	RerollTask
)
