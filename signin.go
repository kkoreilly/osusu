package main

import (
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type signIn struct {
	app.Compo
	user   User
	status string
}

func (s *signIn) Render() app.UI {
	return &Page{
		ID:                     "sign-in",
		Title:                  "Sign In",
		Description:            "Sign into MealRec",
		AuthenticationRequired: false,
		TitleElement:           "Sign In",
		Elements: []app.UI{
			app.Form().ID("sign-in-page-form").Class("form").OnSubmit(s.OnSubmit).Body(
				NewTextInput("sign-in-page-username", "Username:", "Username", true, &s.user.Username),
				NewTextInput("sign-in-page-password", "Password:", "Password", false, &s.user.Password).SetType("password"),
				NewCheckboxChip("sign-in-page-remember-me", "Remember Me", true, &s.user.RememberMe),
				app.Div().ID("sign-in-page-action-button-row").Class("action-button-row").Body(
					app.A().ID("sign-in-page-cancel").Class("action-button", "secondary-action-button").Href("/").Text("Cancel").Title("Cancel the Sign In and Navigate to the Start"),
					app.Input().ID("sign-in-page-submit").Class("action-button", "primary-action-button").Name("submit").Type("submit").Value("Sign In"),
				),
			),
			app.P().ID("sign-in-page-status").Class("status-text").Text(s.status),
		},
	}
}

func (s *signIn) OnSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	s.status = "Loading..."

	ctx.Defer(func(ctx app.Context) {
		user, err := SignInAPI.Call(s.user)
		if err != nil {
			s.status = err.Error()
			s.Update()
			return
		}
		// if no error, we are now authenticated
		authenticated = time.Now()
		SetCurrentUser(user, ctx)
		ctx.Navigate("/people")
	})
}
