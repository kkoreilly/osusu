package main

import (
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Entries is a slice of multiple entries
type Entries []Entry

// MissingData returns whether the given person is missing data in any of the given entries
func (e Entries) MissingData(person Person) bool {
	for _, entry := range e {
		if entry.MissingData(person) {
			return true
		}
	}
	return false
}

type entries struct {
	app.Compo
	person  Person
	meal    Meal
	options Options
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
			e.options = GetOptions(ctx)
			entries, err := GetEntriesForMealAPI.Call(e.meal.ID)
			if err != nil {
				log.Println(err)
				return
			}
			sort.Slice(entries, func(i, j int) bool {
				return entries[i].Date.After(entries[j].Date)
			})
			e.entries = entries
		},
		TitleElement: "Entries for " + e.meal.Name,
		Elements: []app.UI{
			app.Div().ID("entries-page-action-button-row").Class("action-button-row").Body(
				app.A().ID("entries-page-back-button").Class("secondary-action-button", "action-button").Href("/meal").Text("Back"),
				app.Button().ID("entries-page-new-button").Class("primary-action-button", "action-button").Text("New").OnClick(e.New),
			),
			app.Div().ID("entries-page-entries-container").Body(
				app.Range(e.entries).Slice(func(i int) app.UI {
					entry := e.entries[i]
					si := strconv.Itoa(i)
					score := entry.Score(e.options)
					colorH := strconv.Itoa((score * 12) / 10)
					scoreText := strconv.Itoa(score)
					missingData := entry.MissingData(e.person)
					return app.Div().ID("entries-page-entry-"+si).Class("entries-page-entry").DataSet("missing-data", missingData).Style("--color-h", colorH).Style("--score-percent", scoreText+"%").
						OnClick(func(ctx app.Context, event app.Event) { e.EntryOnClick(ctx, event, entry) }).Body(
						app.Span().ID("entries-page-entry-date"+si).Class("entries-page-entry-date").Text(entry.Date.Format("Monday, January 2, 2006")),
						app.Span().ID("entries-page-entry-score"+si).Class("entries-page-entry-score").Text(scoreText),
					)
				}),
			),
		},
	}
}

func (e *entries) EntryOnClick(ctx app.Context, event app.Event, entry Entry) {
	SetCurrentEntry(entry, ctx)
	ctx.Navigate("/entry")
}

func (e *entries) New(ctx app.Context, event app.Event) {
	newEntry := Entry{
		UserID:      GetCurrentUser(ctx).ID,
		MealID:      e.meal.ID,
		Date:        time.Now(),
		Type:        "Dinner",
		Source:      "Cooking",
		Cost:        PersonMap{e.person.ID: 50},
		Effort:      PersonMap{e.person.ID: 50},
		Healthiness: PersonMap{e.person.ID: 50},
		Taste:       PersonMap{e.person.ID: 50},
	}
	// if there are previous entries, copy the values from the latest (the entries are already sorted with the latest first by OnNav)
	if len(e.entries) > 0 {
		previousEntry := e.entries[0]
		newEntry = Entry{
			UserID:      GetCurrentUser(ctx).ID,
			MealID:      e.meal.ID,
			Date:        time.Now(),
			Type:        previousEntry.Type,
			Source:      previousEntry.Source,
			Cost:        previousEntry.Cost,
			Effort:      previousEntry.Effort,
			Healthiness: previousEntry.Healthiness,
			Taste:       previousEntry.Taste,
		}
	}

	entry, err := CreateEntryAPI.Call(newEntry)
	if err != nil {
		log.Println(err)
		return
	}
	SetCurrentEntry(entry, ctx)
	ctx.Navigate("/entry")
}
