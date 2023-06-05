// Classify determines the cuisines and categories of the recipes.
package main

import (
	"log"
	"sort"
	"time"

	"github.com/emer/emergent/decoder"
	"github.com/kkoreilly/osusu/osusu"
)

const (
	// Threshold is how many times a word has to occur for it to be included in the recipe words slice]
	// Was found using the following results from a simple error test with n = 3:
	// Threshold = 0: 58%
	// 10: 60%
	// 50: 58%
	// 100: 59%
	// 500: 57%
	// 1_000: 48%
	// 5_000: 32%
	// 10_000: 45%
	// 50_000: 90%
	// 100_000: 105%
	Threshold = 10
	// Rounds is how many times the neural network is run
	Rounds = 1000
)

func main() {
	log.Println("Starting Classify")
	t := time.Now()
	osusu.InitRecipeConstants()
	recipes, err := osusu.LoadRecipes()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("load recipes", time.Since(t))
	recipes = recipes.ConsolidateCategories()
	recipes = recipes.ConsolidateCuisines()
	log.Println("consolidate categories and cuisines", time.Since(t))
	oldMap, oldCount := recipes.CountCuisines()
	words := GetRecipeWords(recipes)
	log.Println("get recipe words", time.Since(t))
	recipes = InferCuisines(recipes, words)
	newMap, newCount := recipes.CountCuisines()
	osusu.RecipeNumberChanges(oldMap, oldCount, newMap, newCount)
	osusu.SaveRecipes(recipes)
}

// GetRecipeWords gets all of the words contained in all of the given recipes
func GetRecipeWords(recipes osusu.Recipes) []string {
	t := time.Now()
	wordMap := map[string]int{}
	for _, recipe := range recipes {

		// Simple error test with n = 3:
		// Ingredients only: 56-60%
		// Ingredients and description: 68%
		// Ingredients and name: 67%
		// Ingredients, description, and name: 72%
		// Therefore ingredients only is best
		// t := time.Now()
		words := []string{}
		for _, ingredient := range recipe.Ingredients {
			words = append(words, osusu.GetWords(ingredient)...)
		}
		// log.Println("get recipe words: inside: get words", time.Since(t))
		for _, word := range words {
			// if word == "33333334326744" {
			// 	log.Println(word, recipe.Ingredients, recipe.URL)
			// }
			wordMap[word]++
		}
		// log.Println("get recipe words: inside: add to word map", time.Since(t))
	}
	log.Println("get recipe words: get word map", time.Since(t))
	res := []string{}
	for word, num := range wordMap {
		if num < Threshold {
			// log.Println(word, num)
			continue
		}
		res = append(res, word)
	}
	log.Println("get recipe words: convert to slice", time.Since(t))
	sort.Slice(res, func(i, j int) bool {
		return wordMap[res[i]] < wordMap[res[j]]
	})
	// for _, word := range res {
	// 	// log.Println(word, wordMap[word])
	// }
	log.Println("get recipe words: sort and print", time.Since(t))
	return res
}

// InferCuisines infers the cuisines of all of the given recipes with them missing using the given recipe words using the SoftMax decoder.
func InferCuisines(recipes osusu.Recipes, words []string) osusu.Recipes {
	ca := decoder.SoftMax{}
	ca.Init(len(osusu.AllCategories), len(words))

	cu := decoder.SoftMax{}
	cu.Init(len(osusu.BaseCuisines), len(words))

	wordMap := map[string]int{}
	for i, word := range words {
		wordMap[word] = i
	}

	categoryMap := map[string]int{}
	for i, category := range osusu.AllCategories {
		categoryMap[category] = i
	}

	cuisineMap := map[string]int{}
	for i, cuisine := range osusu.BaseCuisines {
		cuisineMap[cuisine] = i
	}

	for i := 0; i < Rounds; i++ {
		log.Println("On Round", i)
		for j, recipe := range recipes {
			ca.Inputs = make([]float32, ca.NInputs)
			cu.Inputs = make([]float32, cu.NInputs)
			words := []string{}
			for _, ingredient := range recipe.Ingredients {
				words = append(words, osusu.GetWords(ingredient)...)
			}
			for _, word := range words {
				i, ok := wordMap[word]
				if ok {
					ca.Inputs[i] = 1
					cu.Inputs[i] = 1
				}
			}
			ca.Forward()
			cu.Forward()

			if len(recipe.Category) != 0 {
				ca.Train(categoryMap[recipe.Category[0]])
			}
			if len(recipe.Cuisine) != 0 {
				cu.Train(cuisineMap[recipe.Cuisine[0]])
			}

			if len(recipe.Category) == 0 && i == Rounds-1 {
				ca.Sort()
				category := osusu.AllCategories[ca.Sorted[0]]
				recipe.Category = []string{category}
				recipes[j] = recipe
				// log.Println(category, "-", recipe.Name)
			}

			if len(recipe.Cuisine) == 0 && i == Rounds-1 {
				cu.Sort()
				cuisine := osusu.BaseCuisines[cu.Sorted[0]]
				recipe.Cuisine = []string{cuisine}
				recipes[j] = recipe
				// log.Println(cuisine, "-", recipe.Name)
			}
		}
	}
	return recipes
}
