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
				TextInput().ID("sign-in-page-username").Label("Username:").Value(&s.user.Username).AutoFocus(true),
				TextInput().ID("sign-in-page-password").Type(passwordInputType).Label("Password:").Value(&s.user.Password),
				ButtonRow().ID("sign-in-page-checkboxes").Buttons(
					CheckboxChip().ID("sign-in-page-show-password").Label("Show Password").Default(false).Value(&s.showPassword),
					CheckboxChip().ID("sign-in-page-remember-me").Label("Remember Me").Default(true).Value(&s.user.RememberMe),
				),
				ButtonRow().ID("sign-in-page").Buttons(
					Button().ID("sign-in-page-cancel").Class("secondary").Icon("cancel").Text("Cancel").OnClick(NavigateEvent("/")),
					Button().ID("sign-in-page-submit").Class("primary").Type("submit").Icon("login").Text("Sign In"),
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
		Navigate("/groups", ctx)
	})
}
