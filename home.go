package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type home struct {
	app.Compo
	meals map[string]Meal
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
			app.Range(h.meals).Map(func(k string) app.UI {
				return app.Div().ID("home-page-meal-" + k).Class("home-page-meal").Text(k).
					OnClick(func(ctx app.Context, e app.Event) { h.MealOnClick(ctx, e, h.meals[k]) })
			}),
		),
	)
}

func (h *home) New(ctx app.Context, e app.Event) {
	SetCurrentMeal(Meal{}, ctx)
	ctx.Navigate("/edit")
}

func (h *home) MealOnClick(ctx app.Context, e app.Event, meal Meal) {
	SetCurrentMeal(meal, ctx)
	ctx.Navigate("/edit")
}

func (h *home) OnNav(ctx app.Context) {
	h.meals = GetMeals(ctx)
	if h.meals == nil {
		h.meals = make(Meals)
		SetMeals(h.meals, ctx)
	}
}
