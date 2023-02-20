package main

import (
	"log"
	"sort"
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type home struct {
	app.Compo
	meals   Meals
	person  Person
	people  People
	options Options
}

func (h *home) Render() app.UI {
	people := make(map[string]bool)
	for k, v := range h.options.People {
		people[strconv.Itoa(k)] = v
	}
	return app.Div().Body(
		app.H1().ID("home-page-title").Class("page-title").Text("Welcome, "+h.person.Name),
		app.Div().ID("home-page-action-button-row").Class("action-button-row").Body(
			app.Button().ID("home-page-options-button").Class("white-action-button", "action-button").Text("Options").OnClick(h.ShowOptions),
			app.Button().ID("home-page-new-button").Class("blue-action-button", "action-button").Text("New").OnClick(h.New),
		),
		app.Hr(),
		app.Div().ID("home-page-meals-container").Body(
			app.Range(h.meals).Slice(func(i int) app.UI {
				meal := h.meals[i]
				si := strconv.Itoa(i)
				score := meal.Score(h.options)
				colorH := strconv.Itoa((score * 12) / 10)
				scoreText := strconv.Itoa(score)
				_, tasteSet := meal.Taste[h.person.ID]
				return app.Div().ID("home-page-meal-"+si).Class("home-page-meal").Style("--color-h", colorH).Style("--score-percent", scoreText+"%").
					OnClick(func(ctx app.Context, e app.Event) { h.MealOnClick(ctx, e, meal) }).Body(
					app.Span().ID("home-page-meal-name-"+si).Class("home-page-meal-name").Text(meal.Name).DataSet("taste-set", tasteSet),
					app.Span().ID("home-page-meal-score-"+si).Class("home-page-meal-score").Text(scoreText),
				)
			}),
		),
		app.Dialog().ID("home-page-options").Body(
			app.Form().ID("home-page-options-form").Class("form").OnSubmit(h.SaveOptions).Body(
				app.Label().ID("home-page-options-taste-label").Class("input-label").For("home-page-options-taste-input").Text("Taste Weight:"),
				app.Input().ID("home-page-options-taste-input").Class("input", "input-range").Type("range").Min(0).Max(100).Value(h.options.TasteWeight),
				app.Label().ID("home-page-options-recency-label").Class("input-label").For("home-page-options-recency-input").Text("Recency Weight:"),
				app.Input().ID("home-page-options-recency-input").Class("input", "input-range").Type("range").Min(0).Max(100).Value(h.options.RecencyWeight),
				app.Label().ID("home-page-options-cost-label").Class("input-label").For("home-page-options-cost-input").Text("Cost Weight:"),
				app.Input().ID("home-page-options-cost-input").Class("input", "input-range").Type("range").Min(0).Max(100).Value(h.options.CostWeight),
				app.Label().ID("home-page-options-effort-label").Class("input-label").For("home-page-options-effort-input").Text("Effort Weight:"),
				app.Input().ID("home-page-options-effort-input").Class("input", "input-range").Type("range").Min(0).Max(100).Value(h.options.EffortWeight),
				app.Label().ID("home-page-options-healthiness-label").Class("input-label").For("home-page-options-healthiness-input").Text("Healthiness Weight:"),
				app.Input().ID("home-page-options-healthiness-input").Class("input", "input-range").Type("range").Min(0).Max(100).Value(h.options.HealthinessWeight),
				app.Label().ID("home-page-options-people-label").Class("input-label").For("home-page-options-people-container").Text("Who Are You Eating With?"),
				app.Div().ID("home-page-options-people-container").Class("chips-container").Body(
					app.Range(h.people).Slice(func(i int) app.UI {
						si := strconv.Itoa(i)
						p := h.people[i]
						checked := h.options.People[p.ID]
						return Chip("home-page-options-person-"+si, "checkbox", "home-page-options-person-"+si, p.Name, checked)
					}),
				),
				app.Label().ID("home-page-options-type-label").Class("input-label").For("home-page-options-type-container").Text("What Type of Meal Are You Having?"),
				app.Div().ID("home-page-options-type-container").Class("chips-container").Body(
					Chip("home-page-options-type-breakfast", "radio", "home-page-options-type", "Breakfast", h.options.Type == "Breakfast"),
					Chip("home-page-options-type-lunch", "radio", "home-page-options-type", "Lunch", h.options.Type == "Lunch"),
					Chip("home-page-options-type-dinner", "radio", "home-page-options-type", "Dinner", h.options.Type == "Dinner"),
				),
				app.Div().ID("home-page-options-action-button-row").Class("action-button-row").Body(
					app.Input().ID("home-page-options-cancel-button").Class("white-action-button", "action-button").Type("button").Value("Cancel").OnClick(h.CancelOptions),
					app.Input().ID("home-page-options-save-button").Class("blue-action-button", "action-button").Type("submit").Value("Save"),
				),
			),
		),
	)
}

// Chip returns a new chip element with the given id, input type, name, value, checked value, and classes
func Chip(id string, inputType string, name string, value string, checked bool, class ...string) app.UI {
	return app.Label().ID(id+"-chip-label").Class("chip-label").For(id+"-chip-input").Body(
		app.Input().ID(id+"-chip-input").Class("chip-input").Type(inputType).Name(name).Checked(checked).Value(value),
		app.Text(value),
	)
}

func (h *home) New(ctx app.Context, e app.Event) {
	meal, err := CreateMealRequest(GetCurrentUser(ctx))
	if err != nil {
		log.Println(err)
		return
	}
	meal.Taste[h.person.ID] = 50
	SetCurrentMeal(meal, ctx)
	ctx.Navigate("/edit")
}

func (h *home) MealOnClick(ctx app.Context, e app.Event, meal Meal) {
	SetCurrentMeal(meal, ctx)
	ctx.Navigate("/edit")
}

func (h *home) OnNav(ctx app.Context) {
	if Authenticate(true, ctx) {
		return
	}
	h.person = GetCurrentPerson(ctx)

	people, err := GetPeopleRequest(GetCurrentUser(ctx))
	if err != nil {
		log.Println(err)
		return
	}
	h.people = people

	h.options = GetOptions(ctx)
	if h.options.People == nil {
		h.options = Options{50, 50, 50, 50, 50, make(map[int]bool), "", []string{}}
	}
	for _, p := range h.people {
		if _, ok := h.options.People[p.ID]; !ok {
			h.options.People[p.ID] = true
		}
	}

	meals, err := GetMealsRequest(GetCurrentUser(ctx))
	if err != nil {
		log.Println(err)
		return
	}
	h.meals = meals
	h.SortMeals()
}

func (h *home) ShowOptions(ctx app.Context, e app.Event) {
	app.Window().GetElementByID("home-page-options").Call("showModal")
}

func (h *home) CancelOptions(ctx app.Context, e app.Event) {
	app.Window().GetElementByID("home-page-options").Call("close")
}

func (h *home) SaveOptions(ctx app.Context, e app.Event) {
	e.PreventDefault()

	h.options.CostWeight = app.Window().GetElementByID("home-page-options-cost-input").Get("valueAsNumber").Int()
	h.options.EffortWeight = app.Window().GetElementByID("home-page-options-effort-input").Get("valueAsNumber").Int()
	h.options.HealthinessWeight = app.Window().GetElementByID("home-page-options-healthiness-input").Get("valueAsNumber").Int()
	h.options.TasteWeight = app.Window().GetElementByID("home-page-options-taste-input").Get("valueAsNumber").Int()
	h.options.RecencyWeight = app.Window().GetElementByID("home-page-options-recency-input").Get("valueAsNumber").Int()
	for i, p := range h.people {
		checked := app.Window().GetElementByID("home-page-options-person-" + strconv.Itoa(i) + "-chip-input").Get("checked").Bool()
		h.options.People[p.ID] = checked
	}

	SetOptions(h.options, ctx)

	app.Window().GetElementByID("home-page-options").Call("close")

	h.SortMeals()
}

func (h *home) SortMeals() {
	sort.Slice(h.meals, func(i, j int) bool {
		return h.meals[i].Score(h.options) > h.meals[j].Score(h.options)
	})
}
