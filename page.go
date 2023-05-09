package main

import (
	"time"
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
	PreOnNavFunc           func(ctx app.Context)
	OnNavFunc              func(ctx app.Context)
	OnClick                func(ctx app.Context, e app.Event)
	TitleElement           string
	SubtitleElement        string
	Elements               []app.UI
	loaded                 bool
	statusText             string
	statusType             StatusType
	updateAvailable        bool
	installAvailable       bool
	user                   User
}

// CurrentPage is the current page the user is on
var CurrentPage *Page

// Render returns the UI of the page based on its attributes
func (p *Page) Render() app.UI {
	// We use current page for some things (account, install, and update buttons) to prevent flashing on page switch.
	// If there is no current page (if we haven't been on a page before), just set it to p.
	if CurrentPage == nil {
		CurrentPage = p
	}
	elements := []app.UI{}
	if p.loaded {
		elements = p.Elements
	}
	width, _ := app.Window().Size()
	smallScreen := width <= 480
	installIcon := "install_desktop"
	if smallScreen {
		installIcon = "install_mobile"
	}
	nameFirstLetter := ""
	accountButtonIcon := "person"
	if len(CurrentPage.user.Name) > 0 {
		nameFirstLetter = string(unicode.ToUpper(rune(CurrentPage.user.Name[0])))
		accountButtonIcon = ""
	}
	return app.Div().ID(p.ID+"-page-container").Class("page-container").DataSet("small-screen", smallScreen).OnClick(p.OnClick).Body(
		app.Header().ID(p.ID+"-page-header").Class("page-header").Body(
			app.Div().ID(p.ID+"-page-top-bar").Class("page-top-bar").Body(
				app.Button().ID(p.ID+"-page-top-bar-icon-link").Class("page-top-bar-icon-button").Type("button").OnClick(NavigateEvent("/")).Title("Navigate to the Osusu Start/Home Page").Body(
					app.Img().ID(p.ID+"-page-top-bar-icon-img").Class("page-top-bar-icon-img").Src("/web/images/icon-192.png"),
					app.If(!smallScreen, app.Span().ID(p.ID+"-page-top-bar-icon-text").Class("page-top-bar-icon-text").Text("Osusu")),
				),
				app.Div().ID(p.ID+"-page-top-bar-buttons").Class("page-top-bar-buttons").Body(
					Button().ID(p.ID+"-page-top-bar-update").Class("top-bar").Icon("update").Text("Update").OnClick(p.UpdateApp).Hidden(!CurrentPage.updateAvailable),
					Button().ID(p.ID+"-page-top-bar-install").Class("top-bar").Icon(installIcon).Text("Install").OnClick(p.InstallApp).Hidden(!CurrentPage.installAvailable),
					Button().ID(p.ID+"-page-top-bar-account").Class("top-bar-account").Icon(accountButtonIcon).Text(nameFirstLetter).OnClick(NavigateEvent("/account")),
				),
			),
		),
		app.Main().ID(p.ID+"-page-main").Class("page-main").Body(
			app.Dialog().ID(p.ID+"-page-status").Class("page-status").DataSet("status-type", p.statusType).Body(
				app.Span().ID(p.ID+"-page-status-text").Class("page-status-text").Text(p.statusText),
				Button().ID(p.ID+"page-status-close-button").Class("page-status-close").Icon("close").OnClick(p.ClosePageStatus),
			),
			app.If(p.TitleElement != "", app.H1().ID(p.ID+"-page-title").Class("page-title").Text(p.TitleElement)),
			app.If(p.SubtitleElement != "", app.P().ID(p.ID+"-page-subtitle").Class("page-subtitle").Text(p.SubtitleElement)),
			app.If(true, elements...),
		),
	)
}

// ShowStatus shows the page status dialog with the given status text with the given status type
func (p *Page) ShowStatus(text string, statusType StatusType) {
	p.statusText = text
	p.statusType = statusType
	app.Window().GetElementByID(p.ID + "-page-status").Call("show")
}

// ShowErrorStatus is a shorthand for ShowStatus(err.Error(), StatusTypeNegative)
func (p *Page) ShowErrorStatus(err error) {
	p.ShowStatus(err.Error(), StatusTypeNegative)
}

// ClosePageStatus closes the page status dialog
func (p *Page) ClosePageStatus(ctx app.Context, event app.Event) {
	app.Window().GetElementByID(p.ID + "-page-status").Call("close")
}

// UpdatePageTitle updates the title of the page to the value set in the page's Title field
func (p *Page) UpdatePageTitle(ctx app.Context) {
	ctx.Page().SetTitle(p.Title + " | Osusu")
}

// OnNav is called when the page is navigated to. It calls the specified PreOnNav function, authenticates the user, sets the title and description, sets the update and install states, and calls the specified OnNav function. The specified PreOnNav function is called before all other steps, including authentication.
func (p *Page) OnNav(ctx app.Context) {
	if p.PreOnNavFunc != nil {
		p.PreOnNavFunc(ctx)
	}
	if Authenticate(p.AuthenticationRequired, ctx) {
		return
	}
	CurrentPage = p

	p.UpdatePageTitle(ctx)
	ctx.Page().SetDescription(p.Description)
	p.updateAvailable = ctx.AppUpdateAvailable()
	p.installAvailable = ctx.IsAppInstallable()
	p.user = CurrentUser(ctx)

	if p.OnNavFunc != nil {
		p.OnNavFunc(ctx)
	}

	ctx.Defer(func(ctx app.Context) {
		app.Window().GetElementByID(p.ID+"-page-main").Get("style").Set("opacity", 1)
	})
	p.loaded = true
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

// ReturnURL returns the state value containing the url to return to after exiting a page that can be accessed from multiple places
func ReturnURL(ctx app.Context) string {
	var returnURL string
	ctx.GetState("returnURL", &returnURL)
	return returnURL
}

// SetReturnURL sets the state value of the url to return to after exiting a page that can be accessed from multiple places
func SetReturnURL(returnURL string, ctx app.Context) {
	ctx.SetState("returnURL", returnURL, app.Persist)
}

// Navigate navigates to the given URL using the given context
func Navigate(url string, ctx app.Context) {
	if url == ctx.Page().URL().Path {
		return
	}
	app.Window().GetElementByID(CurrentPage.ID+"-page-main").Get("style").Set("opacity", 0)
	ctx.Defer(func(ctx app.Context) {
		time.Sleep(250 * time.Millisecond)
		ctx.Navigate(url)
	})
}

// NavigateEvent returns an event handler that navigates the user to the given URL
func NavigateEvent(url string) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		Navigate(url, ctx)
	}
}

// ReturnToReturnURL returns the user to the url to return to after exiting a page that can be accessed from multiple places.
// The event value is not used, but it allows this function be used as an event handler.
func ReturnToReturnURL(ctx app.Context, e app.Event) {
	Navigate(ReturnURL(ctx), ctx)
}

// Back navigates to the previous page in history. The context and event values are not used, but they allow this function to be used as an event handler
func Back(ctx app.Context, e app.Event) {
	app.Window().Get("history").Call("back")
}
