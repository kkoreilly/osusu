package client

import (
	"github.com/kkoreilly/osusu/page"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func Start() {
	app.Route("/", &page.Start{})
	app.Route("/signin", &page.SignIn{})
	app.Route("/signup", &page.SignUp{})
	app.Route("/groups", &page.GroupsPage{})
	app.Route("/group", &page.GroupPage{})
	app.Route("/search", &page.Search{})
	app.Route("/home", &page.Home{})
	app.Route("/search", &page.Home{})
	app.Route("/history", &page.Home{})
	app.Route("/discover", &page.Home{})
	app.Route("/meal", &page.MealPage{})
	app.Route("/recipe", &page.RecipePage{})
	app.Route("/entries", &page.EntriesPage{})
	app.Route("/entry", &page.EntryPage{})
	app.Route("/account", &page.Account{})
	app.RouteWithRegexp("/join.*", &page.Join{})

	app.RunWhenOnBrowser()
}
