package main

import (
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// TemporarySessionLength is how long the user can be logged in without checking whether they are authenticated
const TemporarySessionLength = 3 * time.Hour

// RememberMeSessionLength is how long the user's session will be saved with remember me
const RememberMeSessionLength = 90 * 24 * time.Hour

// authenticated is when, if ever, the user has already been authenticated in this session of the app. This information is used to skip unnecessary additional authentication requests in the same session.
var authenticated time.Time

// Authenticate checks whether the user is signed in and takes an action or takes no action based on that. It returns whether the calling function should return.
// If required is set to true, Auth does nothing if the user is signed in and redirects the user to the sign in page otherwise.
// If required is set to false, Auth redirects the user to the home page if the user is signed in, and does nothing otherwise.
func Authenticate(required bool, ctx app.Context) bool {
	ok := time.Since(authenticated) < TemporarySessionLength
	if !ok {
		user := CurrentUser(ctx)
		if user.Session != "" {
			_, err := AuthenticateSessionAPI.Call(user)
			if err == nil {
				ok = true
				authenticated = time.Now()
			}
		}
	}
	switch {
	case required && !ok:
		ctx.Navigate("/signin")
	case !required && ok:
		ctx.Navigate("/home")
	default:
		return false
	}
	return true
}
