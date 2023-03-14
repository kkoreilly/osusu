package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type signIn struct {
	app.Compo
	user   User
	status string
}

func (s *signIn) Render() app.UI {
	return app.Div().Body(
		app.H1().ID("sign-in-page-title").Class("page-title").Text("Sign In"),
		app.Form().ID("sign-in-page-form").Class("form").OnSubmit(s.OnSubmit).Body(
			NewTextInput("sign-in-page-username", "Username:", "Username", true, &s.user.Username),
			NewTextInput("sign-in-page-password", "Password:", "Password", false, &s.user.Password).SetType("password"),
			NewCheckboxChip("sign-in-page-remember-me", "Remember Me", true, &s.user.RememberMe),
			app.Div().ID("sign-in-page-action-button-row").Class("action-button-row").Body(
				app.A().ID("sign-in-page-cancel").Class("action-button", "white-action-button").Href("/").Text("Cancel"),
				app.Input().ID("sign-in-page-submit").Class("action-button", "blue-action-button").Name("submit").Type("submit").Value("Sign In"),
			),
		),
		app.P().ID("sign-in-page-status").Class("status-text").Text(s.status),
	)
}

func (s *signIn) OnNav(ctx app.Context) {
	if Authenticate(false, ctx) {
		return
	}
}

func (s *signIn) OnSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	s.status = "Loading..."

	ctx.Defer(func(ctx app.Context) {
		user, err := SignInRequest(s.user)
		if err != nil {
			s.status = err.Error()
			s.Update()
			return
		}
		SetCurrentUser(user, ctx)
		ctx.Navigate("/people")
	})

}
