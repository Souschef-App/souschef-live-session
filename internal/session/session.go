package session

import (
	"fmt"
	"souschef/data"
	"sync"
)

type Helper struct {
	TaskID string
}

type Session struct {
	IsRunning bool
	HostID    string
	Helpers   map[string]*Helper
	Recipes   []data.Recipe
	mu        sync.Mutex
}

var Live *Session

func (s *Session) Start(userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if userID != s.HostID {
		return fmt.Errorf("only the session host can start the session")
	} else if s.IsRunning {
		return fmt.Errorf("session already started")
	}

	s.IsRunning = true
	fmt.Println("Live session started!")

	return nil
}

func (s *Session) Stop(userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if userID != s.HostID {
		return fmt.Errorf("only the session host can stop the session")
	} else if !s.IsRunning {
		return fmt.Errorf("session has not started")
	}

	s.IsRunning = false
	fmt.Println("Live session stopped.")

	// TODO: Reset to default

	return nil
}

func (s *Session) AddHelper(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Helpers[userID] = &Helper{}
}

func (s *Session) RemoveHelper(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO: Handle unfinished tasks

	delete(s.Helpers, userID)
}

func (s *Session) CompleteTask(userID string) (*data.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	helper, exist := s.Helpers[userID]
	if !exist {
		return nil, fmt.Errorf("user not found")
	}

	// TODO:
	// 1. Get users task (might have error: "task completion desync" where user has no assigned task)
	task := &data.Task{}
	// 2. Complete task in algorithm
	// 3. Remove task from helper
	helper.TaskID = ""

	return task, nil
}

func (s *Session) AssignTask(userID string) (*data.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	helper, exist := s.Helpers[userID]
	if !exist {
		return nil, fmt.Errorf("user not found")
	} else if helper.TaskID != "" {
		return nil, fmt.Errorf("user already assigned task")
	}

	// TODO: Get task from algorithm
	task := &data.Task{}

	helper.TaskID = task.ID
	task.AssigneeID = userID

	return task, nil
}
