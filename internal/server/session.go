package server

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	Status bool
}

type Session struct {
	HostID    string
	IsRunning bool
	Clients   []Client
	mu        sync.Mutex
}

func (s *Session) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.IsRunning = true
}

func (s *Session) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.IsRunning = false
}
