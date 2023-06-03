package page

import (
	"strings"

	"github.com/kkoreilly/osusu/api"
	"github.com/kkoreilly/osusu/compo"
	"github.com/kkoreilly/osusu/osusu"
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

type Join struct {
	app.Compo
	groupCode string
	user      osusu.User
}

func (j *Join) Render() app.UI {
	return &compo.Page{
		ID:                     "join",
		Title:                  "Join Group",
		Description:            "Join a group",
		AuthenticationRequired: true,
		PreOnNavFunc: func(ctx app.Context) {
			joinURL := ctx.Page().URL().String()
			SetJoinURL(joinURL, ctx)
			split := strings.Split(joinURL, "/")
			groupCode := split[len(split)-1]
			if groupCode == "join" {
				groupCode = ""
			}
			j.groupCode = groupCode

		},
		OnNavFunc: func(ctx app.Context) {
			j.user = osusu.CurrentUser(ctx)
		},
		TitleElement: "Join Group",
		Elements: []app.UI{
			app.Form().ID("join-page-form").Class("form").OnSubmit(j.OnSubmit).Body(
				compo.TextInput().ID("join-page-form-join-url").Label("Join Code:").Value(&j.groupCode).AutoFocus(true),
				compo.ButtonRow().ID("join-page").Buttons(
					compo.Button().ID("join-page-cancel").Class("secondary").Icon("cancel").Text("Cancel").OnClick(compo.ReturnToReturnURL),
					compo.Button().ID("join-page-join").Class("primary").Type("submit").Icon("group").Text("Join Group"),
				),
			),
		},
	}
}

func (j *Join) OnSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	split := strings.Split(j.groupCode, "/")
	j.groupCode = split[len(split)-1]

	group, err := api.JoinGroup.Call(osusu.GroupJoin{GroupCode: j.groupCode, UserID: j.user.ID})
	if err != nil {
		compo.CurrentPage.ShowErrorStatus(err)
		return
	}
	osusu.SetCurrentGroup(group, ctx)
	ctx.Navigate("/home")
}
