package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func main() {
	app.Route("/", &start{})
	app.Route("/signin", &signIn{})
	app.Route("/signup", &signUp{})
	app.Route("/home", &home{})
	app.Route("/edit", &edit{})
	app.Route("/recommendations", &recommendations{})
	app.Route("/people", &people{})

	// Once the routes set up, the next thing to do is to either launch the app
	// or the server that serves the app.
	//
	// When executed on the client-side, the RunWhenOnBrowser() function
	// launches the app,  starting a loop that listens for app events and
	// executes client instructions. Since it is a blocking call, the code below
	// it will never be executed.
	//
	// When executed on the server-side, RunWhenOnBrowser() does nothing, which
	// lets room for server implementation without the need for precompiling
	// instructions.
	app.RunWhenOnBrowser()

	// Finally, launching the server that serves the app is done by using the Go
	// standard HTTP package.
	//
	// The Handler is an HTTP handler that serves the client and all its
	// required resources to make it work into a web browser. Here it is
	// configured to handle requests with a path that starts with "/".
	http.Handle("/", &app.Handler{
		Name:        "MealRec",
		Title:       "MealRec",
		Icon:        app.Icon{Default: "/web/images/icon-192.png", Large: "/web/images/icon-512.png"},
		Description: "An app for getting recommendations on what meals to eat",
		Styles:      []string{"https://fonts.googleapis.com/css?family=Roboto", "/web/css/global.css", "/web/css/start.css", "/web/css/signinup.css", "/web/css/home.css", "/web/css/edit.css", "/web/css/people.css"},
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
