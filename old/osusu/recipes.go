package osusu

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kkoreilly/osusu/util/file"
	"github.com/kkoreilly/osusu/util/mat"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// A Recipe is an external recipe that can be used for new meal recommendations
type Recipe struct {
	Name              string
	URL               string
	Description       string
	Image             string
	Author            string
	DatePublished     time.Time
	DateModified      time.Time
	Category          []string
	Cuisine           []string
	Ingredients       []string
	TotalTime         string
	PrepTime          string
	CookTime          string
	TotalTimeDuration time.Duration `json:"-"`
	PrepTimeDuration  time.Duration `json:"-"`
	CookTimeDuration  time.Duration `json:"-"`
	Yield             int
	RatingValue       float64 // out of 5
	RatingCount       int
	RatingScore       int `json:"-"`
	RatingWeight      int `json:"-"`
	Nutrition         Nutrition
	Source            string `json:"-"`
	BaseScoreIndex    Score  `json:"-"` // index score values for base information about a recipe (using info like calories, time, ingredients, etc)
	BaseScore         Score  // percentile values of BaseScoreIndex
	Score             Score
}

// Recipes is a slice of multiple recipes
type Recipes []Recipe

// Nutrition represents the nutritional information of a recipe
type Nutrition struct {
	Calories       int // unit: Calories (kcal)
	Carbohydrate   int // g
	Cholesterol    int // mg
	Fiber          int // g
	Protein        int // g
	Fat            int // g
	SaturatedFat   int // g
	UnsaturatedFat int // g
	Sodium         int // mg
	Sugar          int // g
}

// CurrentRecipe returns the current recipe value from local storage
func CurrentRecipe(ctx app.Context) Recipe {
	var recipe Recipe
	ctx.GetState("currentRecipe", &recipe)
	return recipe
}

// SetCurrentRecipe sets the current recipe value in local storage
func SetCurrentRecipe(recipe Recipe, ctx app.Context) {
	ctx.SetState("currentRecipe", recipe, app.Persist)
}

// LoadRecipes gets all of the recipes from the given file
func LoadRecipes(path string) (Recipes, error) {
	recipes := Recipes{}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &recipes)
	if err != nil {
		return nil, err
	}
	for i, recipe := range recipes {
		// load durations
		if recipe.TotalTime != "" {
			recipe.TotalTimeDuration, err = time.ParseDuration(recipe.TotalTime)
			if err != nil {
				log.Println(fmt.Errorf("error loading total time duration: %w", err))
			}
		}
		if recipe.CookTime != "" {
			recipe.CookTimeDuration, err = time.ParseDuration(recipe.CookTime)
			if err != nil {
				log.Println(fmt.Errorf("error loading cook time duration: %w", err))
			}
		}
		if recipe.PrepTime != "" {
			recipe.PrepTimeDuration, err = time.ParseDuration(recipe.PrepTime)
			if err != nil {
				log.Println(fmt.Errorf("error loading prep time duration: %w", err))
			}
		}
		recipe.RatingScore = int(100 * recipe.RatingValue / 5)
		recipe.RatingWeight = recipe.RatingCount
		if recipe.RatingWeight < 50 {
			recipe.RatingWeight = 50
		}
		if recipe.RatingWeight > 150 {
			recipe.RatingWeight = 150
		}
		// unescape strings because they were escaped to encode in json
		recipe.Name = html.UnescapeString(recipe.Name)
		recipe.Description = html.UnescapeString(recipe.Description)
		for i, ingredient := range recipe.Ingredients {
			ingredient = html.UnescapeString(ingredient)
			ingredient = strings.ReplaceAll(ingredient, "0.33333334326744", "1/3")
			ingredient = strings.ReplaceAll(ingredient, "0.66666668653488", "1/6")
			recipe.Ingredients[i] = ingredient
		}
		recipes[i] = recipe
	}
	return recipes, nil
}

