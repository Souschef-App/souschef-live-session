package session

import (
	"sort"
	"souschef/data"
)

type TaskManager struct {
	TaskRegistry  map[string]*data.Task // taskID → Task
	AssignedTasks map[string]*data.Task
	Dependants    map[string][]string
	QueuedTasks   []*data.Task
}

func CreateTaskManager() *TaskManager {
	return &TaskManager{
		TaskRegistry:  make(map[string]*data.Task),
		AssignedTasks: make(map[string]*data.Task),
		Dependants:    make(map[string][]string),
		QueuedTasks:   []*data.Task{},
	}
}

func (t *TaskManager) AllTasksCompleted() bool {
	for _, task := range t.TaskRegistry {
		if !task.Completed {
			return false
		}
	}

	return true
}

func (t *TaskManager) Init(recipes []*data.Recipe) {
	// Create task registry & preload task queue
	for _, recipe := range recipes {
		for _, task := range recipe.Tasks {
			tempTask := task
			t.TaskRegistry[tempTask.ID] = &tempTask
			if len(tempTask.Dependencies) == 0 {
				t.QueuedTasks = append(t.QueuedTasks, &tempTask)
			}
		}
	}

	// Create reverse dependencies table
	// i.e: Find all tasks that depend on "task_x"
	for taskID, task := range t.TaskRegistry {
		if len(task.Dependencies) > 0 {
			for _, dependantID := range task.Dependencies {
				if _, ok := t.Dependants[dependantID]; !ok {
					// Create key if doesn't exist
					t.Dependants[dependantID] = []string{}
				}

				t.Dependants[dependantID] = append(t.Dependants[dependantID], taskID)
			}
		}
	}
}

func (t *TaskManager) GetTask(skill data.SkillLevel) *data.Task {
	for i, task := range t.QueuedTasks {
		if int(task.Difficulty) <= int(skill) {
			// Remove from queued tasks
			t.QueuedTasks = append(t.QueuedTasks[:i], t.QueuedTasks[i+1:]...)

			// Add to assigned tasks
			t.AssignedTasks[task.ID] = task

			return task
		}
	}

	return nil
}

// Marks a task as completed and removes it from the assigned tasks.
// It checks if any dependent tasks can now be queued and adds them to the queue.
//
// Parameters:
// - taskID: The unique identifier of the task to complete.
//
// Returns:
// - The completed task or nil if the task does not exist in the assigned tasks.
func (t *TaskManager) CompleteTask(taskID string) *data.Task {
	// 1. Mark task as completed
	task, exist := t.AssignedTasks[taskID]
	if !exist {
		return nil
	}

	task.Completed = true

	// 2. Remove task from assigned tasks
	delete(t.AssignedTasks, taskID)

	// 3. Check if any dependent tasks can now be queued
	var oldQueueSize = len(t.QueuedTasks)
	dependants, exist := t.Dependants[task.ID]
	if exist {
		for _, taskID := range dependants {
			task, exist := t.TaskRegistry[taskID]
			if exist && t.hasUncompletedDeps(task) {
				t.QueuedTasks = append(t.QueuedTasks, task)
			}
		}
	}

	var queueModified = len(t.QueuedTasks) != oldQueueSize

	// 4. Re-sort queued tasks based on duration (Optional)
	if queueModified {
		sort.Slice(t.QueuedTasks, func(i, j int) bool {
			return t.QueuedTasks[i].Duration > t.QueuedTasks[j].Duration
		})
	}

	return task
}

func (t *TaskManager) ReassignTask(taskID string, skill data.SkillLevel) *data.Task {
	oldTask := t.AssignedTasks[taskID]
	newTask := t.GetTask(skill)

	// Remove task from assigned tasks
	if newTask != nil {
		delete(t.AssignedTasks, oldTask.ID)
		t.QueuedTasks = append(t.QueuedTasks, oldTask)
		return newTask
	}

	return oldTask
}

func (t *TaskManager) UnassignTask(taskID string) {
	task, exist := t.AssignedTasks[taskID]
	if exist {
		delete(t.AssignedTasks, taskID)
		t.QueuedTasks = append(t.QueuedTasks, task)
	}
}

func (t *TaskManager) hasUncompletedDeps(task *data.Task) bool {
	for _, taskID := range task.Dependencies {
		dependancyTask, exist := t.TaskRegistry[taskID]
		if exist && !dependancyTask.Completed {
			return false
		}
	}

	return true
}
