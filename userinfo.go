package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// SaveUsername saves the user's username
func SaveUsername(username string, ctx app.Context) {
	ctx.SetState("username", username)
}

// GetUsername gets the user's username
func GetUsername(ctx app.Context) string {
	var username string
	ctx.GetState("username", &username)
	return username
}
