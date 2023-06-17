package page

import (
	"time"

	"github.com/kkoreilly/osusu/api"
	"github.com/kkoreilly/osusu/compo"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type SignUp struct {
	app.Compo
	user         osusu.User
	showPassword bool
}

func (s *SignUp) Render() app.UI {
	passwordInputType := "password"
	if s.showPassword {
		passwordInputType = "text"
	}
	return &compo.Page{
		ID:                     "sign-up",
		Title:                  "Sign Up",
		Description:            "Sign up to Osusu",
		AuthenticationRequired: false,
		TitleElement:           "Sign Up",
		Elements: []app.UI{
			app.Form().ID("sign-up-page-form").Class("form").OnSubmit(s.OnSubmit).Body(
				compo.TextInput().ID("sign-up-page-username").Label("Username:").Value(&s.user.Username).AutoFocus(true),
				compo.TextInput().ID("sign-up-page-password").Type(passwordInputType).Label("Password:").Value(&s.user.Password),
				compo.TextInput().ID("sign-up-page-name").Label("Name:").Value(&s.user.Name),
				&compo.ButtonRow{ID: "sign-up-page-checkboxes", Buttons: []app.UI{
					compo.CheckboxChip().ID("sign-up-page-show-password").Label("Show Password").Default(false).Value(&s.showPassword),
					compo.CheckboxChip().ID("sign-up-page-remember-me").Label("Remember Me").Default(true).Value(&s.user.RememberMe),
				}},
				&compo.ButtonRow{ID: "sign-up-page", Buttons: []app.UI{
					&compo.Button{ID: "sign-up-page-cancel", Class: "secondary", Icon: "cancel", Text: "Cancel", OnClick: compo.NavigateEvent("/")},
					&compo.Button{ID: "sign-up-page-submit", Class: "primary", Type: "submit", Icon: "app_registration", Text: "Sign Up"},
				}},
			),
		},
	}
}

func (s *SignUp) OnSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	compo.CurrentPage.ShowStatus("Loading...", osusu.StatusTypeNeutral)

	ctx.Defer(func(ctx app.Context) {
		user, err := api.SignUp.Call(s.user)
		if err != nil {
			compo.CurrentPage.ShowErrorStatus(err)
			s.Update()
			return
		}
		// if no error, we are now authenticated
		compo.Authenticated = time.Now()
		osusu.SetCurrentUser(user, ctx)
		compo.Navigate("/groups", ctx)
	})
}
