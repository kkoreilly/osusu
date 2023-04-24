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
		},
	}
}
