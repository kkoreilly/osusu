// Classify determines the cuisines and categories of the recipes.
package main

import (
	"log"
	"sort"

	"github.com/emer/emergent/decoder"
	"github.com/emer/etable/eplot"
	"github.com/emer/etable/etable"
	"github.com/emer/etable/etensor"
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
	// Epochs is how many times the neural network is run
	Epochs = 1000
)

var ErrorTable *etable.Table

var ErrorPlot *eplot.Plot2D

var vp *gi.Viewport2D

func main() {
	log.Println("Starting Classify")
	gimain.Main(func() {
		mainrun()
	})
}

func mainrun() {
	win := gi.NewMainWindow("osusu-classify", "Osusu Classify Recipes", 1024, 768)
	vp = win.WinViewport2D()
	updt := vp.UpdateStart()
	mfr := win.SetMainFrame()

	ErrorTable = etable.NewTable("error-table")
	ErrorTable.AddCol(etensor.NewFloat64([]int{1}, nil, nil), "epoch")
	ErrorTable.AddCol(etensor.NewFloat64([]int{1}, nil, nil), "categoryError")
	ErrorTable.AddCol(etensor.NewFloat64([]int{1}, nil, nil), "cuisineError")
	ErrorTable.AddRows(Epochs)

	ErrorPlot = eplot.AddNewPlot2D(mfr, "error-plot")
	ErrorPlot.SetTable(ErrorTable)
	ErrorPlot.Params.FmMetaMap(ErrorTable.MetaData)
	ErrorPlot.Params.Title = "Classify"
	ErrorPlot.Params.XAxisCol = "epoch"
	ErrorPlot.Params.XAxisLabel = "Epoch"
	ErrorPlot.Params.YAxisLabel = "Error"
	ErrorPlot.Params.Points = true

	ErrorPlot.ColParams("epoch").On = true
	ErrorPlot.ColParams("categoryError").On = true
	// ErrorPlot.ColParams("categoryError").Range.Max = 1
	// ErrorPlot.ColParams("categoryError").Range.FixMax = true
	ErrorPlot.ColParams("cuisineError").On = true
	// ErrorPlot.ColParams("cuisineError").Range.Max = 1
	// ErrorPlot.ColParams("cuisineError").Range.FixMax = true
	ErrorPlot.SetStretchMax()

	// ErrorPlot.Update()

	go Classify()

	vp.UpdateEndNoSig(updt)
	win.StartEventLoop()
}

// Classify loads, classifies, and saves the recipes
func Classify() {
	log.Println("Loading Recipes")
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

	log.Println("Consolidating Categories and Cuisines")

	recipes = recipes.ConsolidateCategories()
	recipes = recipes.ConsolidateCuisines()

	log.Println("Getting Recipe Words")

	oldMap, oldCount := recipes.CountCuisines()
	words, recipeWordsMap := GetRecipeWords(recipes)
	log.Println("Classifying Recipes")
	Infer(recipes, words, recipeWordsMap)

	newMap, newCount := recipes.CountCuisines()
	osusu.RecipeNumberChanges(oldMap, oldCount, newMap, newCount)
	osusu.SaveRecipes(recipes, "web/data/recipes.json")
}

// GetRecipeWords gets all of the words contained in all of the given recipes and returns a slice of all of them and a map of the words for each recipe index
func GetRecipeWords(recipes osusu.Recipes) ([]string, map[int][]string) {
	wordMap := map[string]int{}
	recipeWordsMap := map[int][]string{}
	for i, recipe := range recipes {

		// Simple error test with n = 3:
		// Ingredients only: 56-60%
		// Ingredients and description: 68%
		// Ingredients and name: 67%
		// Ingredients, description, and name: 72%
		// Therefore ingredients only is best
		words := append(osusu.GetWords(recipe.Name), osusu.GetWords(recipe.Description)...)
		for _, ingredient := range recipe.Ingredients {
			words = append(words, osusu.GetWords(ingredient)...)
		}
		if recipeWordsMap[i] == nil {
			recipeWordsMap[i] = []string{}
		}
		for _, word := range words {
			// if word == "33333334326744" {
			// 	log.Println(word, recipe.Ingredients, recipe.URL)
			// }
			wordMap[word]++
			recipeWordsMap[i] = append(recipeWordsMap[i], word)
		}
	}
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
	sort.Slice(words, func(i, j int) bool {
		return wordMap[words[i]] < wordMap[words[j]]
	})
	for _, word := range words {
		log.Println(word, wordMap[word])
	}
	return words, recipeWordsMap
}

// Infer infers the categories and cuisines of all of the given recipes with them missing using the given recipe words using the SoftMax decoder.
func Infer(recipes osusu.Recipes, words []string, recipeWordsMap map[int][]string) osusu.Recipes {
	ca := decoder.SoftMax{} // category decoder
	ca.Init(len(osusu.AllCategories), len(words))

	cu := decoder.SoftMax{} // cuisine decoder
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

	for i := 0; i < Epochs; i++ {
		log.Println("On Epoch", i)
		caNum := 0
		caNumRight := 0
		cuNum := 0
		cuNumRight := 0
		for j, recipe := range recipes {
			for i := range ca.Inputs {
				ca.Inputs[i] = 0
			}
			for i := range cu.Inputs {
				cu.Inputs[i] = 0
			}
			words := recipeWordsMap[j]
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
				caNum++
			}
			if len(recipe.Cuisine) != 0 {
				cu.Train(cuisineMap[recipe.Cuisine[0]])
				cuNum++
			}

			ca.Sort()
			category := osusu.AllCategories[ca.Sorted[0]]
			if len(recipe.Category) != 0 && category == recipe.Category[0] {
				caNumRight++
			}

			if len(recipe.Category) == 0 && i == Epochs-1 {
				recipe.Category = []string{category}
				recipes[j] = recipe
			}

			cu.Sort()
			cuisine := osusu.BaseCuisines[cu.Sorted[0]]
			if len(recipe.Cuisine) != 0 && cuisine == recipe.Cuisine[0] {
				cuNumRight++
			}
			if len(recipe.Cuisine) == 0 && i == Epochs-1 {
				recipe.Cuisine = []string{cuisine}
				recipes[j] = recipe
			}
		}
		caError := 1 - (float64(caNumRight) / float64(caNum))
		cuError := 1 - (float64(cuNumRight) / float64(cuNum))
		// log.Println(caError, cuError)
		ErrorTable.SetCellFloat("epoch", i, float64(i))
		ErrorTable.SetCellFloat("categoryError", i, caError)
		ErrorTable.SetCellFloat("cuisineError", i, cuError)
		ErrorPlot.GoUpdatePlot()
	}
	return recipes
}
