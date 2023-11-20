package data

type TaskStatus int

const (
	Unassigned TaskStatus = iota
	InProgress
	Background
	Completed
)

type FeedAction int

const (
	Assignment FeedAction = iota
	Completion
	Reroll
)

type SkillLevel int

const (
	Beginner SkillLevel = iota
	Intermediate
	Expert
)

type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
)

type OccasionType int

const (
	Home OccasionType = iota
	Professional
	Educational
)

type CookingUnit int

const (
	None CookingUnit = iota
	Ounces
	Pounds
	Grams
	Kilograms
	Teaspoons
	Tablespoons
	Cups
	Pints
	Quarts
	Gallons
	Mililiters
	Liters
)
