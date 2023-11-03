package data

import "time"

type WelcomeSnapshot struct {
	Users    []*User         `json:"users"`
	Livefeed []*FeedSnapshot `json:"livefeed"`
}

type FeedSnapshot struct {
	User      *User      `json:"user"`
	Task      *Task      `json:"task"`
	Status    TaskStatus `json:"status"`
	Timestamp time.Time  `json:"timestamp"`
}

type User struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	SkillLevel SkillLevel `json:"skillLevel"`
	TaskID     string     `json:"taskID"`
}

type MealPlan struct {
	ID       string       `json:"id"`
	HostID   string       `json:"hostID"`
	Occasion OccasionType `json:"occasion"`
	Recipes  []*Recipe    `json:"recipes"`
}

type Recipe struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Duration    float64       `json:"duration"`
	Difficulty  Difficulty    `json:"difficulty"`
	NumServings int           `json:"numServings"`
	Tasks       []Task        `json:"tasks"`
	Ingredient  []Ingredient  `json:"ingredients"`
	Kitchenware []Kitchenware `json:"kitchenware"`
}

type Task struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Duration     float64       `json:"duration"`
	Difficulty   Difficulty    `json:"difficulty"`
	Priority     int           `json:"priority"`
	Dependencies []string      `json:"dependencies"`
	Ingredients  []Ingredient  `json:"ingredients"`
	Kitchenware  []Kitchenware `json:"kitchenware"`
	Completed    bool
}

type Ingredient struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Quantity float64     `json:"quantity"`
	Unit     CookingUnit `json:"unit"`
}

type Kitchenware struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}
