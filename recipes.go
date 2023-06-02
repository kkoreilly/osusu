package main

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

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"golang.org/x/exp/constraints"
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

// LoadRecipes gets all of the recipes from the recipes.json file
func LoadRecipes() (Recipes, error) {
	recipes := Recipes{}
	b, err := os.ReadFile("web/data/recipes.json")
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
			recipe.Ingredients[i] = html.UnescapeString(ingredient)
		}
		recipes[i] = recipe
	}
	return recipes, nil
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
	score.Taste = int(100*r.RatingValue) + Min(r.RatingCount, 500)                // rating value combined with rating count, higher = better rated = better
	score.Recency = int(r.DatePublished.Unix()/3600 + r.DateModified.Unix()/3600) // hours since 1970 for date published and modified, higher = more recent = better
	return score
}

// Min returns the minimum value of the given two values
func Min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// Max returns the maximum value of the given two values
func Max[T constraints.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
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
	categoryMap := map[string][]string{
		"Dinner":     {"Dinner", "Entree", "Pasta"},
		"Drink":      {"Drink", "Beverage", "Cocktail", "Coffee"},
		"Dessert":    {"Dessert", "Cake", "Candy", "Pie"},
		"Lunch":      {"Lunch", "Sandwich"},
		"Ingredient": {"Ingredient", "Bread", "Condiment", "Jam / Jelly", "Sauce", "Spice Mix"},
		"Appetizer":  {"Appetizer", "Salad", "Soup"},
		"Side":       {"Side Dish"},
		"Breakfast":  {"Breakfast"},
		"Snack":      {"Snack"},
		"Brunch":     {"Brunch"},
	}
	for i, recipe := range r {
		// need unique categories so use map to prevent duplicates
		categories := map[string]bool{}
		for _, category := range recipe.Category {
			for k, v := range categoryMap {
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
	for _, cuisine := range ignoredCuisines {
		ignoredCuisinesMap[cuisine] = true
	}
	for i, recipe := range r {
		// need unique cuisines so use map to prevent duplicates
		cuisines := map[string]bool{}
		// if len(recipe.Cuisine) == 0 {
		// 	log.Println(recipe.Name, recipe.Description)
		// }
		for _, cuisine := range recipe.Cuisine {
			got := false
			for k, v := range cuisineToCuisineMap {
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
		for mainCuisine, subCuisines := range cuisineToCuisineMap {
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
		// for cuisine, words := range cuisineWords {
		// 	for _, word := range words {
		// 		if strings.Contains(recipe.Name, word) || strings.Contains(recipe.Description, word) {
		// 			cuisines[cuisine] = true
		// 		}
		// 		for _, ingredient := range recipe.Ingredients {
		// 			if strings.Contains(ingredient, word) {
		// 				cuisines[cuisine] = true
		// 			}
		// 		}
		// 	}
		// }
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
		return newMap[diff[i]]-oldMap[diff[i]] < newMap[diff[j]]-oldMap[diff[j]]
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

// // FixRecipeTimes returns the given recipes with durations formatted correctly
// func FixRecipeTimes(recipes Recipes) Recipes {
// 	for i, recipe := range recipes {
// 		var prepDuration, cookDuration, totalDuration time.Duration
// 		if recipe.PrepTime != "" {
// 			prepDuration, _ = ParseDuration(recipe.PrepTime)
// 			recipe.PrepTime = prepDuration.String()
// 		}
// 		if recipe.CookTime != "" {
// 			cookDuration, _ = ParseDuration(recipe.CookTime)
// 			recipe.CookTime = cookDuration.String()
// 		}
// 		if recipe.TotalTime != "" {
// 			totalDuration, _ = ParseDuration(recipe.TotalTime)
// 			recipe.TotalTime = totalDuration.String()
// 		} else {
// 			recipe.TotalTime = (prepDuration + cookDuration).String()
// 		}
// 		recipes[i] = recipe
// 	}
// 	return recipes
// }

// // sourceAccuracy is a map with rough accuracy estimates for every source
// var sourceAccuracy = map[string]int{
// 	"bbcfood":              50,
// 	"elanaspantry":         100,
// 	"lovefood":             100,
// 	"delishhh":             100,
// 	"thevintagemixer":      100,
// 	"backtoherroots":       100,
// 	"cookieandkate":        0,
// 	"jamieoliver":          80,
// 	"paninihappy":          100,
// 	"bunkycooks":           0,
// 	"steamykitchen":        100,
// 	"chow":                 0,
// 	"seriouseats":          90,
// 	"thelittlekitchen":     100,
// 	"williamssonoma":       0,
// 	"whatsgabycooking":     60,
// 	"cookincanuck":         100,
// 	"eatthelove":           100,
// 	"naturallyella":        0,
// 	"aspicyperspective":    0,
// 	"food":                 0,
// 	"pickypalate":          100,
// 	"thepioneerwoman":      100,
// 	"foodnetwork":          100,
// 	"epicurious":           0,
// 	"tastykitchen":         100,
// 	"biggirlssmallkitchen": 70,
// 	"bonappetit":           80,
// 	"allrecipes":           30,
// 	"browneyedbaker":       90,
// 	"101cookbooks":         100,
// 	"bbcgoodfood":          90,
// 	"smittenkitchen":       100,
// }

// // EstimateValid estimates what number of the given recipes are valid using the source accuracy map
// func EstimateValid(recipes Recipes) int {
// 	sum := 0
// 	for _, recipe := range recipes {
// 		accuracy := sourceAccuracy[recipe.Source]
// 		sum += accuracy
// 	}
// 	return sum / 100
// }

// // RemoveInvalidRecipes returns the given recipes with all recipes that return 404s removed
// func RemoveInvalidRecipes(recipes Recipes) Recipes {
// 	// amount valid and invalid per source
// 	invalid := map[string]int{}
// 	valid := map[string]int{}
// 	total := map[string]int{}
// 	res := Recipes{}
// 	// t := time.Now()
// 	for i, recipe := range recipes {
// 		if total[recipe.Source] >= 10 {
// 			continue
// 		}
// 		log.Println(total[recipe.Source], recipe.Source)
// 		resp, err := http.Get(recipe.URL)
// 		if err != nil {
// 			log.Println("error fetching recipe in remove invalid recipes:", err)
// 			invalid[recipe.Source]++
// 			total[recipe.Source]++
// 			continue
// 		}
// 		if resp.StatusCode != 200 {
// 			// log.Println("bad status code:", resp.StatusCode, "with source", recipe.Source)
// 			invalid[recipe.Source]++
// 			total[recipe.Source]++
// 			continue
// 		}
// 		// log.Println("valid recipe found with source", recipe.Source)
// 		valid[recipe.Source]++
// 		total[recipe.Source]++
// 		res = append(res, recipe)
// 		if i != 0 {
// 			// log.Println("percent valid:", strconv.Itoa(100*len(res)/i)+"%", "total:", i, "valid:", len(res), "invalid:", i-len(res))
// 			// log.Println("total time:", time.Since(t), "time per:", time.Since(t)/time.Duration(i), "estimated total:", time.Duration(len(recipes))*time.Since(t)/time.Duration((i)))
// 		}
// 	}
// 	for source, val := range total {
// 		nInvalid, nValid := invalid[source], valid[source]
// 		log.Println("Keep:", 100*nValid/val > 70, "Percent Valid:", strconv.Itoa(100*nValid/val)+"%", "Source:", source, "Number Invalid:", nInvalid, "Number Valid:", nValid, "Total Number:", val)
// 	}
// 	return res
// }

// SaveRecipes writes the given recipes to the web/data/newrecipes.json file
func SaveRecipes(recipes Recipes) error {
	jsonData, err := json.Marshal(recipes)
	if err != nil {
		return err
	}
	err = os.WriteFile("web/data/newrecipes.json", jsonData, 0666)
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
	UsedSources  map[string]bool
	N            int // the iteration that we are on (ie: n = 3 for the fourth time we are getting recipes for the same options), we return 100 new meals each time
}

// RecommendRecipes returns a list of recommended new recipes based on the given recommend recipes data
func RecommendRecipes(data RecommendRecipesData) Recipes {
	if len(allRecipes) == 0 {
		return Recipes{}
	}
	scores := []Score{}
	indices := []int{}
	numSkipped := 0
	recipes := []Recipe{}
	for i, recipe := range allRecipes {
		// if we have already used this recipe URL (and by extension this recipe), skip to prevent duplicates
		if data.UsedSources[recipe.URL] {
			numSkipped++
			continue
		}
		// check that at least one category satisfies at least one type option
		gotCategory := false
		if !gotCategory {
			for _, recipeCategory := range recipe.Category {
				if data.Options.Category[recipeCategory] {
					gotCategory = true
					break
				}
			}
			if !gotCategory {
				numSkipped++
				continue
			}
		}

		// check that at least one cuisine satisfies at least one cuisine option
		gotCuisine := false
		for _, recipeCuisine := range recipe.Cuisine {
			for optionsCuisine, val := range data.Options.Cuisine {
				if val && recipeCuisine == optionsCuisine {
					gotCuisine = true
					break
				}
			}
			if gotCuisine {
				break
			}
		}
		if !gotCuisine {
			numSkipped++
			continue
		}
		// need to subtract num skipped so that i stays in line with the number of items in the indices slice
		i -= numSkipped
		// only append to indices and recipes if it matches conditions above
		indices = append(indices, i)
		recipes = append(recipes, recipe)
		words := GetWords(recipe.Name)
		recipeScores := []ScoreWeight{}
		// use map to track unique matches and simple int to count total
		matches := map[string]bool{}
		numMatches := 0
		for j, word := range words {
			score, ok := data.WordScoreMap[word]
			if ok {
				matches[word] = true
				numMatches++
				// earlier words more important
				recipeScores = append(recipeScores, ScoreWeight{Score: score, Weight: len(words) - j})
			}
		}
		// numUniqueMatches := len(matches)

		score := AverageScore(recipeScores)
		// importance of user info is based on how much of it their is, capped at 200 words, which is the same as base info weight
		score = AverageScore([]ScoreWeight{{Score: score, Weight: Min(len(data.WordScoreMap), 200)}, {Score: recipe.BaseScore, Weight: 200}})
		score.Total = score.ComputeTotal(data.Options)

		scores = append(scores, score)
	}
	// if we got nothing, bail now to prevent slice bounds error later
	if len(indices) == 0 {
		return Recipes{}
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
		recipe := recipes[i]
		recipe.Score = scores[i]
		res = append(res, recipe)
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
	separators := []rune{' ', ',', '.', '(', ')', '+', '–', '—'}
	// use map for easier access
	separatorsMap := map[rune]bool{}
	for _, separator := range separators {
		separatorsMap[separator] = true
	}
	ignoredWords := []string{"a", "an", "the", "and", "or", "with", "to", "from", "about", "above", "across", "against", "along", "at", "but"}
	// use map for easier access
	ignoredWordsMap := map[string]bool{}
	for _, word := range ignoredWords {
		ignoredWordsMap[word] = true
	}
	for _, r := range text {
		if ignoredWordsMap[curStr] {
			curStr = ""
			continue
		}
		if separatorsMap[r] {
			if curStr != "" {
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
