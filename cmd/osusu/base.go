package main

import (
	"github.com/kkoreilly/osusu/auth"
	"goki.dev/gi/v2/gi"
	"goki.dev/goosi/events"
)

func base(sc *gi.Scene) {
	gi.NewLabel(sc).SetType(gi.LabelHeadlineLarge).SetText("Osusu")
	gi.NewLabel(sc).SetText("An app for getting recommendations on what meals to eat in a group based on the ratings of each member of the group, and the cost, effort, healthiness, and recency of the meal.")

	brow := gi.NewLayout(sc)
	gi.NewButton(brow).SetText("Sign in").OnClick(func(e events.Event) {
		signIn(sc)
	})
	gi.NewButton(brow).SetText("Sign up").OnClick(func(e events.Event) {
		signUp(sc)
	})
	auth.GoogleButton(brow)
}
