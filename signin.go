package main

import (
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type signIn struct {
	app.Compo
	user         User
	showPassword bool
}

func (s *signIn) Render() app.UI {
	passwordInputType := "password"
	if s.showPassword {
		passwordInputType = "text"
	}
	return &Page{
		ID:                     "sign-in",
		Title:                  "Sign In",
		Description:            "Sign into Osusu",
		AuthenticationRequired: false,
		TitleElement:           "Sign In",
		Elements: []app.UI{
			app.Form().ID("sign-in-page-form").Class("form").OnSubmit(s.OnSubmit).Body(
				NewTextInput("sign-in-page-username", "Username:", "Username", true, &s.user.Username),
				NewTextInput("sign-in-page-password", "Password:", "Password", false, &s.user.Password).SetType(passwordInputType),
				app.Div().ID("sign-in-page-checkboxes").Class("action-button-row").Body(
					NewCheckboxChip("sign-in-page-show-password", "Show Password", false, &s.showPassword),
					NewCheckboxChip("sign-in-page-remember-me", "Remember Me", true, &s.user.RememberMe),
				),
				app.Div().ID("sign-in-page-action-button-row").Class("action-button-row").Body(
					app.A().ID("sign-in-page-cancel").Class("action-button", "secondary-action-button").Href("/").Text("Cancel").Title("Cancel the Sign In and Navigate to the Start"),
					app.Button().ID("sign-in-page-submit").Class("action-button", "primary-action-button").Name("submit").Type("submit").Text("Sign In"),
				),
			),
		},
	}
}

func (s *signIn) OnSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	CurrentPage.ShowStatus("Loading...", StatusTypeNeutral)

	ctx.Defer(func(ctx app.Context) {
		user, err := SignInAPI.Call(s.user)
		if err != nil {
			CurrentPage.ShowErrorStatus(err)
			s.Update()
			return
		}
		// if no error, we are now authenticated
		authenticated = time.Now()
		SetCurrentUser(user, ctx)
		ctx.Navigate("/groups")
	})
}
