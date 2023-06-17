// Package page provides types for all of the pages of the app
package page

import (
	"time"

	"github.com/kkoreilly/osusu/api"
	"github.com/kkoreilly/osusu/compo"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Account struct {
	app.Compo
	group osusu.Group
	user  osusu.User
}

func (a *Account) Render() app.UI {
	viewGroupIcon := "visibility"
	viewGroupText := "View Group"
	if a.user.ID == a.group.Owner {
		viewGroupIcon = "edit"
		viewGroupText = "Edit Group"
	}
	return &compo.Page{
		ID:                     "account",
		Title:                  "Account",
		Description:            "View and change account information",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			compo.SetReturnURL("/account", ctx)
			a.group = osusu.CurrentGroup(ctx)
			a.user = osusu.CurrentUser(ctx)
		},
		TitleElement:    "Account",
		SubtitleElement: "You are currently signed into " + a.user.Username + " with the name " + a.user.Name + " and the group " + a.group.Name + ".",
		Elements: []app.UI{
			&compo.ButtonRow{ID: "account-page-top", Buttons: []app.UI{
				&compo.Button{ID: "account-page-sign-out", Class: "danger", Icon: "logout", Text: "Sign Out", OnClick: a.InitialSignOut},
				&compo.Button{ID: "account-page-change-group", Class: "secondary", Icon: "group", Text: "Change Group", OnClick: compo.NavigateEvent("/groups")},
				&compo.Button{ID: "account-page-view-group", Class: "primary", Icon: viewGroupIcon, Text: viewGroupText, OnClick: a.ViewGroup},
			}},
			app.H2().ID("account-page-user-info-subtitle").Text("Change User Information:"),
			app.Form().ID("account-page-user-info-form").Class("form").OnSubmit(a.ChangeUserInfo).Body(
				compo.TextInput().ID("account-page-username").Label("Username:").Value(&a.user.Username),
				compo.TextInput().ID("account-page-name").Label("Name:").Value(&a.user.Name),
				&compo.Button{ID: "account-page-user-info-save", Class: "primary", Type: "submit", Icon: "save", Text: "Save"},
			),
			app.H2().ID("account-page-password-subtitle").Text("Change Password:"),
			app.Form().ID("account-page-password-form").Class("form").OnSubmit(a.ChangePassword).Body(
				compo.TextInput().ID("account-page-password").Type("password").Label("Password:").Value(&a.user.Password),
				&compo.Button{ID: "account-page-password-save", Class: "tertiary", Type: "submit", Icon: "save", Text: "Save"},
			),
			app.Dialog().ID("account-page-confirm-sign-out").Class("modal").Body(
				app.P().ID("account-page-confirm-sign-out-text").Class("confirm-delete-text").Text("Are you sure you want to sign out?"),
				&compo.ButtonRow{ID: "account-page-confirm-sign-out", Buttons: []app.UI{
					&compo.Button{ID: "account-page-confirm-sign-out-sign-out", Class: "danger", Icon: "logout", Text: "Yes, Sign Out", OnClick: a.ConfirmSignOut},
					&compo.Button{ID: "account-page-confirm-sign-out-cancel", Class: "secondary", Icon: "cancel", Text: "No, Cancel", OnClick: a.CancelSignOut},
				}},
			),
		},
	}
}

func (a *Account) InitialSignOut(ctx app.Context, e app.Event) {
	e.PreventDefault()
	app.Window().GetElementByID("account-page-confirm-sign-out").Call("showModal")
}

func (a *Account) ConfirmSignOut(ctx app.Context, e app.Event) {
	e.PreventDefault()
	user := osusu.CurrentUser(ctx)
	if user.Session != "" {
		_, err := api.SignOut.Call(osusu.CurrentUser(ctx))
		if err != nil {
			compo.CurrentPage.ShowErrorStatus(err)
			return
		}
	}
	// if no error, we are no longer authenticated
	compo.Authenticated = time.UnixMilli(0)
	ctx.LocalStorage().Del("currentUser")
	ctx.LocalStorage().Del("currentGroup")

	compo.Navigate("/signin", ctx)
}

func (a *Account) CancelSignOut(ctx app.Context, e app.Event) {
	e.PreventDefault()
	app.Window().GetElementByID("account-page-confirm-sign-out").Call("close")
}

func (a *Account) ChangeUserInfo(ctx app.Context, e app.Event) {
	e.PreventDefault()
	_, err := api.UpdateUserInfo.Call(a.user)
	if err != nil {
		compo.CurrentPage.ShowErrorStatus(err)
		return
	}
	compo.CurrentPage.ShowStatus("User Info Updated!", osusu.StatusTypePositive)
	osusu.SetCurrentUser(a.user, ctx)
}

func (a *Account) ChangePassword(ctx app.Context, e app.Event) {
	e.PreventDefault()

	compo.CurrentPage.ShowStatus("Loading...", osusu.StatusTypeNeutral)

	ctx.Defer(func(ctx app.Context) {
		_, err := api.UpdatePassword.Call(a.user)
		if err != nil {
			compo.CurrentPage.ShowErrorStatus(err)
			a.Update()
			a.user.Password = ""
			return
		}
		compo.CurrentPage.ShowStatus("Password Updated!", osusu.StatusTypePositive)
		a.Update()
		a.user.Password = ""
		osusu.SetCurrentUser(a.user, ctx)
	})
}

func (a *Account) ViewGroup(ctx app.Context, e app.Event) {
	osusu.SetIsGroupNew(false, ctx)
	compo.Navigate("/group", ctx)
}
