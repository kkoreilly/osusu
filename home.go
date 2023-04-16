package main

import (
	"log"
	"sort"
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type home struct {
	app.Compo
	meals              Meals
	entriesForEachMeal map[int64]Entries // entries for each meal id
	user               User
	person             Person
	people             People
	options            Options
	peopleOptions      map[string]bool
}

func (h *home) Render() app.UI {
	peopleString := []string{}
	for _, p := range h.people {
		peopleString = append(peopleString, p.Name)
	}
	// need to copy to separate array from because append modifies the underlying array
	var cuisines = make([]string, len(h.user.Cuisines))
	copy(cuisines, h.user.Cuisines)
	return &Page{
		ID:                     "home",
		Title:                  "Home",
		Description:            "Osusu home",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			h.user = GetCurrentUser(ctx)
			cuisines, err := GetUserCuisinesAPI.Call(h.user.ID)
			if err != nil {
				CurrentPage.ShowErrorStatus(err)
				return
			}
			h.user.Cuisines = cuisines
			SetCurrentUser(h.user, ctx)

			h.person = GetCurrentPerson(ctx)

			people, err := GetPeopleAPI.Call(h.user.ID)
			if err != nil {
				CurrentPage.ShowErrorStatus(err)
				return
			}
			h.people = people

			h.options = GetOptions(ctx)
			if h.options.People == nil {
				h.options = DefaultOptions(h.user)
			}
			h.options = h.options.RemoveInvalidCuisines(h.user.Cuisines)
			SetOptions(h.options, ctx)
			h.peopleOptions = make(map[string]bool)
			for _, p := range h.people {
				if _, ok := h.options.People[p.ID]; !ok {
					h.options.People[p.ID] = true
				}
				h.peopleOptions[p.Name] = h.options.People[p.ID]
			}

			meals, err := GetMealsAPI.Call(h.user.ID)
			if err != nil {
				CurrentPage.ShowErrorStatus(err)
				return
			}
			h.meals = meals

			entries, err := GetEntriesAPI.Call(h.user.ID)
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
		TitleElement: "Welcome, " + h.person.Name,
		Elements: []app.UI{
			app.Div().ID("home-page-action-button-row").Class("action-button-row").Body(
				app.Button().ID("home-page-new-button").Class("secondary-action-button", "action-button").Text("New Meal").OnClick(h.New),
				app.Button().ID("home-page-options-button").Class("primary-action-button", "action-button").Text("Search").OnClick(h.ShowOptions),
			),
			app.Div().ID("home-page-meals-container").Body(
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
					missingData := entries.MissingData(h.person)
					return app.Button().ID("home-page-meal-"+si).Class("home-page-meal").DataSet("missing-data", missingData).Style("--color-h", colorH).Style("--score-percent", scoreText+"%").
						OnClick(func(ctx app.Context, e app.Event) { h.MealOnClick(ctx, e, meal) }).Body(
						app.Span().ID("home-page-meal-name-"+si).Class("home-page-meal-name").Text(meal.Name),
						app.Span().ID("home-page-meal-score-"+si).Class("home-page-meal-score").Text(scoreText),
					)
				}),
			),
			app.Dialog().ID("home-page-options").OnClick(h.OptionsOnClick).Body(
				app.Form().ID("home-page-options-form").Class("form").OnSubmit(h.SaveOptions).OnClick(h.OptionsFormOnClick).Body(
					NewRadioChips("home-page-options-type", "What meal are you eating?", "Dinner", &h.options.Type, mealTypes...),
					NewCheckboxChips("home-page-options", "Who are you eating with?", map[string]bool{}, &h.peopleOptions, peopleString...),
					NewCheckboxChips("home-page-options-source", "What meal sources are okay?", map[string]bool{"Cooking": true, "Dine-In": true, "Takeout": true}, &h.options.Source, mealSources...),
					NewCheckboxChips("home-page-options-cuisine", "What cuisines are okay?", map[string]bool{"American": true}, &h.options.Cuisine, append(cuisines, "+")...).SetOnChange(h.CuisinesOnChange),
					newCuisinesDialog("home-page", h.CuisinesDialogOnSave),
					NewRangeInput("home-page-options-taste", "How important is taste?", &h.options.TasteWeight),
					NewRangeInput("home-page-options-recency", "How important is recency?", &h.options.RecencyWeight),
					NewRangeInput("home-page-options-cost", "How important is cost?", &h.options.CostWeight),
					NewRangeInput("home-page-options-effort", "How important is effort?", &h.options.EffortWeight),
					NewRangeInput("home-page-options-healthiness", "How important is healthiness?", &h.options.HealthinessWeight),

					app.Div().ID("home-page-options-action-button-row").Class("action-button-row").Body(
						app.Input().ID("home-page-options-cancel-button").Class("secondary-action-button", "action-button").Type("button").Value("Cancel").OnClick(h.CancelOptions),
						app.Input().ID("home-page-options-save-button").Class("primary-action-button", "action-button").Type("submit").Value("Search"),
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
	log.Println("options form on click")
	// cancel the closing of the dialog if they actually
	e.Call("stopPropagation")
}

func (h *home) New(ctx app.Context, e app.Event) {
	meal, err := CreateMealAPI.Call(GetCurrentUser(ctx).ID)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	SetCurrentMeal(meal, ctx)
	ctx.Navigate("/meal")
}

func (h *home) MealOnClick(ctx app.Context, e app.Event, meal Meal) {
	SetCurrentMeal(meal, ctx)
	ctx.Navigate("/meal")
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

	for _, p := range h.people {
		h.options.People[p.ID] = h.peopleOptions[p.Name]
	}

	SetOptions(h.options, ctx)

	app.Window().GetElementByID("home-page-options").Call("close")

	h.SortMeals()
}

func (h *home) CuisinesOnChange(ctx app.Context, event app.Event, val string) {
	if val == "+" {
		h.options.Cuisine[val] = false
		event.Get("target").Set("checked", false)
		app.Window().GetElementByID("home-page-cuisines-dialog").Call("showModal")
	}
}

func (h *home) CuisinesDialogOnSave(ctx app.Context, event app.Event) {
	h.user = GetCurrentUser(ctx)
	h.options = h.options.RemoveInvalidCuisines(h.user.Cuisines)
}

func (h *home) SortMeals() {
	sort.Slice(h.meals, func(i, j int) bool {
		// prioritize meals with missing data, then score
		mealI := h.meals[i]
		entriesI := h.entriesForEachMeal[mealI.ID]
		iMissingData := entriesI.MissingData(h.person)
		mealJ := h.meals[j]
		entriesJ := h.entriesForEachMeal[mealJ.ID]
		jMissingData := entriesJ.MissingData(h.person)
		if iMissingData && !jMissingData {
			return true
		}
		if !iMissingData && jMissingData {
			return false
		}
		return h.meals[i].Score(h.entriesForEachMeal[h.meals[i].ID], h.options) > h.meals[j].Score(h.entriesForEachMeal[h.meals[j].ID], h.options)
	})
}
