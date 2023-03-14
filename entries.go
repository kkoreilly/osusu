package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type entries struct {
	app.Compo
	meal   Meal
	person Person
}

func (e *entries) Render() app.UI {
	return app.Div().Body(
		app.H1().ID("entries-page-title").Class("page-title").Text("Entries for "+e.meal.Name),
		app.Div().ID("entries-page-action-button-row").Class("action-button-row").Body(
			app.A().ID("entries-page-back-button").Class("white-action-button", "action-button").Href("/edit").Text("Back"),
		),
		app.Hr(),
	)
}

func (e *entries) OnNav(ctx app.Context) {
	if Authenticate(true, ctx) {
		return
	}
	e.person = GetCurrentPerson(ctx)
	e.meal = GetCurrentMeal(ctx)
}
