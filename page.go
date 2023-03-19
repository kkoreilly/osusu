package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Page struct {
	app.Compo
	ID                     string
	Title                  string
	Description            string
	AuthenticationRequired bool
	OnNavFunc              func(ctx app.Context)
	TitleElement           string
	Elements               []app.UI
}

func (p *Page) Render() app.UI {
	width, _ := app.Window().Size()
	smallScreen := width <= 480
	return app.Div().ID(p.ID + "-page-container").Class("page-container").Body(append(
		[]app.UI{
			app.Div().ID(p.ID+"-page-top-bar").Class("page-top-bar").Body(
				app.A().ID(p.ID+"-page-top-bar-icon-link").Class("page-top-bar-icon-link").Href("/").Body(
					app.Img().ID(p.ID+"-page-top-bar-icon-img").Class("page-top-bar-icon-img").Src("/web/images/icon-192.png"),
					app.If(!smallScreen, app.Span().ID(p.ID+"-page-top-bar-icon-text").Class("page-top-bar-icon-text").Text("MealRec")),
				),
				app.If(p.TitleElement != "", app.H1().ID(p.ID+"-page-title").Class("page-title").Text(p.TitleElement)),
				app.Div().ID(p.ID+"-page-top-bar-buttons").Class("page-top-bar-buttons").Body(
					app.Button().ID(p.ID+"-page-top-bar-update-button").Class("page-top-bar-button", "page-top-bar-update-button").Text("Update"),
					app.Button().ID(p.ID+"-page-top-bar-install-button").Class("page-top-bar-button", "page-top-bar-install-button").Text("Install"),
				),
			),
		},
		p.Elements...,
	)...)
}

func (p *Page) OnNav(ctx app.Context) {
	if Authenticate(p.AuthenticationRequired, ctx) {
		return
	}
	ctx.Page().SetTitle("MealRec | " + p.Title)
	ctx.Page().SetDescription(p.Description)
	if p.OnNavFunc != nil {
		p.OnNavFunc(ctx)
	}
}
