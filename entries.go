package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type entries struct {
	app.Compo
	meal   Meal
	person Person
}

func (e *entries) Render() app.UI {
	return &Page{
		ID:                     "entries",
		Title:                  "Entries",
		Description:            "View meal entries",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			e.person = GetCurrentPerson(ctx)
			e.meal = GetCurrentMeal(ctx)
		},
		TitleElement: "Entries for " + e.meal.Name,
		Elements: []app.UI{
			app.Div().ID("entries-page-action-button-row").Class("action-button-row").Body(
				app.A().ID("entries-page-back-button").Class("white-action-button", "action-button").Href("/edit").Text("Back"),
			),
			app.Hr(),
		},
	}
}
