package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gocolly/colly/v2"
)

// AllRecipesData contains the data from the allrecipes.com recipe information json text
type AllRecipesData struct {
	Type             []any `json:"@type"`
	Name             string
	Description      string
	Image            map[string]any
	RecipeCategory   []any
	RecipeCuisine    []any
	RecipeIngredient []any
	TotalTime        string
	PrepTime         string
	CookTime         string
	AggregateRating  map[string]any
	Nutrition        map[string]any
}

// ScrapeAllRecipes crawls allrecipes.com and saves the recipes it scrapes to web/data/allrecipes.csv
func ScrapeAllRecipes() {
	c := colly.NewCollector(colly.AllowedDomains("www.allrecipes.com"))

	c.OnXML("//sitemapindex/sitemap/loc", func(e *colly.XMLElement) {
		e.Request.Visit(e.Text)
	})

	total := 0
	c.OnXML("//urlset/url/loc", func(e *colly.XMLElement) {
		e.Request.Visit(e.Text)
	})

	c.OnHTML(`script[type="application/ld+json"]`, func(e *colly.HTMLElement) {
		// all info is stored in the first element of an array for some reason, so need to use slice
		initialData := []AllRecipesData{}
		err := json.Unmarshal([]byte(e.Text), &initialData)
		if err != nil {
			log.Println("error unmarshaling data for", e.Request.URL.String()+":", err)
			return
		}
		data := initialData[0]
		recipe := Recipe{
			Name:        data.Name,
			URL:         e.Request.URL.String(),
			Description: data.Description,
			Image:       ForceType[string](data.Image["url"]),
			Category:    ForceSliceType[string](data.RecipeCategory),
			Cuisine:     ForceSliceType[string](data.RecipeCuisine),
			Ingredients: ForceSliceType[string](data.RecipeCategory),
			TotalTime:   FormatDuration(data.TotalTime),
			PrepTime:    FormatDuration(data.PrepTime),
			CookTime:    FormatDuration(data.CookTime),
			RatingValue: ForceTypeFunc(data.AggregateRating["ratingValue"], func(value string) (float64, error) { return strconv.ParseFloat(value, 64) }),
			RatingCount: ForceTypeFunc(data.AggregateRating["ratingCount"], strconv.Atoi),
			Nutrition: Nutrition{
				Calories: ForceTypeFunc(data.Nutrition["calories"], AtoiIgnore),
			},
		}
		log.Println("got recipe", recipe)
		// // check if it is actually a recipe
		// types, ok := data["@type"].([]any)
		// if !ok {
		// 	log.Println("discarded", e.Request.URL, "because couldn't load type")
		// 	return
		// }
		// gotType := false
		// for _, typ := range types {
		// 	if typ == "Recipe" {
		// 		gotType = true
		// 	}
		// }
		// if !gotType {
		// 	log.Println("discarded", e.Request.URL, "because type", types, "does not contain Recipe")
		// 	return
		// }
		// recipe := Recipe{}
		// recipe.URL = e.Request.URL.String()
		// recipe.Name, ok = data["name"].(string)
		// if !ok {
		// 	log.Println("discarded", e.Request.URL, "because couldn't load name")
		// 	return
		// }
		// recipe.Description, ok = data["description"].(string)
		// if !ok {
		// 	log.Println("discarded", e.Request.URL, "because couldn't load description")
		// 	return
		// }
		// // image is map with url as object, so need to do this
		// image, ok := data["image"].(map[string]any)
		// if !ok {
		// 	log.Println("discarded", e.Request.URL, "because couldn't load image")
		// 	return
		// }
		// recipe.Image, ok = image["url"].(string)
		// if !ok {
		// 	log.Println("discarded", e.Request.URL, "because couldn't load image URL")
		// 	return
		// }
		// // categories are a slice of any, so we need to loop over to get in terms of slice of strings
		// categories, ok := data["recipeCategory"].([]any)
		// if !ok {
		// 	log.Println("discarded", e.Request.URL, "because couldn't load category")
		// 	return
		// }
		// recipe.Category = []string{}
		// for _, category := range categories {
		// 	categoryString, ok := category.(string)
		// 	if !ok {
		// 		continue
		// 	}
		// 	recipe.Category = append(recipe.Category, categoryString)
		// }
		// // cuisines are also a slice of any, so we need to loop over to get in terms of slice of strings
		// cuisines, ok := data["recipeCuisine"].([]any)
		// if !ok {
		// 	log.Println("discarded", e.Request.URL, "because couldn't load cuisine")
		// 	return
		// }
		// recipe.Cuisine = []string{}
		// for _, cuisine := range cuisines {
		// 	cuisineString, ok := cuisine.(string)
		// 	if !ok {
		// 		continue
		// 	}
		// 	recipe.Cuisine = append(recipe.Cuisine, cuisineString)
		// }
		// // for ingredients, once again need to loop over slice of any to get slice of strings
		// ingredients, ok := data["recipeIngredient"].([]any)
		// if !ok {
		// 	log.Println("discarded", e.Request.URL, "because couldn't load ingredients")
		// 	return
		// }
		// recipe.Ingredients = []string{}
		// for _, ingredient := range ingredients {
		// 	ingredientString, ok := ingredient.(string)
		// 	if !ok {
		// 		continue
		// 	}
		// 	recipe.Ingredients = append(recipe.Ingredients, ingredientString)
		// }
		// recipe.TotalTime, ok = data["totalTime"].(string)
		// if !ok {
		// 	log.Println("discarded", e.Request.URL, "because couldn't load total time")
		// 	return
		// }
		// recipe.TotalTime = FormatDuration(recipe.TotalTime)
		// recipe.PrepTime, ok = data["prepTime"].(string)
		// if !ok {
		// 	log.Println("discarded", e.Request.URL, "because couldn't load prep time")
		// 	return
		// }
		// recipe.PrepTime = FormatDuration(recipe.PrepTime)
		// recipe.CookTime, ok = data["cookTime"].(string)
		// if !ok {
		// 	log.Println("discarded", e.Request.URL, "because couldn't load cook time")
		// 	return
		// }
		// recipe.CookTime = FormatDuration(recipe.CookTime)
		// rating, ok := data["aggregateRating"].(map[string]any)
		// if !ok {
		// 	log.Println("discarded", e.Request.URL, "because couldn't load rating")
		// 	return
		// }
		// recipe.RatingValue, ok = rating["ratingValue"].(string)
		// if !ok {
		// 	log.Println("discarded", e.Request.URL, "because couldn't load rating value")
		// 	return
		// }
		// recipe.RatingCount, ok = rating["ratingCount"].(string)
		// if !ok {
		// 	log.Println("discarded", e.Request.URL, "because couldn't load rating count")
		// 	return
		// }
		// nutrition, ok := data["nutrition"].(map[string]any)
		// if !ok {
		// 	log.Println("discarded", e.Request.URL, "because couldn't load nutrition information")
		// 	return
		// }
		// calories, ok := nutrition["calories"].(string)
		// if !ok {
		// 	log.Println("discarded", e.Request.URL, "because couldn't load calories")
		// 	return
		// }
		// recipe.Nutrition.Calories, err = ParseInt(calories)
		// if err != nil {
		// 	log.Println("discarded", e.Request.URL, "because error with parsing calories:", err)
		// 	return
		// }
		// log.Println("got recipe:", recipe)
	})

	c.Visit("https://www.allrecipes.com/sitemap.xml")
	fmt.Println(total)
}

// ForceType converts the given value to the given type and returns the zero value of the given type if the value can not be converted
func ForceType[T any](value any) T {
	// need to catch ok value to prevent panic even if we don't use it
	val, _ := value.(T)
	return val
}

// ForceTypeFunc converts the given value to the given type T by first converting the value to type I and then calling the given function with that value. If the function returns an error, ForceTypeFunc returns the zero value of the given type.
func ForceTypeFunc[T, I any](value any, fun func(value I) (T, error)) T {
	val := ForceType[I](value)
	res, _ := fun(val)
	return res
}

// ForceSliceType converts the given slice value to a slice of the given type and uses a zero value of the given type for every element that can not be converted
func ForceSliceType[T any](value []any) []T {
	res := make([]T, len(value))
	for i, v := range value {
		// need to catch ok value to prevent panic even if we don't use it
		val, _ := v.(T)
		res[i] = val
	}
	return res
}
