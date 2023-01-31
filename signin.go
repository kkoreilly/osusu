package main

import (
	"log"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type signIn struct {
	app.Compo
}

func (s *signIn) Render() app.UI {
	return app.Div().Body(
		app.H1().ID("sign-in-page-title").Class("page-title").Text("Sign In"),
		app.Form().ID("sign-in-page-form").OnSubmit(s.OnSubmit).Body(
			app.Input().ID("sign-in-page-username").Class("sign-in-page-input").Name("username").Type("text").Placeholder("Username"),
			app.Input().ID("sign-in-page-password").Class("sign-in-page-input").Name("password").Type("password").Placeholder("Password"),
			app.Div().ID("sign-in-page-action-button-row").Class("action-button-row").Body(
				app.A().ID("sign-in-page-cancel").Class("action-button", "white-action-button").Href("/").Text("Cancel"),
				app.Input().ID("sign-in-page-submit").Class("action-button", "blue-action-button").Name("submit").Type("submit").Value("Sign In"),
			),
		),
	)
}

func (s *signIn) OnSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()
	username := app.Window().GetElementByID("sign-in-page-username").Get("value")
	password := app.Window().GetElementByID("sign-in-page-password").Get("value")
	log.Println("Username:", username, "Password:", password)
}
