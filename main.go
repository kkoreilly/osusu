package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func main() {
	app.Route("/", &start{})
	app.Route("/signin", &signIn{})
	app.Route("/signup", &signUp{})
	app.Route("/home", &home{})
	app.Route("/meal", &meal{})
	app.Route("/entries", &entries{})
	app.Route("/entry", &entry{})
	app.Route("/people", &people{})
	app.Route("/person", &person{})

	app.RunWhenOnBrowser()

	startServer()
}
