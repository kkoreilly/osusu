package main

import (
	"errors"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/kkoreilly/osusu/osusu"
	"goki.dev/gi/v2/gi"
	"goki.dev/kid"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
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
		var oldUser osusu.User
		err = osusu.DB.First(&oldUser, "email = ?", user.Email).Error
		// if we already have a user with the same email, we don't need to make a new account
		if err == nil {
			curUser = &oldUser
			home()
			return
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			gi.NewDialog(brow).Title("Error checking for existing user").Prompt(err.Error()).Run()
			return
		}
		err = osusu.DB.Create(user).Error
		if err != nil {
			gi.NewDialog(brow).Title("Error creating new user").Prompt(err.Error()).Run()
		}
		curUser = user
		home()
	})
}
