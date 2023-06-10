package page

import (
	"strconv"
	"time"

	"github.com/kkoreilly/osusu/api"
	"github.com/kkoreilly/osusu/compo"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Group struct {
	app.Compo
	group           osusu.Group
	isGroupNew      bool
	user            osusu.User
	isOwner         bool
	members         osusu.Users
	joinLinkClicked bool
}

func (g *Group) Render() app.UI {
	titleText := "Edit Group"
	saveButtonIcon := "save"
	saveButtonText := "Save"
	if g.isGroupNew {
		titleText = "Create Group"
		saveButtonIcon = "add"
		saveButtonText = "Create"
	}
	cancelButtonIcon := "cancel"
	cancelButtonText := "Cancel"
	if !g.isOwner {
		cancelButtonIcon = "arrow_back"
		cancelButtonText = "Back"
		titleText = "View Group"
	}
	joinLinkText := "https://osusu.fly.dev/join/" + g.group.Code
	if g.joinLinkClicked {
		joinLinkText = "Link Copied!"
	}
	return &compo.Page{
		ID:                     "group",
		Title:                  titleText,
		Description:            "View, edit, and select a group",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			g.group = osusu.CurrentGroup(ctx)
			g.isGroupNew = osusu.IsGroupNew(ctx)
			if g.isOwner {
				if g.isOwner {
					titleText = "Edit Group"
				}
				if g.isGroupNew {
					titleText = "Create Group"
				}
				compo.CurrentPage.Title = titleText
				compo.CurrentPage.UpdatePageTitle(ctx)
			}
			g.user = osusu.CurrentUser(ctx)
			g.isOwner = g.isGroupNew || g.user.ID == g.group.Owner
			users, err := api.GetUsers.Call(g.group.Members)
			if err != nil {
				compo.CurrentPage.ShowErrorStatus(err)
				return
			}
			g.members = users
		},
		TitleElement: titleText,
		Elements: []app.UI{
			app.Form().ID("group-page-form").Class("form").OnSubmit(g.OnSubmit).Body(
				app.If(g.isOwner,
					compo.TextInput().ID("group-page-name").Label("Name:").Value(&g.group.Name).AutoFocus(true),
				).Else(
					app.Span().ID("group-page-name").Text("Name: "+g.group.Name),
				),
				app.If(!g.isGroupNew,
					app.Span().ID("group-page-join-link-text").Text("Join Link (click to copy):"),
					app.Span().ID("group-page-join-link").Text(joinLinkText).DataSet("clicked", g.joinLinkClicked).OnClick(g.JoinLinkOnClick),
					app.Span().ID("group-page-members-label").Text("Group Members:"),
					app.Table().ID("group-page-members-table").Body(
						app.THead().ID("group-page-members-table-header").Body(
							app.Tr().ID("group-page-members-table-header-row").Body(
								app.Th().ID("group-page-members-table-header-name").Text("Name:"),
								app.Th().ID("group-page-members-table-header-username").Text("Username:"),
							),
						),
						app.TBody().Body(
							app.Range(g.members).Slice(func(i int) app.UI {
								user := g.members[i]
								si := strconv.Itoa(i)
								isOwner := user.ID == g.group.Owner
								return app.Tr().ID("group-page-member-"+si).Class("group-page-member").DataSet("is-owner", isOwner).Body(
									app.Td().ID("group-page-member-name-"+si).Class("group-page-member-name").Text(user.Name),
									app.Td().ID("group-page-member-username-"+si).Class("group-page-member-username").Text(user.Username),
								)
							}),
						),
					),
				),
				compo.ButtonRow().ID("group-page").Buttons(
					compo.Button().ID("group-page-delete").Class("danger").Icon("delete").Text("Delete").Hidden(!g.isOwner || g.isGroupNew),
					compo.Button().ID("group-page-cancel").Class("secondary").Icon(cancelButtonIcon).Text(cancelButtonText).OnClick(compo.ReturnToReturnURL),
					compo.Button().ID("group-page-save").Class("primary").Type("submit").Icon(saveButtonIcon).Text(saveButtonText).Hidden(!g.isOwner),
				),
			),
		},
	}
}

func (g *Group) JoinLinkOnClick(ctx app.Context, e app.Event) {
	if g.joinLinkClicked || app.Window().Get("navigator").Get("clipboard").Truthy() {
		e.PreventDefault()
		g.joinLinkClicked = true
		app.Window().Get("navigator").Get("clipboard").Call("writeText", "https://osusu.fly.dev/join/"+g.group.Code)
		ctx.Defer(func(ctx app.Context) {
			time.Sleep(1 * time.Second)
			g.joinLinkClicked = false
			g.Update()
		})
	}
}

func (g *Group) OnSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	if g.isGroupNew {
		g.group.Owner = g.user.ID
		g.group.Members = []int64{g.user.ID}
		group, err := api.CreateGroup.Call(g.group)
		if err != nil {
			compo.CurrentPage.ShowErrorStatus(err)
			return
		}
		g.group = group
		osusu.SetCurrentGroup(g.group, ctx)
		compo.ReturnToReturnURL(ctx, e)
		return
	}
	_, err := api.UpdateGroup.Call(g.group)
	if err != nil {
		compo.CurrentPage.ShowErrorStatus(err)
		return
	}
	osusu.SetCurrentGroup(g.group, ctx)
	compo.ReturnToReturnURL(ctx, e)
}
