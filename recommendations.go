package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type recommendations struct {
	app.Compo
}

func (r *recommendations) Render() app.UI {
	return app.Div().Body(
		app.H1().ID("recommendations-page-title").Class("page-title").Text("Recommendations"),
	)
}
