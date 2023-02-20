package main

import (
	"log"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type edit struct {
	app.Compo
	meal   Meal
	person Person
}

func (e *edit) Render() app.UI {
	return app.Div().Body(
		app.H1().ID("edit-page-title").Class("page-title").Text("Edit Meal"),
		app.Form().ID("edit-page-form").Class("form").OnSubmit(e.OnSubmit).Body(
			app.Label().ID("edit-page-name-label").Class("input-label").For("edit-page-name-input").Text("Name:"),
			app.Input().ID("edit-page-name-input").Class("input").Type("text").Placeholder("Meal Name").AutoFocus(true).Value(e.meal.Name),
			app.Label().ID("edit-page-last-done-label").Class("input-label").For("edit-page-last-done-input").Text("Last Done:"),
			app.Input().ID("edit-page-last-done-input").Class("input").Type("date").Value(e.meal.LastDone.Format("2006-01-02")),
			app.Label().ID("edit-page-type-label").Class("input-label").For("edit-page-type-inputs-container").Text("Type:"),
			app.Div().ID("edit-page-type-inputs-container").Class("chips-container").Body(
				Chip("edit-page-type-breakfast", "radio", "edit-page-type", "Breakfast", e.meal.Type == "Breakfast"),
				Chip("edit-page-type-lunch", "radio", "edit-page-type", "Lunch", e.meal.Type == "Lunch"),
				Chip("edit-page-type-dinner", "radio", "edit-page-type", "Dinner", e.meal.Type == "Dinner"),
			),
			app.Label().ID("edit-page-source-label").Class("input-label").For("edit-page-source-inputs-container").Text("Source:"),
			app.Div().ID("edit-page-source-inputs-container").Class("chips-container").Body(
				Chip("edit-page-source-cooking", "radio", "edit-page-source", "Cooking", e.meal.Source == "Cooking"),
				Chip("edit-page-source-restaurant", "radio", "edit-page-source", "Dine-In", e.meal.Source == "Dine-In"),
				Chip("edit-page-source-takeout", "radio", "edit-page-source", "Takeout", e.meal.Source == "Takeout"),
			),
			app.Label().ID("edit-page-taste-label").Class("input-label").For("edit-page-taste-input").Text("Your Rating:"),
			app.Input().ID("edit-page-taste-input").Class("input", "input-range").Type("range").Min(0).Max(100).Value(e.meal.Taste[e.person.ID]),
			app.Hr().ID("edit-page-hr"),
			app.Label().ID("edit-page-cost-label").Class("input-label").For("edit-page-cost-input").Text("Cost:"),
			app.Input().ID("edit-page-cost-input").Class("input", "input-range").Type("range").Min(0).Max(100).Value(e.meal.Cost),
			app.Label().ID("edit-page-effort-label").Class("input-label").For("edit-page-effort-input").Text("Effort:"),
			app.Input().ID("edit-page-effort-input").Class("input", "input-range").Type("range").Min(0).Max(100).Value(e.meal.Effort),
			app.Label().ID("edit-page-healthiness-label").Class("input-label").For("edit-page-healthiness-input").Text("Healthiness:"),
			app.Input().ID("edit-page-healthiness-input").Class("input", "input-range").Type("range").Min(0).Max(100).Value(e.meal.Healthiness),
			app.Div().ID("edit-page-action-button-row").Class("action-button-row").Body(
				app.Input().ID("edit-page-delete-button").Class("action-button", "red-action-button").Type("button").Value("Delete").OnClick(e.InitialDelete),
				app.A().ID("edit-page-cancel-button").Class("action-button", "white-action-button").Href("/home").Text("Cancel"),
				app.Input().ID("edit-page-save-button").Class("action-button", "blue-action-button").Type("submit").Value("Save"),
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
	if e.meal.Type == "" {
		e.meal.Type = "Dinner"
	}
	if e.meal.Source == "" {
		e.meal.Source = "Cooking"
	}
}

func (e *edit) OnSubmit(ctx app.Context, event app.Event) {
	event.PreventDefault()

	e.meal.Name = app.Window().GetElementByID("edit-page-name-input").Get("value").String()
	e.meal.Cost = app.Window().GetElementByID("edit-page-cost-input").Get("valueAsNumber").Int()
	e.meal.Effort = app.Window().GetElementByID("edit-page-effort-input").Get("valueAsNumber").Int()
	e.meal.Healthiness = app.Window().GetElementByID("edit-page-healthiness-input").Get("valueAsNumber").Int()
	e.meal.Taste[e.person.ID] = app.Window().GetElementByID("edit-page-taste-input").Get("valueAsNumber").Int()

	e.meal.LastDone = time.UnixMilli(int64((app.Window().GetElementByID("edit-page-last-done-input").Get("valueAsNumber").Int()))).UTC()
	elem := app.Window().Get("document").Call("querySelector", `input[name="edit-page-type"]:checked`)
	if !elem.IsNull() {
		e.meal.Type = elem.Get("value").String()
	}
	elem = app.Window().Get("document").Call("querySelector", `input[name="edit-page-source"]:checked`)
	if !elem.IsNull() {
		e.meal.Source = elem.Get("value").String()
	}
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
