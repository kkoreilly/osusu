package main

import (
	"strconv"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// A Group is a group of users that can determine what to eat together
type Group struct {
	ID       int64
	Owner    int64
	Code     string
	Name     string
	Members  []int64
	Cuisines []string
}

// GroupJoin is the struct that contains data used to have a person join a group
type GroupJoin struct {
	GroupCode string
	UserID    int64
}

// CurrentGroup returns the current group state value
func CurrentGroup(ctx app.Context) Group {
	var group Group
	ctx.GetState("currentGroup", &group)
	return group
}

// SetCurrentGroup sets the current group state value to the given group
func SetCurrentGroup(group Group, ctx app.Context) {
	ctx.SetState("currentGroup", group, app.Persist)
}

// IsGroupNew gets whether the current group is a new group
func IsGroupNew(ctx app.Context) bool {
	var isGroupNew bool
	ctx.GetState("isGroupNew", &isGroupNew)
	return isGroupNew
}

// SetIsGroupNew sets whether the current group is a new group
func SetIsGroupNew(isGroupNew bool, ctx app.Context) {
	ctx.SetState("isGroupNew", isGroupNew, app.Persist)
}

type group struct {
	app.Compo
	group           Group
	isGroupNew      bool
	user            User
	isOwner         bool
	members         Users
	joinLinkClicked bool
}

func (g *group) Render() app.UI {
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
	return &Page{
		ID:                     "group",
		Title:                  titleText,
		Description:            "View, edit, and select a group",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			g.group = CurrentGroup(ctx)
			g.isGroupNew = IsGroupNew(ctx)
			if g.isOwner {
				if g.isOwner {
					titleText = "Edit Group"
				}
				if g.isGroupNew {
					titleText = "Create Group"
				}
				CurrentPage.Title = titleText
				CurrentPage.UpdatePageTitle(ctx)
			}
			g.user = CurrentUser(ctx)
			g.isOwner = g.isGroupNew || g.user.ID == g.group.Owner
			users, err := GetUsersAPI.Call(g.group.Members)
			if err != nil {
				CurrentPage.ShowErrorStatus(err)
				return
			}
			g.members = users
		},
		TitleElement: titleText,
		Elements: []app.UI{
			app.Form().ID("group-page-form").Class("form").OnSubmit(g.OnSubmit).Body(
				app.If(g.isOwner,
					TextInput().ID("group-page-name").Label("Name:").Value(&g.group.Name).AutoFocus(true),
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
				ButtonRow().ID("group-page").Buttons(
					Button().ID("group-page-delete").Class("danger").Icon("delete").Text("Delete").Hidden(!g.isOwner || g.isGroupNew),
					Button().ID("group-page-cancel").Class("secondary").Icon(cancelButtonIcon).Text(cancelButtonText).OnClick(ReturnToReturnURL),
					Button().ID("group-page-save").Class("primary").Type("submit").Icon(saveButtonIcon).Text(saveButtonText).Hidden(!g.isOwner),
				),
			),
		},
	}
}

func (g *group) JoinLinkOnClick(ctx app.Context, e app.Event) {
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

func (g *group) OnSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	if g.isGroupNew {
		g.group.Owner = g.user.ID
		g.group.Members = []int64{g.user.ID}
		group, err := CreateGroupAPI.Call(g.group)
		if err != nil {
			CurrentPage.ShowErrorStatus(err)
			return
		}
		g.group = group
		SetCurrentGroup(g.group, ctx)
		ReturnToReturnURL(ctx, e)
		return
	}
	_, err := UpdateGroupAPI.Call(g.group)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	SetCurrentGroup(g.group, ctx)
	ReturnToReturnURL(ctx, e)
}
