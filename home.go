package main

import (
	"sort"
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type home struct {
	app.Compo
	group              Group
	user               User
	users              []User
	meals              Meals
	entriesForEachMeal map[int64]Entries // entries for each meal id
	options            Options
	usersOptions       map[string]bool
	cuisinesInUse      map[string]bool // key is cuisine name, value is if it's being used
	currentMeal        Meal            // the current selected meal for the context menu
}

func (h *home) Render() app.UI {
	usersStrings := []string{}
	for _, u := range h.users {
		usersStrings = append(usersStrings, u.Name)
	}
	// // need to copy to separate array from because append modifies the underlying array
	// var cuisines = make([]string, len(h.user.Cuisines))
	// copy(cuisines, h.user.Cuisines)
	cuisines := []string{}
	for cuisine, val := range h.cuisinesInUse {
		if val {
			cuisines = append(cuisines, cuisine)
		}
	}
	return &Page{
		ID:                     "home",
		Title:                  "Home",
		Description:            "View, sort, and filter your meals.",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			SetReturnURL("/home", ctx)
			h.group = GetCurrentGroup(ctx)
			if h.group.Name == "" {
				ctx.Navigate("/groups")
			}
			h.user = GetCurrentUser(ctx)
			cuisines, err := GetUserCuisinesAPI.Call(h.user.ID)
			if err != nil {
				CurrentPage.ShowErrorStatus(err)
				return
			}
			h.user.Cuisines = cuisines
			SetCurrentUser(h.user, ctx)

			users, err := GetUsersAPI.Call(h.group.Members)
			if err != nil {
				CurrentPage.ShowErrorStatus(err)
				return
			}
			h.users = users

			// people, err := GetPeopleAPI.Call(h.user.ID)
			// if err != nil {
			// 	CurrentPage.ShowErrorStatus(err)
			// 	return
			// }
			// h.people = people

			h.options = GetOptions(ctx)
			if h.options.Users == nil {
				h.options = DefaultOptions(h.user)
			}
			h.options = h.options.RemoveInvalidCuisines(h.user.Cuisines)
			SetOptions(h.options, ctx)
			h.usersOptions = make(map[string]bool)
			for _, p := range h.users {
				if _, ok := h.options.Users[p.ID]; !ok {
					h.options.Users[p.ID] = true
				}
				h.usersOptions[p.Name] = h.options.Users[p.ID]
			}

			meals, err := GetMealsAPI.Call(h.group.ID)
			if err != nil {
				CurrentPage.ShowErrorStatus(err)
				return
			}
			h.meals = meals

			h.cuisinesInUse = map[string]bool{}
			for _, meal := range h.meals {
				for _, cuisine := range meal.Cuisine {
					h.cuisinesInUse[cuisine] = true
					// if the user has not yet set whether or not to allow this cuisine (if it is new), automatically set it to true
					_, ok := h.options.Cuisine[cuisine]
					if !ok {
						h.options.Cuisine[cuisine] = true
					}
				}
			}

			entries, err := GetEntriesAPI.Call(h.group.ID)
			if err != nil {
				CurrentPage.ShowErrorStatus(err)
				return
			}
			h.entriesForEachMeal = make(map[int64]Entries)
			for _, entry := range entries {
				entries := h.entriesForEachMeal[entry.MealID]
				if entries == nil {
					entries = Entries{}
				}
				entries = append(entries, entry)
				h.entriesForEachMeal[entry.MealID] = entries
			}
			h.SortMeals()
		},
		OnClick:      h.PageOnClick,
		TitleElement: "Welcome, " + h.user.Name,
		Elements: []app.UI{
			app.Div().ID("home-page-action-button-row").Class("action-button-row").Body(
				app.Button().ID("home-page-new-button").Class("secondary-action-button", "action-button").Text("New Meal").OnClick(h.New),
				app.Button().ID("home-page-options-button").Class("primary-action-button", "action-button").Text("Search").OnClick(h.ShowOptions),
			),
			app.Table().ID("home-page-meals-table").Body(
				app.THead().ID("home-page-meals-table-header").Body(
					app.Tr().ID("home-page-meals-table-header-row").Body(
						app.Th().ID("home-page-meals-table-header-name").Text("Name"),
						app.Th().ID("home-page-meals-table-header-score").Text("Score"),
					),
				),
				app.TBody().ID("home-page-meals-table-body").Body(
					app.Range(h.meals).Slice(func(i int) app.UI {
						meal := h.meals[i]
						si := strconv.Itoa(i)
						entries := h.entriesForEachMeal[meal.ID]

						// check if at least one cuisine satisfies a cuisine requirement
						gotCuisine := false
						for _, mealCuisine := range meal.Cuisine {
							for optionCuisine, value := range h.options.Cuisine {
								if value && mealCuisine == optionCuisine {
									gotCuisine = true
								}
							}
						}
						if !gotCuisine {
							return app.Text("")
						}

						// check if at least one entry satisfies the type and source requirements if there is at least one entry.
						if len(entries) > 0 {
							gotType := false
							gotSource := false
							for _, entry := range entries {
								if entry.Type == h.options.Type {
									gotType = true
								}
								if h.options.Source[entry.Source] {
									gotSource = true
								}
							}
							if !(gotType && gotSource) {
								return app.Text("")
							}
						}

						score := meal.Score(entries, h.options)
						colorH := strconv.Itoa((score * 12) / 10)
						scoreText := strconv.Itoa(score)
						missingData := entries.MissingData(h.user)
						isCurrentMeal := meal.ID == h.currentMeal.ID
						return app.Tr().ID("home-page-meal-"+si).Class("home-page-meal").DataSet("missing-data", missingData).DataSet("current-meal", isCurrentMeal).Style("--color-h", colorH).Style("--score-percent", scoreText+"%").
							OnClick(func(ctx app.Context, e app.Event) { h.MealOnClick(ctx, e, meal) }).Body(
							app.Td().ID("home-page-meal-name-"+si).Class("home-page-meal-name").Text(meal.Name),
							app.Td().ID("home-page-meal-score-"+si).Class("home-page-meal-score").Text(scoreText),
						)
					}),
				),
			),
			app.Dialog().ID("home-page-meal-dialog").OnClick(h.MealDialogOnClick).Body(
				app.Button().ID("home-page-meal-dialog-new-entry-button").Class("action-button", "primary-action-button").Text("New Entry").OnClick(h.NewEntryOnClick),
				app.Button().ID("home-page-meal-dialog-view-entries-button").Class("action-button", "secondary-action-button").Text("View Entries").OnClick(h.ViewEntriesOnClick),
				app.Button().ID("home-page-meal-dialog-edit-meal-button").Class("action-button", "tertiary-action-button").Text("Edit Meal").OnClick(h.EditMealOnClick),
				app.Button().ID("home-page-meal-dialog-delete-meal-button").Class("action-button", "danger-action-button").Text("Delete Meal").OnClick(h.DeleteMealOnClick),
			),

			app.Dialog().ID("home-page-confirm-delete-meal").Body(
				app.P().ID("home-page-confirm-delete-meal-text").Class("confirm-delete-text").Text("Are you sure you want to delete this meal?"),
				app.Div().ID("home-page-confirm-delete-meal-action-button-row").Class("action-button-row").Body(
					app.Button().ID("home-page-confirm-delete-meal-delete").Class("action-button", "danger-action-button").Text("Yes, Delete").OnClick(h.ConfirmDeleteMealOnClick),
					app.Button().ID("home-page-confirm-delete-meal-cancel").Class("action-button", "secondary-action-button").Text("No, Cancel").OnClick(h.CancelDeleteMealOnClick),
				),
			),
			app.Dialog().ID("home-page-options").OnClick(h.OptionsOnClick).Body(
				app.Form().ID("home-page-options-form").Class("form").OnSubmit(h.SaveOptions).OnClick(h.OptionsFormOnClick).Body(
					NewRadioChips("home-page-options-type", "What meal are you eating?", "Dinner", &h.options.Type, mealTypes...),
					NewCheckboxChips("home-page-options", "Who are you eating with?", map[string]bool{}, &h.usersOptions, usersStrings...),
					NewCheckboxChips("home-page-options-source", "What meal sources are okay?", map[string]bool{"Cooking": true, "Dine-In": true, "Takeout": true}, &h.options.Source, mealSources...),
					NewCheckboxChips("home-page-options-cuisine", "What cuisines are okay?", map[string]bool{"American": true}, &h.options.Cuisine, cuisines...),
					newCuisinesDialog("home-page", h.CuisinesDialogOnSave),
					NewRangeInput("home-page-options-taste", "How important is taste?", &h.options.TasteWeight),
					NewRangeInput("home-page-options-recency", "How important is recency?", &h.options.RecencyWeight),
					NewRangeInput("home-page-options-cost", "How important is cost?", &h.options.CostWeight),
					NewRangeInput("home-page-options-effort", "How important is effort?", &h.options.EffortWeight),
					NewRangeInput("home-page-options-healthiness", "How important is healthiness?", &h.options.HealthinessWeight),
					app.Div().ID("home-page-options-action-button-row").Class("action-button-row").Body(
						app.Button().ID("home-page-options-cancel-button").Class("secondary-action-button", "action-button").Type("button").Text("Cancel").OnClick(h.CancelOptions),
						app.Button().ID("home-page-options-save-button").Class("primary-action-button", "action-button").Type("submit").Text("Search"),
					),
				),
			),
		},
	}
}

