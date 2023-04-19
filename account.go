package main

import (
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type account struct {
	app.Compo
	group Group
	user  User
}

func (a *account) Render() app.UI {
	viewGroupText := "View Group"
	if a.user.ID == a.group.Owner {
		viewGroupText = "Edit Group"
	}
	return &Page{
		ID:                     "account",
		Title:                  "Account",
		Description:            "View and change account information",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			SetReturnURL("/account", ctx)
			a.group = GetCurrentGroup(ctx)
			a.user = GetCurrentUser(ctx)
		},
		TitleElement:    "Account",
		SubtitleElement: "You are currently signed into " + a.user.Username + " with the name " + a.user.Name + " and the group " + a.group.Name + ".",
		Elements: []app.UI{
			app.Div().ID("acount-page-top-action-button-row").Class("action-button-row").Body(
				app.Button().ID("account-page-sign-out-button").Class("danger-action-button", "action-button").Text("Sign Out").OnClick(a.InitialSignOut),
				app.A().ID("account-page-change-group-button").Class("secondary-action-button", "action-button").Href("/groups").Text("Change Group"),
				app.Button().ID("account-page-view-group-button").Class("primary-action-button", "action-button").Text(viewGroupText).OnClick(a.ViewGroup),
			),
			app.Form().ID("account-page-username-form").Class("form").OnSubmit(a.ChangeUsername).Body(
				NewTextInput("account-page-username", "Change Your Username:", "", false, &a.user.Username),
				app.Div().ID("account-page-username-action-button-row").Class("action-button-row").Body(
					app.Button().ID("account-page-username-save-button").Class("primary-action-button", "action-button").Type("submit").Text("Save Username"),
				),
			),
			app.Form().ID("account-page-password-form").Class("form").OnSubmit(a.ChangePassword).Body(
				NewTextInput("account-page-password", "Change Your Password:", "••••••••", false, &a.user.Password).SetType("password"),
				app.Div().ID("account-page-password-action-button-row").Class("action-button-row").Body(
					app.Button().ID("account-page-password-save-button").Class("tertiary-action-button", "action-button").Type("submit").Text("Save Password"),
				),
			),
			app.Dialog().ID("account-page-confirm-sign-out").Body(
				app.P().ID("account-page-confirm-sign-out-text").Class("confirm-delete-text").Text("Are you sure you want to sign out?"),
				app.Div().ID("account-page-confirm-sign-out-action-button-row").Class("action-button-row").Body(
					app.Button().ID("account-page-confirm-sign-out-sign-out").Class("action-button", "danger-action-button").Text("Yes, Sign Out").OnClick(a.ConfirmSignOut),
					app.Button().ID("account-page-confirm-sign-out-cancel").Class("action-button", "secondary-action-button").Text("No, Cancel").OnClick(a.CancelSignOut),
				),
			),
		},
	}
}

func (a *account) InitialSignOut(ctx app.Context, e app.Event) {
	e.PreventDefault()
	app.Window().GetElementByID("account-page-confirm-sign-out").Call("showModal")
}

func (a *account) ConfirmSignOut(ctx app.Context, e app.Event) {
	e.PreventDefault()
	user := GetCurrentUser(ctx)
	if user.Session != "" {
		_, err := SignOutAPI.Call(GetCurrentUser(ctx))
		if err != nil {
			CurrentPage.ShowErrorStatus(err)
			return
		}
	}
	// if no error, we are no longer authenticated
	authenticated = time.UnixMilli(0)
	ctx.LocalStorage().Del("currentUser")
	ctx.LocalStorage().Del("currentGroup")

	ctx.Navigate("/signin")
}

func (a *account) CancelSignOut(ctx app.Context, e app.Event) {
	e.PreventDefault()
	app.Window().GetElementByID("account-page-confirm-sign-out").Call("close")
}

func (a *account) ChangeUsername(ctx app.Context, e app.Event) {
	e.PreventDefault()
	_, err := UpdateUsernameAPI.Call(a.user)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	CurrentPage.ShowStatus("Username Updated!", StatusTypePositive)
	SetCurrentUser(a.user, ctx)
}

func (a *account) ChangePassword(ctx app.Context, e app.Event) {
	e.PreventDefault()

	CurrentPage.ShowStatus("Loading...", StatusTypeNeutral)

	ctx.Defer(func(ctx app.Context) {
		_, err := UpdatePasswordAPI.Call(a.user)
		if err != nil {
			CurrentPage.ShowErrorStatus(err)
			a.Update()
			a.user.Password = ""
			return
		}
		CurrentPage.ShowStatus("Password Updated!", StatusTypePositive)
		a.Update()
		a.user.Password = ""
		SetCurrentUser(a.user, ctx)
	})
}

func (a *account) ViewGroup(ctx app.Context, e app.Event) {
	SetIsGroupNew(false, ctx)
	ctx.Navigate("/group")
}