// SaveRecipes writes the given recipes to the given file after safely moving any file currently there to web/data/recipes.json.old
func SaveRecipes(recipes Recipes, path string) error {
	if file.Exists(path) {
		file.RenameSafe(path, "web/data/recipes.json.old")
	}
	jsonData, err := json.Marshal(recipes)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, jsonData, 0666)
	return err
}

// ComputeBaseScores returns the recipes with the base score for each recipe computed
func (r Recipes) ComputeBaseScores() Recipes {
	r = r.ComputeBaseScoreIndices()
	len := len(r)
	// we sort recipes by the base score indices on each metric and then loop over to find the percentile for each recipe on each metric and use that for the base score
	sort.Slice(r, func(i, j int) bool {
		// higher is worse, so put greater values first
		return r[i].BaseScoreIndex.Cost > r[j].BaseScoreIndex.Cost
	})
	for i := range r {
		r[i].BaseScore.Cost = Percentile(i, len)
	}
	sort.Slice(r, func(i, j int) bool {
		// higher is worse, so put greater values first
		return r[i].BaseScoreIndex.Effort > r[j].BaseScoreIndex.Effort
	})
	for i := range r {
		r[i].BaseScore.Effort = Percentile(i, len)
	}
	sort.Slice(r, func(i, j int) bool {
		// higher is worse, so put greater values first
		return r[i].BaseScoreIndex.Healthiness > r[j].BaseScoreIndex.Healthiness
	})
	for i := range r {
		r[i].BaseScore.Healthiness = Percentile(i, len)
	}
	sort.Slice(r, func(i, j int) bool {
		// higher is better, so put greater values last
		return r[i].BaseScoreIndex.Taste < r[j].BaseScoreIndex.Taste
	})
	for i := range r {
		r[i].BaseScore.Taste = Percentile(i, len)
	}
	sort.Slice(r, func(i, j int) bool {
		// higher is better, so put greater values last
		return r[i].BaseScoreIndex.Recency < r[j].BaseScoreIndex.Recency
	})
	for i := range r {
		r[i].BaseScore.Recency = Percentile(i, len)
	}
	return r
}

// Percentile returns the percentile of the element at the given index position in a sorted slice of the given length, normalized to range between 0 and 100
func Percentile(index, length int) int {
	return (100*index + length/2) / length
}

// ComputeBaseScoreIndices returns the recipes with the base score index for each recipe computed
func (r Recipes) ComputeBaseScoreIndices() Recipes {
	for i, recipe := range r {
		recipe.BaseScoreIndex = recipe.ComputeBaseScoreIndex()
		r[i] = recipe
	}
	return r
}

// ComputeBaseScoreIndex computes and returns the base score index for the given recipe
func (r Recipe) ComputeBaseScoreIndex() Score {
	score := Score{}
	score.Cost = len(r.Ingredients) // just length of ingredients, obviously can be improved to actually look at ingredients, but in general cost will increase with number of ingredients, higher = more expensive = worse
	// use generic total time duration of one hour if it isn't defined
	if r.TotalTimeDuration == 0 {
		r.TotalTimeDuration = time.Hour
	}
	score.Effort = len(r.Ingredients) + int(r.TotalTimeDuration.Minutes()) // use combination of number of ingredients and total time, higher = more effort = worse
	// avoid div by 0
	if r.Nutrition.Protein == 0 {
		score.Healthiness = r.Nutrition.Sugar * 10
	} else {
		score.Healthiness = 100 * r.Nutrition.Sugar / r.Nutrition.Protein // ratio of sugar to protein, higher = more sugar = worse
	}
	score.Taste = int(100*r.RatingValue) + mat.Min(r.RatingCount, 500)            // rating value combined with rating count, higher = better rated = better
	score.Recency = int(r.DatePublished.Unix()/3600 + r.DateModified.Unix()/3600) // hours since 1970 for date published and modified, higher = more recent = better
	return score
}

// GetWords returns all of the words contained within the name, description, and ingredients of the recipe
func (r Recipe) GetWords() []string {
	words := GetWords(r.Name)
	words = append(words, GetWords(r.Description)...)
	for _, ingredient := range r.Ingredients {
		words = append(words, GetWords(ingredient)...)
	}
	return words
}

