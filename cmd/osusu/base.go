package main

import (
	"fmt"

	"github.com/kkoreilly/osusu/auth"
	"goki.dev/gi/v2/gi"
	"golang.org/x/oauth2"
)

func base(sc *gi.Scene) {
	gi.NewLabel(sc).SetType(gi.LabelHeadlineLarge).SetText("Osusu")
	gi.NewLabel(sc).SetText("An app for getting recommendations on what meals to eat in a group based on the ratings of each member of the group, and the cost, effort, healthiness, and recency of the meal.")

	brow := gi.NewLayout(sc)
	auth.Buttons(brow, func(token *oauth2.Token) {
		fmt.Println("got token", token)
	})
}
