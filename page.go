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
	SubtitleElement        string
	Elements               []app.UI
	updateAvailable        bool
	installAvailable       bool
}

func (p *Page) Render() app.UI {
	width, _ := app.Window().Size()
	smallScreen := width <= 480
	return app.Div().ID(p.ID + "-page-container").Class("page-container").Body([]app.UI{
		app.Header().ID(p.ID+"-page-header").Class("page-header").Body(
			app.Div().ID(p.ID+"-page-top-bar").Class("page-top-bar").Body(
				app.A().ID(p.ID+"-page-top-bar-icon-link").Class("page-top-bar-icon-link").Href("/").Title("Navigate to the MealRec Start/Home Page").Body(
					app.Img().ID(p.ID+"-page-top-bar-icon-img").Class("page-top-bar-icon-img").Src("/web/images/icon-192.png"),
					app.If(!smallScreen, app.Span().ID(p.ID+"-page-top-bar-icon-text").Class("page-top-bar-icon-text").Text("MealRec")),
				),
				app.Div().ID(p.ID+"-page-top-bar-buttons").Class("page-top-bar-buttons").Body(
					app.If(p.updateAvailable, app.Button().ID(p.ID+"-page-top-bar-update-button").Class("page-top-bar-button", "page-top-bar-update-button").Text("Update").Title("Update to the Latest Version of MealRec").OnClick(p.UpdateApp)),
					app.If(p.installAvailable, app.Button().ID(p.ID+"-page-top-bar-install-button").Class("page-top-bar-button", "page-top-bar-install-button").Text("Install").Title("Install MealRec to Your Device").OnClick(p.InstallApp)),
				),
			),
			app.If(p.TitleElement != "", app.H1().ID(p.ID+"-page-title").Class("page-title").Text(p.TitleElement)),
			app.If(p.SubtitleElement != "", app.P().ID(p.ID+"-page-subtitle").Class("page-subtitle").Text(p.SubtitleElement)),
		),
		app.Main().ID(p.ID + "-page-main").Class("page-main").Body(
			p.Elements...,
		),
		app.Footer().ID(p.ID+"-page-footer").Class("page-footer").Body(
			app.Span().Text("Copyright © 2023, MealRec"),
			app.A().Href("https://www.flaticon.com/free-icons/pizza").Title("pizza icons").Text("Pizza icons created by Freepik - Flaticon"),
		),
	}...)
}

func (p *Page) OnNav(ctx app.Context) {
	if Authenticate(p.AuthenticationRequired, ctx) {
		return
	}
	ctx.Page().SetTitle("MealRec | " + p.Title)
	ctx.Page().SetDescription(p.Description)
	p.updateAvailable = ctx.AppUpdateAvailable()
	p.installAvailable = ctx.IsAppInstallable()
	if p.OnNavFunc != nil {
		p.OnNavFunc(ctx)
	}
}

func (p *Page) OnAppUpdate(ctx app.Context) {
	p.updateAvailable = ctx.AppUpdateAvailable()
}

func (p *Page) OnAppInstallChange(ctx app.Context) {
	p.installAvailable = ctx.IsAppInstallable()
}

func (p *Page) UpdateApp(ctx app.Context, e app.Event) {
	ctx.Reload()
}

func (p *Page) InstallApp(ctx app.Context, e app.Event) {
	ctx.ShowAppInstallPrompt()
}
