package main

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

// A Recipe is an external recipe that can be used for new meal recommendations
type Recipe struct {
	Name        string
	Ingredients string
	URL         string
	Image       string
	TotalTime   string
	PrepTime    string
	CookTime    string
	Description string
}

// Recipes is a slice of multiple recipes
type Recipes []Recipe

// GetRecipes gets all of the recipes from the recipes.json file
func GetRecipes() (Recipes, error) {
	recipes := Recipes{}
	b, err := os.ReadFile("web/recipes.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &recipes)
	if err != nil {
		return nil, err
	}
	return recipes, nil
}

// FixRecipeTimes returns the given recipes with durations formatted correctly
func FixRecipeTimes(recipes Recipes) Recipes {
	for i, recipe := range recipes {
		var prepDuration, cookDuration, totalDuration time.Duration
		if recipe.PrepTime != "" {
			prepDuration, _ = ParseDuration(recipe.PrepTime)
			recipe.PrepTime = prepDuration.String()
		}
		if recipe.CookTime != "" {
			cookDuration, _ = ParseDuration(recipe.CookTime)
			recipe.CookTime = cookDuration.String()
		}
		if recipe.TotalTime != "" {
			totalDuration, _ = ParseDuration(recipe.TotalTime)
			recipe.TotalTime = totalDuration.String()
		} else {
			recipe.TotalTime = (prepDuration + cookDuration).String()
		}
		recipes[i] = recipe
	}
	return recipes
}

// ParseDuration parses an ISO 8601 duration (the format used in the recipes source)
func ParseDuration(duration string) (time.Duration, error) {
	// remove PT at start
	if len(duration) > 2 {
		duration = duration[2:]
	}
	// change H, M, etc to h, m, etc
	duration = strings.ToLower(duration)
	return time.ParseDuration(duration)
}
