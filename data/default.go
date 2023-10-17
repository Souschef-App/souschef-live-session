package data

var DefaultRecipe = Recipe{
	ID:          "123",
	Name:        "Savory Party Bread",
	Duration:    1500,
	NumServings: 8,
	Difficulty:  Easy,
	Tasks: []Task{
		{
			ID:           "task-3",
			Title:        "Warm Pizza Sauce",
			Description:  "Warm pizza sauce in a microwave-safe bowl until hot, about 45 seconds.",
			Duration:     1.0,
			Difficulty:   Easy,
			Priority:     0,
			Dependencies: []string{},
			Ingredients: []Ingredient{
				{
					ID:       "ingredient-2",
					Name:     "Pizza Sauce",
					Quantity: 0.5,
					Unit:     Cups,
				},
			},
			Kitchenware: []Kitchenware{
				{
					ID:       "kw-3",
					Name:     "Microwave-Safe Bowl",
					Quantity: 1,
				},
				{
					ID:       "kw-4",
					Name:     "Microwave",
					Quantity: 1,
				},
			},
		},
		{
			ID:           "task-6",
			Title:        "Slice Pizza",
			Description:  "Cut the baked pizza crosswise into slices. Serve and enjoy.",
			Duration:     5.0,
			Difficulty:   Easy,
			Priority:     0,
			Dependencies: []string{"task-5"},
			Ingredients:  []Ingredient{},
			Kitchenware: []Kitchenware{
				{
					ID:       "kw-2",
					Name:     "Knife",
					Quantity: 1,
				},
			},
		},
		{
			ID:           "task-1",
			Title:        "Preheat Oven",
			Description:  "Preheat the oven to 375 degrees F (190 degrees C).",
			Duration:     10.0,
			Difficulty:   Easy,
			Priority:     1,
			Dependencies: []string{},
			Ingredients:  []Ingredient{},
			Kitchenware: []Kitchenware{
				{
					ID:       "kw-1",
					Name:     "Oven",
					Quantity: 1,
				},
			},
		},
		{
			ID:           "task-2",
			Title:        "Prepare Baguette",
			Description:  "Split the French baguette in half lengthwise.",
			Duration:     2.0,
			Difficulty:   Medium,
			Priority:     0,
			Dependencies: []string{},
			Ingredients: []Ingredient{
				{
					ID:       "ingredient-1",
					Name:     "French Baguette",
					Quantity: 1.0,
					Unit:     None,
				},
			},
			Kitchenware: []Kitchenware{
				{
					ID:       "kw-2",
					Name:     "Knife",
					Quantity: 1,
				},
			},
		},
		{
			ID:           "task-5",
			Title:        "Bake Pizza",
			Description:  "Place the assembled pizza onto a baking tray and bake in the preheated oven until cheese is golden, about 18 minutes. Optional step: turn on the oven’s broiler, set a rack 6 inches below the heating element, and broil the pizza for a deeper color, 1 to 2 minutes.",
			Duration:     20.0,
			Difficulty:   Easy,
			Priority:     0,
			Dependencies: []string{"task-1", "task-4"},
			Ingredients:  []Ingredient{},
			Kitchenware: []Kitchenware{
				{
					ID:       "kw-5",
					Name:     "Baking Tray",
					Quantity: 1,
				},
			},
		},
		{
			ID:           "task-4",
			Title:        "Assemble Pizza",
			Description:  "Spread pizza sauce on the baguette, sprinkle on cheese, and scatter sausage and pepperoni on top as desired.",
			Duration:     5.0,
			Difficulty:   Hard,
			Priority:     0,
			Dependencies: []string{"task-3", "task-2", "task-4.2", "task-4.3"},
			Ingredients: []Ingredient{
				{
					ID:       "ingredient-3",
					Name:     "Mozzarella Cheese",
					Quantity: 0.66,
					Unit:     Cups,
				},
				{
					ID:       "ingredient-4",
					Name:     "Sausage",
					Quantity: 2.0,
					Unit:     Ounces,
				},
				{
					ID:       "ingredient-5",
					Name:     "Mini Pepperoni Slices",
					Quantity: 2.0,
					Unit:     Ounces,
				},
			},
			Kitchenware: []Kitchenware{},
		},
		{
			ID:           "task-4.2",
			Title:        "Grate Cheese",
			Description:  "Grate Mozzarella cheese for the pizza topping.",
			Duration:     3.0,
			Difficulty:   Medium,
			Dependencies: []string{},
			Ingredients: []Ingredient{
				{
					ID:       "ingredient-3",
					Name:     "Mozzarella Cheese",
					Quantity: 0.66,
					Unit:     Cups,
				},
			},
			Kitchenware: []Kitchenware{
				{
					ID:       "kw-6",
					Name:     "Grater",
					Quantity: 1,
				},
			},
		},
		{
			ID:           "task-4.3",
			Title:        "Slice Sausage",
			Description:  "Slice the sausage for the pizza.",
			Duration:     2.0,
			Difficulty:   Medium,
			Dependencies: []string{},
			Ingredients: []Ingredient{
				{
					ID:       "ingredient-4",
					Name:     "Sausage",
					Quantity: 1,
					Unit:     None,
				},
			},
			Kitchenware: []Kitchenware{
				{
					ID:       "kw-2",
					Name:     "Knife",
					Quantity: 1,
				},
			},
		},
	},
	Ingredient: []Ingredient{
		{
			ID:       "ingredient-1",
			Name:     "French Baguette",
			Quantity: 1.0,
			Unit:     None,
		},
		{
			ID:       "ingredient-2",
			Name:     "Pizza Sauce",
			Quantity: 0.5,
			Unit:     Cups,
		},
		{
			ID:       "ingredient-3",
			Name:     "Mozzarella Cheese",
			Quantity: 0.66,
			Unit:     Cups,
		},
		{
			ID:       "ingredient-4",
			Name:     "Sausage",
			Quantity: 2.0,
			Unit:     Ounces,
		},
		{
			ID:       "ingredient-5",
			Name:     "Mini Pepperoni Slices",
			Quantity: 2.0,
			Unit:     Ounces,
		},
	},
	Kitchenware: []Kitchenware{
		{
			ID:       "kw-1",
			Name:     "Oven",
			Quantity: 1,
		},
		{
			ID:       "kw-2",
			Name:     "Knife",
			Quantity: 2,
		},
		{
			ID:       "kw-3",
			Name:     "Microwave-Safe Bowl",
			Quantity: 1,
		},
		{
			ID:       "kw-4",
			Name:     "Microwave",
			Quantity: 1,
		},
		{
			ID:       "kw-5",
			Name:     "Baking Tray",
			Quantity: 1,
		},
		{
			ID:       "kw-6",
			Name:     "Grater",
			Quantity: 1,
		},
	},
}
