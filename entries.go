package main

import (
	"sort"
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Entries is a slice of multiple entries
type Entries []Entry

// // MissingData returns whether the given user is missing data in any of the given entries
// func (e Entries) MissingData(user User) bool {
// 	for _, entry := range e {
// 		if entry.MissingData(user) {
// 			return true
// 		}
// 	}
// 	return false
// }

type entries struct {
	app.Compo
	user    User
	meal    Meal
	options Options
	entries Entries
}

func (e *entries) Render() app.UI {
	// width, _ := app.Window().Size()
	// smallScreen := width <= 480
	return &Page{
		ID:                     "entries",
		Title:                  "Entries",
		Description:            "View meal entries",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			SetReturnURL("/entries", ctx)
			e.user = CurrentUser(ctx)
			e.meal = CurrentMeal(ctx)
			e.options = GetOptions(ctx)
			// recency doesn't matter for entries
			e.options.RecencyWeight = 0
			entries, err := GetEntriesForMealAPI.Call(e.meal.ID)
			if err != nil {
				CurrentPage.ShowErrorStatus(err)
				return
			}
			e.entries = entries
			sort.Slice(entries, func(i, j int) bool {
				// // prioritize entries with missing data, then date
				// entryI := e.entries[i]
				// iMissingData := entryI.MissingData(e.user)
				// entryJ := e.entries[j]
				// jMissingData := entryJ.MissingData(e.user)
				// if iMissingData && !jMissingData {
				// 	return true
				// }
				// if !iMissingData && jMissingData {
				// 	return false
				// }
				return entries[i].Date.After(entries[j].Date)
			})
		},
		TitleElement: "Entries for " + e.meal.Name,
		Elements: []app.UI{
			ButtonRow().ID("entries-page").Buttons(
				Button().ID("entries-page-back").Class("secondary").Icon("arrow_back").Text("Back").OnClick(NavigateEvent("/home")),
				Button().ID("entries-page-new").Class("primary").Icon("add").Text("New").OnClick(e.New),
			),
			app.Div().ID("entries-page-entries-container").Class("meal-images-container").Body(
				app.Range(e.entries).Slice(func(i int) app.UI {
					si := strconv.Itoa(i)
					entry := e.entries[i]
					score := entry.Score(e.options)
					return MealImage().ID("entries-page-entry" + si).Class("entries-page-entry").MainText(entry.Date.Format("Monday, January 2, 2006")).SecondaryText(entry.Type + " • " + entry.Source).Score(score).OnClick(func(ctx app.Context, event app.Event) { e.EntryOnClick(ctx, event, entry) }).OnClickScope(entry.ID)
				}),
			),
			// app.Table().ID("entries-page-entries-table").Body(
			// 	app.THead().ID("entries-page-entries-table-header").Body(
			// 		app.Tr().ID("entries-page-entries-table-header-row").Body(
			// 			app.Th().ID("entries-page-entries-table-header-date").Text("Date"),
			// 			app.Th().ID("entries-page-entries-table-header-total").Class("table-header-score").Text("Total"),
			// 			app.Th().ID("entries-page-entries-table-header-taste").Class("table-header-score").Text("Taste"),
			// 			app.Th().ID("entries-page-entries-table-header-cost").Class("table-header-score").Text("Cost"),
			// 			app.Th().ID("entries-page-entries-table-header-effort").Class("table-header-score").Text("Effort"),
			// 			app.Th().ID("entries-page-entries-table-header-healthiness").Class("table-header-score").Text("Health"),
			// 			app.If(!smallScreen,
			// 				app.Th().ID("entries-page-entries-table-header-type").Text("Type"),
			// 				app.Th().ID("entries-page-entries-table-header-source").Text("Source"),
			// 				app.Th().ID("entries-page-entries-table-header-people").Text("People"),
			// 			),
			// 		),
			// 	),
			// 	app.TBody().ID("entries-page-entries-table-body").Body(
			// 		app.Range(e.entries).Slice(func(i int) app.UI {
			// 			entry := e.entries[i]
			// 			si := strconv.Itoa(i)
			// 			score := entry.Score(e.options)
			// 			colorH := strconv.Itoa((score.Total * 12) / 10)
			// 			scoreText := strconv.Itoa(score.Total)

			// 			// missingData := entry.MissingData(e.user)
			// 			return app.Tr().ID("entries-page-entry-"+si).Class("entries-page-entry").Style("--color-h", colorH).Style("--score", scoreText+"%").
			// 				OnClick(func(ctx app.Context, event app.Event) { e.EntryOnClick(ctx, event, entry) }).Body(
			// 				app.Td().ID("entries-page-entry-date-"+si).Class("entries-page-entry-date").Text(entry.Date.Format("Jan 2, 2006")),
			// 				MealScore("entries-page-entry-total-"+si, "entries-page-entry-total", score.Total, "Total"),
			// 				MealScore("entries-page-entry-taste-"+si, "entries-page-entry-taste", score.Taste, "Taste"),
			// 				MealScore("entries-page-entry-cost-"+si, "entries-page-entry-cost", score.Cost, "Cost"),
			// 				MealScore("entries-page-entry-effort-"+si, "entries-page-entry-effort", score.Effort, "Effort"),
			// 				MealScore("entries-page-entry-healthiness-"+si, "entries-page-entry-healthiness", score.Healthiness, "Healthiness"),
			// 				app.If(!smallScreen,
			// 					app.Td().ID("entries-page-entry-type-"+si).Class("entries-page-entry-type").Text(entry.Type),
			// 					app.Td().ID("entries-page-entry-source-"+si).Class("entries-page-entry-source").Text(entry.Source),
			// 					app.Td().ID("entries-page-entry-people-"+si).Class("entries-page-entry-people").Text(ListString(e.options.UsersList(entry.Cost))),
			// 				),
			// 			)
			// 		}),
			// 	),
			// ),
		},
	}
}

func (e *entries) EntryOnClick(ctx app.Context, event app.Event, entry Entry) {
	SetIsEntryNew(false, ctx)
	SetCurrentEntry(entry, ctx)
	Navigate("/entry", ctx)
}

func (e *entries) New(ctx app.Context, event app.Event) {
	entry := NewEntry(CurrentGroup(ctx), e.user, e.meal, e.entries)
	SetIsEntryNew(true, ctx)
	SetCurrentEntry(entry, ctx)
	Navigate("/entry", ctx)
}
