package session

import (
	"fmt"
	"souschef/data"
	"sync"
)

type Session struct {
	IsRunning   bool
	HostID      string
	Helpers     map[string]*data.Helper
	Recipes     []data.Recipe
	TaskManager TaskManager
	mu          sync.Mutex
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
	s.TaskManager.Init(s.Recipes)
	fmt.Println("Live session started!")

	return nil
}

func (s *Session) Stop(userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if userID != s.HostID {
		return fmt.Errorf("only the session host can stop the session")
	}

	return s.internalStop()
}

func (s *Session) internalStop() error {
	if !s.IsRunning {
		return fmt.Errorf("session has not started")
	}

	s.IsRunning = false
	fmt.Println("Live session stopped.")

	// TODO: Reset algorithm to default

	return nil
}

func (s *Session) AddHelper(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Helpers[userID] = &data.Helper{
		Skill: data.Expert,
	}
}

func (s *Session) RemoveHelper(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	helper := s.Helpers[userID]
	if helper.TaskID != "" {
		s.TaskManager.UnassignTask(helper.TaskID)
	}

	delete(s.Helpers, userID)
}

func (s *Session) CompleteTask(userID string) (*data.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.IsRunning {
		return nil, fmt.Errorf("session has not started")
	}

	helper, exist := s.Helpers[userID]
	if !exist {
		return nil, fmt.Errorf("user not found")
	}

	task := s.TaskManager.CompleteTask(helper.TaskID)
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}

	// TODO: Refine this idea
	if s.TaskManager.AllTasksCompleted() {
		s.internalStop()
	}

	helper.TaskID = ""

	return task, nil
}

func (s *Session) RerollTask(userID string) (*data.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.IsRunning {
		return nil, fmt.Errorf("session has not started")
	}

	helper, exist := s.Helpers[userID]
	if !exist {
		return nil, fmt.Errorf("user not found")
	}

	newTask := s.TaskManager.ReassignTask(helper.TaskID, helper.Skill)
	helper.TaskID = ""
	if newTask != nil {
		helper.TaskID = newTask.ID
	}

	return newTask, nil
}

func (s *Session) AssignTask(userID string) (*data.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.IsRunning {
		return nil, fmt.Errorf("session has not started")
	}

	helper, exist := s.Helpers[userID]
	if !exist {
		return nil, fmt.Errorf("user not found")
	} else if helper.TaskID != "" {
		return nil, fmt.Errorf("user already assigned task")
	}

	task := s.TaskManager.GetTask(helper.Skill)
	if task != nil {
		helper.TaskID = task.ID
	}

	return task, nil // task can be nil, meaning no suitable task
}
