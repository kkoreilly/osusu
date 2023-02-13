package main

import (
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type people struct {
	app.Compo
	people []string
}

func (p *people) Render() app.UI {
	return app.Div().Body(
		app.H1().ID("people-page-title").Class("page-title").Text("Who Are You?"),
		app.Div().ID("people-page-action-button-row").Class("action-button-row").Body(
			app.Button().ID("people-page-new-person").Class("action-button", "blue-action-button").Text("New Person").OnClick(p.New),
		),
		app.Div().ID("people-page-people-container").Body(
			app.Range(p.people).Slice(func(i int) app.UI {
				return app.Div().ID("people-page-person-" + strconv.Itoa(i)).Class("people-page-person").Text(p.people[i]).
					OnClick(func(ctx app.Context, e app.Event) { p.PersonOnClick(ctx, e, i) })
			}),
		),
	)
}

func (p *people) OnNav(ctx app.Context) {
	p.people = GetCurrentUser(ctx).People
}

func (p *people) New(ctx app.Context, e app.Event) {
	p.people = append(p.people, "A Person")
}

func (p *people) PersonOnClick(ctx app.Context, e app.Event, i int) {
	ctx.Navigate("/person")
}
