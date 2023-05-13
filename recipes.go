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

	// only used when unmarshalling json
	TotalTime string
	PrepTime  string
	CookTime  string

	// used later after loading json
	TotalDuration time.Duration
	PrepDuration  time.Duration
	CookDuration  time.Duration
}

// Recipes is a slice of multiple recipes
type Recipes []Recipe

// GetRecipes gets all of the recipes from the recipes.json file and validates them
func GetRecipes() (Recipes, error) {
	recipies := Recipes{}
	b, err := os.ReadFile("web/recipes.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &recipies)
	if err != nil {
		return nil, err
	}
	for i, recipe := range recipies {
		if recipe.PrepTime != "" {
			prepDuration, _ := ParseDuration(recipe.PrepTime)
			recipe.PrepDuration = prepDuration
		}
		if recipe.CookTime != "" {
			cookDuration, _ := ParseDuration(recipe.CookTime)
			recipe.CookDuration = cookDuration
		}
		if recipe.TotalTime != "" {
			totalDuration, _ := ParseDuration(recipe.TotalTime)
			recipe.TotalDuration = totalDuration
		} else {
			recipe.TotalDuration = recipe.PrepDuration + recipe.CookDuration
		}
		recipies[i] = recipe
	}
	return recipies, nil
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
