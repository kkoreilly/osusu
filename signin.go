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
			// app.Label().ID("sign-in-page-username-label").Class("input-label").For("sign-in-page-username").Text("Username:"),
			// app.Input().ID("sign-in-page-username").Class("input", "sign-in-page-input").Name("username").Type("text").Placeholder("Username").AutoFocus(true),
			// app.Label().ID("sign-in-page-password-label").Class("input-label").For("sign-in-page-password").Text("Password:"),
			// app.Input().ID("sign-in-page-password").Class("input", "sign-in-page-input").Name("password").Type("password").Placeholder("Password"),
			NewCheckboxChip("sign-in-page-remember-me", "Remember Me", true, &s.user.RememberMe),
			// Chip("sign-in-page-remember-me", "checkbox", "sign-in-page-remember-me", "Remember Me", true),
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
		// username := app.Window().GetElementByID("sign-in-page-username").Get("value").String()
		// password := app.Window().GetElementByID("sign-in-page-password").Get("value").String()
		// s.user.RememberMe = app.Window().GetElementByID("sign-in-page-remember-me-chip-input").Get("checked").Bool()

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
