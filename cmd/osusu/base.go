package main

import (
	"errors"
	"path/filepath"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/kkoreilly/osusu/osusu"
	"goki.dev/gi/v2/gi"
	"goki.dev/girl/styles"
	"goki.dev/kid"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

func base(b *gi.Body) {
	b.Style(func(s *styles.Style) {
		s.Justify.Content = styles.Center
		s.Align.Content = styles.Center
		s.Align.Items = styles.Center
		s.Text.Align = styles.Center
	})

	gi.NewLabel(b).SetType(gi.LabelDisplayLarge).SetText("Osusu")
	gi.NewLabel(b).SetType(gi.LabelTitleLarge).SetText("An app for getting recommendations on what meals to eat in a group based on the ratings of each member of the group, and the cost, effort, healthiness, and recency of the meal.")

	fun := func(token *oauth2.Token, userInfo *oidc.UserInfo) {
		user := &osusu.User{}
		err := userInfo.Claims(&user)
		if err != nil {
			gi.ErrorDialog(b, err)
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
			gi.ErrorDialog(b, err)
			return
		}
		err = osusu.DB.Create(user).Error
		if err != nil {
			gi.ErrorDialog(b, err)
		}
		curUser = user
		home()
	}
	kid.Buttons(b, &kid.ButtonsConfig{
		SuccessFunc: fun,
		TokenFile: func(provider, email string) string {
			return filepath.Join(gi.AppPrefsDir(), provider+"-token.json")
		},
	}).Style(func(s *styles.Style) {
		s.Grow.Set(0, 0)
	})
}

/*
func loadSession(b *gi.Body) {
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
		gi.ErrorDialog(b, err).Run()
		return
	}
	// sessions expire after 2 weeks
	if time.Since(session.CreatedAt) > 2*7*24*time.Hour {
		err := osusu.DB.Delete(session).Error
		if err != nil {
			gi.ErrorDialog(b, err).Run()
		}
		return
	}
	curUser = &session.User
	home()
}

func saveSession(b *gi.Body) {
	bs := make([]byte, 16)
	rand.Read(bs)
	token := hex.EncodeToString(bs)
	session := &osusu.Session{
		UserID: curUser.ID,
		Token:  token,
	}
	err := osusu.DB.Create(session).Error
	if err != nil {
		gi.ErrorDialog(b, err).Run()
		return
	}
	err = os.WriteFile(filepath.Join(goosi.TheApp.AppPrefsDir(), "sessionToken.json"), []byte(token), 0666)
	if err != nil {
		gi.ErrorDialog(b, err).Run()
	}
}
*/
