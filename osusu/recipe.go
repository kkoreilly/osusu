package osusu

import (
	"fmt"
	"html"
	"strings"
	"time"
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
	Category          []string   `view:"-"`
	CategoryFlag      Categories `json:"-" label:"Category"`
	Cuisine           []string   `view:"-"`
	CuisineFlag       Cuisines   `json:"-" label:"Cuisine"`
	Ingredients       []string
	TotalTime         string        `view:"-"`
	PrepTime          string        `view:"-"`
	CookTime          string        `view:"-"`
	TotalTimeDuration time.Duration `json:"-" label:"Total time" viewif:"TotalTime!=\"\""`
	PrepTimeDuration  time.Duration `json:"-" label:"Prep time" viewif:"PrepTime!=\"\""`
	CookTimeDuration  time.Duration `json:"-" label:"Cook time" viewif:"CookTime!=\"\""`
	Yield             int
	RatingValue       float64 `view:"slider" min:"0" max:"5"`
	RatingCount       int
	RatingScore       int `view:"-" json:"-"`
	RatingWeight      int `view:"-" json:"-"`
	Nutrition         Nutrition
	Source            string `json:"-"`
	BaseScoreIndex    Score  `view:"-" json:"-"` // index score values for base information about a recipe (using info like calories, time, ingredients, etc)
	BaseScore         Score  `view:"-"`          // percentile values of BaseScoreIndex
	Score             Score  `view:"-"`
}

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

// Init initializes computed values in the recipe after it has been loaded.
func (r *Recipe) Init() error {
	var err error
	// load durations
	if r.TotalTime != "" {
		r.TotalTimeDuration, err = time.ParseDuration(r.TotalTime)
		if err != nil {
			return fmt.Errorf("error loading total time duration: %w", err)
		}
	}
	if r.CookTime != "" {
		r.CookTimeDuration, err = time.ParseDuration(r.CookTime)
		if err != nil {
			return fmt.Errorf("error loading cook time duration: %w", err)
		}
	}
	if r.PrepTime != "" {
		r.PrepTimeDuration, err = time.ParseDuration(r.PrepTime)
		if err != nil {
			return fmt.Errorf("error loading prep time duration: %w", err)
		}
	}
	r.RatingScore = int(100 * r.RatingValue / 5)
	r.RatingWeight = r.RatingCount
	if r.RatingWeight < 50 {
		r.RatingWeight = 50
	}
	if r.RatingWeight > 150 {
		r.RatingWeight = 150
	}
	// unescape strings because they were escaped to encode in json
	r.Name = html.UnescapeString(r.Name)
	r.Description = html.UnescapeString(r.Description)
	for i, ingredient := range r.Ingredients {
		ingredient = html.UnescapeString(ingredient)
		ingredient = strings.ReplaceAll(ingredient, "0.33333334326744", "1/3")
		ingredient = strings.ReplaceAll(ingredient, "0.66666668653488", "1/6")
		r.Ingredients[i] = ingredient
	}
	return nil
}

// Text returns all of the text associated with the recipe as one string.
// It is intended to be used as text encoding model data, so it should
// not be presented to end-users.
func (r *Recipe) Text() string {
	return strings.Join([]string{r.Name, r.Description, strings.Join(r.Ingredients, "\n")}, "\n")
}
