package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type edit struct {
	app.Compo
	meal  Meal
	meals Meals
}

func (e *edit) Render() app.UI {
	return app.Div().Body(
		app.H1().ID("edit-page-title").Class("page-title").Text("Edit"),
		app.Form().ID("edit-page-form").Class("form").OnSubmit(e.OnSubmit).Body(
			app.Label().ID("edit-page-name-label").Class("input-label", "edit-page-input-label").For("edit-page-name-input").Text("Name:"),
			app.Input().ID("edit-page-name-input").Class("input", "edit-page-input").Type("text").Placeholder("Name"),
			app.Label().ID("edit-page-cost-label").Class("input-label", "edit-page-input-label").For("edit-page-cost-input").Text("Cost:"),
			app.Input().ID("edit-page-cost-input").Class("input", "input-range", "edit-page-input").Type("range").Min(0).Max(100),
			app.Label().ID("edit-page-effort-label").Class("input-label", "edit-page-input-label").For("edit-page-effort-input").Text("Effort:"),
			app.Input().ID("edit-page-effort-input").Class("input", "input-range", "edit-page-input").Type("range").Min(0).Max(100),
			app.Label().ID("edit-page-healthiness-label").Class("input-label", "edit-page-input-label").For("edit-page-healthiness-input").Text("Healthiness:"),
			app.Input().ID("edit-page-healthiness-input").Class("input", "input-range", "edit-page-input").Type("range").Min(0).Max(100),
			app.Div().ID("edit-page-action-button-row").Class("action-button-row").Body(
				app.A().ID("edit-page-cancel-button").Class("action-button", "white-action-button").Href("/home").Text("Cancel"),
				app.Input().ID("edit-page-save-button").Class("action-button", "blue-action-button").Type("submit").Value("Save"),
			),
		),
	)
}

func (e *edit) OnNav(ctx app.Context) {
	e.meals = GetMeals(ctx)
	if e.meals == nil {
		e.meals = make(Meals)
		SetMeals(e.meals, ctx)
	}
	e.meal = GetCurrentMeal(ctx)
	app.Window().GetElementByID("edit-page-name-input").Set("value", e.meal.Name)
	app.Window().GetElementByID("edit-page-cost-input").Set("value", e.meal.Cost)
	app.Window().GetElementByID("edit-page-effort-input").Set("value", e.meal.Effort)
	app.Window().GetElementByID("edit-page-healthiness-input").Set("value", e.meal.Healthiness)
}

func (e *edit) OnSubmit(ctx app.Context, event app.Event) {
	event.PreventDefault()
	delete(e.meals, e.meal.Name)

	e.meal.Name = app.Window().GetElementByID("edit-page-name-input").Get("value").String()
	e.meal.Cost = app.Window().GetElementByID("edit-page-cost-input").Get("valueAsNumber").Int()
	e.meal.Effort = app.Window().GetElementByID("edit-page-effort-input").Get("valueAsNumber").Int()
	e.meal.Healthiness = app.Window().GetElementByID("edit-page-healthiness-input").Get("valueAsNumber").Int()

	e.meals[e.meal.Name] = e.meal
	SetMeals(e.meals, ctx)
	ctx.Navigate("/home")
}