// CountCategories returns how many of each category there are in the given recipes and how many recipes have any category
func (r Recipes) CountCategories() (map[string]int, int) {
	res := map[string]int{}
	total := 0
	for _, recipe := range r {
		if len(recipe.Category) != 0 {
			total++
		}
		for _, category := range recipe.Category {
			res[category]++
		}
	}
	return res, total
}

// CountCuisines returns how many of each cuisine there are in the given recipes and how many recipes have any cuisine
func (r Recipes) CountCuisines() (map[string]int, int) {
	res := map[string]int{}
	total := 0
	for _, recipe := range r {
		if len(recipe.Cuisine) != 0 {
			total++
		}
		for _, cuisine := range recipe.Cuisine {
			res[cuisine]++
		}
	}
	return res, total
}

// ConsolidateCategories consolidates the categories of the given recipes into a more concise set
func (r Recipes) ConsolidateCategories() Recipes {
	for i, recipe := range r {
		// need unique categories so use map to prevent duplicates
		categories := map[string]bool{}
		for _, category := range recipe.Category {
			for k, v := range CategoryToCategoryMap {
				for _, mapCategory := range v {
					if mapCategory == category {
						categories[k] = true
					}
				}
			}
		}
		newCategory := []string{}
		for category := range categories {
			newCategory = append(newCategory, category)
		}
		r[i].Category = newCategory
	}
	return r
}

// ConsolidateCuisines consolidates the cuisines of the given recipes into a more concise set
func (r Recipes) ConsolidateCuisines() Recipes {
	// convert into map for easier and quicker access
	ignoredCuisinesMap := map[string]bool{}
	for _, cuisine := range IgnoredCuisines {
		ignoredCuisinesMap[cuisine] = true
	}
	for i, recipe := range r {
		// need unique cuisines so use map to prevent duplicates
		cuisines := map[string]bool{}
		for _, cuisine := range recipe.Cuisine {
			got := false
			for k, v := range CuisineToCuisineMap {
				for _, mapCuisine := range v {
					if mapCuisine == cuisine {
						cuisines[k] = true
						got = true
					}
				}
			}
			if !got && !ignoredCuisinesMap[cuisine] {
				log.Println("uncaught", cuisine)
			}
		}
		for mainCuisine, subCuisines := range CuisineToCuisineMap {
			for _, cuisine := range subCuisines {
				cuisineLower := strings.ToLower(cuisine)
				if strings.Contains(recipe.Name, cuisine) || strings.Contains(recipe.Name, cuisineLower) || strings.Contains(recipe.Description, cuisine) || strings.Contains(recipe.Description, cuisineLower) {
					cuisines[mainCuisine] = true
				}
				for _, ingredient := range recipe.Ingredients {
					if strings.Contains(ingredient, cuisine) || strings.Contains(ingredient, cuisineLower) {
						cuisines[mainCuisine] = true
					}
				}
			}
		}
		newCuisine := []string{}
		for cuisine := range cuisines {
			newCuisine = append(newCuisine, cuisine)
		}
		r[i].Cuisine = newCuisine
	}
	return r
}

