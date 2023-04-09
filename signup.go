package main

import (
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type signUp struct {
	app.Compo
	user         User
	showPassword bool
}

func (s *signUp) Render() app.UI {
	passwordInputType := "password"
	if s.showPassword {
		passwordInputType = "text"
	}
	return &Page{
		ID:                     "sign-up",
		Title:                  "Sign Up",
		Description:            "Sign up to Osusu",
		AuthenticationRequired: false,
		TitleElement:           "Sign Up",
		Elements: []app.UI{
			app.Form().ID("sign-up-page-form").Class("form").OnSubmit(s.OnSubmit).Body(
				NewTextInput("sign-up-page-username", "Username:", "Username", true, &s.user.Username),
				NewTextInput("sign-up-page-password", "Password:", "Password", false, &s.user.Password).SetType(passwordInputType),
				app.Div().ID("sign-up-page-checkboxes").Class("action-button-row").Body(
					NewCheckboxChip("sign-up-page-show-password", "Show Password", false, &s.showPassword),
					NewCheckboxChip("sign-up-page-remember-me", "Remember Me", true, &s.user.RememberMe),
				),
				app.Div().ID("sign-up-page-action-button-row").Class("action-button-row").Body(
					app.A().ID("sign-up-page-cancel").Class("action-button", "secondary-action-button").Href("/").Text("Cancel"),
					app.Input().ID("sign-up-page-submit").Class("action-button", "primary-action-button").Name("submit").Type("submit").Value("Sign Up"),
				),
			),
		},
	}
}

func (s *signUp) OnSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	CurrentPage.ShowStatus("Loading...", StatusTypeNeutral)

	ctx.Defer(func(ctx app.Context) {
		user, err := SignUpAPI.Call(s.user)
		if err != nil {
			CurrentPage.ShowStatus(err.Error(), StatusTypeNegative)
			s.Update()
			return
		}
		// if no error, we are now authenticated
		authenticated = time.Now()
		SetCurrentUser(user, ctx)
		ctx.Navigate("/people")
	})
}
