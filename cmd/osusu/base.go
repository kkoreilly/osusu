package main

import (
	"errors"
	"path/filepath"

	"cogentcore.org/core/base/auth"
	"cogentcore.org/core/core"
	"cogentcore.org/core/styles"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/kkoreilly/osusu/osusu"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

func base(b *core.Body) {
	b.Styler(func(s *styles.Style) {
		s.Justify.Content = styles.Center
		s.Align.Content = styles.Center
		s.Align.Items = styles.Center
		s.Text.Align = styles.Center
	})

	core.NewText(b).SetType(core.TextDisplayLarge).SetText("Osusu")
	core.NewText(b).SetType(core.TextTitleLarge).SetText("An app for getting recommendations on what meals to eat in a group based on the ratings of each member of the group, and the cost, effort, healthiness, and recency of the meal.")

	fun := func(token *oauth2.Token, userInfo *oidc.UserInfo) {
		user := &osusu.User{}
		err := userInfo.Claims(&user)
		if err != nil {
			core.ErrorDialog(b, err)
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
			core.ErrorDialog(b, err)
			return
		}
		err = osusu.DB.Create(user).Error
		if err != nil {
			core.ErrorDialog(b, err)
		}
		curUser = user
		home()
	}
	auth.Buttons(b, &auth.ButtonsConfig{
		SuccessFunc: fun,
		TokenFile: func(provider, email string) string {
			return filepath.Join(core.TheApp.AppDataDir(), provider+"-token.json")
		},
	}).Styler(func(s *styles.Style) {
		s.Grow.Set(0, 0)
	})
}

/*
func loadSession(b *core.Body) {
	token, err := os.ReadFile(filepath.Join(core.TheApp.AppDataDir(), "sessionToken.json"))
	if err != nil {
		return
	}
	session := &osusu.Session{}
	err = osusu.DB.Preload("User").First(session, "token = ?", string(token)).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	if err != nil {
		core.ErrorDialog(b, err).Run()
		return
	}
	// sessions expire after 2 weeks
	if time.Since(session.CreatedAt) > 2*7*24*time.Hour {
		err := osusu.DB.Delete(session).Error
		if err != nil {
			core.ErrorDialog(b, err).Run()
		}
		return
	}
	curUser = &session.User
	home()
}

func saveSession(b *core.Body) {
	bs := make([]byte, 16)
	rand.Read(bs)
	token := hex.EncodeToString(bs)
	session := &osusu.Session{
		UserID: curUser.ID,
		Token:  token,
	}
	err := osusu.DB.Create(session).Error
	if err != nil {
		core.ErrorDialog(b, err).Run()
		return
	}
	err = os.WriteFile(filepath.Join(core.TheApp.AppDataDir(), "sessionToken.json"), []byte(token), 0666)
	if err != nil {
		core.ErrorDialog(b, err).Run()
	}
}
*/
