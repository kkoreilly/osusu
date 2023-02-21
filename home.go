package main

import (
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type home struct {
	app.Compo
	meals         Meals
	person        Person
	people        People
	options       Options
	peopleOptions map[string]bool
	// test          map[string]bool
}

func (h *home) Render() app.UI {
	peopleString := []string{}
	for _, p := range h.people {
		peopleString = append(peopleString, p.Name)
	}
	return app.Div().Body(
		// NewCheckboxChips("test", "What fruits do you like?", map[string]bool{"Banana": true, "Peach": true}, &h.test, "Apple", "Banana", "Mango", "Peach"),
		app.H1().ID("home-page-title").Class("page-title").Text("Welcome, "+h.person.Name),
		app.Div().ID("home-page-action-button-row").Class("action-button-row").Body(
			app.Button().ID("home-page-sign-out-button").Class("red-action-button", "action-button").Text("Sign Out").OnClick(h.InitialSignOut),
			app.Button().ID("home-page-options-button").Class("white-action-button", "action-button").Text("Options").OnClick(h.ShowOptions),
			app.Button().ID("home-page-new-button").Class("blue-action-button", "action-button").Text("New").OnClick(h.New),
		),
		app.Hr(),
		app.Div().ID("home-page-meals-container").Body(
			app.Range(h.meals).Slice(func(i int) app.UI {
				meal := h.meals[i]
				if meal.Type != h.options.Type || !h.options.Source[meal.Source] {
					return app.Text("")
				}
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
				NewRangeInput("home-page-options-taste", "How important is taste?", &h.options.TasteWeight),
				NewRangeInput("home-page-options-recency", "How important is recency?", &h.options.RecencyWeight),
				NewRangeInput("home-page-options-cost", "How important is cost?", &h.options.CostWeight),
				NewRangeInput("home-page-options-effort", "How important is effort?", &h.options.EffortWeight),
				NewRangeInput("home-page-options-healthiness", "How important is healthiness?", &h.options.HealthinessWeight),
				// &Input[int]{ID: "home-page-options-taste", Label: "How important is taste?", InputClass: "input-range", Type: "range", Value: &h.options.TasteWeight, ValueFunc: ValueFuncInt},
				// &Input[int]{ID: "home-page-options-recency", Label: "How important is recency?", InputClass: "input-range", Type: "range", Value: &h.options.RecencyWeight, ValueFunc: ValueFuncInt},
				// &Input[int]{ID: "home-page-options-cost", Label: "How important is cost?", InputClass: "input-range", Type: "range", Value: &h.options.CostWeight, ValueFunc: ValueFuncInt},
				// &Input[int]{ID: "home-page-options-effort", Label: "How important is effort?", InputClass: "input-range", Type: "range", Value: &h.options.EffortWeight, ValueFunc: ValueFuncInt},
				// &Input[int]{ID: "home-page-options-healthiness", Label: "How important is healthiness?", InputClass: "input-range", Type: "range", Value: &h.options.HealthinessWeight, ValueFunc: ValueFuncInt},
				// app.Label().ID("home-page-options-taste-label").Class("input-label").For("home-page-options-taste-input").Text("How important is taste?"),
				// app.Input().ID("home-page-options-taste-input").Class("input", "input-range").Type("range").Min(0).Max(100).Value(h.options.TasteWeight),
				// app.Label().ID("home-page-options-recency-label").Class("input-label").For("home-page-options-recency-input").Text("How important is recency?"),
				// app.Input().ID("home-page-options-recency-input").Class("input", "input-range").Type("range").Min(0).Max(100).Value(h.options.RecencyWeight),
				// app.Label().ID("home-page-options-cost-label").Class("input-label").For("home-page-options-cost-input").Text("How important is cost?"),
				// app.Input().ID("home-page-options-cost-input").Class("input", "input-range").Type("range").Min(0).Max(100).Value(h.options.CostWeight),
				// app.Label().ID("home-page-options-effort-label").Class("input-label").For("home-page-options-effort-input").Text("How important is effort?"),
				// app.Input().ID("home-page-options-effort-input").Class("input", "input-range").Type("range").Min(0).Max(100).Value(h.options.EffortWeight),
				// app.Label().ID("home-page-options-healthiness-label").Class("input-label").For("home-page-options-healthiness-input").Text("How important is healthiness?"),
				// app.Input().ID("home-page-options-healthiness-input").Class("input", "input-range").Type("range").Min(0).Max(100).Value(h.options.HealthinessWeight),
				NewCheckboxChips("home-page-options", "Who are you eating with?", map[string]bool{}, &h.peopleOptions, peopleString...),
				// app.Label().ID("home-page-options-people-label").Class("input-label").For("home-page-options-people-container").Text("Who are you eating with?"),
				// app.Div().ID("home-page-options-people-container").Class("chips-container").Body(
				// 	app.Range(h.people).Slice(func(i int) app.UI {
				// 		si := strconv.Itoa(i)
				// 		p := h.people[i]
				// 		checked := h.options.People[p.ID]
				// 		return Chip("home-page-options-person-"+si, "checkbox", "home-page-options-person-"+si, p.Name, checked)
				// 	}),
				// ),
				NewRadioChips("home-page-options-type", "What meal are you eating?", "Dinner", &h.options.Type, mealTypes...),
				// app.Label().ID("home-page-options-type-label").Class("input-label").For("home-page-options-type-container").Text("What meal are you eating?"),
				// app.Div().ID("home-page-options-type-container").Class("chips-container").Body(
				// 	Chip("home-page-options-type-breakfast", "radio", "home-page-options-type", "Breakfast", h.options.Type == "Breakfast"),
				// 	Chip("home-page-options-type-lunch", "radio", "home-page-options-type", "Lunch", h.options.Type == "Lunch"),
				// 	Chip("home-page-options-type-dinner", "radio", "home-page-options-type", "Dinner", h.options.Type == "Dinner"),
				// ),
				NewCheckboxChips("home-page-options-source", "What meal sources are okay?", map[string]bool{"Cooking": true, "Dine-In": true, "Takeout": true}, &h.options.Source, mealSources...),
				// app.Label().ID("home-page-options-source-label").Class("input-label").For("home-page-options-source-container").Text("What meal sources are okay?"),
				// app.Div().ID("home-page-options-source-container").Class("chips-container").Body(
				// 	Chip("home-page-options-source-cooking", "checkbox", "home-page-options-source", "Cooking", h.options.Source["Cooking"]),
				// 	Chip("home-page-options-source-dine-in", "checkbox", "home-page-options-source", "Dine-In", h.options.Source["Dine-In"]),
				// 	Chip("home-page-options-source-takeout", "checkbox", "home-page-options-source", "Takeout", h.options.Source["Takeout"]),
				// ),
				app.Div().ID("home-page-options-action-button-row").Class("action-button-row").Body(
					app.Input().ID("home-page-options-cancel-button").Class("white-action-button", "action-button").Type("button").Value("Cancel").OnClick(h.CancelOptions),
					app.Input().ID("home-page-options-save-button").Class("blue-action-button", "action-button").Type("submit").Value("Save"),
				),
			),
		),
		app.Dialog().ID("home-page-confirm-sign-out").Body(
			app.P().ID("home-page-confirm-sign-out-text").Text("Are you sure you want to sign out?"),
			app.Div().ID("home-page-confirm-sign-out-action-button-row").Class("action-button-row").Body(
				app.Button().ID("home-page-confirm-sign-out-sign-out").Class("action-button", "red-action-button").Text("Yes, Sign Out").OnClick(h.ConfirmSignOut),
				app.Button().ID("edit-page-confirm-sign-out-cancel").Class("action-button", "white-action-button").Text("No, Cancel").OnClick(h.CancelSignOut),
			),
		),
	)
}

// // Chip returns a new chip element with the given id, input type, name, value, checked value, and classes
// func Chip(id string, inputType string, name string, value string, checked bool, class ...string) app.UI {
// 	return app.Label().ID(id+"-chip-label").Class("chip-label").For(id+"-chip-input").Body(
// 		app.Input().ID(id+"-chip-input").Class("chip-input").Type(inputType).Name(name).Checked(checked).Value(value),
// 		app.Text(value),
// 	)
// }

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
		h.options = Options{50, 50, 50, 50, 50, make(map[int]bool), "Dinner", map[string]bool{"Cooking": true, "Dine-In": true, "Takeout": true}}
	}
	for _, p := range h.people {
		if _, ok := h.options.People[p.ID]; !ok {
			h.options.People[p.ID] = true
		}
		h.peopleOptions[p.Name] = h.options.People[p.ID]
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
	// log.Println(h.test)
}

func (h *home) CancelOptions(ctx app.Context, e app.Event) {
	app.Window().GetElementByID("home-page-options").Call("close")
}

func (h *home) SaveOptions(ctx app.Context, e app.Event) {
	e.PreventDefault()

	// h.options.CostWeight = app.Window().GetElementByID("home-page-options-cost-input").Get("valueAsNumber").Int()
	// h.options.EffortWeight = app.Window().GetElementByID("home-page-options-effort-input").Get("valueAsNumber").Int()
	// h.options.HealthinessWeight = app.Window().GetElementByID("home-page-options-healthiness-input").Get("valueAsNumber").Int()
	// h.options.TasteWeight = app.Window().GetElementByID("home-page-options-taste-input").Get("valueAsNumber").Int()
	// h.options.RecencyWeight = app.Window().GetElementByID("home-page-options-recency-input").Get("valueAsNumber").Int()
	// for i, p := range h.people {
	// 	checked := app.Window().GetElementByID("home-page-options-person-" + strconv.Itoa(i) + "-chip-input").Get("checked").Bool()
	// 	h.options.People[p.ID] = checked
	// }
	for _, p := range h.people {
		// log.Println(p.ID, p.Name, h.peopleOptions[p.Name])
		h.options.People[p.ID] = h.peopleOptions[p.Name]
	}
	// elem := app.Window().Get("document").Call("querySelector", `input[name="home-page-options-type"]:checked`)
	// if !elem.IsNull() {
	// 	h.options.Type = elem.Get("value").String()
	// }
	// h.options.Source["Cooking"] = app.Window().GetElementByID("home-page-options-source-cooking-chip-input").Get("checked").Bool()
	// h.options.Source["Dine-In"] = app.Window().GetElementByID("home-page-options-source-dine-in-chip-input").Get("checked").Bool()
	// h.options.Source["Takeout"] = app.Window().GetElementByID("home-page-options-source-takeout-chip-input").Get("checked").Bool()

	SetOptions(h.options, ctx)

	app.Window().GetElementByID("home-page-options").Call("close")

	h.SortMeals()
}

func (h *home) SortMeals() {
	sort.Slice(h.meals, func(i, j int) bool {
		return h.meals[i].Score(h.options) > h.meals[j].Score(h.options)
	})
}

func (h *home) InitialSignOut(ctx app.Context, event app.Event) {
	event.PreventDefault()
	app.Window().GetElementByID("home-page-confirm-sign-out").Call("showModal")
}

func (h *home) ConfirmSignOut(ctx app.Context, event app.Event) {
	event.PreventDefault()
	user := GetCurrentUser(ctx)
	if user.Session != "" {
		err := SignOutRequest(GetCurrentUser(ctx))
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		authenticated = time.UnixMilli(0)
	}
	ctx.LocalStorage().Del("currentUser")

	ctx.Navigate("/signin")
}

func (h *home) CancelSignOut(ctx app.Context, event app.Event) {
	event.PreventDefault()
	app.Window().GetElementByID("home-page-confirm-sign-out").Call("close")
}
