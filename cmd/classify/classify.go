// Classify determines the cuisines and categories of the recipes.
package main

import (
	"flag"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/emer/etable/eplot"
	"github.com/emer/etable/etable"
	"github.com/emer/etable/etensor"
	"github.com/goki/gi/gi"
	"github.com/goki/gi/gimain"
	"github.com/kkoreilly/gobp"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/kkoreilly/osusu/util/file"
)

const (
	// Threshold is how many times a word has to occur for it to be included in the recipe words slice
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
	ErrorTable.AddCol(etensor.NewFloat64([]int{1}, nil, nil), "categoryTestError")
	ErrorTable.AddCol(etensor.NewFloat64([]int{1}, nil, nil), "cuisineTestError")
	ErrorTable.AddCol(etensor.NewFloat64([]int{1}, nil, nil), "categoryAverageSSE")
	ErrorTable.AddCol(etensor.NewFloat64([]int{1}, nil, nil), "cuisineAverageSSE")
	ErrorTable.AddCol(etensor.NewFloat64([]int{1}, nil, nil), "categoryTrainError")
	ErrorTable.AddCol(etensor.NewFloat64([]int{1}, nil, nil), "cuisineTrainError")
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
	ErrorPlot.ColParams("categoryTestError").On = true
	ErrorPlot.ColParams("categoryTestError").Lbl = "Category Test Error"
	ErrorPlot.ColParams("cuisineTestError").On = true
	ErrorPlot.ColParams("cuisineTestError").Lbl = "Cuisine Test Error"
	ErrorPlot.ColParams("categoryAverageSSE").On = true
	ErrorPlot.ColParams("categoryAverageSSE").Lbl = "Category Average SSE"
	ErrorPlot.ColParams("cuisineAverageSSE").On = true
	ErrorPlot.ColParams("cuisineAverageSSE").Lbl = "Cuisine Average SSE"
	ErrorPlot.ColParams("categoryTrainError").On = true
	ErrorPlot.ColParams("categoryTrainError").Lbl = "Category Train Error"
	ErrorPlot.ColParams("cuisineTrainError").On = true
	ErrorPlot.ColParams("cuisineTrainError").Lbl = "Cuisine Train Error"
	ErrorPlot.SetStretchMax()

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
		words := append(osusu.GetWords(recipe.Name), osusu.GetWords(recipe.Description)...)
		for _, ingredient := range recipe.Ingredients {
			words = append(words, osusu.GetWords(ingredient)...)
		}
		if recipeWordsMap[i] == nil {
			recipeWordsMap[i] = []string{}
		}
		for _, word := range words {
			wordMap[word]++
			recipeWordsMap[i] = append(recipeWordsMap[i], word)
		}
	}
	words := []string{}
	for word, num := range wordMap {
		if num < Threshold {
			continue
		}
		words = append(words, word)
	}
	// remove words that don't meet the threshold from recipeWordsMap also
	for i, words := range recipeWordsMap {
		newWords := []string{}
		for _, word := range words {
			if wordMap[word] >= Threshold {
				newWords = append(newWords, word)
			}
		}
		recipeWordsMap[i] = newWords
	}
	// sort.Slice(words, func(i, j int) bool {
	// 	return wordMap[words[i]] < wordMap[words[j]]
	// })
	// for _, word := range words {
	// 	log.Println(word, wordMap[word])
	// }
	return words, recipeWordsMap
}

