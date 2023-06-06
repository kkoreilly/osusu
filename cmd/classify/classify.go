// Classify determines the cuisines and categories of the recipes.
package main

import (
	"log"
	"sort"
	"time"

	"github.com/emer/emergent/decoder"
	"github.com/goki/gi/gi"
	"github.com/goki/gi/gimain"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/kkoreilly/osusu/util/file"
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
	gimain.Main(func() {
		mainrun()
	})
}

func mainrun() {
	log.Println("Starting Classify")
	CreateUI()
	t := time.Now()
	osusu.InitRecipeConstants()
	var recipes osusu.Recipes
	var err error
	// if we have already classified (so initial is moved to old), get from initial (scraper) because we can't classify already classified recipes.
	// otherwise, load from main (which should also be initial/scraper)
	if file.Exists("web/data/recipes.json.old") {
		recipes, err = osusu.LoadRecipes("web/data/recipes.json.old")
	} else {
		recipes, err = osusu.LoadRecipes("web/data/recipes.json")
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println("load recipes", time.Since(t))
	recipes = recipes.ConsolidateCategories()
	recipes = recipes.ConsolidateCuisines()
	log.Println("consolidate categories and cuisines", time.Since(t))
	oldMap, oldCount := recipes.CountCuisines()
	words, recipeWordsMap := GetRecipeWords(recipes)
	log.Println(len(words), len(recipeWordsMap), len(recipes))
	log.Println("get recipe words", time.Since(t))
	recipes = InferCuisines(recipes, words, recipeWordsMap)
	newMap, newCount := recipes.CountCuisines()
	osusu.RecipeNumberChanges(oldMap, oldCount, newMap, newCount)
	osusu.SaveRecipes(recipes, "web/data/recipes.json")
}

// CreateUI creates the graphical user interface to display plots
func CreateUI() {
	win := gi.NewMainWindow("osusu-classify", "Osusu Classify Recipes", 1024, 768)
	vp := win.WinViewport2D()
	updt := vp.UpdateStart()
	mfr := win.SetMainFrame()

	title := gi.AddNewLabel(mfr, "classify", "Classify")
	title.SetProp("text-align", "center")
	title.SetStretchMaxWidth()

	vp.UpdateEndNoSig(updt)
	win.GoStartEventLoop()
}

// GetRecipeWords gets all of the words contained in all of the given recipes and returns a slice of all of them and a map of the words for each recipe index
func GetRecipeWords(recipes osusu.Recipes) ([]string, map[int][]string) {
	t := time.Now()
	wordMap := map[string]int{}
	recipeWordsMap := map[int][]string{}
	for i, recipe := range recipes {

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
		if recipeWordsMap[i] == nil {
			recipeWordsMap[i] = []string{}
		}
		for _, word := range words {
			// if word == "33333334326744" {
			// 	log.Println(word, recipe.Ingredients, recipe.URL)
			// }
			wordMap[word]++
			recipeWordsMap[i] = append(recipeWordsMap[i], word) // could probably save a couple milliseconds by only adding if the word meets the threshold and is included in the words list, but not worth it now
		}
		// log.Println("get recipe words: inside: add to word map", time.Since(t))
	}
	log.Println("get recipe words: get word map", time.Since(t))
	words := []string{}
	for word, num := range wordMap {
		if num < Threshold {
			// log.Println(word, num)
			continue
		}
		words = append(words, word)
	}
	// totalOld := 0
	// totalNew := 0
	// remove words that don't meet the threshold from recipeWordsMap also
	for i, words := range recipeWordsMap {
		// totalOld += len(words)
		newWords := []string{}
		for _, word := range words {
			if wordMap[word] >= Threshold {
				newWords = append(newWords, word)
			}
		}
		recipeWordsMap[i] = newWords
		// totalNew += len(newWords)
	}
	// log.Println("old", totalOld, "new", totalNew, "diff", totalOld-totalNew)
	log.Println("get recipe words: convert to slice", time.Since(t))
	sort.Slice(words, func(i, j int) bool {
		return wordMap[words[i]] < wordMap[words[j]]
	})
	for _, word := range words {
		log.Println(word, wordMap[word])
	}
	log.Println("get recipe words: sort and print", time.Since(t))
	return words, recipeWordsMap
}

// InferCuisines infers the cuisines of all of the given recipes with them missing using the given recipe words using the SoftMax decoder.
func InferCuisines(recipes osusu.Recipes, words []string, recipeWordsMap map[int][]string) osusu.Recipes {
	t := time.Now()
	ca := decoder.SoftMax{}
	ca.Init(len(osusu.AllCategories), len(words))

	cu := decoder.SoftMax{}
	cu.Init(len(osusu.BaseCuisines), len(words))
	log.Println("infer cuisines: decoder init", time.Since(t))

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

	log.Println("infer cuisines: map setup", time.Since(t))

	for i := 0; i < Rounds; i++ {
		log.Println("On Round", i)
		// t := time.Now()
		for j, recipe := range recipes {
			// t := time.Now()
			// ca.Inputs = make([]float32, ca.NInputs)
			// cu.Inputs = make([]float32, cu.NInputs)
			for i := range ca.Inputs {
				ca.Inputs[i] = 0
			}
			for i := range cu.Inputs {
				cu.Inputs[i] = 0
			}
			// log.Println("infer cuisines: round: recipe: make inputs", time.Since(t))
			// words := []string{}
			// for _, ingredient := range recipe.Ingredients {
			// 	words = append(words, osusu.GetWords(ingredient)...)
			// }

			words := recipeWordsMap[j]
			// log.Println("infer cuisines: round: get words", time.Since(t))
			for _, word := range words {
				i, ok := wordMap[word]
				if ok {
					ca.Inputs[i] = 1
					cu.Inputs[i] = 1
				}
			}
			// log.Println("infer cuisines: round: set inputs", time.Since(t))
			ca.Forward()
			cu.Forward()
			// log.Println("infer cuisines: round: forward", time.Since(t))
			if len(recipe.Category) != 0 {
				ca.Train(categoryMap[recipe.Category[0]])
			}
			if len(recipe.Cuisine) != 0 {
				cu.Train(cuisineMap[recipe.Cuisine[0]])
			}
			// log.Println("infer cuisines: round: train", time.Since(t))

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
			// log.Println("infer cuisines: round", time.Since(t))
		}
		// log.Println("infer cuisines: round", time.Since(t))
	}
	return recipes
}
