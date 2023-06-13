package page

import (
	"sort"
	"strconv"
	"time"

	"github.com/kkoreilly/osusu/api"
	"github.com/kkoreilly/osusu/compo"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/kkoreilly/osusu/util/list"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Search struct {
	app.Compo
	group       osusu.Group
	user        osusu.User
	meals       osusu.Meals
	mealScores  map[int64]osusu.Score
	entries     osusu.Entries
	mealEntries map[int64]osusu.Entries
	options     osusu.Options
	currentMeal osusu.Meal
}

func (s *Search) Render() app.UI {
	return &compo.Page{
		ID:                     "search",
		Title:                  "Search",
		Description:            "Search for the best meals to eat given your current circumstances",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			compo.SetReturnURL("/search", ctx)
			s.group = osusu.CurrentGroup(ctx)
			if s.group.Name == "" {
				compo.Navigate("/groups", ctx)
			}
			s.user = osusu.CurrentUser(ctx)
			cuisines, err := api.GetGroupCuisines.Call(s.group.ID)
			if err != nil {
				compo.CurrentPage.ShowErrorStatus(err)
				return
			}
			s.group.Cuisines = cuisines
			osusu.SetCurrentGroup(s.group, ctx)

			s.options = osusu.GetOptions(ctx)
			if s.options.Users == nil {
				s.options = osusu.DefaultOptions(s.group)
			}

			meals, err := api.GetMeals.Call(s.group.ID)
			if err != nil {
				compo.CurrentPage.ShowErrorStatus(err)
				return
			}
			s.meals = meals

			entries, err := api.GetEntries.Call(s.group.ID)
			if err != nil {
				compo.CurrentPage.ShowErrorStatus(err)
				return
			}
			s.entries = entries
			sort.Slice(s.entries, func(i, j int) bool {
				return s.entries[i].Date.After(s.entries[j].Date)
			})
			s.mealEntries = make(map[int64]osusu.Entries)
			for _, entry := range s.entries {
				entries := s.mealEntries[entry.MealID]
				if entries == nil {
					entries = osusu.Entries{}
				}
				entries = append(entries, entry)
				s.mealEntries[entry.MealID] = entries
			}
			s.mealScores = make(map[int64]osusu.Score)
			s.SortMeals()

			compo.CurrentPage.AddOnClick(s.PageOnClick)
		},
		TitleElement:    "Search",
		SubtitleElement: "Search for the best meals to eat given your current circumstances",
		Elements: []app.UI{
			compo.ButtonRow().ID("search-page").Buttons(
				compo.Button().ID("search-page-new").Class("secondary").Icon("add").Text("New Meal").OnClick(s.NewMeal),
				compo.Button().ID("search-page-search").Class("primary").Icon("search").Text("Search").OnClick(s.ShowOptions),
			),
			compo.QuickOptions().ID("search-page").Options(&s.options).Group(s.group).Meals(s.meals).OnSave(func(ctx app.Context, e app.Event) { s.SortMeals() }),
			app.Div().ID("search-page-meals-container").Class("meal-images-container").Body(
				app.Range(s.meals).Slice(func(i int) app.UI {
					si := strconv.Itoa(i)
					meal := s.meals[i]
					entries := s.mealEntries[meal.ID]

					// check if at least one category satisfies a category option (or there are no categories and unset is an option)
					gotCategory := len(meal.Category) == 0 && s.options.Category["Unset"]
					if !gotCategory {
						for _, mealCategory := range meal.Category {
							for optionCategory, value := range s.options.Category {
								if value && mealCategory == optionCategory {
									gotCategory = true
									break
								}
							}
							if gotCategory {
								break
							}
						}
						if !gotCategory {
							return app.Text("")
						}
					}

					// check if at least one cuisine satisfies a cuisine option (or there are no cuisines and unset is an option)
					gotCuisine := len(meal.Cuisine) == 0 && s.options.Cuisine["Unset"]
					if !gotCuisine {
						for _, mealCuisine := range meal.Cuisine {
							for optionCuisine, value := range s.options.Cuisine {
								if value && mealCuisine == optionCuisine {
									gotCuisine = true
									break
								}
							}
							if gotCuisine {
								break
							}
						}
						if !gotCuisine {
							return app.Text("")
						}
					}

					// check if at least one entry satisfies the source requirements if there is at least one entry.
					if len(entries) > 0 {
						gotSource := false
						for _, entry := range entries {
							if s.options.Source[entry.Source] {
								gotSource = true
							}
						}
						if !gotSource {
							return app.Text("")
						}
					}
					score := s.mealScores[meal.ID]
					isCurrentMeal := meal.ID == s.currentMeal.ID

					// only put • between category and cuisine if both exist
					secondaryText := ""
					if len(meal.Category) != 0 && len(meal.Cuisine) != 0 {
						secondaryText = list.Slice(meal.Category) + " • " + list.Slice(meal.Cuisine)
					} else {
						secondaryText = list.Slice(meal.Category) + list.Slice(meal.Cuisine)
					}
					return compo.MealImage().ID("search-page-meal-" + si).Class("search-page-meal").Selected(isCurrentMeal).Img(meal.Image).MainText(meal.Name).SecondaryText(secondaryText).Score(score).OnClick(func(ctx app.Context, e app.Event) { s.MealOnClick(ctx, e, meal) }).OnClickScope(meal.ID)
				}),
			),
			app.Dialog().ID("search-page-meal-dialog").OnClick(s.MealDialogOnClick).Body(
				compo.Button().ID("search-page-meal-dialog-new-entry").Class("primary").Icon("add").Text("New Entry").OnClick(s.NewEntry),
				compo.Button().ID("search-page-meal-dialog-view-entries").Class("secondary").Icon("visibility").Text("View Entries").OnClick(s.ViewEntries),
				compo.Button().ID("search-page-meal-dialog-edit-meal").Class("secondary").Icon("edit").Text("Edit Meal").OnClick(s.EditMeal),
			),
			compo.Options().ID("search-page").Options(&s.options).OnSave(func(ctx app.Context, e app.Event) { s.SortMeals() }),
		},
	}
}

