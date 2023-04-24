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
	viewGroupIcon := "visibility"
	viewGroupText := "View Group"
	if a.user.ID == a.group.Owner {
		viewGroupIcon = "edit"
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
			ButtonRow().ID("account-page-top").Buttons(
				Button().ID("account-page-sign-out").Class("danger").Icon("logout").Text("Sign Out").OnClick(a.InitialSignOut),
				Button().ID("account-page-change-group").Class("secondary").Icon("group").Text("Change Group").OnClick(NavigateEvent("/groups")),
				Button().ID("account-page-view-group").Class("primary").Icon(viewGroupIcon).Text(viewGroupText).OnClick(a.ViewGroup),
			),
			app.H2().ID("account-page-user-info-subtitle").Text("Change User Information:"),
			app.Form().ID("account-page-user-info-form").Class("form").OnSubmit(a.ChangeUserInfo).Body(
				TextInput().ID("account-page-username").Label("Username:").Value(&a.user.Username),
				TextInput().ID("account-page-name").Label("Name:").Value(&a.user.Name),
				Button().ID("account-page-user-info-save").Class("primary").Type("submit").Icon("save").Text("Save"),
			),
			app.H2().ID("account-page-password-subtitle").Text("Change Password:"),
			app.Form().ID("account-page-password-form").Class("form").OnSubmit(a.ChangePassword).Body(
				TextInput().ID("account-page-password").Type("password").Label("Password:").Value(&a.user.Password),
				Button().ID("account-page-password-save").Class("tertiary").Type("submit").Icon("save").Text("Save"),
			),
			app.Dialog().ID("account-page-confirm-sign-out").Class("modal").Body(
				app.P().ID("account-page-confirm-sign-out-text").Class("confirm-delete-text").Text("Are you sure you want to sign out?"),
				ButtonRow().ID("account-page-confirm-sign-out").Buttons(
					Button().ID("account-page-confirm-sign-out-sign-out").Class("danger").Icon("logout").Text("Yes, Sign Out").OnClick(a.ConfirmSignOut),
					Button().ID("account-page-confirm-sign-out-cancel").Class("secondary").Icon("cancel").Text("No, Cancel").OnClick(a.CancelSignOut),
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

	Navigate("/signin", ctx)
}

func (a *account) CancelSignOut(ctx app.Context, e app.Event) {
	e.PreventDefault()
	app.Window().GetElementByID("account-page-confirm-sign-out").Call("close")
}

func (a *account) ChangeUserInfo(ctx app.Context, e app.Event) {
	e.PreventDefault()
	_, err := UpdateUserInfoAPI.Call(a.user)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	CurrentPage.ShowStatus("User Info Updated!", StatusTypePositive)
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
	Navigate("/group", ctx)
}