// InferCuisines infers the cuisines of the recipes in the given recipes that don't have a cuisine set.
// It uses the values of the recipes with cuisines already set to do this.
func (r Recipes) InferCuisines(numRecipesPerCuisine map[string]int) Recipes {
	wordCuisineMap := map[string](map[string]int){} // map[word](map[cuisine]numTimes){}
	cuisineNumWords := map[string]int{}
	// get word cuisine map
	for _, recipe := range r {
		// won't add to map if no cuisine
		if len(recipe.Cuisine) == 0 {
			continue
		}
		words := GetWords(recipe.Name)
		for _, word := range words {
			if wordCuisineMap[word] == nil {
				wordCuisineMap[word] = map[string]int{}
			}
			for _, cuisine := range recipe.Cuisine {
				wordCuisineMap[word][cuisine]++
				cuisineNumWords[cuisine]++
			}
		}
	}
	// log.Println(wordCuisineMap)
	// use it to infer cuisines
	for i, recipe := range r {
		// not needed if we already have cuisine
		if len(recipe.Cuisine) != 0 {
			continue
		}
		cuisineMap := map[string]int{}
		words := append(GetWords(recipe.Name), GetWords(recipe.Description)...)
		for _, ingredient := range recipe.Ingredients {
			words = append(words, GetWords(ingredient)...)
		}
		for _, word := range words {
			sum := 0
			for _, value := range wordCuisineMap[word] {
				sum += value
			}
			for cuisine, value := range wordCuisineMap[word] {
				cuisineMap[cuisine] += 1000 * value / sum
			}
		}
		highestCuisine := ""
		highestValue := 0.0
		for cuisine, value := range cuisineMap {
			weightedValue := 1000000 * float64(value)
			if weightedValue > highestValue {
				highestCuisine = cuisine
				highestValue = weightedValue
			}
		}

		recipe.Cuisine = []string{highestCuisine}
		r[i] = recipe
	}
	return r
}

// RecipeNumberChanges logs information about the changes from the given old recipe counts to the given new recipe counts
func RecipeNumberChanges(oldMap map[string]int, oldCount int, newMap map[string]int, newCount int) {
	diff := []string{}
	for name := range newMap {
		diff = append(diff, name)
	}
	sort.Slice(diff, func(i, j int) bool {
		return 100*newMap[diff[i]]/oldMap[diff[i]] < 100*newMap[diff[j]]/oldMap[diff[j]]
	})
	for _, name := range diff {
		difference := newMap[name] - oldMap[name]
		log.Println(name, "Difference:", difference)
		if oldMap[name] == 0 {
			continue
		}
		percent := 100 * newMap[name] / oldMap[name]
		log.Println(name, "Percent Difference:", strconv.Itoa(percent)+"%")
	}
	fmt.Println("") // get line space
	log.Println("Total Difference:", newCount-oldCount)
	totalPercentDiff := 100 * newCount / oldCount
	log.Println("Total Percent Difference:", strconv.Itoa(totalPercentDiff)+"%")
	log.Println("Median Difference:", newMap[diff[len(diff)/2]]-oldMap[diff[len(diff)/2]])
	sort.Slice(diff, func(i, j int) bool {
		return 100*newMap[diff[i]]/oldMap[diff[i]] < 100*newMap[diff[j]]/oldMap[diff[j]]
	})
	medianPercentDiff := 100 * newMap[diff[len(diff)/2]] / oldMap[diff[len(diff)/2]]
	log.Println("Median Percent Difference:", strconv.Itoa(medianPercentDiff)+"%")
	log.Println("Error:", strconv.Itoa(int(math.Abs(float64(medianPercentDiff-totalPercentDiff))))+"%")
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

// FormatDuration converts an ISO 8601 duration to a Go-formatted duration
func FormatDuration(duration string) string {
	// remove PT at start
	if len(duration) > 2 {
		duration = duration[2:]
	}
	// change H, M, etc to h, m, etc
	duration = strings.ToLower(duration)
	return duration
}

// ParseDuration parses an ISO 8601 duration (the format used in the recipes source)
func ParseDuration(duration string) (time.Duration, error) {
	duration = FormatDuration(duration)
	return time.ParseDuration(duration)
}

// GetWords gets all of the words contained within the given text
func GetWords(text string) []string {
	res := []string{}
	curStr := ""
	for _, r := range text {
		if WordSeparatorsMap[r] {
			if curStr != "" && !IgnoredWordsMap[curStr] {
				res = append(res, curStr)
			}
			curStr = ""
			continue
		}
		curStr = string(append([]rune(curStr), r))
	}
	if curStr != "" && !IgnoredWordsMap[curStr] {
		res = append(res, curStr)
	}
	return res
}
