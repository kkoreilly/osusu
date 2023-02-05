package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

// Person is a struct that represents the data of a person within an account
type Person struct {
	Name string
}

// People is a slice that represents multiple people
type People []Person

// SetPeople sets the people that are in the current account
func SetPeople(people People, ctx app.Context) {
	ctx.SetState("people", people, app.Persist)
}

// GetPeople gets the people that are in the current account
func GetPeople(ctx app.Context) People {
	var people People
	ctx.GetState("people", &people)
	return people
}
