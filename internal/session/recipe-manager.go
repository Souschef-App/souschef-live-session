package session

import (
	"sort"
	"souschef/data"
)

type RecipeManager struct {
	Recipes  []*TaskManager
	Registry map[string]*data.Task
}

func CreateRecipeManager(recipes []*data.Recipe) *RecipeManager {
	rm := &RecipeManager{
		Recipes:  []*TaskManager{},
		Registry: make(map[string]*data.Task),
	}

	// Populate recipe registry
	for _, recipe := range recipes {
		taskManager := CreateTaskManager(recipe)
		rm.Recipes = append(rm.Recipes, taskManager)
		for taskID, task := range taskManager.Registry {
			rm.Registry[taskID] = task
		}
	}

	return rm
}

func (r *RecipeManager) AllTasksCompleted() bool {
	for _, recipe := range r.Recipes {
		if recipe.Progress != 1 {
			return false
		}
	}

	return true
}

func (r *RecipeManager) GetTask(skill data.SkillLevel) *data.Task {
	// Sort recipes by progress (lowest to highest)
	sort.Slice(r.Recipes, func(i, j int) bool {
		return r.Recipes[i].Progress < r.Recipes[j].Progress
	})

	for _, recipe := range r.Recipes {
		task := recipe.FindEligibleTask(skill)
		if task != nil {
			return task
		}
	}

	return nil
}

func (r *RecipeManager) CompleteTask(taskID string) *data.Task {
	for _, recipe := range r.Recipes {
		if recipe.CompleteTask(taskID) {
			return recipe.Registry[taskID]
		}
	}

	return nil
}

func (r *RecipeManager) ReassignTask(taskID string, skill data.SkillLevel) (*data.Task, *data.Task) {
	for _, recipe := range r.Recipes {
		oldTask, exist := recipe.Registry[taskID]
		if exist {
			newTask := recipe.FindEligibleTask(skill)
			if newTask != nil {
				recipe.UnassignTask(taskID)
				return oldTask, newTask
			}

			return oldTask, nil
		}
	}

	return nil, nil
}

func (r *RecipeManager) UnassignTask(taskID string) {
	for _, recipe := range r.Recipes {
		_, exist := recipe.Registry[taskID]
		if exist {
			recipe.UnassignTask(taskID)
		}
	}
}
