package main

import (
	"log"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type signUp struct {
	app.Compo
}

func (s *signUp) Render() app.UI {
	return app.Div().Body(
		app.H1().ID("sign-up-page-title").Class("page-title").Text("Sign Up"),
		app.Form().ID("sign-up-page-form").OnSubmit(s.OnSubmit).Body(
			app.Input().ID("sign-up-page-username").Class("sign-up-page-input").Type("text").Placeholder("Username"),
			app.Input().ID("sign-up-page-password").Class("sign-up-page-input").Type("password").Placeholder("Password"),
			app.Div().ID("sign-up-page-action-button-row").Class("action-button-row").Body(
				app.A().ID("sign-up-page-cancel").Class("action-button", "white-action-button").Href("/").Text("Cancel"),
				app.Input().ID("sign-up-page-submit").Class("action-button", "blue-action-button").Name("submit").Type("submit").Value("Sign Up"),
			)),
	)
}

func (s *signUp) OnSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()
	username := app.Window().GetElementByID("sign-up-page-username").Get("value")
	password := app.Window().GetElementByID("sign-up-page-password").Get("value")
	log.Println("Username:", username, "Password:", password)
}
