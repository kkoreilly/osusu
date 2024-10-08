package osusu

import (
	"fmt"
	"html"
	"strings"
	"time"

	"gonum.org/v1/gonum/stat"
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
	Category          []string   `display:"-"`
	CategoryFlag      Categories `json:"-" label:"Category"`
	Cuisine           []string   `display:"-"`
	CuisineFlag       Cuisines   `json:"-" label:"Cuisine"`
	Ingredients       []string
	TotalTime         string        `display:"-"`
	PrepTime          string        `display:"-"`
	CookTime          string        `display:"-"`
	TotalTimeDuration time.Duration `json:"-" label:"Total time" viewif:"TotalTime!=\"\""`
	PrepTimeDuration  time.Duration `json:"-" label:"Prep time" viewif:"PrepTime!=\"\""`
	CookTimeDuration  time.Duration `json:"-" label:"Cook time" viewif:"CookTime!=\"\""`
	Yield             int
	RatingValue       float64 `display:"slider" min:"0" max:"5"`
	RatingCount       int
	RatingScore       int `display:"-" json:"-"`
	RatingWeight      int `display:"-" json:"-"`
	Nutrition         Nutrition
	Source            string `json:"-"`
	// index score values for base information about a recipe (using info like calories, time, ingredients, etc)
	BaseScoreIndex Score `json:"-"`
	// percentile values of BaseScoreIndex
	BaseScore Score
	// keyed by meal ID
	TextEncodingScores map[uint]float32 `json:"-"`
	// index score values for text encoding based scores
	EncodingScoreIndex Score `json:"-"`
	// percentile values of EncodingScoreIndex
	EncodingScore Score `json:"-"`
	Score         Score `json:"-"`
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

// ComputeBaseScoreIndex computes and sets the base score index for the given recipe
func (r *Recipe) ComputeBaseScoreIndex() {
	r.BaseScoreIndex = Score{}
	// just length of ingredients, obviously can be improved to actually look at ingredients, but in general cost will increase with number of ingredients, higher = more expensive = worse
	r.BaseScoreIndex.Cost = -len(r.Ingredients)
	// use generic total time duration of one hour if it isn't defined
	if r.TotalTimeDuration == 0 {
		r.TotalTimeDuration = time.Hour
	}
	// use combination of number of ingredients and total time, higher = more effort = worse
	r.BaseScoreIndex.Effort = -len(r.Ingredients) - int(r.TotalTimeDuration.Minutes())
	// avoid div by 0
	if r.Nutrition.Protein == 0 {
		r.BaseScoreIndex.Healthiness = -r.Nutrition.Sugar * 10
	} else {
		// ratio of sugar to protein, higher = more sugar = worse
		r.BaseScoreIndex.Healthiness = -100 * r.Nutrition.Sugar / r.Nutrition.Protein
	}
	// rating value combined with rating count, higher = better rated = better
	r.BaseScoreIndex.Taste = int(100*r.RatingValue) + min(r.RatingCount, 500)
	// hours since 1970 for date published and modified, higher = more recent = better
	r.BaseScoreIndex.Recency = int(r.DatePublished.Unix()/3600 + r.DateModified.Unix()/3600)
}

// ComputeNormScores computes the normalized base and
// text encoding scores for each recipe. The base and text
// encoding score indices already need to be computed.
// The normalized scores are the z-scores centered at 50
// and bounded at roughly 0 and 100.
func ComputeNormScores(recipes []*Recipe) {
	len := len(recipes)
	// we sort recipes by the base score indices on each metric and then loop over to find the percentile for each recipe on each metric and use that for the base score
	doCompute := func(scoreValue func(s *Score) *int, indexScoreObject, percentileScoreObject func(r *Recipe) *Score) {
		scores := make([]float64, len)
		for i, recipe := range recipes {
			scores[i] = float64(*scoreValue(indexScoreObject(recipe)))
		}
		mean, std := stat.MeanStdDev(scores, nil)
		for i, recipe := range recipes {
			z := (scores[i] - mean) / std
			if std == 0 {
				z = 0
			}
			sc := (50 * z / 3) + 50
			*scoreValue(percentileScoreObject(recipe)) = int(sc)
		}
	}
	compute := func(scoreValue func(s *Score) *int) {
		doCompute(scoreValue,
			func(r *Recipe) *Score { return &r.BaseScoreIndex },
			func(r *Recipe) *Score { return &r.BaseScore },
		)
		doCompute(scoreValue,
			func(r *Recipe) *Score { return &r.EncodingScoreIndex },
			func(r *Recipe) *Score { return &r.EncodingScore },
		)
	}
	compute(func(s *Score) *int { return &s.Taste })
	compute(func(s *Score) *int { return &s.Recency })
	compute(func(s *Score) *int { return &s.Cost })
	compute(func(s *Score) *int { return &s.Effort })
	compute(func(s *Score) *int { return &s.Healthiness })
}
