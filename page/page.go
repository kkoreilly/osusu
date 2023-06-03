// Package page provides the page type and types for all of the pages of the app
package page

import (
	"net/url"
	"strconv"
	"time"
	"unicode"

	"github.com/kkoreilly/osusu/compo"
	"github.com/kkoreilly/osusu/osusu"
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
	onClick                []app.EventHandler
	TitleElement           string
	SubtitleElement        string
	Elements               []app.UI
	loaded                 bool
	statusText             string
	statusType             osusu.StatusType
	updateAvailable        bool
	installAvailable       bool
	user                   osusu.User
	mode                   string
	dialogElements         app.Value // all of the elements that should be closed on page click
}

// CurrentPage is the current page the user is on
var CurrentPage *Page

// authenticated is when, if ever, the user has already been authenticated in this session of the app. This information is used to skip unnecessary additional authentication requests in the same session.
var authenticated time.Time

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
	// installIcon := "install_desktop"
	// if smallScreen {
	// 	installIcon = "install_mobile"
	// }
	// homeButtonText := "Home"
	// searchText := "Search"
	// historyText := "History"
	// discoverText := "Discover"
	// installText := "Install"
	// updateText := "Update"
	// if smallScreen {
	// 	// homeButtonText = ""
	// 	searchText, historyText, discoverText, installText, updateText = "", "", "", "", ""
	// }
	nameFirstLetter := ""
	accountButtonIcon := "person"
	if len(CurrentPage.user.Name) > 0 {
		nameFirstLetter = string(unicode.ToUpper(rune(CurrentPage.user.Name[0])))
		accountButtonIcon = ""
	}
	// we display nav bar on top on big screen and on bottom on small screen, so we need to have it as separate thing that is easy to insert
	navBar := app.Div().ID(p.ID+"-page-nav-bar").Class("page-nav-bar").Body(
		// Button().ID(p.ID+"-page-nav-bar-home").Class("nav-bar").Icon("home").Text(homeButtonText).OnClick(NavigateEvent("/")),
		// app.Dialog().ID(p.ID+"-page-menu").Class("page-menu").Body(
		// 	Button().ID(p.ID+"-home").Class("page-menu").Icon("home").Text("Home").OnClick(NavigateEvent("/")),
		// ),
		// app.Button().ID(p.ID+"-page-nav-bar-icon-link").Class("page-nav-bar-icon-button").Type("button").OnClick(NavigateEvent("/")).Title("Navigate to the Osusu Start/Home Page").Body(
		// 	app.Img().ID(p.ID+"-page-nav-bar-icon-img").Class("page-nav-bar-icon-img").Src("/web/images/icon-192.png"),
		// 	app.If(!smallScreen, app.Span().ID(p.ID+"-page-nav-bar-icon-text").Class("page-nav-bar-icon-text").Text("Osusu")),
		// ),
		compo.Button().ID(p.ID+"page-nav-bar-search").Class("open-"+strconv.FormatBool(CurrentPage.mode == "Search")+" nav-bar").Icon("search").Text("Search").OnClick(p.NavBarOnClick("/search", "Search")),
		compo.Button().ID(p.ID+"page-nav-bar-discover").Class("open-"+strconv.FormatBool(CurrentPage.mode == "Discover")+" nav-bar").Icon("explore").Text("Discover").OnClick(p.NavBarOnClick("/discover", "Discover")),
		compo.Button().ID(p.ID+"page-nav-bar-history").Class("open-"+strconv.FormatBool(CurrentPage.mode == "History")+" nav-bar").Icon("history").Text("History").OnClick(p.NavBarOnClick("/history", "History")),
		// app.If(false, Button().ID(p.ID+"-page-nav-bar-update").Class("nav-bar").Icon("update").Text(updateText).OnClick(p.UpdateApp).Hidden(!CurrentPage.updateAvailable),
		// 	Button().ID(p.ID+"-page-nav-bar-install").Class("nav-bar").Icon(installIcon).Text(installText).OnClick(p.InstallApp).Hidden(!CurrentPage.installAvailable),
		// 	Button().ID(p.ID+"-page-nav-bar-account").Class("nav-bar-account").Icon(accountButtonIcon).Text(nameFirstLetter).OnClick(NavigateEvent("/account"))),

		// app.Div().ID(p.ID+"-page-nav-bar-buttons").Class("page-nav-bar-buttons").Body(),
	)
	return app.Div().ID(p.ID+"-page-container").Class("page-container").DataSet("small-screen", smallScreen).OnClick(p.OnClickEvent).Body(
		app.Header().ID(p.ID+"-page-header").Class("page-header").Body(
		// app.Div().ID(p.ID+"-page-top-bar").Class("page-top-bar", "page-nav-bar").Body(
		// 	Button().ID(p.ID+"-page-top-bar-update").Class("nav-bar").Icon("update").Text("Update").OnClick(p.UpdateApp).Hidden(!CurrentPage.updateAvailable),
		// 	Button().ID(p.ID+"-page-top-bar-install").Class("nav-bar").Icon(installIcon).Text("Install").OnClick(p.InstallApp).Hidden(!CurrentPage.installAvailable),
		// ),
		),
		app.Main().ID(p.ID+"-page-main").Class("page-main").Body(
			app.Dialog().ID(p.ID+"-page-status").Class("page-status").DataSet("status-type", p.statusType).Body(
				app.Span().ID(p.ID+"-page-status-text").Class("page-status-text").Text(p.statusText),
				compo.Button().ID(p.ID+"page-status-close-button").Class("page-status-close").Icon("close").OnClick(p.ClosePageStatus),
			),
			compo.Button().ID(p.ID+"-page-top-bar-account").Class("top-bar-account").Icon(accountButtonIcon).Text(nameFirstLetter).OnClick(NavigateEvent("/account")),
			app.If(p.TitleElement != "", app.H1().ID(p.ID+"-page-title").Class("page-title").Text(p.TitleElement)),
			app.If(p.SubtitleElement != "", app.P().ID(p.ID+"-page-subtitle").Class("page-subtitle").Text(p.SubtitleElement)),
			app.If(true, elements...),
		),
		navBar,
	)
}

