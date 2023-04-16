package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type cuisinesDialog struct {
	app.Compo
	ID             string
	OnSave         func(ctx app.Context, e app.Event)
	user           User
	cuisines       map[string]bool
	newCuisineName string
}

func (c *cuisinesDialog) Render() app.UI {
	return app.Dialog().ID(c.ID+"-cuisines-dialog").Class("cuisines-dialog").Body(
		app.Form().ID(c.ID+"-cuisines-dialog-new-cuisine-form").Class("form").OnSubmit(c.NewCuisine).Body(
			NewTextInput(c.ID+"-cuisines-dialog-new-name", "New Cuisine Name:", "New Cuisine Name", false, &c.newCuisineName),
			app.Input().ID(c.ID+"-cuisines-dialog-new-button").Class("action-button", "tertiary-action-button").Type("submit").Value("Make New Cuisine"),
		),
		app.Form().ID(c.ID+"-cuisines-dialog-form").Class("form").OnSubmit(c.Save).Body(
			NewCheckboxChips(c.ID+"-cuisines-dialog-checkbox-chips", "What cuisine options should be available?", map[string]bool{}, &c.cuisines, c.user.Cuisines...),
			app.Div().ID(c.ID+"-cuisines-dialog-action-button-row").Class("action-button-row").Body(
				app.Input().ID(c.ID+"-cuisines-dialog-cancel-button").Class("action-button", "secondary-action-button").Type("button").Value("Cancel").OnClick(c.Cancel),
				app.Input().ID(c.ID+"-cuisines-dialog-save-button").Class("action-button", "primary-action-button").Type("submit").Value("Save"),
			),
		),
	)
}

func newCuisinesDialog(id string, onSave func(ctx app.Context, e app.Event)) *cuisinesDialog {
	return &cuisinesDialog{ID: id, OnSave: onSave}
}

func (c *cuisinesDialog) OnNav(ctx app.Context) {
	c.user = GetCurrentUser(ctx)
	for _, cuisine := range c.user.Cuisines {
		c.cuisines[cuisine] = true
	}
}

func (c *cuisinesDialog) NewCuisine(ctx app.Context, e app.Event) {
	e.PreventDefault()

	input := app.Window().GetElementByID(c.ID + "-cuisines-dialog-new-name-input")
	name := input.Get("value").String()
	c.cuisines[name] = true
	c.user.Cuisines = append(c.user.Cuisines, name)
	input.Call("blur")
	ctx.Defer(func(ctx app.Context) {
		input.Set("value", "")
	})
}

func (c *cuisinesDialog) Cancel(ctx app.Context, e app.Event) {
	app.Window().GetElementByID(c.ID + "-cuisines-dialog").Call("close")
}

func (c *cuisinesDialog) Save(ctx app.Context, e app.Event) {
	e.PreventDefault()

	c.user.Cuisines = []string{}
	for cuisine, value := range c.cuisines {
		if value {
			c.user.Cuisines = append(c.user.Cuisines, cuisine)
		}
	}
	_, err := UpdateUserCuisinesAPI.Call(c.user)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	SetCurrentUser(c.user, ctx)
	app.Window().GetElementByID(c.ID + "-cuisines-dialog").Call("close")
	c.OnSave(ctx, e)
}
