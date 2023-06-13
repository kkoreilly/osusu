package client

import (
	"github.com/kkoreilly/osusu/page"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func Start() {
	app.Route("/", &page.Start{})
	app.Route("/signin", &page.SignIn{})
	app.Route("/signup", &page.SignUp{})

	app.Route("/groups", &page.Groups{})
	app.Route("/group", &page.Group{})
	app.RouteWithRegexp("/join.*", &page.Join{})
	app.Route("/account", &page.Account{})

	app.Route("/search", &page.Search{})
	app.Route("/discover", &page.Discover{})
	app.Route("/history", &page.History{})

	app.Route("/meal", &page.Meal{})
	app.Route("/recipe", &page.Recipe{})
	app.Route("/entries", &page.Entries{})
	app.Route("/entry", &page.Entry{})

	app.RunWhenOnBrowser()
}
