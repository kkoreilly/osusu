package page

import (
	"sort"
	"strconv"

	"github.com/kkoreilly/osusu/api"
	"github.com/kkoreilly/osusu/compo"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type EntriesPage struct {
	app.Compo
	user    osusu.User
	meal    osusu.Meal
	options osusu.Options
	entries osusu.Entries
}

func (e *EntriesPage) Render() app.UI {
	// width, _ := app.Window().Size()
	// smallScreen := width <= 480
	return &compo.Page{
		ID:                     "entries",
		Title:                  "Entries",
		Description:            "View meal entries",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			compo.SetReturnURL("/entries", ctx)
			e.user = osusu.CurrentUser(ctx)
			e.meal = osusu.CurrentMeal(ctx)
			e.options = osusu.GetOptions(ctx)
			// recency doesn't matter for entries
			e.options.RecencyWeight = 0
			entries, err := api.GetEntriesForMeal.Call(e.meal.ID)
			if err != nil {
				compo.CurrentPage.ShowErrorStatus(err)
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
			compo.ButtonRow().ID("entries-page").Buttons(
				compo.Button().ID("entries-page-back").Class("secondary").Icon("arrow_back").Text("Back").OnClick(compo.NavigateEvent("/home")),
				compo.Button().ID("entries-page-new").Class("primary").Icon("add").Text("New").OnClick(e.New),
			),
			app.Div().ID("entries-page-entries-container").Class("meal-images-container").Body(
				app.Range(e.entries).Slice(func(i int) app.UI {
					si := strconv.Itoa(i)
					entry := e.entries[i]
					score := entry.Score(e.options)
					return compo.MealImage().ID("entries-page-entry" + si).Class("entries-page-entry").MainText(entry.Date.Format("Monday, January 2, 2006")).SecondaryText(entry.Type + " • " + entry.Source).Score(score).OnClick(func(ctx app.Context, event app.Event) { e.EntryOnClick(ctx, event, entry) }).OnClickScope(entry.ID)
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

func (e *EntriesPage) EntryOnClick(ctx app.Context, event app.Event, entry osusu.Entry) {
	osusu.SetIsEntryNew(false, ctx)
	osusu.SetCurrentEntry(entry, ctx)
	compo.Navigate("/entry", ctx)
}

func (e *EntriesPage) New(ctx app.Context, event app.Event) {
	entry := osusu.NewEntry(osusu.CurrentGroup(ctx), e.user, e.meal, e.entries)
	osusu.SetIsEntryNew(true, ctx)
	osusu.SetCurrentEntry(entry, ctx)
	compo.Navigate("/entry", ctx)
}
