package session

import (
	"fmt"
	"souschef/data"
	"souschef/internal/utils"
	"sync"
	"time"
)

// TODO: Keep track of live feed
type Session struct {
	IsRunning   bool
	HostID      string
	Recipes     []*data.Recipe
	Livefeed    []*data.FeedSnapshot
	TaskManager *TaskManager
	Observable  *utils.Observable
	mu          sync.Mutex
}

func CreateSession(mealplan data.MealPlan) *Session {
	return &Session{
		IsRunning:   false,
		HostID:      mealplan.HostID,
		Recipes:     mealplan.Recipes,
		Livefeed:    []*data.FeedSnapshot{},
		TaskManager: CreateTaskManager(),
		Observable:  utils.CreateObservable(),
	}
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

func (s *Session) CompleteTask(user *data.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.IsRunning {
		return fmt.Errorf("session has not started")
	}

	// 1. Mark task as completed
	completedTask := s.TaskManager.CompleteTask(user.TaskID)
	if completedTask == nil {
		return fmt.Errorf("task completed not found")
	}

	// 2. Unassign user's task
	user.TaskID = ""
	s.recordSnapshot(user, completedTask, data.Completed)

	// TODO: Refine this idea
	if s.TaskManager.AllTasksCompleted() {
		s.internalStop()
	}

	return nil
}

func (s *Session) RerollTask(user *data.User) (*data.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.IsRunning {
		return nil, fmt.Errorf("session has not started")
	}

	var userSkillLevel = user.SkillLevel
	if user.ID == s.HostID {
		userSkillLevel = data.Expert
	}

	newTask := s.TaskManager.ReassignTask(user.TaskID, userSkillLevel)
	if newTask != nil && user.TaskID != newTask.ID {
		user.TaskID = newTask.ID
		s.recordSnapshot(user, newTask, data.Rerolled)
	}

	return newTask, nil
}

func (s *Session) AssignTask(user *data.User) (*data.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.IsRunning {
		return nil, fmt.Errorf("session has not started")
	}

	if user.TaskID != "" {
		return nil, fmt.Errorf("user already assigned task")
	}

	var userSkillLevel = user.SkillLevel
	if user.ID == s.HostID {
		userSkillLevel = data.Expert
	}

	task := s.TaskManager.GetTask(userSkillLevel)
	if task != nil {
		user.TaskID = task.ID
		s.recordSnapshot(user, task, data.Assigned)
	}

	return task, nil // task can be nil when no suitable task found
}

func (s *Session) recordSnapshot(user *data.User, task *data.Task, status data.TaskStatus) {
	feedSnapshot := &data.FeedSnapshot{
		User:      user,
		Task:      task,
		Status:    status,
		Timestamp: time.Now(),
	}

	// prepend
	s.Livefeed = append([]*data.FeedSnapshot{feedSnapshot}, s.Livefeed...)
	s.Observable.NotifyObservers(feedSnapshot)
}
