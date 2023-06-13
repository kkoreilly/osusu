package page

import (
	"strconv"

	"github.com/kkoreilly/osusu/api"
	"github.com/kkoreilly/osusu/compo"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Groups struct {
	app.Compo
	groups osusu.Groups
}

func (g *Groups) Render() app.UI {
	return &compo.Page{
		ID:                     "groups",
		Title:                  "Groups",
		Description:            "Select your group",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			compo.SetReturnURL("/groups", ctx)
			groups, err := api.GetGroups.Call(osusu.CurrentUser(ctx).ID)
			if err != nil {
				compo.CurrentPage.ShowErrorStatus(err)
			}
			g.groups = groups
		},
		TitleElement: "Select a Group",
		Elements: []app.UI{
			compo.ButtonRow().ID("groups-page").Buttons(
				compo.Button().ID("groups-page-new-group").Class("secondary").Icon("add").Text("New Group").OnClick(g.New),
				compo.Button().ID("groups-page-join-group").Class("primary").Icon("group").Text("Join Group").OnClick(compo.NavigateEvent("/join")),
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

func (g *Groups) New(ctx app.Context, e app.Event) {
	osusu.SetIsGroupNew(true, ctx)
	osusu.SetCurrentGroup(osusu.Group{}, ctx)
	compo.Navigate("/group", ctx)
}

func (g *Groups) GroupOnClick(ctx app.Context, e app.Event, group osusu.Group) {
	osusu.SetIsGroupNew(false, ctx)
	osusu.SetCurrentGroup(group, ctx)
	compo.Navigate("/search", ctx)
}