func (s *Search) NewMeal(ctx app.Context, e app.Event) {
	osusu.SetIsMealNew(true, ctx)
	osusu.SetCurrentMeal(osusu.Meal{}, ctx)
	compo.Navigate("/meal", ctx)
}

func (s *Search) PageOnClick(ctx app.Context, e app.Event) {
	// close meal dialog on page click (this will be stopped by another event if someone clicks on the meal dialog itself)
	mealDialog := app.Window().GetElementByID("search-page-meal-dialog")
	// need to check for null because people can click on page before dialog exists
	if !mealDialog.IsNull() {
		mealDialog.Call("close")
		s.currentMeal = osusu.Meal{}
	}
}

func (s *Search) MealOnClick(ctx app.Context, e app.Event, meal osusu.Meal) {
	e.Call("stopPropagation")
	s.currentMeal = meal
	osusu.SetCurrentMeal(meal, ctx)
	dialog := app.Window().GetElementByID("search-page-meal-dialog")
	if dialog.Get("open").Bool() {
		ctx.Dispatch(func(ctx app.Context) {
			dialog.Call("close")
		})
		ctx.Defer(func(ctx app.Context) {
			time.Sleep(250 * time.Millisecond)
			s.UpdateMealDialogPosition(ctx, e, dialog)
			dialog.Call("show")
		})
		return
	}
	s.UpdateMealDialogPosition(ctx, e, dialog)
	dialog.Call("show")
}

func (s *Search) UpdateMealDialogPosition(ctx app.Context, e app.Event, dialog app.Value) {
	pageX, pageY := e.Get("pageX").Int(), e.Get("pageY").Int()
	clientX, clientY := e.Get("clientX").Int(), e.Get("clientY").Int()
	clientWidth, clientHeight := dialog.Get("clientWidth").Int(), dialog.Get("clientHeight").Int()
	pageWidth, pageHeight := app.Window().Size()
	translateX, translateY := "0%", "0%"
	if clientX+clientWidth >= pageWidth {
		translateX = "-100%"
	}
	if clientY+clientHeight >= pageHeight {
		translateY = "-100%"
	}
	dialog.Get("style").Set("top", strconv.Itoa(pageY)+"px")
	dialog.Get("style").Set("left", strconv.Itoa(pageX)+"px")
	dialog.Get("style").Set("transform", "translate("+translateX+", "+translateY)
}

func (s *Search) MealDialogOnClick(ctx app.Context, e app.Event) {
	// stop the meal dialog from being closed by the page on click event
	e.Call("stopPropagation")
}

func (s *Search) NewEntry(ctx app.Context, e app.Event) {
	entry := osusu.NewEntry(s.group, s.user, s.currentMeal, s.mealEntries[s.currentMeal.ID])
	osusu.SetIsEntryNew(true, ctx)
	osusu.SetCurrentEntry(entry, ctx)
	compo.Navigate("/entry", ctx)
}

func (s *Search) ViewEntries(ctx app.Context, e app.Event) {
	compo.Navigate("/entries", ctx)
}

func (s *Search) EditMeal(ctx app.Context, e app.Event) {
	osusu.SetIsMealNew(false, ctx)
	compo.Navigate("/meal", ctx)
}

func (s *Search) ShowOptions(ctx app.Context, e app.Event) {
	app.Window().GetElementByID("search-page-options").Call("showModal")
}

func (s *Search) SortMeals() {
	for _, meal := range s.meals {
		s.mealScores[meal.ID] = meal.Score(s.mealEntries[meal.ID], s.options)
	}
	sort.Slice(s.meals, func(i, j int) bool {
		return s.mealScores[s.meals[i].ID].Total > s.mealScores[s.meals[j].ID].Total
	})
}
