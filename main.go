package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func main() {
	app.Route("/", &start{})
	app.Route("/signin", &signIn{})
	app.Route("/signup", &signUp{})
	app.Route("/home", &home{})
	app.Route("/edit", &edit{})
	app.Route("/recommendations", &recommendations{})
	app.Route("/people", &people{})
	app.Route("/person", &person{})

	app.RunWhenOnBrowser()

	startServer()
}
