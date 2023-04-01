package main

import (
	"log"
	"strconv"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Entries []Entry

type entries struct {
	app.Compo
	person  Person
	meal    Meal
	entries Entries
}

func (e *entries) Render() app.UI {
	return &Page{
		ID:                     "entries",
		Title:                  "Entries",
		Description:            "View meal entries",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			e.person = GetCurrentPerson(ctx)
			e.meal = GetCurrentMeal(ctx)
			entries, err := GetEntriesAPI.Call(e.meal.ID)
			if err != nil {
				log.Println(err)
				return
			}
			e.entries = entries
		},
		TitleElement: "Entries for " + e.meal.Name,
		Elements: []app.UI{
			app.Div().ID("entries-page-action-button-row").Class("action-button-row").Body(
				app.A().ID("entries-page-back-button").Class("secondary-action-button", "action-button").Href("/meal").Text("Back"),
				app.Button().ID("entries-page-new-button").Class("primary-action-button", "action-button").Text("New").OnClick(e.New),
			),
			app.Hr(),
			app.Div().ID("entries-page-entries-container").Body(
				app.Range(e.entries).Slice(func(i int) app.UI {
					entry := e.entries[i]
					si := strconv.Itoa(i)
					return app.Div().ID("entries-page-entry-" + si).Class("entries-page-entry").
						OnClick(func(ctx app.Context, event app.Event) { e.EntryOnClick(ctx, event, entry) }).Body(
						app.Span().ID("entries-page-entry-date-text" + si).Class("entries-page-entry-date-text").Text(entry.Date.Format("Monday, January 2, 2006")),
					)
				}),
			),
			app.Hr(),
		},
	}
}

func (e *entries) EntryOnClick(ctx app.Context, event app.Event, entry Entry) {
	SetCurrentEntry(entry, ctx)
	ctx.Navigate("/entry")
}

func (e *entries) New(ctx app.Context, event app.Event) {
	entry, err := CreateEntryAPI.Call(Entry{
		UserID:      GetCurrentUser(ctx).ID,
		MealID:      e.meal.ID,
		Date:        time.Now(),
		Type:        "Dinner",
		Source:      "Cooking",
		Cost:        PersonMap{e.person.ID: 50},
		Effort:      PersonMap{e.person.ID: 50},
		Healthiness: PersonMap{e.person.ID: 50},
		Taste:       PersonMap{e.person.ID: 50},
	})
	if err != nil {
		log.Println(err)
		return
	}
	SetCurrentEntry(entry, ctx)
	ctx.Navigate("/entry")
}
