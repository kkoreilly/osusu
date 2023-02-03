package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

// Meal is a struct that represents the data of a meal
type Meal struct {
	Name        string
	Cost        int
	Effort      int
	Healthiness int
}

// Meals ia a map that represents the meals
type Meals map[string]Meal

// SetCurrentMeal sets the current meal state value to the given meal, using the given context
func SetCurrentMeal(meal Meal, ctx app.Context) {
	ctx.SetState("currentMeal", meal, app.Persist)
}

// GetCurrentMeal gets and returns the current meal state value, using the given context
func GetCurrentMeal(ctx app.Context) Meal {
	var meal Meal
	ctx.GetState("currentMeal", &meal)
	return meal
}

// SetMeals sets the meals state value to the given meals, using the given context
func SetMeals(meals Meals, ctx app.Context) {
	ctx.SetState("meals", meals, app.Persist)
}

// GetMeals gets and returns the meals state value, using the given context
func GetMeals(ctx app.Context) Meals {
	var meals Meals
	ctx.GetState("meals", &meals)
	return meals
}
