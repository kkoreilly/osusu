package main

import (
	"log"
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type home struct {
	app.Compo
	meals Meals
}

func (h *home) Render() app.UI {
	return app.Div().Body(
		app.H1().ID("home-page-title").Class("page-title").Text("Home"),
		app.Div().ID("home-page-action-button-row").Class("action-button-row").Body(
			app.A().ID("home-page-rec-button").Class("white-action-button", "action-button").Href("/recommendations").Text("Recommendations"),
			app.Button().ID("home-page-new-button").Class("blue-action-button", "action-button").Text("New").OnClick(h.New),
		),
		app.Hr(),
		app.Div().ID("home-page-meals-container").Body(
			app.Range(h.meals).Slice(func(i int) app.UI {
				meal := h.meals[i]
				return app.Div().ID("home-page-meal-" + strconv.Itoa(i)).Class("home-page-meal").Text(meal.Name).
					OnClick(func(ctx app.Context, e app.Event) { h.MealOnClick(ctx, e, meal) })
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
	h.meals = meals
}
