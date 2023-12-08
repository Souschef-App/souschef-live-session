package session

import (
	"sort"
	"souschef/data"
	"time"
)

type TaskManager struct {
	Registry        map[string]*data.Task
	Dependants      map[string][]string
	AssignableTasks []*data.Task
	Progress        float32
	completedCount  int
}

func CreateTaskManager(recipe *data.Recipe) *TaskManager {
	tm := &TaskManager{
		Registry:        make(map[string]*data.Task),
		Dependants:      make(map[string][]string),
		AssignableTasks: []*data.Task{},
		Progress:        0.0,
		completedCount:  0,
	}

	// Populate Registry & AssignableTasks
	for index, task := range recipe.Tasks {
		tempTask := task
		tempTask.Order = index
		tempTask.Status = data.Unassigned // SAFETY
		tm.Registry[tempTask.ID] = &tempTask
		if len(tempTask.Dependencies) == 0 {
			tm.AssignableTasks = append(tm.AssignableTasks, &tempTask)
		}
	}

	// Populate Dependants by creating a reverse dependency map
	// i.e: Find all tasks that depend on task_x
	for taskID, task := range tm.Registry {
		if len(task.Dependencies) > 0 {
			for _, dependant := range task.Dependencies {
				if _, ok := tm.Dependants[dependant.DependencyID]; !ok {
					// Create entry if doesn't exist
					tm.Dependants[dependant.DependencyID] = []string{}
				}

				tm.Dependants[dependant.DependencyID] = append(tm.Dependants[dependant.DependencyID], taskID)
			}
		}
	}

	return tm
}

func (t *TaskManager) FindEligibleTask(skill data.SkillLevel) *data.Task {
	for i, task := range t.AssignableTasks {
		if int(task.Difficulty) <= int(skill) {
			if t.markTaskAssigned(task) {
				// Remove
				t.AssignableTasks = append(t.AssignableTasks[:i], t.AssignableTasks[i+1:]...)
			}

			return task
		}
	}

	return nil
}

func (t *TaskManager) CompleteTask(taskID string) bool {
	task, exist := t.Registry[taskID]
	if !exist {
		return false
	}

	// Task can only be completed if it is either in progress or in the background
	if task.Status != data.InProgress && task.Status != data.Background {
		return false
	}

	if task.IsBackground && task.Status == data.InProgress {
		task.Status = data.Background
		task.Timestamp = time.Now()
		return true
	}

	task.Status = data.Completed
	task.Timestamp = time.Now()

	t.completedCount += 1
	t.calculateProgress()
	t.tryAddingDepsToAssignable(taskID)

	return true
}

func (t *TaskManager) UnassignTask(taskID string) {
	task, exist := t.Registry[taskID]
	if exist && task.Status == data.InProgress {
		task.Status = data.Unassigned
		task.Timestamp = time.Now()
		t.AssignableTasks = append(t.AssignableTasks, task)
	}
}

// PRIVATE METHODS

func (t *TaskManager) tryAddingDepsToAssignable(taskID string) {
	var queueModified = false

	dependants, exist := t.Dependants[taskID]
	if exist {
		for _, taskID := range dependants {
			task, exist := t.Registry[taskID]
			// SAFETY CHECK: Task should always be unassigned, but we check anyways
			if exist && task.Status == data.Unassigned && !t.hasUncompletedDeps(task) {
				t.AssignableTasks = append(t.AssignableTasks, task)
				queueModified = true
			}
		}
	}

	// Sort based on task order
	if queueModified {
		sort.Slice(t.AssignableTasks, func(i, j int) bool {
			return t.AssignableTasks[i].Order > t.AssignableTasks[j].Order
		})
	}
}

func (t *TaskManager) markTaskAssigned(task *data.Task) bool {
	if task.Status == data.Completed || task.Status == data.Background {
		return false
	}

	// Task is either unassigned or in-progress (reroll)
	task.Status = data.InProgress
	task.Timestamp = time.Now()
	return true
}

func (t *TaskManager) hasUncompletedDeps(task *data.Task) bool {
	for _, taskDependency := range task.Dependencies {
		dependancyTask, exist := t.Registry[taskDependency.DependencyID]
		if exist && dependancyTask.Status != data.Completed {
			return true
		}
	}

	return false
}

func (r *TaskManager) calculateProgress() {
	numTask := len(r.Registry)

	if numTask > 0 {
		r.Progress = float32(r.completedCount) / float32(numTask)
	}
}
