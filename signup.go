package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type signUp struct {
	app.Compo
}

func (s *signUp) Render() app.UI {
	return app.Div().Body(
		app.H1().ID("sign-up-page-title").Class("page-title").Text("Sign Up"),
		app.Form().ID("sign-up-page-form").OnSubmit(s.OnSubmit).Body(
			app.Input().ID("sign-up-page-username").Class("sign-up-page-input").Type("text").Placeholder("Username"),
			app.Input().ID("sign-up-page-password").Class("sign-up-page-input").Type("password").Placeholder("Password"),
			app.Input().ID("sign-up-page-submit").Class("blue-action-button").Type("submit").Value("Sign Up"),
		),
	)
}

func (s *signUp) OnSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()
}
