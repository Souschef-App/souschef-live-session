package main

import (
	"fmt"
	"souschef/data"
	"souschef/internal/server"
	"souschef/internal/session"
)

func main() {
	fmt.Println("Initializing...")

	// TODO: Kafka message:
	// 1. Request websocket IP & mealplan
	requestResolved := make(chan struct{})

	go func() {
		var mealplan = data.MealPlan{
			ID:       "123",
			HostID:   "123",
			Occasion: data.Home,
			Recipes:  []data.Recipe{data.DefaultRecipe},
		}

		session.Live = &session.Session{
			IsRunning: false,
			Helpers:   make(map[string]*data.Helper),
			HostID:    mealplan.HostID,
			Recipes:   mealplan.Recipes,
			TaskManager: session.TaskManager{
				QueuedTasks:   []*data.Task{},
				TaskRegistry:  make(map[string]*data.Task),
				AssignedTasks: make(map[string]*data.Task),
				Dependants:    make(map[string][]string),
			},
		}

		close(requestResolved)
	}()

	<-requestResolved

	session.Live.Start("123")
	server.StartWebSocket(":8080")

	// TODO: Kafka message:
	// 1. Notify "ready for client connections"
}
