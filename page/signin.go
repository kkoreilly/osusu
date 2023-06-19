package page

import (
	"time"

	"github.com/kkoreilly/osusu/api"
	"github.com/kkoreilly/osusu/compo"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type SignIn struct {
	app.Compo
	user osusu.User
}

func (s *SignIn) Render() app.UI {
	return &compo.Page{
		ID:                     "sign-in",
		Title:                  "Sign In",
		Description:            "Sign into Osusu",
		AuthenticationRequired: false,
		TitleElement:           "Sign In",
		Elements: []app.UI{
			app.Form().ID("sign-in-page-form").Class("form").OnSubmit(s.OnSubmit).Body(
				compo.TextInput(&compo.Input[string]{ID: "sign-in-page-username", Label: "Username:", Value: &s.user.Username, AutoFocus: true}),
				compo.PasswordInput(&compo.Input[string]{ID: "sign-in-page-password", Label: "Password:", Value: &s.user.Password}),
				&compo.ButtonRow{ID: "sign-in-page-checkboxes", Buttons: []app.UI{
					&compo.CheckboxChip{ID: "sign-in-page-remember-me", Label: "Remember Me", Default: true, Value: &s.user.RememberMe},
				}},
				&compo.ButtonRow{ID: "sign-in-page", Buttons: []app.UI{
					&compo.Button{ID: "sign-in-page-cancel", Class: "secondary", Icon: "cancel", Text: "Cancel", OnClick: compo.NavigateEvent("/")},
					&compo.Button{ID: "sign-in-page-submit", Class: "primary", Type: "submit", Icon: "login", Text: "Sign In"},
				}},
			),
		},
	}
}

func (s *SignIn) OnSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	compo.CurrentPage.ShowStatus("Loading...", osusu.StatusTypeNeutral)

	ctx.Defer(func(ctx app.Context) {
		user, err := api.SignIn.Call(s.user)
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