// Infer infers the categories and cuisines of all of the given recipes with them missing using the given recipe words using the SoftMax decoder.
func Infer(recipes osusu.Recipes, words []string, recipeWordsMap map[int][]string) osusu.Recipes {
	numHiddenLayers := flag.Int("layers", 2, "number of hidden layers")
	numHiddenUnits := flag.Int("units", 500, "number of hidden units")
	flag.Parse()
	ca := gobp.NewNetwork(len(words), len(osusu.AllCategories), *numHiddenLayers, *numHiddenUnits) // category decoder
	ca.OutputActivationFunc = gobp.SoftMax
	ca.LearningRate = 0.05

	cu := gobp.NewNetwork(len(words), len(osusu.BaseCuisines), *numHiddenLayers, *numHiddenUnits) // cuisine decoder
	cu.OutputActivationFunc = gobp.SoftMax
	cu.LearningRate = 0.05

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

	// indices with category set
	caSet := []int{}
	// indices with category unset
	caUnset := []int{}

	// indices with cuisine set
	cuSet := []int{}
	// indices with cuisine unset
	cuUnset := []int{}

	for i, recipe := range recipes {
		if len(recipe.Category) == 0 {
			caUnset = append(caUnset, i)
		} else {
			caSet = append(caSet, i)
		}
		if len(recipe.Cuisine) == 0 {
			cuUnset = append(cuUnset, i)
		} else {
			cuSet = append(cuSet, i)
		}
	}

	rand.Seed(time.Now().UnixNano())

	caPerm := rand.Perm(len(caSet))
	caTrain := map[int]bool{}
	for i, val := range caPerm {
		if i < 9*len(caPerm)/10 {
			caTrain[val] = true
		}
	}

	cuPerm := rand.Perm(len(cuSet))
	cuTrain := map[int]bool{}
	for i, val := range cuPerm {
		if i < 9*len(cuPerm)/10 {
			cuTrain[val] = true
		}
	}

	log.Println(len(caSet), len(caUnset), len(cuSet), len(cuUnset))
	log.Println(len(caPerm), len(caTrain), len(cuPerm), len(cuTrain))

	recipeInputs := make([][]float32, len(recipes))

	for i := range recipes {
		inputs := make([]float32, len(words))
		words := recipeWordsMap[i]
		for _, word := range words {
			j, ok := wordMap[word]
			if ok {
				inputs[j] = 1
			}
		}
		recipeInputs[i] = inputs
	}

	for e := 0; e < Epochs; e++ {
		log.Println("Starting Epoch", e)
		var caNum, caNumRight, caTrainNum, caTrainNumRight, cuNum, cuNumRight, cuTrainNum, cuTrainNumRight int
		var caTotalSSE, cuTotalSSE float32
		var totalTime time.Duration
		for i, recipe := range recipes {
			if i%1000 == 0 {
				log.Println("Starting Recipe", i, totalTime)
				totalTime = 0
			}

			ca.Inputs = recipeInputs[i]
			cu.Inputs = recipeInputs[i]

			t := time.Now()

			ca.Forward()
			cu.Forward()

			if caTrain[i] {
				for j := 0; j < len(osusu.AllCategories); j++ {
					ca.Targets[j] = 0
				}
				for _, category := range recipe.Category {
					ca.Targets[categoryMap[category]] = 1
				}

				sse := ca.Back()
				if math.IsNaN(float64(sse)) || math.IsInf(float64(sse), 0) {
					log.Println("infinite or NaN category sse: epoch", e, "recipe", i, "sse", sse)
				} else {
					caTotalSSE += sse
				}
			}
			if cuTrain[i] {
				for j := 0; j < len(osusu.BaseCuisines); j++ {
					cu.Targets[j] = 0
				}
				for _, cuisine := range recipe.Cuisine {
					cu.Targets[cuisineMap[cuisine]] = 1
				}
				sse := cu.Back()
				if math.IsNaN(float64(sse)) || math.IsInf(float64(sse), 0) {
					log.Println("infinite or NaN cuisine sse: epoch", e, "recipe", i, "sse", sse)
				} else {
					cuTotalSSE += sse
				}
			}
			totalTime += time.Since(t)

			outputs := ca.Outputs()
			highestIdx := -1
			var highestVal float32 = -1000000
			for j, output := range outputs {
				if output > highestVal {
					highestIdx = j
					highestVal = output
				}
			}
			category := osusu.AllCategories[highestIdx]
			// only if we have something to test -- we always test and set normal nums if test and train nums if train
			if len(recipe.Category) != 0 {
				if caTrain[i] {
					caTrainNum++
				} else {
					caNum++
				}
				for _, recipeCategory := range recipe.Category {
					if category == recipeCategory {
						if caTrain[i] {
							caTrainNumRight++
						} else {
							caNumRight++
						}
						break
					}
				}
			}

			// only capture if on last round and not already set
			if len(recipe.Category) == 0 && e == Epochs-1 {
				recipe.Category = []string{category}
				recipes[i] = recipe
			}

			outputs = cu.Outputs()
			highestIdx = -1
			highestVal = 0
			for j, output := range outputs {
				if output > highestVal {
					highestIdx = j
					highestVal = output
				}
			}
			cuisine := osusu.BaseCuisines[highestIdx]
			// only if we have something to test -- we always test and set normal nums if test and train nums if train
			if len(recipe.Cuisine) != 0 {
				if cuTrain[i] {
					cuTrainNum++
				} else {
					cuNum++
				}
				for _, recipeCuisine := range recipe.Cuisine {
					if cuisine == recipeCuisine {
						if cuTrain[i] {
							cuTrainNumRight++
						} else {
							cuNumRight++
						}
						break
					}
				}
			}
			// only capture if on last round and not already set
			if len(recipe.Cuisine) == 0 && e == Epochs-1 {
				recipe.Cuisine = []string{cuisine}
				recipes[i] = recipe
			}
		}
		caAverageSSE := float64(caTotalSSE) / float64(len(caTrain))
		cuAverageSSE := float64(cuTotalSSE) / float64(len(cuTrain))
		caError := 1 - (float64(caNumRight) / float64(caNum))
		cuError := 1 - (float64(cuNumRight) / float64(cuNum))
		caTrainError := 1 - (float64(caTrainNumRight) / float64(caTrainNum))
		cuTrainError := 1 - (float64(cuTrainNumRight) / float64(cuTrainNum))
		log.Println(e, caAverageSSE, cuAverageSSE, caError, cuError, caTrainError, cuTrainError)
		ErrorTable.SetCellFloat("epoch", e, float64(e))
		ErrorTable.SetCellFloat("categoryTestError", e, caError)
		ErrorTable.SetCellFloat("cuisineTestError", e, cuError)
		ErrorTable.SetCellFloat("categoryAverageSSE", e, caAverageSSE)
		ErrorTable.SetCellFloat("cuisineAverageSSE", e, cuAverageSSE)
		ErrorTable.SetCellFloat("categoryTrainError", e, caTrainError)
		ErrorTable.SetCellFloat("cuisineTrainError", e, cuTrainError)
		ErrorPlot.GoUpdatePlot()
	}
	return recipes
}
