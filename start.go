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
			app.Div().ID("start-page-action-button-row").Class("action-button-row").Body(
				app.Button().ID("start-page-sign-up").Class("action-button", "secondary-action-button").Type("button").OnClick(NavigateEvent("/signup")).Text("Sign Up").Title("Navigate to the Sign Up Page"),
				app.Button().ID("start-page-sign-in").Class("action-button", "primary-action-button").Type("button").OnClick(NavigateEvent("/signin")).Text("Sign In").Title("Navigate to the Sign In Page"),
			),
		},
	}
}
