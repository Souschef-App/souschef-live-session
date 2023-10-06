package data

type MealPlan struct {
	ID       string       `json:"id"`
	HostID   string       `json:"host_id"`
	Occasion OccasionType `json:"occasion"`
	Recipes  []Recipe     `json:"recipes"`
}

type Recipe struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Duration    float64       `json:"duration"`
	Difficulty  int           `json:"difficulty"`
	NumServings int           `json:"num_servings"`
	Tasks       []Task        `json:"tasks"`
	Ingredient  []Ingredient  `json:"ingredients"`
	Kitchenware []Kitchenware `json:"kitchenware"`
}

type Task struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Duration     float64       `json:"duration"`
	Difficulty   int           `json:"difficulty"`
	Priority     int           `json:"priority"`
	Dependencies []Task        `json:"dependencies"`
	Ingredients  []Ingredient  `json:"ingredients"`
	Kitchenware  []Kitchenware `json:"kitchenware"`
	AssigneeID   string
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
