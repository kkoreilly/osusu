package main

import (
	"log"

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
	return app.Dialog().ID(c.ID+"-cuisines-dialog").Class("cuisines-dialog", "modal").Body(
		app.Form().ID(c.ID+"-cuisines-dialog-new-cuisine-form").Class("form").OnSubmit(c.NewCuisine).Body(
			TextInput().ID(c.ID+"-cuisines-dialog-new-name").Label("New Cuisine Name:").Value(&c.newCuisineName),
			Button().ID(c.ID+"-cuisines-dialog-new").Class("tertiary").Type("submit").Icon("add").Text("Create New Cuisine"),
		),
		app.Form().ID(c.ID+"-cuisines-dialog-form").Class("form").OnSubmit(c.Save).Body(
			CheckboxChips().ID(c.ID+"-cuisines-dialog-chips").Label("What cuisine options should be available?").Value(&c.cuisines).Options(c.user.Cuisines...),
			ButtonRow().ID(c.ID+"-cuisines-dialog").Buttons(
				Button().ID(c.ID+"-cuisines-dialog-cancel").Class("secondary").Icon("cancel").Text("Cancel").OnClick(c.Cancel),
				Button().ID(c.ID+"-cuisines-dialog-save").Class("primary").Type("submit").Icon("save").Text("Save"),
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
	log.Println(c.user, c.cuisines)
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
