package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type person struct {
	app.Compo
}

func (p *person) Render() app.UI {
	return app.Div().Body(
		app.H1().ID("person-page-title").Class("page-title").Text("Person"),
	)
}
