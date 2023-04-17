package main

import (
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Meal is a struct that represents the data of a meal
type Meal struct {
	ID          int64
	UserID      int64
	Name        string
	Description string
	Cuisine     []string
}

// Meals is a slice that represents multiple meals
type Meals []Meal

// Score produces a score from 0 to 100 for the meal based on the given entries and options
func (m Meal) Score(entries Entries, options Options) int {
	entriesSum := 0
	var latestDate time.Time
	for _, entry := range entries {
		entriesSum += entry.Score(options)
		if entry.Date.After(latestDate) {
			latestDate = entry.Date
		}
	}
	recencyScore := int(2 * time.Now().Truncate(time.Hour*24).UTC().Sub(latestDate) / (time.Hour * 24))
	if recencyScore > 100 {
		recencyScore = 100
	}
	// add up all of the weights except recency and multiply all of the scores except for recency by them to make the other weights affect how much recency matters
	weightsToal := options.CostWeight + options.EffortWeight + options.HealthinessWeight + options.TasteWeight
	sum := weightsToal*entriesSum + options.RecencyWeight*recencyScore
	den := weightsToal*len(entries) + options.RecencyWeight
	if den == 0 {
		return 0
	}
	return sum / den
}

// RemoveInvalidCuisines returns the the meal with all invalid cuisines removed, using the given cuisine options
func (m Meal) RemoveInvalidCuisines(cuisines []string) Meal {
	res := []string{}
	for _, mealCuisine := range m.Cuisine {
		for _, cuisineOption := range cuisines {
			if mealCuisine == cuisineOption {
				res = append(res, mealCuisine)
			}
		}
	}
	m.Cuisine = res
	return m
}

// SetCurrentMeal sets the current meal state value to the given meal, using the given context
func SetCurrentMeal(meal Meal, ctx app.Context) {
	ctx.SetState("currentMeal", meal, app.Persist)
}

// GetCurrentMeal gets and returns the current meal state value, using the given context
func GetCurrentMeal(ctx app.Context) Meal {
	var meal Meal
	ctx.GetState("currentMeal", &meal)
	return meal
}

var (
	mealTypes   = []string{"Breakfast", "Lunch", "Dinner"}
	mealSources = []string{"Cooking", "Dine-In", "Takeout"}
	// mealCuisines = []string{"American", "Chinese", "Indian", "Italian", "Japanese", "Korean", "Mexican", "+"}
)

type meal struct {
	app.Compo
	user    User
	meal    Meal
	person  Person
	cuisine map[string]bool
}

func (m *meal) Render() app.UI {
	// need to copy to separate array from because append modifies the underlying array
	var cuisines = make([]string, len(m.user.Cuisines))
	copy(cuisines, m.user.Cuisines)
	return &Page{
		ID:                     "meal",
		Title:                  "Edit Meal",
		Description:            "Edit a meal",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			m.user = GetCurrentUser(ctx)
			m.person = GetCurrentPerson(ctx)
			m.meal = GetCurrentMeal(ctx)

			cuisines, err := GetUserCuisinesAPI.Call(m.user.ID)
			if err != nil {
				CurrentPage.ShowErrorStatus(err)
				return
			}
			m.user.Cuisines = cuisines
			SetCurrentUser(m.user, ctx)

			m.meal = m.meal.RemoveInvalidCuisines(m.user.Cuisines)

			if m.meal.Cuisine == nil {
				m.meal.Cuisine = []string{"American"}
			}

			m.cuisine = make(map[string]bool)
			for _, cuisine := range m.meal.Cuisine {
				m.cuisine[cuisine] = true
			}
		},
		TitleElement: "Edit Meal",
		Elements: []app.UI{
			app.Form().ID("meal-page-form").Class("form").OnSubmit(m.OnSubmit).Body(
				NewTextInput("meal-page-name", "What is the name of this meal?", "Meal Name", true, &m.meal.Name),
				NewTextarea("meal-page-description", "Description/Notes:", "Meal description/notes", false, &m.meal.Description),
				NewCheckboxChips("meal-page-cuisine", "Cuisines:", map[string]bool{"American": true}, &m.cuisine, append(cuisines, "+")...).SetOnChange(m.CuisinesOnChange),
				newCuisinesDialog("meal-page", m.CuisinesDialogOnSave),
				app.Div().ID("meal-page-action-button-row").Class("action-button-row").Body(
					app.Input().ID("meal-page-delete-button").Class("action-button", "danger-action-button").Type("button").Value("Delete").OnClick(m.InitialDelete),
					app.A().ID("meal-page-cancel-button").Class("action-button", "secondary-action-button").Href("/home").Text("Cancel"),
					// app.Input().ID("meal-page-entries-button").Class("action-button", "tertiary-action-button").Type("button").Value("View Entries").OnClick(m.ViewEntries),
					app.Input().ID("meal-page-save-button").Class("action-button", "primary-action-button").Type("submit").Value("Save"),
				),
			),

			app.Dialog().ID("meal-page-confirm-delete").Body(
				app.P().ID("meal-page-confirm-delete-text").Class("confirm-delete-text").Text("Are you sure you want to delete this meal?"),
				app.Div().ID("meal-page-confirm-delete-action-button-row").Class("action-button-row").Body(
					app.Button().ID("meal-page-confirm-delete-delete").Class("action-button", "danger-action-button").Text("Yes, Delete").OnClick(m.ConfirmDelete),
					app.Button().ID("meal-page-confirm-delete-cancel").Class("action-button", "secondary-action-button").Text("No, Cancel").OnClick(m.CancelDelete),
				),
			),
		},
	}
}

func (m *meal) CuisinesOnChange(ctx app.Context, event app.Event, val string) {
	if val == "+" {
		m.cuisine[val] = false
		event.Get("target").Set("checked", false)
		app.Window().GetElementByID("meal-page-cuisines-dialog").Call("showModal")
	}
}

func (m *meal) CuisinesDialogOnSave(ctx app.Context, event app.Event) {
	m.user = GetCurrentUser(ctx)
	m.meal = m.meal.RemoveInvalidCuisines(m.user.Cuisines)
}

func (m *meal) OnSubmit(ctx app.Context, event app.Event) {
	event.PreventDefault()

	m.meal.Cuisine = []string{}
	for cuisine, value := range m.cuisine {
		if value {
			m.meal.Cuisine = append(m.meal.Cuisine, cuisine)
		}
	}

	_, err := UpdateMealAPI.Call(m.meal)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	SetCurrentMeal(m.meal, ctx)

	ctx.Navigate("/home")
}

func (m *meal) ViewEntries(ctx app.Context, event app.Event) {
	event.PreventDefault()

	m.meal.Cuisine = []string{}
	for cuisine, value := range m.cuisine {
		if value {
			m.meal.Cuisine = append(m.meal.Cuisine, cuisine)
		}
	}

	_, err := UpdateMealAPI.Call(m.meal)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	SetCurrentMeal(m.meal, ctx)

	ctx.Navigate("/entries")
}

func (m *meal) InitialDelete(ctx app.Context, event app.Event) {
	event.PreventDefault()
	app.Window().GetElementByID("meal-page-confirm-delete").Call("showModal")
}

func (m *meal) ConfirmDelete(ctx app.Context, event app.Event) {
	event.PreventDefault()

	_, err := DeleteMealAPI.Call(m.meal.ID)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	SetCurrentMeal(Meal{}, ctx)

	ctx.Navigate("/home")
}

func (m *meal) CancelDelete(ctx app.Context, event app.Event) {
	event.PreventDefault()
	app.Window().GetElementByID("meal-page-confirm-delete").Call("close")
}
