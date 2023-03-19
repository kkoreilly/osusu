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
			h.SortMeals()
		},
		TitleElement: "Welcome, " + h.person.Name,
		Elements: []app.UI{
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
					NewCheckboxChips("home-page-options", "Who are you eating with?", map[string]bool{}, &h.peopleOptions, peopleString...),
					NewRadioChips("home-page-options-type", "What meal are you eating?", "Dinner", &h.options.Type, mealTypes...),
					NewCheckboxChips("home-page-options-source", "What meal sources are okay?", map[string]bool{"Cooking": true, "Dine-In": true, "Takeout": true}, &h.options.Source, mealSources...),
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
		},
	}
}

func (h *home) New(ctx app.Context, e app.Event) {
	meal, err := CreateMealAPI.Call(GetCurrentUser(ctx).ID)
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
		_, err := SignOutAPI.Call(GetCurrentUser(ctx))
		if err != nil {
			log.Println(err)
			return
		}
	}
	// if no error, we are no longer authenticated
	authenticated = time.UnixMilli(0)
	ctx.LocalStorage().Del("currentUser")

	ctx.Navigate("/signin")
}

func (h *home) CancelSignOut(ctx app.Context, event app.Event) {
	event.PreventDefault()
	app.Window().GetElementByID("home-page-confirm-sign-out").Call("close")
}
