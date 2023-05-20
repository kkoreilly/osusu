package main

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
	"time"
)

// A Recipe is an external recipe that can be used for new meal recommendations
type Recipe struct {
	Name string
	// Ingredients string
	URL         string
	Image       string
	TotalTime   string
	PrepTime    string
	CookTime    string
	Score       Score
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

// SaveRecipes writes the given recipes to the web/newrecipes.json file
func SaveRecipes(recipes Recipes) error {
	jsonData, err := json.Marshal(recipes)
	if err != nil {
		return err
	}
	err = os.WriteFile("web/newrecipes.json", jsonData, 0666)
	return err
}

// GenerateWordMap represents a map with words as keys and all of the recipes that contain them as values
func GenerateWordMap(recipes Recipes) map[string]Recipes {
	res := map[string]Recipes{}
	for _, recipe := range recipes {
		words := GetWords(recipe.Name)
		for _, word := range words {
			if res[word] == nil {
				res[word] = Recipes{}
			}
			res[word] = append(res[word], recipe)
		}
	}
	return res
}

// RecommendRecipesData is the data used in a recommend recipes call
type RecommendRecipesData struct {
	WordScoreMap map[string]Score
	Options      Options
	N            int // the iteration that we are on (ie: n = 3 for the fourth time we are getting recipes for the same options), we return 100 new meals each time
}

// RecommendRecipes returns a list of recommended new recipes based on the given recommend recipes data
func RecommendRecipes(data RecommendRecipesData) Recipes {
	if len(allRecipes) == 0 {
		return Recipes{}
	}
	scores := []Score{}
	indices := []int{}
	for i, recipe := range allRecipes {
		indices = append(indices, i)
		words := GetWords(recipe.Name)
		recipeScores := []Score{}
		// use map to track unique matches and simple int to count total
		matches := map[string]bool{}
		numMatches := 0
		for _, word := range words {
			score, ok := data.WordScoreMap[word]
			if ok {
				matches[word] = true
				numMatches++
				recipeScores = append(recipeScores, score)
			}
		}
		numUniqueMatches := len(matches)

		score := AverageScore(recipeScores)
		score.Total += numMatches * numUniqueMatches

		scores = append(scores, score)
	}
	sort.Slice(indices, func(i, j int) bool {
		return scores[indices[i]].Total > scores[indices[j]].Total
	})
	res := Recipes{}
	upper := (data.N + 1) * 100
	lower := data.N * 100
	if upper >= len(indices) {
		upper = len(indices) - 1
	}
	if lower >= len(indices) {
		lower = len(indices) - 1
	}
	for _, i := range indices[lower:upper] {
		recipe := allRecipes[i]
		recipe.Score = scores[i]
		res = append(res, recipe)
	}
	return res
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

// GetWords gets all of the words contained within the given text
func GetWords(text string) []string {
	res := []string{}
	curStr := ""
	for _, r := range text {
		if r == ' ' || r == ',' || r == '.' || r == '(' || r == ')' || r == '+' || r == '–' || r == '—' {
			if curStr != "" && curStr != "and" {
				res = append(res, curStr)
				curStr = ""
			}
			continue
		}
		curStr = string(append([]rune(curStr), r))
	}
	if curStr != "" {
		res = append(res, curStr)
	}
	return res
}