func (h *home) OptionsOnClick(ctx app.Context, e app.Event) {
	// if the options dialog on click event is triggered, close the options because the dialog includes the whole page and a separate event will cancel this if they actually clicked on the dialog
	h.SaveOptions(ctx, e)
}

func (h *home) OptionsFormOnClick(ctx app.Context, e app.Event) {
	// cancel the closing of the dialog if they actually
	e.Call("stopPropagation")
}

func (h *home) New(ctx app.Context, e app.Event) {
	SetIsMealNew(true, ctx)
	SetCurrentMeal(Meal{}, ctx)
	ctx.Navigate("/meal")
}

func (h *home) PageOnClick(ctx app.Context, e app.Event) {
	// close meal dialog on page click (this will be stopped by another event if someone clicks on the meal dialog itself)
	app.Window().GetElementByID("home-page-meal-dialog").Call("close")
	h.currentMeal = Meal{}
}

func (h *home) MealOnClick(ctx app.Context, e app.Event, meal Meal) {
	e.Call("stopPropagation")
	h.currentMeal = meal
	SetCurrentMeal(meal, ctx)
	dialog := app.Window().GetElementByID("home-page-meal-dialog")
	dialog.Call("show")
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

	// ctx.Navigate("/meal")
}

func (h *home) MealDialogOnClick(ctx app.Context, e app.Event) {
	// stop the meal dialog from being closed by the page on click event
	e.Call("stopPropagation")
}

