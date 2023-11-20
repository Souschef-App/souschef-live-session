package session

import (
	"fmt"
	"souschef/data"
	"souschef/internal/utils"
	"sync"
	"time"
)

type Session struct {
	IsRunning     bool
	HostID        string
	Livefeed      []*data.FeedSnapshot
	RecipeManager *RecipeManager
	Observable    *utils.Observable
	mu            sync.Mutex
}

var Live *Session

func CreateSession(mealplan data.MealPlan) *Session {
	return &Session{
		IsRunning:     false,
		HostID:        mealplan.HostID,
		Livefeed:      []*data.FeedSnapshot{},
		RecipeManager: CreateRecipeManager(mealplan.Recipes),
		Observable:    utils.CreateObservable(),
	}
}

// PUBLIC METHODS

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
	}

	return s.shutdown()
}

func (s *Session) CompleteTask(user *data.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.IsRunning {
		return fmt.Errorf("session has not started")
	}

	// 1. Mark task as completed
	completedTask := s.RecipeManager.CompleteTask(user.TaskID)
	if completedTask == nil {
		return fmt.Errorf("unable to complete task '%s': the task was not found or is not currently in progress", user.TaskID)
	}

	// 2. Unassign user's task
	user.TaskID = ""
	s.recordSnapshot(user, completedTask, data.Completion)

	// 3. Check if all tasks are completed, i.e. session over
	if s.RecipeManager.AllTasksCompleted() {
		s.shutdown()
	}

	return nil
}

// If error is nil, task is guaranteed to not be nil
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

	oldTask, newTask := s.RecipeManager.ReassignTask(user.TaskID, userSkillLevel)
	if oldTask == nil && newTask == nil {
		return nil, fmt.Errorf("failed to reroll task")
	}

	s.recordSnapshot(user, oldTask, data.Reroll)

	if newTask == nil {
		s.recordSnapshot(user, oldTask, data.Assignment)
		return oldTask, nil
	} else {
		user.TaskID = newTask.ID
		s.recordSnapshot(user, newTask, data.Assignment)
		return newTask, nil
	}
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

	task := s.RecipeManager.GetTask(userSkillLevel)
	if task != nil {
		user.TaskID = task.ID
		s.recordSnapshot(user, task, data.Assignment)
	}

	return task, nil // task can be nil when no suitable task found
}

// PRIVATE METHODS

func (s *Session) shutdown() error {
	if !s.IsRunning {
		return fmt.Errorf("session has not started")
	}

	s.IsRunning = false
	fmt.Println("Live session stopped.")

	return nil
}

func (s *Session) recordSnapshot(user *data.User, task *data.Task, status data.FeedAction) {
	feedSnapshot := &data.FeedSnapshot{
		User:      user,
		Task:      task,
		Action:    status,
		Timestamp: time.Now(),
	}

	// prepend
	s.Livefeed = append([]*data.FeedSnapshot{feedSnapshot}, s.Livefeed...)
	s.Observable.NotifyObservers(feedSnapshot)
}
