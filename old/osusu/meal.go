package osusu

import (
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Meal is a struct that represents the data of a meal
type Meal struct {
	ID          int64
	GroupID     int64
	Name        string
	Description string
	Source      string
	Image       string
	Category    []string
	Cuisine     []string
}

// Meals is a slice that represents multiple meals
type Meals []Meal

// Score produces a score object for the meal based on the given entries and options
func (m Meal) Score(entries Entries, options Options) Score {
	scores := []ScoreWeight{}
	var latestDate time.Time
	for _, entry := range entries {
		// importance is hours since 1970--makes more recent ones more important
		scores = append(scores, ScoreWeight{Score: entry.Score(options), Weight: int(entry.Date.Unix() / 3600)})
		if entry.Date.After(latestDate) {
			latestDate = entry.Date
		}
	}
	score := AverageScore(scores)
	recencyScore := int(2 * time.Now().Truncate(time.Hour*24).UTC().Sub(latestDate) / (time.Hour * 24))
	if recencyScore > 100 {
		recencyScore = 100
	}
	score.Recency = recencyScore
	return score
}

func (m Meal) LatestDate(entries Entries) time.Time {
	var latestDate time.Time
	for _, entry := range entries {
		if entry.Date.After(latestDate) {
			latestDate = entry.Date
		}
	}
	return latestDate
}

// RemoveInvalidCuisines returns the the meal with all invalid cuisines removed, using the given cuisine options
func (m Meal) RemoveInvalidCuisines(cuisines []string) Meal {
	res := []string{}
	for _, mealCuisine := range m.Cuisine {
		for _, cuisineOption := range cuisines {
			if mealCuisine == cuisineOption {
				res = append(res, mealCuisine)
			}
		}
	}
	m.Cuisine = res
	return m
}

// CurrentMeal gets and returns the current meal state value, using the given context
func CurrentMeal(ctx app.Context) Meal {
	var meal Meal
	ctx.GetState("currentMeal", &meal)
	return meal
}

// SetCurrentMeal sets the current meal state value to the given meal, using the given context
func SetCurrentMeal(meal Meal, ctx app.Context) {
	ctx.SetState("currentMeal", meal, app.Persist)
}

// IsMealNew gets the state value specifying whether the current meal is new
func IsMealNew(ctx app.Context) bool {
	var isMealNew bool
	ctx.GetState("isMealNew", &isMealNew)
	return isMealNew
}

// SetIsMealNew sets the state value specifying whether the current meal is new
func SetIsMealNew(isMealNew bool, ctx app.Context) {
	ctx.SetState("isMealNew", isMealNew, app.Persist)
}
