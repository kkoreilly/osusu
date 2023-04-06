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
	entriesForEachMeal map[int]Entries // entries for each meal id
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
	return &Page{
		ID:                     "home",
		Title:                  "Home",
		Description:            "MealRec home",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			h.person = GetCurrentPerson(ctx)

			people, err := GetPeopleAPI.Call(GetCurrentUser(ctx).ID)
			if err != nil {
				log.Println(err)
				return
			}
			h.people = people

			h.options = GetOptions(ctx)
			if h.options.People == nil {
				h.options = Options{50, 50, 50, 50, 50, make(map[int]bool), "Dinner", map[string]bool{"Cooking": true, "Dine-In": true, "Takeout": true}}
			}
			h.peopleOptions = make(map[string]bool)
			for _, p := range h.people {
				if _, ok := h.options.People[p.ID]; !ok {
					h.options.People[p.ID] = true
				}
				h.peopleOptions[p.Name] = h.options.People[p.ID]
			}

			meals, err := GetMealsAPI.Call(GetCurrentUser(ctx).ID)
			if err != nil {
				log.Println(err)
				return
			}
			h.meals = meals

			entries, err := GetEntriesAPI.Call(GetCurrentUser(ctx).ID)
			if err != nil {
				log.Println(err)
				return
			}
			h.entriesForEachMeal = make(map[int]Entries)
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
				app.Button().ID("home-page-options-button").Class("secondary-action-button", "action-button").Text("Options").OnClick(h.ShowOptions),
				app.Button().ID("home-page-new-button").Class("primary-action-button", "action-button").Text("New").OnClick(h.New),
			),
			app.Div().ID("home-page-meals-container").Body(
				app.Range(h.meals).Slice(func(i int) app.UI {
					meal := h.meals[i]
					si := strconv.Itoa(i)
					entries := h.entriesForEachMeal[meal.ID]

					// check if at least one entry satisfies the type and source requirements.
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

					score := meal.Score(entries, h.options)
					colorH := strconv.Itoa((score * 12) / 10)
					scoreText := strconv.Itoa(score)
					missingData := entries.MissingData(h.person)
					return app.Div().ID("home-page-meal-"+si).Class("home-page-meal").DataSet("missing-data", missingData).Style("--color-h", colorH).Style("--score-percent", scoreText+"%").
						OnClick(func(ctx app.Context, e app.Event) { h.MealOnClick(ctx, e, meal) }).Body(
						app.Span().ID("home-page-meal-name-"+si).Class("home-page-meal-name").Text(meal.Name),
						app.Span().ID("home-page-meal-score-"+si).Class("home-page-meal-score").Text(scoreText),
					)
				}),
			),
			app.Dialog().ID("home-page-options").Body(
				app.Form().ID("home-page-options-form").Class("form").OnSubmit(h.SaveOptions).Body(
					NewRangeInput("home-page-options-taste", "How important is taste?", &h.options.TasteWeight),
					NewRangeInput("home-page-options-recency", "How important is recency?", &h.options.RecencyWeight),
					NewRangeInput("home-page-options-cost", "How important is cost?", &h.options.CostWeight),
					NewRangeInput("home-page-options-effort", "How important is effort?", &h.options.EffortWeight),
					NewRangeInput("home-page-options-healthiness", "How important is healthiness?", &h.options.HealthinessWeight),
					NewCheckboxChips("home-page-options", "Who are you eating with?", map[string]bool{}, &h.peopleOptions, peopleString...),
					NewRadioChips("home-page-options-type", "What meal are you eating?", "Dinner", &h.options.Type, mealTypes...),
					NewCheckboxChips("home-page-options-source", "What meal sources are okay?", map[string]bool{"Cooking": true, "Dine-In": true, "Takeout": true}, &h.options.Source, mealSources...),
					app.Div().ID("home-page-options-action-button-row").Class("action-button-row").Body(
						app.Input().ID("home-page-options-cancel-button").Class("secondary-action-button", "action-button").Type("button").Value("Cancel").OnClick(h.CancelOptions),
						app.Input().ID("home-page-options-save-button").Class("primary-action-button", "action-button").Type("submit").Value("Save"),
					),
				),
			),
		},
	}
}

func (h *home) New(ctx app.Context, e app.Event) {
	meal, err := CreateMealAPI.Call(GetCurrentUser(ctx).ID)
	if err != nil {
		log.Println(err)
		return
	}
	// meal.Taste[h.person.ID] = 50
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

func (h *home) SortMeals() {
	sort.Slice(h.meals, func(i, j int) bool {
		return h.meals[i].Score(h.entriesForEachMeal[h.meals[i].ID], h.options) > h.meals[j].Score(h.entriesForEachMeal[h.meals[j].ID], h.options)
	})
}