func (h *home) NewEntryOnClick(ctx app.Context, e app.Event) {
	entry := NewEntry(h.group, h.user, h.currentMeal, h.entriesForEachMeal[h.currentMeal.ID])
	SetIsEntryNew(true, ctx)
	SetCurrentEntry(entry, ctx)
	ctx.Navigate("/entry")
}

func (h *home) ViewEntriesOnClick(ctx app.Context, e app.Event) {
	ctx.Navigate("/entries")
}

func (h *home) EditMealOnClick(ctx app.Context, e app.Event) {
	SetIsMealNew(false, ctx)
	ctx.Navigate("/meal")
}

func (h *home) DeleteMealOnClick(ctx app.Context, e app.Event) {
	e.PreventDefault()
	app.Window().GetElementByID("home-page-confirm-delete-meal").Call("showModal")
}

func (h *home) ConfirmDeleteMealOnClick(ctx app.Context, e app.Event) {
	e.PreventDefault()

	_, err := DeleteMealAPI.Call(h.currentMeal.ID)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	SetCurrentMeal(Meal{}, ctx)
	meals, err := GetMealsAPI.Call(h.user.ID)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	h.meals = meals
	app.Window().GetElementByID("home-page-confirm-delete-meal").Call("close")
}

func (h *home) CancelDeleteMealOnClick(ctx app.Context, e app.Event) {
	e.PreventDefault()
	app.Window().GetElementByID("home-page-confirm-delete-meal").Call("close")
}

func (h *home) ShowOptions(ctx app.Context, e app.Event) {
	app.Window().GetElementByID("home-page-options").Call("showModal")
}

func (h *home) CancelOptions(ctx app.Context, e app.Event) {
	h.options = GetOptions(ctx)
	app.Window().GetElementByID("home-page-options").Call("close")
}

func (h *home) SaveOptions(ctx app.Context, e app.Event) {
	e.PreventDefault()

	for _, u := range h.users {
		h.options.Users[u.ID] = h.usersOptions[u.Name]
	}

	SetOptions(h.options, ctx)

	app.Window().GetElementByID("home-page-options").Call("close")

	h.SortMeals()
}

// func (h *home) CuisinesOnChange(ctx app.Context, event app.Event, val string) {
// 	if val == "+" {
// 		h.options.Cuisine[val] = false
// 		event.Get("target").Set("checked", false)
// 		app.Window().GetElementByID("home-page-cuisines-dialog").Call("showModal")
// 	}
// }

func (h *home) CuisinesDialogOnSave(ctx app.Context, event app.Event) {
	h.user = GetCurrentUser(ctx)
	h.options = h.options.RemoveInvalidCuisines(h.user.Cuisines)
}

func (h *home) SortMeals() {
	sort.Slice(h.meals, func(i, j int) bool {
		// prioritize meals with missing data, then score
		mealI := h.meals[i]
		entriesI := h.entriesForEachMeal[mealI.ID]
		iMissingData := entriesI.MissingData(h.user)
		mealJ := h.meals[j]
		entriesJ := h.entriesForEachMeal[mealJ.ID]
		jMissingData := entriesJ.MissingData(h.user)
		if iMissingData && !jMissingData {
			return true
		}
		if !iMissingData && jMissingData {
			return false
		}
		return h.meals[i].Score(h.entriesForEachMeal[h.meals[i].ID], h.options) > h.meals[j].Score(h.entriesForEachMeal[h.meals[j].ID], h.options)
	})
}
