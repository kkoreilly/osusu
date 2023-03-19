package main

import (
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Meal is a struct that represents the data of a meal
type Meal struct {
	ID          int
	Name        string
	Description string
	Cost        int
	Effort      int
	Healthiness int
	Taste       map[int]int // key is person id, value is taste rating
	Type        string
	Source      string
	LastDone    time.Time
	UserID      int
}

// Meals is a slice that represents multiple meals
type Meals []Meal

// Score produces a score from 0 to 100 for the meal based on its attributes and the given options
func (m Meal) Score(options Options) int {
	// average of all attributes
	var tasteSum int
	for i, v := range m.Taste {
		use := options.People[i]
		// invert the person's rating if they are not participating
		if use {
			tasteSum += v
		} else {
			tasteSum += 100 - v
		}
	}
	recencyScore := int(2 * time.Now().Truncate(time.Hour*24).UTC().Sub(m.LastDone) / (time.Hour * 24))
	if recencyScore > 100 {
		recencyScore = 100
	}
	sum := options.CostWeight*(100-m.Cost) + options.EffortWeight*(100-m.Effort) + options.HealthinessWeight*m.Healthiness + options.TasteWeight*tasteSum + options.RecencyWeight*recencyScore
	den := options.CostWeight + options.EffortWeight + options.HealthinessWeight + len(m.Taste)*options.TasteWeight + options.RecencyWeight
	if den == 0 {
		return 0
	}
	return sum / den
}

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
