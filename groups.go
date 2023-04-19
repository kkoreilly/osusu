package main

import (
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Groups is a slice of groups
type Groups []Group

type groups struct {
	app.Compo
	groups Groups
}

func (g *groups) Render() app.UI {
	return &Page{
		ID:                     "groups",
		Title:                  "Groups",
		Description:            "Select your group",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			SetReturnURL("/groups", ctx)
			groups, err := GetGroupsAPI.Call(GetCurrentUser(ctx).ID)
			if err != nil {
				CurrentPage.ShowErrorStatus(err)
			}
			g.groups = groups
		},
		TitleElement: "Select a Group",
		Elements: []app.UI{
			app.Div().ID("groups-page-action-button-row").Class("action-button-row").Body(
				app.Button().ID("groups-page-new-group").Class("action-button", "primary-action-button").Text("New Group").OnClick(g.New),
			),
			app.Div().ID("groups-page-groups-container").Body(
				app.Range(g.groups).Slice(func(i int) app.UI {
					return app.Button().ID("groups-page-group-" + strconv.Itoa(i)).Class("groups-page-group").Text(g.groups[i].Name).
						OnClick(func(ctx app.Context, e app.Event) { g.GroupOnClick(ctx, e, g.groups[i]) })
				}),
			),
		},
	}
}

func (g *groups) New(ctx app.Context, e app.Event) {
	SetIsGroupNew(true, ctx)
	SetCurrentGroup(Group{}, ctx)
	ctx.Navigate("/group")
}

func (g *groups) GroupOnClick(ctx app.Context, e app.Event, group Group) {
	SetIsGroupNew(false, ctx)
	SetCurrentGroup(group, ctx)
	ctx.Navigate("/home")
}
