package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type signIn struct {
	app.Compo
	err string
}

func (s *signIn) Render() app.UI {
	return app.Div().Body(
		app.H1().ID("sign-in-page-title").Class("page-title").Text("Sign In"),
		app.Form().ID("sign-in-page-form").Class("form").OnSubmit(s.OnSubmit).Body(
			app.Input().ID("sign-in-page-username").Class("input", "sign-in-page-input").Name("username").Type("text").Placeholder("Username"),
			app.Input().ID("sign-in-page-password").Class("input", "sign-in-page-input").Name("password").Type("password").Placeholder("Password"),
			app.Div().ID("sign-in-page-action-button-row").Class("action-button-row").Body(
				app.A().ID("sign-in-page-cancel").Class("action-button", "white-action-button").Href("/").Text("Cancel"),
				app.Input().ID("sign-in-page-submit").Class("action-button", "blue-action-button").Name("submit").Type("submit").Value("Sign In"),
			),
		),
		app.P().ID("sign-in-page-error").Class("error-text").Text(s.err),
	)
}

func (s *signIn) OnSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	username := app.Window().GetElementByID("sign-in-page-username").Get("value").String()
	password := app.Window().GetElementByID("sign-in-page-password").Get("value").String()
	user := User{Username: username, Password: password}

	err := SignInRequest(user)
	if err != nil {
		s.err = err.Error()
		return
	}
	ctx.Navigate("/people")
}
