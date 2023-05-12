package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type start struct {
	app.Compo
}

func (s *start) Render() app.UI {
	return &Page{
		ID:                     "start",
		Title:                  "Start",
		Description:            "Osusu is an app for getting recommendations on what meals to eat in a group based on the ratings of each member of the group, and the cost, effort, healthiness, and recency of the meal.",
		AuthenticationRequired: false,
		TitleElement:           "Welcome to Osusu!",
		SubtitleElement:        "Osusu is an app for getting recommendations on what meals to eat in a group based on the ratings of each member of the group, and the cost, effort, healthiness, and recency of the meal.",
		Elements: []app.UI{
			ButtonRow().ID("start-page").Buttons(
				Button().ID("start-page-sign-up").Class("secondary").Icon("app_registration").Text("Sign Up").OnClick(NavigateEvent("/signup")),
				Button().ID("start-page-sign-in").Class("primary").Icon("login").Text("Sign In").OnClick(NavigateEvent("/signin")),
			),
			app.Div().ID("start-page-info").Body(
				StartPageInfos([]startPageInfo{
					{id: "recommendations", title: "Get Recommendations", body: "Get a ranked list of meal recommendations scored on various metrics."},
					{id: "options", title: "Customize Options", body: "Customize your options and get recommendations that satisfy your needs."},
					{id: "everyone", title: "Works for Everyone", body: "Find meals that satisfy everyone's preferences and constraints, no matter how picky people are."},
					{id: "history", title: "Track History", body: "Track when you eat meals and how their quality changes over time."},
				}),
			),
		},
	}
}

type startPageInfo struct {
	id    string
	title string
	body  string
}

func StartPageInfos(infos []startPageInfo) app.UI {
	return app.Div().ID("start-page-infos").Body(
		app.Range(infos).Slice(func(i int) app.UI {
			info := infos[i]
			return app.Div().ID("start-page-info-container-"+info.id).Class("start-page-info-container").Body(
				app.H2().ID("start-page-info-title-"+info.id).Class("start-page-info-title").Text(info.title),
				app.P().ID("start-page-info-body-"+info.id).Class("start-page-info-body").Text(info.body),
			)
		}),
	)
}
