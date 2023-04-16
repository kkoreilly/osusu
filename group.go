package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

// A Group is a group of users that can determine what to eat together
type Group struct {
	ID      int64
	Owner   int64
	Code    string
	Name    string
	Members []int64
}

// GroupJoin is the struct that contains data used to have a person join a group
type GroupJoin struct {
	GroupCode string
	UserID    int64
}

// SetCurrentGroup sets the current group state value to the given group
func SetCurrentGroup(group Group, ctx app.Context) {
	ctx.SetState("currentGroup", group, app.Persist)
}

// GetCurrentGroup returns the current group state value
func GetCurrentGroup(ctx app.Context) Group {
	var group Group
	ctx.GetState("currentGroup", &group)
	return group
}

// SetIsGroupNew sets whether the current group is a new group
func SetIsGroupNew(isGroupNew bool, ctx app.Context) {
	ctx.SetState("isGroupNew", isGroupNew, app.Persist)
}

// GetIsGroupNew gets whether the current group is a new group
func GetIsGroupNew(ctx app.Context) bool {
	var isGroupNew bool
	ctx.GetState("isGroupNew", &isGroupNew)
	return isGroupNew
}

type group struct {
	app.Compo
	group      Group
	isGroupNew bool
	user       User
	isOwner    bool
}

func (g *group) Render() app.UI {
	titleText := "Edit Group"
	if g.isGroupNew {
		titleText = "Create Group"
	}
	cancelButtonText := "Cancel"
	if !g.isOwner {
		cancelButtonText = "Back"
		titleText = "View Group"
	}
	return &Page{
		ID:                     "group",
		Title:                  titleText,
		Description:            "View, edit, and select a group",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			g.group = GetCurrentGroup(ctx)
			g.isGroupNew = GetIsGroupNew(ctx)
			if g.isGroupNew || !g.isOwner {
				CurrentPage.Title = titleText
				CurrentPage.UpdatePageTitle(ctx)
			}
			g.user = GetCurrentUser(ctx)
			g.isOwner = g.isGroupNew || g.user.ID == g.group.Owner
			g.isGroupNew = false
			g.isOwner = false
		},
		TitleElement: titleText,
		Elements: []app.UI{
			app.Form().ID("group-page-form").Class("form").Body(
				app.If(g.isOwner,
					NewTextInput("group-page-name", "Name:", "Group Name", true, &g.group.Name),
				).Else(
					app.Span().ID("group-page-name").Text("Name: "+g.group.Name),
				),
				app.If(!g.isGroupNew,
					app.Span().ID("group-page-join-link-text").Body(
						app.Text("Join Link: "),
						app.A().ID("group-page-join-link").Href("https://osusu.fly.dev/join/"+g.group.Code).Text("osusu.fly.dev/join/"+g.group.Code),
					),
				),
				app.Div().ID("group-page-action-button-row").Class("action-button-row").Body(
					app.If(g.isOwner && !g.isGroupNew,
						app.Button().ID("group-page-delete-button").Class("action-button", "danger-action-button").Text("Delete"),
					),
					app.Button().ID("grou-page-cancel-button").Class("action-button", "secondary-action-button").Text(cancelButtonText),
					app.If(g.isOwner && !g.isGroupNew,
						app.Button().ID("group-page-save-button").Class("action-button", "primary-action-button").Text("Save"),
					),
					app.If(g.isGroupNew,
						app.Button().ID("group-page-create-button").Class("action-button", "primary-action-button").Text("Create"),
					),
				),
			),
		},
	}
}
