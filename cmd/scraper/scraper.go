package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"unicode"

	"github.com/gocolly/colly/v2"
	"github.com/kkoreilly/osusu/osusu"
)

// AllRecipesData contains the data from the allrecipes.com recipe information json text
type AllRecipesData struct {
	Type             []string `json:"@type"`
	Name             string
	Description      string
	Image            map[string]any
	Author           []map[string]string
	DatePublished    string
	DateModified     string
	RecipeCategory   []string
	RecipeCuisine    []string
	RecipeIngredient []string
	TotalTime        string
	PrepTime         string
	CookTime         string
	RecipeYield      []string
	AggregateRating  map[string]string
	Nutrition        map[string]string
}

// ScrapeAllRecipes crawls allrecipes.com and saves the recipes it scrapes to web/data/allrecipes.csv
func ScrapeAllRecipes() {
	log.Println("Scraping allrecipes.com")
	c := colly.NewCollector(colly.AllowedDomains("www.allrecipes.com"))

	recipes := osusu.Recipes{}

	c.OnXML("//sitemapindex/sitemap/loc", func(e *colly.XMLElement) {
		e.Request.Visit(e.Text)
	})

	c.OnXML("//urlset/url/loc", func(e *colly.XMLElement) {
		e.Request.Visit(e.Text)
	})

	c.OnHTML(`script[type="application/ld+json"]`, func(e *colly.HTMLElement) {
		// all info is stored in the first element of an array for some reason, so need to use slice
		initialData := []AllRecipesData{}
		err := json.Unmarshal([]byte(e.Text), &initialData)
		if err != nil {
			log.Println(fmt.Errorf("error unmarshaling data for %v: %w", e.Request.URL.String(), err))
			return
		}
		data := initialData[0]
		// check if it is actually recipe
		gotType := false
		for _, typ := range data.Type {
			if typ == "Recipe" {
				gotType = true
			}
		}
		if !gotType {
			log.Println("discarding", e.Request.URL, "because no type Recipe in", data.Type)
			return
		}
		author := ""
		for i, v := range data.Author {
			author += v["name"]
			if i != len(data.Author)-1 {
				author += " "
			}
		}
		yield := ""
		if len(data.RecipeYield) > 0 {
			yield = data.RecipeYield[0]
		}
		datePublished, _ := time.Parse(time.RFC3339, data.DatePublished)
		dateModified, _ := time.Parse(time.RFC3339, data.DateModified)
		// NOTE: the reason we force many things and ignore errors is that some recipes have some data missing, and if they do, we still want to use them and we don't care about the errors
		recipe := osusu.Recipe{
			Name:          data.Name,
			URL:           e.Request.URL.String(),
			Description:   data.Description,
			Image:         ForceType[string](data.Image["url"]),
			Author:        author,
			DatePublished: datePublished,
			DateModified:  dateModified,
			Category:      data.RecipeCategory,
			Cuisine:       data.RecipeCuisine,
			Ingredients:   data.RecipeIngredient,
			TotalTime:     osusu.FormatDuration(data.TotalTime),
			PrepTime:      osusu.FormatDuration(data.PrepTime),
			CookTime:      osusu.FormatDuration(data.CookTime),
			Yield:         ForceFunc(yield, AtoiIgnore),
			RatingValue:   ForceFunc(data.AggregateRating["ratingValue"], ParseFloat64),
			RatingCount:   ForceFunc(data.AggregateRating["ratingCount"], AtoiIgnore),
			Nutrition: osusu.Nutrition{
				Calories:       ForceFunc(data.Nutrition["calories"], AtoiIgnore),
				Carbohydrate:   ForceFunc(data.Nutrition["carbohydrateContent"], AtoiIgnore),
				Cholesterol:    ForceFunc(data.Nutrition["cholesterolContent"], AtoiIgnore),
				Fiber:          ForceFunc(data.Nutrition["fiberContent"], AtoiIgnore),
				Protein:        ForceFunc(data.Nutrition["proteinContent"], AtoiIgnore),
				Fat:            ForceFunc(data.Nutrition["fatContent"], AtoiIgnore),
				SaturatedFat:   ForceFunc(data.Nutrition["saturatedFatContent"], AtoiIgnore),
				UnsaturatedFat: ForceFunc(data.Nutrition["unsaturatedFatContent"], AtoiIgnore),
				Sodium:         ForceFunc(data.Nutrition["sodiumContent"], AtoiIgnore),
				Sugar:          ForceFunc(data.Nutrition["sugarContent"], AtoiIgnore),
			},
		}
		recipes = append(recipes, recipe)
		log.Println("total recipes:", len(recipes), "\ngot new recipe:", recipe)
	})

	c.Visit("https://www.allrecipes.com/sitemap.xml")

	data, err := json.Marshal(recipes)
	if err != nil {
		log.Println(fmt.Errorf("failed to save crawled recipes as json: %w", err))
	}
	err = os.WriteFile("web/data/allrecipes.json", data, 0666)
	if err != nil {
		log.Println(fmt.Errorf("failed to save crawled recipes to file: %w", err))
	}
}

// ForceType converts the given value to the given type and returns the zero value of the given type if the value can not be converted
func ForceType[T any](value any) T {
	// need to catch ok value to prevent panic even if we don't use it
	val, _ := value.(T)
	return val
}

// ForceFunc calls the given function with the given value and returns the result, ignoring any error.
func ForceFunc[I, O any](value I, fun func(value I) (O, error)) O {
	res, _ := fun(value)
	return res
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

// ParseFloat64 parses the given string as a float64
func ParseFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
