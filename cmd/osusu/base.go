package main

import (
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/kkoreilly/osusu/osusu"
	"goki.dev/gi/v2/gi"
	"goki.dev/kid"
	"golang.org/x/oauth2"
)

func base(sc *gi.Scene) {
	gi.NewLabel(sc).SetType(gi.LabelHeadlineLarge).SetText("Osusu")
	gi.NewLabel(sc).SetText("An app for getting recommendations on what meals to eat in a group based on the ratings of each member of the group, and the cost, effort, healthiness, and recency of the meal.")

	brow := gi.NewLayout(sc)
	kid.Buttons(brow, func(token *oauth2.Token, userInfo *oidc.UserInfo) {
		user := &osusu.User{}
		err := userInfo.Claims(&user)
		if err != nil {
			gi.NewDialog(brow).Title("Error getting user info").Prompt(err.Error()).Run()
			return
		}
		fmt.Println(user)
		osusu.DB.Create(user)
	})
}
