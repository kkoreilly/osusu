package main

import (
	"encoding/json"
	"html"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// A Recipe is an external recipe that can be used for new meal recommendations
type Recipe struct {
	Name         string
	URL          string
	Description  string
	Image        string
	Category     []string
	Cuisine      []string
	Ingredients  []string
	TotalTime    string
	PrepTime     string
	CookTime     string
	RatingValue  float64 // out of 5
	RatingCount  int
	RatingScore  int
	RatingWeight int
	Nutrition    Nutrition
	Source       string
	Score        Score
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

// AtoiIgnore parses the given int string into an int, ignoring all non-numeric characters in the string
func AtoiIgnore(nutrition string) (int, error) {
	runes := []rune{}
	for _, r := range nutrition {
		if unicode.IsDigit(r) {
			runes = append(runes, r)
		}
	}
	return strconv.Atoi(string(runes))
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
	b, err := os.ReadFile("web/recipes.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &recipes)
	if err != nil {
		return nil, err
	}
	for i, recipe := range recipes {
		if recipe.RatingValue == "" {
			continue
		}
		ratingValue, err := strconv.ParseFloat(recipe.RatingValue, 64)
		if err != nil {
			return nil, err
		}
		ratingCount, err := strconv.Atoi(recipe.RatingCount)
		if err != nil {
			return nil, err
		}
		recipe.RatingScore = int(100 * ratingValue / 5)
		recipe.RatingWeight = ratingCount
		if recipe.RatingWeight < 50 {
			recipe.RatingWeight = 50
		}
		if recipe.RatingWeight > 150 {
			recipe.RatingWeight = 150
		}
		// unescape strings because they were escaped to encode in json
		recipe.Name = html.UnescapeString(recipe.Name)
		recipe.Description = html.UnescapeString(recipe.Description)
		// recipe.Ingredients = html.UnescapeString(recipe.Ingredients)
		recipes[i] = recipe
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

// sourceAccuracy is a map with rough accuracy estimates for every source
var sourceAccuracy = map[string]int{
	"bbcfood":              50,
	"elanaspantry":         100,
	"lovefood":             100,
	"delishhh":             100,
	"thevintagemixer":      100,
	"backtoherroots":       100,
	"cookieandkate":        0,
	"jamieoliver":          80,
	"paninihappy":          100,
	"bunkycooks":           0,
	"steamykitchen":        100,
	"chow":                 0,
	"seriouseats":          90,
	"thelittlekitchen":     100,
	"williamssonoma":       0,
	"whatsgabycooking":     60,
	"cookincanuck":         100,
	"eatthelove":           100,
	"naturallyella":        0,
	"aspicyperspective":    0,
	"food":                 0,
	"pickypalate":          100,
	"thepioneerwoman":      100,
	"foodnetwork":          100,
	"epicurious":           0,
	"tastykitchen":         100,
	"biggirlssmallkitchen": 70,
	"bonappetit":           80,
	"allrecipes":           30,
	"browneyedbaker":       90,
	"101cookbooks":         100,
	"bbcgoodfood":          90,
	"smittenkitchen":       100,
}

// EstimateValid estimates what number of the given recipes are valid using the source accuracy map
func EstimateValid(recipes Recipes) int {
	sum := 0
	for _, recipe := range recipes {
		accuracy := sourceAccuracy[recipe.Source]
		sum += accuracy
	}
	return sum / 100
}

// RemoveInvalidRecipes returns the given recipes with all recipes that return 404s removed
func RemoveInvalidRecipes(recipes Recipes) Recipes {
	// amount valid and invalid per source
	invalid := map[string]int{}
	valid := map[string]int{}
	total := map[string]int{}
	res := Recipes{}
	// t := time.Now()
	for i, recipe := range recipes {
		if total[recipe.Source] >= 10 {
			continue
		}
		log.Println(total[recipe.Source], recipe.Source)
		resp, err := http.Get(recipe.URL)
		if err != nil {
			log.Println("error fetching recipe in remove invalid recipes:", err)
			invalid[recipe.Source]++
			total[recipe.Source]++
			continue
		}
		if resp.StatusCode != 200 {
			// log.Println("bad status code:", resp.StatusCode, "with source", recipe.Source)
			invalid[recipe.Source]++
			total[recipe.Source]++
			continue
		}
		// log.Println("valid recipe found with source", recipe.Source)
		valid[recipe.Source]++
		total[recipe.Source]++
		res = append(res, recipe)
		if i != 0 {
			// log.Println("percent valid:", strconv.Itoa(100*len(res)/i)+"%", "total:", i, "valid:", len(res), "invalid:", i-len(res))
			// log.Println("total time:", time.Since(t), "time per:", time.Since(t)/time.Duration(i), "estimated total:", time.Duration(len(recipes))*time.Since(t)/time.Duration((i)))
		}
	}
	for source, val := range total {
		nInvalid, nValid := invalid[source], valid[source]
		log.Println("Keep:", 100*nValid/val > 70, "Percent Valid:", strconv.Itoa(100*nValid/val)+"%", "Source:", source, "Number Invalid:", nInvalid, "Number Valid:", nValid, "Total Number:", val)
	}
	return res
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
		score.Total = (100*score.Total + recipe.RatingWeight*recipe.RatingScore) / (100 + recipe.RatingWeight)

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
