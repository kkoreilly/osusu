package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/kkoreilly/osusu/osusu"
	"goki.dev/gi/v2/gi"
	"goki.dev/goosi"
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
			gi.ErrorDialog(brow, err).Run()
			return
		}
		var oldUser osusu.User
		err = osusu.DB.First(&oldUser, "email = ?", user.Email).Error
		// if we already have a user with the same email, we don't need to make a new account
		if err == nil {
			curUser = &oldUser
			saveSession(sc)
			home()
			return
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			gi.ErrorDialog(brow, err).Run()
			return
		}
		err = osusu.DB.Create(user).Error
		if err != nil {
			gi.ErrorDialog(brow, err).Run()
		}
		curUser = user
		saveSession(sc)
		home()
	})
}

func loadSession(sc *gi.Scene) {
	token, err := os.ReadFile(filepath.Join(goosi.TheApp.AppPrefsDir(), "sessionToken.json"))
	if err != nil {
		return
	}
	session := &osusu.Session{}
	err = osusu.DB.Preload("User").First(session, "token = ?", string(token)).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	if err != nil {
		gi.ErrorDialog(sc, err).Run()
		return
	}
	// sessions expire after 2 weeks
	if time.Since(session.CreatedAt) > 2*7*24*time.Hour {
		err := osusu.DB.Delete(session).Error
		if err != nil {
			gi.ErrorDialog(sc, err).Run()
		}
		return
	}
	curUser = &session.User
	home()
}

func saveSession(sc *gi.Scene) {
	b := make([]byte, 16)
	rand.Read(b)
	token := hex.EncodeToString(b)
	session := &osusu.Session{
		UserID: curUser.ID,
		Token:  token,
	}
	err := osusu.DB.Create(session).Error
	if err != nil {
		gi.ErrorDialog(sc, err).Run()
		return
	}
	err = os.WriteFile(filepath.Join(goosi.TheApp.AppPrefsDir(), "sessionToken.json"), []byte(token), 0666)
	if err != nil {
		gi.ErrorDialog(sc, err).Run()
	}
}
