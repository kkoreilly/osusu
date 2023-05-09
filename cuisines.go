package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type cuisinesDialogCompo struct {
	app.Compo
	ID             string
	OnSave         func(ctx app.Context, e app.Event)
	group          Group
	cuisines       map[string]bool
	newCuisineName string
}

func (c *cuisinesDialogCompo) Render() app.UI {
	return app.Dialog().ID(c.ID+"-cuisines-dialog").Class("cuisines-dialog", "modal").Body(
		app.Form().ID(c.ID+"-cuisines-dialog-form").Class("form").OnSubmit(c.NewCuisine).Body(
			TextInput().ID(c.ID+"-cuisines-dialog-name").Label("New Cuisine Name:").Value(&c.newCuisineName),
			ButtonRow().ID(c.ID+"-cuisines-dialog-button-row").Buttons(
				Button().ID(c.ID+"-cuisines-dialog-cancel").Class("secondary").Icon("cancel").Text("Cancel").OnClick(c.Cancel),
				Button().ID(c.ID+"-cuisines-dialog-new").Class("primary").Type("submit").Icon("add").Text("Create New Cuisine"),
			),
		),
	)
}

func cuisinesDialog(id string, onSave func(ctx app.Context, e app.Event)) *cuisinesDialogCompo {
	return &cuisinesDialogCompo{ID: id, OnSave: onSave}
}

func (c *cuisinesDialogCompo) OnNav(ctx app.Context) {
	c.group = CurrentGroup(ctx)
	for _, cuisine := range c.group.Cuisines {
		c.cuisines[cuisine] = true
	}
}

func (c *cuisinesDialogCompo) NewCuisine(ctx app.Context, e app.Event) {
	e.PreventDefault()

	input := app.Window().GetElementByID(c.ID + "-cuisines-dialog-name-input")
	name := input.Get("value").String()
	c.cuisines[name] = true
	c.group.Cuisines = append(c.group.Cuisines, name)
	input.Call("blur")
	ctx.Defer(func(ctx app.Context) {
		input.Set("value", "")
	})
}

func (c *cuisinesDialogCompo) Cancel(ctx app.Context, e app.Event) {
	app.Window().GetElementByID(c.ID + "-cuisines-dialog").Call("close")
}

func (c *cuisinesDialogCompo) Save(ctx app.Context, e app.Event) {
	e.PreventDefault()

	c.group.Cuisines = []string{}
	for cuisine, value := range c.cuisines {
		if value {
			c.group.Cuisines = append(c.group.Cuisines, cuisine)
		}
	}
	_, err := UpdateGroupCuisinesAPI.Call(c.group)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	SetCurrentGroup(c.group, ctx)
	app.Window().GetElementByID(c.ID + "-cuisines-dialog").Call("close")
	c.OnSave(ctx, e)
}
