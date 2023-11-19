package osusu

import "time"

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
	TotalTime         string        `view:"-"`
	PrepTime          string        `view:"-"`
	CookTime          string        `view:"-"`
	TotalTimeDuration time.Duration `json:"-" label:"Total time"`
	PrepTimeDuration  time.Duration `json:"-" label:"Prep time"`
	CookTimeDuration  time.Duration `json:"-" label:"Cook time"`
	Yield             int
	RatingValue       float64 `view:"slider" min:"0" max:"5"`
	RatingCount       int
	RatingScore       int `json:"-"`
	RatingWeight      int `json:"-"`
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
