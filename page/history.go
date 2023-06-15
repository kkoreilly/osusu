package page

import (
	"sort"
	"strconv"

	"github.com/kkoreilly/osusu/api"
	"github.com/kkoreilly/osusu/compo"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/kkoreilly/osusu/util/cond"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type History struct {
	app.Compo
	group           osusu.Group
	user            osusu.User
	meals           osusu.Meals
	entries         osusu.Entries
	entryMeals      map[int64]osusu.Meal // the meal each entry is associated with
	entryScores     map[int64]osusu.Score
	showEntries     map[int64]bool // whether each entry is shown
	numEntriesShown int
	options         osusu.Options
}

func (h *History) Render() app.UI {
	return &compo.Page{
		ID:                     "history",
		Title:                  "History",
		Description:            "View the history of what meals you've eaten and how they were",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			compo.SetReturnURL("/history", ctx)
			h.group = osusu.CurrentGroup(ctx)
			if h.group.Name == "" {
				compo.Navigate("/groups", ctx)
			}
			h.user = osusu.CurrentUser(ctx)
			cuisines, err := api.GetGroupCuisines.Call(h.group.ID)
			if err != nil {
				compo.CurrentPage.ShowErrorStatus(err)
				return
			}
			h.group.Cuisines = cuisines
			osusu.SetCurrentGroup(h.group, ctx)

			h.options = osusu.GetOptions(ctx)
			if h.options.Users == nil {
				h.options = osusu.DefaultOptions(h.group)
			}

			meals, err := api.GetMeals.Call(h.group.ID)
			if err != nil {
				compo.CurrentPage.ShowErrorStatus(err)
				return
			}
			h.meals = meals

			entries, err := api.GetEntries.Call(h.group.ID)
			if err != nil {
				compo.CurrentPage.ShowErrorStatus(err)
				return
			}
			h.entries = entries

			h.entryMeals = make(map[int64]osusu.Meal)
			for _, entry := range h.entries {
				for _, meal := range h.meals {
					if entry.MealID == meal.ID {
						h.entryMeals[entry.ID] = meal
						break
					}
				}
			}

			h.entryScores = make(map[int64]osusu.Score)
			h.showEntries = make(map[int64]bool)
			h.SortEntries()
		},
		TitleElement:    "History",
		SubtitleElement: "View the history of what meals you've eaten and how they were",
		Elements: []app.UI{
			compo.ButtonRow().ID("history-page").Buttons(
				compo.Button().ID("history-page-search").Class("primary").Icon("search").Text("Search").OnClick(h.ShowOptions),
			),
			compo.QuickOptions().ID("history-page").Options(&h.options).Group(h.group).Meals(h.meals).OnSave(func(ctx app.Context, e app.Event) { h.SortEntries() }),
			app.P().ID("history-page-no-entries-shown").Class("centered-text").Text(cond.IfElse(len(h.entries) == 0, "You have not created any entries yet. Please try adding a new entry by navigating to the Search page, selecting a meal, and pressing the New Entry button.", "No entries satisfy your filters. Please try changing them or adding a new entry by navigating to the Search page, selecting a meal, and pressing the New Entry button.")).Hidden(h.numEntriesShown != 0),
			app.Div().ID("history-page-entries-container").Class("meal-images-container").Body(
				app.Range(h.entries).Slice(func(i int) app.UI {
					si := strconv.Itoa(i)
					entry := h.entries[i]

					if !h.showEntries[entry.ID] {
						return app.Text("")
					}

					score := h.entryScores[entry.ID]
					entryMeal := h.entryMeals[entry.ID]
					secondaryText := entryMeal.Name
					if entry.Category != "" {
						secondaryText += " • " + entry.Category
					}
					if entry.Source != "" {
						secondaryText += " • " + entry.Source
					}
					return compo.MealImage().ID("history-page-entry-" + si).Class("history-page-entry").Img(entryMeal.Image).MainText(entry.Date.Format("Monday, January 2, 2006")).SecondaryText(secondaryText).Score(score).OnClick(func(ctx app.Context, e app.Event) { h.EntryOnClick(ctx, e, entry) }).OnClickScope(entry.ID)
				}),
			),
			// recency is irrelevant for history
			compo.Options().ID("history-page").Options(&h.options).Exclude("recency").OnSave(func(ctx app.Context, e app.Event) { h.SortEntries() }),
		},
	}
}

func (h *History) ShowOptions(ctx app.Context, e app.Event) {
	app.Window().GetElementByID("history-page-options").Call("showModal")
}

func (h *History) EntryOnClick(ctx app.Context, e app.Event, entry osusu.Entry) {
	osusu.SetIsEntryNew(false, ctx)
	osusu.SetCurrentEntry(entry, ctx)
	compo.Navigate("/entry", ctx)
}

func (h *History) SortEntries() {
	h.numEntriesShown = 0
	for _, entry := range h.entries {
		entryMeal := h.entryMeals[entry.ID]
		gotCategory := false
		for category, val := range h.options.Category {
			if val && category == entry.Category {
				gotCategory = true
				break
			}
		}
		if !gotCategory {
			h.showEntries[entry.ID] = false
			continue
		}
		gotCuisine := false
		for _, mealCuisine := range entryMeal.Cuisine {
			for optionCuisine, val := range h.options.Cuisine {
				if val && mealCuisine == optionCuisine {
					gotCuisine = true
					break
				}
			}
			if gotCuisine {
				break
			}
		}
		if !gotCuisine {
			h.showEntries[entry.ID] = false
			continue
		}
		gotSource := false
		for source, val := range h.options.Source {
			if val && source == entry.Source {
				gotSource = true
				break
			}
		}
		if !gotSource {
			h.showEntries[entry.ID] = false
			continue
		}
		h.showEntries[entry.ID] = true
		h.numEntriesShown++
		h.entryScores[entry.ID] = entry.Score(h.options)
	}
	sort.Slice(h.entries, func(i, j int) bool {
		return h.entries[i].Date.After(h.entries[j].Date)
	})
}
