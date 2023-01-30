package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// start is a component that displays a simple "Hello World!". A component is a
// customizable, independent, and reusable UI element. It is created by
// embedding app.Compo into a struct.
type start struct {
	app.Compo
}

// The Render method is where the component appearance is defined. Here, a
// "Hello World!" is displayed as a heading.
func (s *start) Render() app.UI {
	return app.Div().Body(
		app.H1().ID("start-page-title").Class("page-title").Text("MealRec"),
		app.P().ID("start-page-subtitle").Text("An app for getting recommendations on what meals to eat in a group based on the ratings of and compatibility with each member of the group."),
		app.Div().ID("start-page-action-button-row").Body(
			app.A().ID("start-page-sign-in").Class("blue-action-button").Href("/signin").Text("Sign In"),
			app.A().ID("start-page-sign-up").Class("blue-action-button").Href("/signup").Text("Sign Up"),
		),
	)
}
