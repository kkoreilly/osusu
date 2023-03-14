package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type signUp struct {
	app.Compo
	user   User
	status string
}

func (s *signUp) Render() app.UI {
	return app.Div().Body(
		app.H1().ID("sign-up-page-title").Class("page-title").Text("Sign Up"),
		app.Form().ID("sign-up-page-form").Class("form").OnSubmit(s.OnSubmit).Body(
			NewTextInput("sign-up-page-username", "Username:", "Username", true, &s.user.Username),
			NewTextInput("sign-up-page-password", "Password:", "Password", false, &s.user.Password).SetType("password"),
			NewCheckboxChip("sign-up-page-remember-me", "Remember Me", true, &s.user.RememberMe),
			app.Div().ID("sign-up-page-action-button-row").Class("action-button-row").Body(
				app.A().ID("sign-up-page-cancel").Class("action-button", "white-action-button").Href("/").Text("Cancel"),
				app.Input().ID("sign-up-page-submit").Class("action-button", "blue-action-button").Name("submit").Type("submit").Value("Sign Up"),
			),
		),
		app.P().ID("sign-up-page-status").Class("status-text").Text(s.status),
	)
}

func (s *signUp) OnNav(ctx app.Context) {
	if Authenticate(false, ctx) {
		return
	}
}

func (s *signUp) OnSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	s.status = "Loading..."

	ctx.Defer(func(ctx app.Context) {
		user, err := CreateUserRequest(s.user)
		if err != nil {
			s.status = err.Error()
			s.Update()
			return
		}
		SetCurrentUser(user, ctx)
		ctx.Navigate("/people")
	})
}
