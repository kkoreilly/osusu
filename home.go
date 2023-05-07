package main

import (
	"sort"
	"strconv"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type home struct {
	app.Compo
	group              Group
	user               User
	users              Users
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
	cuisines := []string{}
	for cuisine, val := range h.cuisinesInUse {
		if val {
			cuisines = append(cuisines, cuisine)
		}
	}
	width, _ := app.Window().Size()
	smallScreen := width <= 480
	return &Page{
		ID:                     "home",
		Title:                  "Home",
		Description:            "View, sort, and filter your meals.",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			SetReturnURL("/home", ctx)
			h.group = CurrentGroup(ctx)
			if h.group.Name == "" {
				Navigate("/groups", ctx)
			}
			h.user = CurrentUser(ctx)
			cuisines, err := GetGroupCuisinesAPI.Call(h.group.ID)
			if err != nil {
				CurrentPage.ShowErrorStatus(err)
				return
			}
			h.group.Cuisines = cuisines
			SetCurrentGroup(h.group, ctx)

			users, err := GetUsersAPI.Call(h.group.Members)
			if err != nil {
				CurrentPage.ShowErrorStatus(err)
				return
			}
			h.users = users

			h.options = GetOptions(ctx)
			if h.options.Users == nil {
				h.options = DefaultOptions(h.group)
			}
			h.options = h.options.RemoveInvalidCuisines(h.group.Cuisines)
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
			for cuisine := range h.options.Cuisine {
				if !h.cuisinesInUse[cuisine] {
					delete(h.options.Cuisine, cuisine)
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
			ButtonRow().ID("home-page").Buttons(
				Button().ID("home-page-new").Class("secondary").Icon("add").Text("New Meal").OnClick(h.NewMeal),
				Button().ID("home-page-search").Class("primary").Icon("search").Text("Search").OnClick(h.ShowOptions),
			),
			app.Table().ID("home-page-meals-table").Body(
				app.THead().ID("home-page-meals-table-header").Body(
					app.Tr().ID("home-page-meals-table-header-row").Body(
						app.Th().ID("home-page-meals-table-header-name").Text("Name"),
						app.Th().ID("home-page-meals-table-header-total").Text("Total"),
						app.Th().ID("home-page-meals-table-header-taste").Text("Taste"),
						app.Th().ID("home-page-meals-table-header-recency").Text("New"),
						app.Th().ID("home-page-meals-table-header-cost").Text("Cost"),
						app.Th().ID("home-page-meals-table-header-effort").Text("Effort"),
						app.Th().ID("home-page-meals-table-header-healthiness").Text("Health"),
						app.If(!smallScreen,
							app.Th().ID("home-page-meals-table-header-cuisines").Text("Cuisines"),
							app.Th().ID("home-page-meals-table-header-description").Text("Description"),
						),
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
						colorH := strconv.Itoa((score.Total * 12) / 10)
						scoreText := strconv.Itoa(score.Total)
						missingData := entries.MissingData(h.user)
						isCurrentMeal := meal.ID == h.currentMeal.ID
						return app.Tr().ID("home-page-meal-"+si).Class("home-page-meal").DataSet("missing-data", missingData).DataSet("current-meal", isCurrentMeal).Style("--color-h", colorH).Style("--score", scoreText+"%").
							OnClick(func(ctx app.Context, e app.Event) { h.MealOnClick(ctx, e, meal) }).Body(
							app.Td().ID("home-page-meal-name-"+si).Class("home-page-meal-name").Text(meal.Name),
							MealScore("home-page-meal-total-"+si, "home-page-meal-total", score.Total),
							MealScore("home-page-meal-taste-"+si, "home-page-meal-taste", score.Taste),
							MealScore("home-page-meal-recency-"+si, "home-page-meal-recency", score.Recency),
							MealScore("home-page-meal-cost-"+si, "home-page-meal-cost", score.Cost),
							MealScore("home-page-meal-effort-"+si, "home-page-meal-effort", score.Effort),
							MealScore("home-page-meal-healthiness-"+si, "home-page-meal-healthiness", score.Healthiness),
							app.If(!smallScreen,
								app.Td().ID("home-page-meal-cuisines-"+si).Class("home-page-meal-cuisines").Text(meal.CuisineString()),
								app.Td().ID("home-page-meal-description-"+si).Class("home-page-meal-description").Text(meal.Description),
							),
						)
					}),
				),
			),
			app.Dialog().ID("home-page-meal-dialog").OnClick(h.MealDialogOnClick).Body(
				Button().ID("home-page-meal-dialog-new-entry").Class("primary").Icon("add").Text("New Entry").OnClick(h.NewEntry),
				Button().ID("home-page-meal-dialog-view-entries").Class("secondary").Icon("visibility").Text("View Entries").OnClick(h.ViewEntries),
				Button().ID("home-page-meal-dialog-edit-meal").Class("tertiary").Icon("edit").Text("Edit Meal").OnClick(h.EditMeal),
				Button().ID("home-page-meal-dialog-delete-meal").Class("danger").Icon("delete").Text("Delete Meal").OnClick(h.DeleteMeal),
			),

			app.Dialog().ID("home-page-confirm-delete-meal").Class("modal").Body(
				app.P().ID("home-page-confirm-delete-meal-text").Class("confirm-delete-text").Text("Are you sure you want to delete this meal?"),
				ButtonRow().ID("home-page-confirm-delete-meal").Buttons(
					Button().ID("home-page-confirm-delete-meal-delete").Class("danger").Icon("delete").Text("Yes, Delete").OnClick(h.ConfirmDeleteMeal),
					Button().ID("home-page-confirm-delete-meal-cancel").Class("secondary").Icon("cancel").Text("No, Cancel").OnClick(h.CancelDeleteMeal),
				),
			),
			app.Dialog().ID("home-page-options").Class("modal").OnClick(h.OptionsOnClick).Body(
				app.Form().ID("home-page-options-form").Class("form").OnSubmit(h.SaveOptions).OnClick(h.OptionsFormOnClick).Body(
					RadioChips().ID("home-page-options-type").Label("What meal are you eating?").Default("Dinner").Value(&h.options.Type).Options(mealTypes...),
					CheckboxChips().ID("home-page-options-users").Label("Who are you eating with?").Value(&h.usersOptions).Options(usersStrings...),
					CheckboxChips().ID("home-page-options-source").Label("What meal sources are okay?").Default(map[string]bool{"Cooking": true, "Dine-In": true, "Takeout": true}).Value(&h.options.Source).Options(mealSources...),
					CheckboxChips().ID("home-page-options-cuisine").Label("What cuisines are okay?").Value(&h.options.Cuisine).Options(cuisines...),
					newCuisinesDialog("home-page", h.CuisinesDialogOnSave),
					RangeInput().ID("home-page-options-taste").Label("How important is taste?").Value(&h.options.TasteWeight),
					RangeInput().ID("home-page-options-recency").Label("How important is recency?").Value(&h.options.RecencyWeight),
					RangeInput().ID("home-page-options-cost").Label("How important is cost?").Value(&h.options.CostWeight),
					RangeInput().ID("home-page-options-effort").Label("How important is effort?").Value(&h.options.EffortWeight),
					RangeInput().ID("home-page-options-healthiness").Label("How important is healthiness?").Value(&h.options.HealthinessWeight),
					ButtonRow().ID("home-page-options").Buttons(
						Button().ID("home-page-options-cancel").Class("secondary").Icon("cancel").Text("Cancel").OnClick(h.CancelOptions),
						Button().ID("home-page-options-save").Class("primary").Type("submit").Icon("search").Text("Search"),
					),
				),
			),
		},
	}
}

// MealScore returns a table cell with a score pie circle containing score information for a meal or entry
func MealScore(id string, class string, score int) app.UI {
	return app.Td().ID(id).Class("meal-score", class).Style("--score", strconv.Itoa(score)).Style("--color-h", strconv.Itoa(score*12/10)).Body(
		app.Div().ID(id+"-circle").Class("meal-score-circle", "pie", class+"-circle").Text(score),
	)
}

func (h *home) OptionsOnClick(ctx app.Context, e app.Event) {
	// if the options dialog on click event is triggered, close the options because the dialog includes the whole page and a separate event will cancel this if they actually clicked on the dialog
	h.SaveOptions(ctx, e)
}

func (h *home) OptionsFormOnClick(ctx app.Context, e app.Event) {
	// cancel the closing of the dialog if they actually
	e.Call("stopPropagation")
}

func (h *home) NewMeal(ctx app.Context, e app.Event) {
	SetIsMealNew(true, ctx)
	SetCurrentMeal(Meal{}, ctx)
	Navigate("/meal", ctx)
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
	if dialog.Get("open").Bool() {
		ctx.Dispatch(func(ctx app.Context) {
			dialog.Call("close")
		})
		ctx.Defer(func(ctx app.Context) {
			time.Sleep(250 * time.Millisecond)
			h.UpdateMealDialogPosition(ctx, e, dialog)
			dialog.Call("show")
		})
		return
	}
	h.UpdateMealDialogPosition(ctx, e, dialog)
	dialog.Call("show")

}

func (h *home) UpdateMealDialogPosition(ctx app.Context, e app.Event, dialog app.Value) {
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

func (h *home) MealDialogOnClick(ctx app.Context, e app.Event) {
	// stop the meal dialog from being closed by the page on click event
	e.Call("stopPropagation")
}

func (h *home) NewEntry(ctx app.Context, e app.Event) {
	entry := NewEntry(h.group, h.user, h.currentMeal, h.entriesForEachMeal[h.currentMeal.ID])
	SetIsEntryNew(true, ctx)
	SetCurrentEntry(entry, ctx)
	Navigate("/entry", ctx)
}

func (h *home) ViewEntries(ctx app.Context, e app.Event) {
	Navigate("/entries", ctx)
}

func (h *home) EditMeal(ctx app.Context, e app.Event) {
	SetIsMealNew(false, ctx)
	Navigate("/meal", ctx)
}

func (h *home) DeleteMeal(ctx app.Context, e app.Event) {
	e.PreventDefault()
	app.Window().GetElementByID("home-page-confirm-delete-meal").Call("showModal")
}

func (h *home) ConfirmDeleteMeal(ctx app.Context, e app.Event) {
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

func (h *home) CancelDeleteMeal(ctx app.Context, e app.Event) {
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

func (h *home) CuisinesDialogOnSave(ctx app.Context, event app.Event) {
	h.user = CurrentUser(ctx)
	h.options = h.options.RemoveInvalidCuisines(h.group.Cuisines)
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
		return h.meals[i].Score(h.entriesForEachMeal[h.meals[i].ID], h.options).Total > h.meals[j].Score(h.entriesForEachMeal[h.meals[j].ID], h.options).Total
	})
}
