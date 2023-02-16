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
			app.Input().ID("edit-page-name-input").Class("input").Type("text").Placeholder("Meal Name").AutoFocus(true),
			app.Label().ID("edit-page-last-done-label").Class("input-label").For("edit-page-last-done-input").Text("Last Done:"),
			app.Input().ID("edit-page-last-done-input").Class("input").Type("date"),
			app.Label().ID("edit-page-taste-label").Class("input-label").For("edit-page-taste-input").Text("Your Rating:"),
			app.Input().ID("edit-page-taste-input").Class("input", "input-range").Type("range").Min(0).Max(100),
			app.Hr().ID("edit-page-hr"),
			app.Label().ID("edit-page-cost-label").Class("input-label").For("edit-page-cost-input").Text("Cost:"),
			app.Input().ID("edit-page-cost-input").Class("input", "input-range").Type("range").Min(0).Max(100),
			app.Label().ID("edit-page-effort-label").Class("input-label").For("edit-page-effort-input").Text("Effort:"),
			app.Input().ID("edit-page-effort-input").Class("input", "input-range").Type("range").Min(0).Max(100),
			app.Label().ID("edit-page-healthiness-label").Class("input-label").For("edit-page-healthiness-input").Text("Healthiness:"),
			app.Input().ID("edit-page-healthiness-input").Class("input", "input-range").Type("range").Min(0).Max(100),
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
	e.person = GetCurrentPerson(ctx)
	e.meal = GetCurrentMeal(ctx)
	app.Window().GetElementByID("edit-page-name-input").Set("value", e.meal.Name)
	app.Window().GetElementByID("edit-page-cost-input").Set("valueAsNumber", e.meal.Cost)
	app.Window().GetElementByID("edit-page-effort-input").Set("valueAsNumber", e.meal.Effort)
	app.Window().GetElementByID("edit-page-healthiness-input").Set("valueAsNumber", e.meal.Healthiness)
	app.Window().GetElementByID("edit-page-taste-input").Set("valueAsNumber", e.meal.Taste[e.person.ID])
	app.Window().GetElementByID("edit-page-last-done-input").Set("valueAsNumber", e.meal.LastDone.UnixMilli())
}

func (e *edit) OnSubmit(ctx app.Context, event app.Event) {
	event.PreventDefault()

	e.meal.Name = app.Window().GetElementByID("edit-page-name-input").Get("value").String()
	e.meal.Cost = app.Window().GetElementByID("edit-page-cost-input").Get("valueAsNumber").Int()
	e.meal.Effort = app.Window().GetElementByID("edit-page-effort-input").Get("valueAsNumber").Int()
	e.meal.Healthiness = app.Window().GetElementByID("edit-page-healthiness-input").Get("valueAsNumber").Int()
	e.meal.Taste[e.person.ID] = app.Window().GetElementByID("edit-page-taste-input").Get("valueAsNumber").Int()
	e.meal.LastDone = time.UnixMilli(int64(app.Window().GetElementByID("edit-page-last-done-input").Get("valueAsDate").Call("getTime").Int())).UTC()

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
