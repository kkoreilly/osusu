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
				TextInput().ID("sign-up-page-username").Label("Username:").Value(&s.user.Username).AutoFocus(true),
				TextInput().ID("sign-up-page-password").Type(passwordInputType).Label("Password:").Value(&s.user.Password),
				TextInput().ID("sign-up-page-name").Label("Name:").Value(&s.user.Name),
				ButtonRow().ID("sign-up-page-checkboxes").Buttons(
					CheckboxChip().ID("sign-up-page-show-password").Label("Show Password").Default(false).Value(&s.showPassword),
					CheckboxChip().ID("sign-up-page-remember-me").Label("Remember Me").Default(true).Value(&s.user.RememberMe),
				),
				ButtonRow().ID("sign-up-page").Buttons(
					Button().ID("sign-up-page-cancel").Class("secondary").Icon("cancel").Text("Cancel").OnClick(NavigateEvent("/")),
					Button().ID("sign-up-page-submit").Class("primary").Type("submit").Icon("app_registration").Text("Sign Up"),
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
			CurrentPage.ShowErrorStatus(err)
			s.Update()
			return
		}
		// if no error, we are now authenticated
		authenticated = time.Now()
		SetCurrentUser(user, ctx)
		Navigate("/groups", ctx)
	})
}
