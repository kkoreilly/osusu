package main

import (
	"log"
	"sort"
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type home struct {
	app.Compo
	meals  Meals
	person Person
}

func (h *home) Render() app.UI {
	return app.Div().Body(
		app.H1().ID("home-page-title").Class("page-title").Text("Welcome Home, "+h.person.Name),
		app.Div().ID("home-page-action-button-row").Class("action-button-row").Body(
			app.A().ID("home-page-rec-button").Class("white-action-button", "action-button").Href("/recommendations").Text("Recommendations"),
			app.Button().ID("home-page-new-button").Class("blue-action-button", "action-button").Text("New").OnClick(h.New),
		),
		app.Hr(),
		app.Div().ID("home-page-meals-container").Body(
			app.Range(h.meals).Slice(func(i int) app.UI {
				meal := h.meals[i]
				si := strconv.Itoa(i)
				score := meal.Score()
				colorH := strconv.Itoa((score * 12) / 10)
				scoreText := strconv.Itoa(score)
				return app.Div().ID("home-page-meal-"+si).Class("home-page-meal").Style("--color-h", colorH).Style("--score-percent", scoreText+"%").
					OnClick(func(ctx app.Context, e app.Event) { h.MealOnClick(ctx, e, meal) }).Body(
					app.P().ID("home-page-meal-name-"+si).Class("home-page-meal-name").Text(meal.Name),
					app.P().ID("home-page-meal-score-"+si).Class("home-page-meal-score").Text(scoreText),
				)
			}),
		),
	)
}

func (h *home) New(ctx app.Context, e app.Event) {
	meal, err := CreateMealRequest(GetCurrentUser(ctx))
	if err != nil {
		log.Println(err)
		return
	}
	SetCurrentMeal(meal, ctx)
	ctx.Navigate("/edit")
}

func (h *home) MealOnClick(ctx app.Context, e app.Event, meal Meal) {
	SetCurrentMeal(meal, ctx)
	ctx.Navigate("/edit")
}

func (h *home) OnNav(ctx app.Context) {
	meals, err := GetMealsRequest(GetCurrentUser(ctx))
	if err != nil {
		log.Println(err)
	}
	sort.Slice(meals, func(i, j int) bool {
		return meals[i].Score() > meals[j].Score()
	})
	h.meals = meals
	h.person = GetCurrentPerson(ctx)
}
