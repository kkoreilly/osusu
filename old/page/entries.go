package page

import (
	"sort"
	"strconv"

	"github.com/kkoreilly/osusu/api"
	"github.com/kkoreilly/osusu/compo"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Entries struct {
	app.Compo
	user    osusu.User
	meal    osusu.Meal
	options osusu.Options
	entries osusu.Entries
}

func (e *Entries) Render() app.UI {
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
				return entries[i].Date.After(entries[j].Date)
			})
		},
		TitleElement: "Entries for " + e.meal.Name,
		Elements: []app.UI{
			&compo.ButtonRow{ID: "entries-page", Buttons: []app.UI{
				&compo.Button{ID: "entries-page-back", Class: "secondary", Icon: "arrow_back", Text: "Back", OnClick: compo.NavigateEvent("/search")},
				&compo.Button{ID: "entries-page-new", Class: "primary", Icon: "add", Text: "New", OnClick: e.New},
			}},
			app.Div().ID("entries-page-entries-container").Class("meal-images-container").Body(
				app.Range(e.entries).Slice(func(i int) app.UI {
					si := strconv.Itoa(i)
					entry := e.entries[i]
					score := entry.Score(e.options)
					return &compo.MealImage{ID: "entries-page-entry" + si, Class: "entries-page-entry", MainText: entry.Date.Format("Monday, January 2, 2006"), SecondaryText: entry.Category + " • " + entry.Source, Score: score, OnClick: func(ctx app.Context, event app.Event) { e.EntryOnClick(ctx, event, entry) }, OnClickScope: []any{entry.ID}}
				}),
			),
		},
	}
}

func (e *Entries) EntryOnClick(ctx app.Context, event app.Event, entry osusu.Entry) {
	osusu.SetIsEntryNew(false, ctx)
	osusu.SetCurrentEntry(entry, ctx)
	compo.Navigate("/entry", ctx)
}

func (e *Entries) New(ctx app.Context, event app.Event) {
	entry := osusu.NewEntry(osusu.CurrentGroup(ctx), e.user, e.meal, e.entries)
	osusu.SetIsEntryNew(true, ctx)
	osusu.SetCurrentEntry(entry, ctx)
	compo.Navigate("/entry", ctx)
}
