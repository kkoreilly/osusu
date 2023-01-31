package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type start struct {
	app.Compo
}

func (s *start) Render() app.UI {
	return app.Div().Body(
		app.H1().ID("start-page-title").Class("page-title").Text("MealRec"),
		app.P().ID("start-page-subtitle").Text("An app for getting recommendations on what meals to eat in a group based on the ratings of and compatibility with each member of the group."),
		app.Div().ID("start-page-action-button-row").Class("action-button-row").Body(
			app.A().ID("start-page-sign-in").Class("action-button", "blue-action-button").Href("/signin").Text("Sign In"),
			app.A().ID("start-page-sign-up").Class("action-button", "blue-action-button").Href("/signup").Text("Sign Up"),
		),
	)
}
