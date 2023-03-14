package main

import (
	"log"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

var (
	mealTypes   = []string{"Breakfast", "Lunch", "Dinner"}
	mealSources = []string{"Cooking", "Dine-In", "Takeout"}
)

type edit struct {
	app.Compo
	meal   Meal
	person Person
}

func (e *edit) Render() app.UI {
	// taste := e.meal.Taste[e.person.ID]
	return app.Div().Body(
		app.H1().ID("edit-page-title").Class("page-title").Text("Edit Meal"),
		app.Form().ID("edit-page-form").Class("form").OnSubmit(e.OnSubmit).Body(
			NewTextInput("edit-page-name", "What is the name of this meal/restaurant?", "Meal Name", true, &e.meal.Name),
			NewTextarea("edit-page-description", "Description/Notes:", "Meal description/notes", false, &e.meal.Description),
			// app.Label().ID("edit-page-last-done-label").Class("input-label").For("edit-page-last-done-input").Text("When did you last eat this?"),
			// app.Input().ID("edit-page-last-done-input").Class("input").Type("date").Value(e.meal.LastDone.Format("2006-01-02")),
			// NewRadioChips("edit-page-type", "What meals can you eat this for?", "Dinner", &e.meal.Type, mealTypes...),
			// NewRadioChips("edit-page-source", "How can you get this?", "Cooking", &e.meal.Source, mealSources...),
			// NewRangeInput("edit-page-taste", "How does this taste?", &taste),
			// NewRangeInput("edit-page-cost", "How much does this cost?", &e.meal.Cost),
			// NewRangeInput("edit-page-effort", "How much effort does this take?", &e.meal.Effort),
			// NewRangeInput("edit-page-healthiness", "How healthy is this?", &e.meal.Healthiness),
			app.Div().ID("edit-page-action-button-row").Class("action-button-row").Body(
				app.Input().ID("edit-page-delete-button").Class("action-button", "red-action-button").Type("button").Value("Delete").OnClick(e.InitialDelete),
				app.A().ID("edit-page-cancel-button").Class("action-button", "white-action-button").Href("/home").Text("Cancel"),
				app.Input().ID("edit-page-save-button").Class("action-button", "blue-action-button").Type("submit").Value("Save"),
				app.A().ID("edit-page-entries-button").Class("action-button", "green-action-button").Href("/entries").Text("View Entries"),
			),
		),
		app.Dialog().ID("edit-page-confirm-delete").Body(
			app.P().ID("edit-page-confirm-delete-text").Text("Are you sure you want to delete this meal?"),
			app.Div().ID("edit-page-confirm-delete-action-button-row").Class("action-button-row").Body(
				app.Button().ID("edit-page-confirm-delete-delete").Class("action-button", "red-action-button").Text("Yes, Delete").OnClick(e.ConfirmDelete),
				app.Button().ID("edit-page-confirm-delete-cancel").Class("action-button", "white-action-button").Text("No, Cancel").OnClick(e.CancelDelete),
			),
		),
	)
}

func (e *edit) OnNav(ctx app.Context) {
	if Authenticate(true, ctx) {
		return
	}
	e.person = GetCurrentPerson(ctx)
	e.meal = GetCurrentMeal(ctx)
}

func (e *edit) OnSubmit(ctx app.Context, event app.Event) {
	event.PreventDefault()

	e.meal.Taste[e.person.ID] = app.Window().GetElementByID("edit-page-taste-input").Get("valueAsNumber").Int()
	e.meal.LastDone = time.UnixMilli(int64((app.Window().GetElementByID("edit-page-last-done-input").Get("valueAsNumber").Int()))).UTC()

	err := UpdateMealRequest(e.meal)
	if err != nil {
		log.Println(err)
		return
	}
	SetCurrentMeal(e.meal, ctx)

	ctx.Navigate("/home")
}

func (e *edit) InitialDelete(ctx app.Context, event app.Event) {
	event.PreventDefault()
	app.Window().GetElementByID("edit-page-confirm-delete").Call("showModal")
}

func (e *edit) ConfirmDelete(ctx app.Context, event app.Event) {
	event.PreventDefault()

	err := DeleteMealRequest(e.meal)
	if err != nil {
		log.Println(err)
		return
	}
	SetCurrentMeal(Meal{}, ctx)

	ctx.Navigate("/home")
}

func (e *edit) CancelDelete(ctx app.Context, event app.Event) {
	event.PreventDefault()
	app.Window().GetElementByID("edit-page-confirm-delete").Call("close")
}
