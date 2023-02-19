package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

// Authenticate checks whether the user is signed in and takes an action or takes no action based on that. It returns whether the calling function should return.
// If required is set to true, Auth does nothing if the user is signed in and redirects the user to the sign in page otherwise.
// If required is set to false, Auth redirects the user to the home page if the user is signed in, and does nothing otherwise.
func Authenticate(required bool, ctx app.Context) bool {
	err := AutoSignIn(ctx)
	switch {
	case required && err != nil:
		ctx.Navigate("/signin")
	case !required && err == nil:
		ctx.Navigate("/home")
	default:
		return false
	}
	return true
}
