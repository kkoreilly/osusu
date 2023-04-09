package main

import (
	"unicode"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Page is the common page structure that all pages have.
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
	statusText             string
	statusType             StatusType
	updateAvailable        bool
	installAvailable       bool
	user                   User
	person                 Person
}

// CurrentPage is the current page the user is on
var CurrentPage *Page

// Render returns the UI of the page based on its attributes
func (p *Page) Render() app.UI {
	width, _ := app.Window().Size()
	smallScreen := width <= 480
	nameFirstLetter := ""
	if len(p.person.Name) > 0 {
		nameFirstLetter = string(unicode.ToUpper(rune(p.person.Name[0])))
	}
	return app.Div().ID(p.ID + "-page-container").Class("page-container").Body([]app.UI{
		app.Header().ID(p.ID+"-page-header").Class("page-header").Body(
			app.Div().ID(p.ID+"-page-top-bar").Class("page-top-bar").Body(
				app.A().ID(p.ID+"-page-top-bar-icon-link").Class("page-top-bar-icon-link").Href("/").Title("Navigate to the Satisfi Start/Home Page").Body(
					app.Img().ID(p.ID+"-page-top-bar-icon-img").Class("page-top-bar-icon-img").Src("/web/images/icon-192.png"),
					app.If(!smallScreen, app.Span().ID(p.ID+"-page-top-bar-icon-text").Class("page-top-bar-icon-text").Text("Satisfi")),
				),

				app.Div().ID(p.ID+"-page-top-bar-buttons").Class("page-top-bar-buttons").Body(
					app.If(p.updateAvailable, app.Button().ID(p.ID+"-page-top-bar-update-button").Class("page-top-bar-button", "page-top-bar-update-button").Text("Update").Title("Update to the Latest Version of Satisfi").OnClick(p.UpdateApp)),
					app.If(p.installAvailable, app.Button().ID(p.ID+"-page-top-bar-install-button").Class("page-top-bar-button", "page-top-bar-install-button").Text("Install").Title("Install Satisfi to Your Device").OnClick(p.InstallApp)),
					app.A().ID(p.ID+"-page-top-bar-account-button").Class("page-top-bar-button", "page-top-bar-account-button").Href("/account").Text("Account").Title("View and Change Account Information").Body(
						app.Span().ID(p.ID+"-page-top-bar-account-button-text").Class("page-top-bar-account-button-text").Text(nameFirstLetter),
					),
				),
			),
			app.Dialog().ID(p.ID+"-page-status").Class("page-status").DataSet("status-type", p.statusType).Body(
				app.Span().ID(p.ID+"-page-status-text").Class("page-status-text").Text(p.statusText),
				app.Button().ID(p.ID+"-page-status-close-button").Class("page-status-close-button").Text("✕").OnClick(p.ClosePageStatus),
			),
			app.If(p.TitleElement != "", app.H1().ID(p.ID+"-page-title").Class("page-title").Text(p.TitleElement)),
			app.If(p.SubtitleElement != "", app.P().ID(p.ID+"-page-subtitle").Class("page-subtitle").Text(p.SubtitleElement)),
		),
		app.Main().ID(p.ID + "-page-main").Class("page-main").Body(
			p.Elements...,
		),
		app.Footer().ID(p.ID+"-page-footer").Class("page-footer").Body(
			app.Span().Text("Copyright © 2023, Satisfi"),
			app.A().Href("https://www.flaticon.com/free-icons/pizza").Title("pizza icons").Text("Pizza icons created by Freepik - Flaticon"),
		),
	}...)
}

// ShowStatus shows the page status dialog with the given status text with the given status type
func (p *Page) ShowStatus(text string, statusType StatusType) {
	p.statusText = text
	p.statusType = statusType
	app.Window().GetElementByID(p.ID + "-page-status").Call("show")
}

// ClosePageStatus closes the page status dialog
func (p *Page) ClosePageStatus(ctx app.Context, event app.Event) {
	app.Window().GetElementByID(p.ID + "-page-status").Call("close")
}

// OnNav is called when the page is navigated to. It authenticates the user, sets the title and description, sets the update and install states, and calls the specified OnNav function.
func (p *Page) OnNav(ctx app.Context) {
	if Authenticate(p.AuthenticationRequired, ctx) {
		return
	}
	CurrentPage = p
	ctx.Page().SetTitle("Satisfi | " + p.Title)
	ctx.Page().SetDescription(p.Description)
	p.updateAvailable = ctx.AppUpdateAvailable()
	p.installAvailable = ctx.IsAppInstallable()
	p.user = GetCurrentUser(ctx)
	p.person = GetCurrentPerson(ctx)
	// if not signed in but person still set, clear person
	if p.user.Username == "" && p.person != (Person{}) {
		p.person = Person{}
		SetCurrentPerson(p.person, ctx)
	}

	if p.OnNavFunc != nil {
		p.OnNavFunc(ctx)
	}
}

// OnAppUpdate is called when the updatability of the app changes
func (p *Page) OnAppUpdate(ctx app.Context) {
	p.updateAvailable = ctx.AppUpdateAvailable()
}

// OnAppInstallChange is called when the installability of the app changes
func (p *Page) OnAppInstallChange(ctx app.Context) {
	p.installAvailable = ctx.IsAppInstallable()
}

// UpdateApp reloads the page to update it after a new version is loaded
func (p *Page) UpdateApp(ctx app.Context, e app.Event) {
	ctx.Reload()
}

// InstallApp shows the app installation prompt
func (p *Page) InstallApp(ctx app.Context, e app.Event) {
	ctx.ShowAppInstallPrompt()
}
