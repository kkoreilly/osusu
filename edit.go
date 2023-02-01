package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type edit struct {
	app.Compo
}

func (e *edit) Render() app.UI {
	return app.Div().Body(
		app.H1().ID("edit-page-title").Class("page-title").Text("Edit"),
		app.Form().ID("edit-page-form").Body(
			app.Label().ID("edit-page-name-label").Class("edit-page-label").For("edit-page-name-input").Text("Name:"),
			app.Input().ID("edit-page-name-input").Class("edit-page-input"),
		),
	)
}

func (e *edit) OnNav(ctx app.Context) {
	app.Window().GetElementByID("edit-page-name-input").Set("value", GetCurrentMeal(ctx).Name)
}