// ShowMenu shows the menu on the page
func (p *Page) ShowMenu(ctx app.Context, e app.Event) {
	app.Window().GetElementByID(p.ID + "-page-menu").Call("show")
}

// ShowStatus shows the page status dialog with the given status text with the given status type
func (p *Page) ShowStatus(text string, statusType osusu.StatusType) {
	p.statusText = text
	p.statusType = statusType
	app.Window().GetElementByID(p.ID + "-page-status").Call("show")
}

// ShowErrorStatus is a shorthand for ShowStatus(err.Error(), StatusTypeNegative)
func (p *Page) ShowErrorStatus(err error) {
	p.ShowStatus(err.Error(), osusu.StatusTypeNegative)
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
	if osusu.Authenticate(p.AuthenticationRequired, ctx) {
		return
	}
	CurrentPage = p

	p.UpdatePageTitle(ctx)
	ctx.Page().SetDescription(p.Description)
	p.updateAvailable = ctx.AppUpdateAvailable()
	p.installAvailable = ctx.IsAppInstallable()
	p.user = osusu.CurrentUser(ctx)
	switch ctx.Page().URL().Path {
	case "/search":
		p.mode = "Search"
	case "/discover":
		p.mode = "Discover"
	case "/history":
		p.mode = "History"
	default:
		p.mode = osusu.GetOptions(ctx).Mode
	}
	p.Update()

	if p.OnNavFunc != nil {
		p.OnNavFunc(ctx)
	}

	ctx.Defer(func(ctx app.Context) {
		app.Window().GetElementByID(p.ID+"-page-main").Get("style").Set("opacity", 1)
		p.dialogElements = app.Window().Get("document").Call("querySelectorAll", ".select, .modal")
	})
	p.loaded = true
}

// OnClickEvent is called when someone clicks on the page
func (p *Page) OnClickEvent(ctx app.Context, e app.Event) {
	// forEachFunc := app.FuncOf(func(this app.Value, args []app.Value) any {
	// 	elem := args[0]
	// 	log.Println(elem)
	// 	elem.Call("close")
	// 	return nil
	// })
	// p.dialogElements.Call("forEach", forEachFunc)
	// log.Println(len(p.OnClick), p.OnClick)
	// log.Println("on click", p.onClick != nil, p.onClick)
	if p.onClick != nil {
		for _, event := range p.onClick {
			// log.Println(event)
			event(ctx, e)
		}
	}
}

// NavBarOnClick returns an event handler for a nav bar button that changes the mode to the given mode and then navigates to the given path
func (p *Page) NavBarOnClick(path string, mode string) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		p.mode = mode
		Navigate(path, ctx)
	}
}

// AddOnClick adds a new on click event handler to the page
func (p *Page) AddOnClick(eventHandler app.EventHandler) {
	if p.onClick == nil {
		p.onClick = []app.EventHandler{}
	}
	p.onClick = append(p.onClick, eventHandler)
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
func Navigate(urlString string, ctx app.Context) {
	// no need to navigate if we are already there (also prevents being stuck on 0 opacity)
	if urlString == ctx.Page().URL().Path {
		return
	}
	urlObject, err := url.Parse(urlString)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	// if external navigation, skip transition
	if urlObject.Host != "" && urlObject.Host != ctx.Page().URL().Host {
		ctx.NavigateTo(urlObject)
		return
	}
	// otherwise, set opacity to 0 and wait 250ms for transition before navigating
	app.Window().GetElementByID(CurrentPage.ID+"-page-main").Get("style").Set("opacity", 0)
	ctx.Defer(func(ctx app.Context) {
		time.Sleep(250 * time.Millisecond)
		ctx.NavigateTo(urlObject)
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
