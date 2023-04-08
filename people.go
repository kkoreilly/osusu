package main

import (
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// People represents the data of multiple people
type People []Person

type people struct {
	app.Compo
	people People
}

func (p *people) Render() app.UI {
	return &Page{
		ID:                     "people",
		Title:                  "People",
		Description:            "Select a person",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			people, err := GetPeopleAPI.Call(GetCurrentUser(ctx).ID)
			if err != nil {
				CurrentPage.ShowStatus(err.Error(), StatusTypeNegative)
				return
			}
			p.people = people
		},
		TitleElement: "Who Are You?",
		Elements: []app.UI{
			app.Div().ID("people-page-action-button-row").Class("action-button-row").Body(
				app.Button().ID("people-page-new-person").Class("action-button", "primary-action-button").Text("New Person").OnClick(p.New),
			),
			app.Div().ID("people-page-people-container").Body(
				app.Range(p.people).Slice(func(i int) app.UI {
					return app.Button().ID("people-page-person-" + strconv.Itoa(i)).Class("people-page-person").Text(p.people[i].Name).
						OnClick(func(ctx app.Context, e app.Event) { p.PersonOnClick(ctx, e, p.people[i]) })
				}),
			),
		},
	}

}

func (p *people) New(ctx app.Context, e app.Event) {
	person, err := CreatePersonAPI.Call(GetCurrentUser(ctx).ID)
	if err != nil {
		CurrentPage.ShowStatus(err.Error(), StatusTypeNegative)
		return
	}
	SetCurrentPerson(person, ctx)
	ctx.Navigate("/person")
}

func (p *people) PersonOnClick(ctx app.Context, e app.Event, person Person) {
	SetCurrentPerson(person, ctx)
	ctx.Navigate("/person")
}

// SetCurrentPerson sets the current person state value
func SetCurrentPerson(person Person, ctx app.Context) {
	ctx.SetState("currentPerson", person, app.Persist)
}

// GetCurrentPerson gets the current person state value
func GetCurrentPerson(ctx app.Context) Person {
	var person Person
	ctx.GetState("currentPerson", &person)
	return person
}
