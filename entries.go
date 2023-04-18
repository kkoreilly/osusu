package main

import (
	"sort"
	"strconv"

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
				CurrentPage.ShowErrorStatus(err)
				return
			}
			e.entries = entries
			sort.Slice(entries, func(i, j int) bool {
				// prioritize entries with missing data, then date
				entryI := e.entries[i]
				iMissingData := entryI.MissingData(e.person)
				entryJ := e.entries[j]
				jMissingData := entryJ.MissingData(e.person)
				if iMissingData && !jMissingData {
					return true
				}
				if !iMissingData && jMissingData {
					return false
				}
				return entries[i].Date.After(entries[j].Date)
			})
		},
		TitleElement: "Entries for " + e.meal.Name,
		Elements: []app.UI{
			app.Div().ID("entries-page-action-button-row").Class("action-button-row").Body(
				app.A().ID("entries-page-back-button").Class("secondary-action-button", "action-button").Href("/home").Text("Back"),
				app.Button().ID("entries-page-new-button").Class("primary-action-button", "action-button").Text("New").OnClick(e.New),
			),
			app.Table().ID("entries-page-entries-table").Body(
				app.THead().ID("entries-page-entries-table-header").Body(
					app.Tr().ID("entries-page-entries-table-header-row").Body(
						app.Th().ID("entries-page-entries-table-header-name").Text("Name"),
						app.Th().ID("entries-page-entries-table-header-score").Text("Score"),
					),
				),
				app.TBody().ID("entries-page-entries-table-body").Body(
					app.Range(e.entries).Slice(func(i int) app.UI {
						entry := e.entries[i]
						si := strconv.Itoa(i)
						score := entry.Score(e.options)
						colorH := strconv.Itoa((score * 12) / 10)
						scoreText := strconv.Itoa(score)
						missingData := entry.MissingData(e.person)
						return app.Tr().ID("entries-page-entry-"+si).Class("entries-page-entry").DataSet("missing-data", missingData).Style("--color-h", colorH).Style("--score-percent", scoreText+"%").
							OnClick(func(ctx app.Context, event app.Event) { e.EntryOnClick(ctx, event, entry) }).Body(
							app.Td().ID("entries-page-entry-date"+si).Class("entries-page-entry-date").Text(entry.Date.Format("1/2/2006")),
							app.Td().ID("entries-page-entry-score"+si).Class("entries-page-entry-score").Text(scoreText),
						)
					}),
				),
			),
		},
	}
}

func (e *entries) EntryOnClick(ctx app.Context, event app.Event, entry Entry) {
	SetIsEntryNew(false, ctx)
	SetCurrentEntry(entry, ctx)
	ctx.Navigate("/entry")
}

func (e *entries) New(ctx app.Context, event app.Event) {
	entry := NewEntry(GetCurrentUser(ctx), e.meal, e.person, e.entries)
	SetIsEntryNew(true, ctx)
	SetCurrentEntry(entry, ctx)
	ctx.Navigate("/entry")
}
