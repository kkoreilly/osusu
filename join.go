package main

import (
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// JoinURL returns the current join URL state value
func JoinURL(ctx app.Context) string {
	var joinURL string
	ctx.GetState("joinURL", &joinURL)
	return joinURL
}

// SetJoinURL sets the current join URL state value
func SetJoinURL(joinURL string, ctx app.Context) {
	ctx.SetState("joinURL", joinURL)
}

type join struct {
	app.Compo
	joinURL string
	user    User
}

func (j *join) Render() app.UI {
	return &Page{
		ID:                     "join",
		Title:                  "Join Group",
		Description:            "Join a group",
		AuthenticationRequired: true,
		PreOnNavFunc: func(ctx app.Context) {
			j.joinURL = ctx.Page().URL().String()
			SetJoinURL(j.joinURL, ctx)
		},
		OnNavFunc: func(ctx app.Context) {
			j.user = CurrentUser(ctx)
		},
		TitleElement: "Join Group",
		Elements: []app.UI{
			app.Form().ID("join-page-form").Class("form").OnSubmit(j.OnSubmit).Body(
				TextInput().ID("join-page-form-join-url").Label("Join URL:").Value(&j.joinURL).AutoFocus(true),
				ButtonRow().ID("join-page").Buttons(
					Button().ID("join-page-cancel").Class("secondary").Icon("cancel").Text("Cancel").OnClick(ReturnToReturnURL),
					Button().ID("join-page-join").Class("primary").Type("submit").Icon("group").Text("Join Group"),
				),
			),
		},
	}
}

func (j *join) OnSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	split := strings.Split(j.joinURL, "/")
	groupCode := split[len(split)-1]
	group, err := JoinGroupAPI.Call(GroupJoin{GroupCode: groupCode, UserID: j.user.ID})
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	SetCurrentGroup(group, ctx)
	ctx.Navigate("/home")
}
