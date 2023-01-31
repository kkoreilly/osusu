package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type home struct {
	app.Compo
}

func (h *home) Render() app.UI {
	return app.H1().ID("home-page-title").Class("page-title").Text("Home")
}
