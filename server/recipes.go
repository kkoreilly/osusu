package server

import (
	"sort"

	"github.com/kkoreilly/osusu/osusu"
	"github.com/kkoreilly/osusu/util/mat"
)

// RecommendRecipesData is the data used in a recommend recipes call
type RecommendRecipesData struct {
	WordScoreMap map[string]osusu.Score
	Options      osusu.Options
	UsedSources  map[string]bool
	N            int // the iteration that we are on (ie: n = 3 for the fourth time we are getting recipes for the same options), we return 100 new meals each time
}

// RecommendRecipes returns a list of recommended new recipes based on the given recommend recipes data
func RecommendRecipes(data RecommendRecipesData) osusu.Recipes {
	if len(AllRecipes) == 0 {
		return osusu.Recipes{}
	}
	scores := []osusu.Score{}
	indices := []int{}
	numSkipped := 0
	recipes := []osusu.Recipe{}
	for i, recipe := range AllRecipes {
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

		words := recipe.GetWords()

		hasExcludedIngredient := false
		for _, word := range words {
			if data.Options.ExcludedIngredients[word] {
				hasExcludedIngredient = true
				break
			}
		}
		if hasExcludedIngredient {
			numSkipped++
			continue
		}
		// need to subtract num skipped so that i stays in line with the number of items in the indices slice
		i -= numSkipped
		// only append to indices and recipes if it matches conditions above
		indices = append(indices, i)
		recipes = append(recipes, recipe)
		recipeScores := []osusu.ScoreWeight{}
		// use map to track unique matches and simple int to count total
		matches := map[string]bool{}
		numMatches := 0
		for j, word := range words {
			score, ok := data.WordScoreMap[word]
			if ok {
				matches[word] = true
				numMatches++
				// earlier words more important
				recipeScores = append(recipeScores, osusu.ScoreWeight{Score: score, Weight: len(words) - j})
			}
		}
		// numUniqueMatches := len(matches)

		score := osusu.AverageScore(recipeScores)
		// importance of user info is based on how much of it their is, capped at 200 words, which is the same as base info weight
		score = osusu.AverageScore([]osusu.ScoreWeight{{Score: score, Weight: mat.Min(len(data.WordScoreMap), 200)}, {Score: recipe.BaseScore, Weight: 200}})
		score.Total = score.ComputeTotal(data.Options)

		scores = append(scores, score)
	}
	// if we got nothing, bail now to prevent slice bounds error later
	if len(indices) == 0 {
		return osusu.Recipes{}
	}
	sort.Slice(indices, func(i, j int) bool {
		return scores[indices[i]].Total > scores[indices[j]].Total
	})
	res := osusu.Recipes{}
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
