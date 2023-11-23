package data

import (
	"time"
)

type WelcomeSnapshot struct {
	Users    []*User          `json:"users"`
	Tasks    map[string]*Task `json:"tasks"`
	Livefeed []*FeedSnapshot  `json:"livefeed"`
}

type FeedSnapshot struct {
	User      *User      `json:"user"`
	Task      *Task      `json:"task"`
	Action    FeedAction `json:"action"`
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
	Dependencies []string      `json:"dependencies"`
	Ingredients  []Ingredient  `json:"ingredients"`
	Kitchenware  []Kitchenware `json:"kitchenware"`
	IsBackground bool          `json:"isBackgroundTask"`
	Status       TaskStatus    `json:"status"`
	Timestamp    time.Time     `json:"timestamp"`
	Order        int           `json:"-"`
}

type Ingredient struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Quantity Fraction    `json:"quantity"`
	Unit     CookingUnit `json:"unit"`
}

type Fraction struct {
	Whole       int64 `json:"whole"`
	Numerator   int64 `json:"numerator"`
	Denominator int64 `json:"denominator"`
}

type Kitchenware struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}
