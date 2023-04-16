package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// User is a struct that represents a user in the users database
type User struct {
	ID         int64
	Username   string
	Password   string
	Cuisines   []string
	Session    string // session id, not part of user in user database, but stored locally
	RememberMe bool   // also not part of user database, but used to transmit whether to save session
}

// SetCurrentUser sets the value of the current user in local storage
func SetCurrentUser(user User, ctx app.Context) {
	if user.RememberMe {
		ctx.SetState("currentUser", user, app.Persist, app.Encrypt, app.ExpiresIn(RememberMeSessionLength))
		return
	}
	ctx.SetState("currentUser", user, app.ExpiresIn(TemporarySessionLength))
}

// GetCurrentUser gets the value of the current user from local storage
func GetCurrentUser(ctx app.Context) User {
	var user User
	ctx.GetState("currentUser", &user)
	return user
}
