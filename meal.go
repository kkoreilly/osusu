package main

import (
	"log"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Meal is a struct that represents the data of a meal
type Meal struct {
	ID          int
	UserID      int
	Name        string
	Description string
	// Cost        int
	// Effort      int
	// Healthiness int
	// Taste       PersonMap
	// Type        string
	// Source      string
	// LastDone    time.Time
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
	sum := 100*entriesSum + options.RecencyWeight*recencyScore
	den := 100*len(entries) + options.RecencyWeight
	if den == 0 {
		return 0
	}
	return sum / den
	// // average of all attributes
	// var tasteSum int
	// for i, v := range m.Taste {
	// 	use := options.People[i]
	// 	// invert the person's rating if they are not participating
	// 	if use {
	// 		tasteSum += v
	// 	} else {
	// 		tasteSum += 100 - v
	// 	}
	// }
	// recencyScore := int(2 * time.Now().Truncate(time.Hour*24).UTC().Sub(m.LastDone) / (time.Hour * 24))
	// if recencyScore > 100 {
	// 	recencyScore = 100
	// }
	// sum := options.CostWeight*(100-m.Cost) + options.EffortWeight*(100-m.Effort) + options.HealthinessWeight*m.Healthiness + options.TasteWeight*tasteSum + options.RecencyWeight*recencyScore
	// den := options.CostWeight + options.EffortWeight + options.HealthinessWeight + len(m.Taste)*options.TasteWeight + options.RecencyWeight
	// if den == 0 {
	// 	return 0
	// }
	// return sum / den
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
)

type meal struct {
	app.Compo
	meal   Meal
	person Person
}

func (m *meal) Render() app.UI {
	// taste := e.meal.Taste[e.person.ID]
	return &Page{
		ID:                     "meal",
		Title:                  "Edit Meal",
		Description:            "Edit a meal",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			m.person = GetCurrentPerson(ctx)
			m.meal = GetCurrentMeal(ctx)
		},
		TitleElement: "Edit Meal",
		Elements: []app.UI{
			app.Form().ID("meal-page-form").Class("form").OnSubmit(m.OnSubmit).Body(
				NewTextInput("meal-page-name", "What is the name of this meal?", "Meal Name", true, &m.meal.Name),
				NewTextarea("meal-page-description", "Description/Notes:", "Meal description/notes", false, &m.meal.Description),
				// app.Label().ID("meal-page-last-done-label").Class("input-label").For("meal-page-last-done-input").Text("When did you last eat this?"),
				// app.Input().ID("meal-page-last-done-input").Class("input").Type("date").Value(e.meal.LastDone.Format("2006-01-02")),
				// NewRadioChips("meal-page-type", "What meals can you eat this for?", "Dinner", &e.meal.Type, mealTypes...),
				// NewRadioChips("meal-page-source", "How can you get this?", "Cooking", &e.meal.Source, mealSources...),
				// NewRangeInput("meal-page-taste", "How does this taste?", &taste),
				// NewRangeInput("meal-page-cost", "How much does this cost?", &e.meal.Cost),
				// NewRangeInput("meal-page-effort", "How much effort does this take?", &e.meal.Effort),
				// NewRangeInput("meal-page-healthiness", "How healthy is this?", &e.meal.Healthiness),
				app.Div().ID("meal-page-action-button-row").Class("action-button-row").Body(
					app.Input().ID("meal-page-delete-button").Class("action-button", "danger-action-button").Type("button").Value("Delete").OnClick(m.InitialDelete),
					app.A().ID("meal-page-cancel-button").Class("action-button", "secondary-action-button").Href("/home").Text("Cancel"),
					app.Input().ID("meal-page-save-button").Class("action-button", "primary-action-button").Type("submit").Value("Save"),
					app.Input().ID("meal-page-entries-button").Class("action-button", "tertiary-action-button").Type("button").Value("View Entries").OnClick(m.ViewEntries),
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

func (m *meal) OnSubmit(ctx app.Context, event app.Event) {
	event.PreventDefault()

	// m.meal.Taste[m.person.ID] = app.Window().GetElementByID("meal-page-taste-input").Get("valueAsNumber").Int()
	// m.meal.LastDone = time.UnixMilli(int64((app.Window().GetElementByID("meal-page-last-done-input").Get("valueAsNumber").Int()))).UTC()

	_, err := UpdateMealAPI.Call(m.meal)
	if err != nil {
		log.Println(err)
		return
	}
	SetCurrentMeal(m.meal, ctx)

	ctx.Navigate("/home")
}

func (m *meal) ViewEntries(ctx app.Context, event app.Event) {
	event.PreventDefault()

	_, err := UpdateMealAPI.Call(m.meal)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return
	}
	SetCurrentMeal(Meal{}, ctx)

	ctx.Navigate("/home")
}

func (m *meal) CancelDelete(ctx app.Context, event app.Event) {
	event.PreventDefault()
	app.Window().GetElementByID("meal-page-confirm-delete").Call("close")
}
