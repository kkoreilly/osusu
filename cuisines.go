package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type cuisinesDialogCompo struct {
	app.Compo
	ID             string
	OnSave         func(ctx app.Context, e app.Event)
	group          Group
	newCuisineName string
}
// newCuisine is the newly created cuisine
var newCuisine string

func (c *cuisinesDialogCompo) Render() app.UI {
	return app.Dialog().ID(c.ID+"-cuisines-dialog").Class("cuisines-dialog", "modal").Body(
		app.Form().ID(c.ID+"-cuisines-dialog-form").Class("form").OnSubmit(c.NewCuisine).Body(
			TextInput().ID(c.ID+"-cuisines-dialog-name").Label("Create New Cuisine:").Value(&c.newCuisineName),
			ButtonRow().ID(c.ID+"-cuisines-dialog-button-row").Buttons(
				Button().ID(c.ID+"-cuisines-dialog-delete").Class("danger").Icon("delete").Text("Delete Unused Cuisines").OnClick(c.DeleteUnusedCuisines),
				Button().ID(c.ID+"-cuisines-dialog-cancel").Class("secondary").Icon("cancel").Text("Cancel").OnClick(c.Cancel),
				Button().ID(c.ID+"-cuisines-dialog-new").Class("primary").Type("submit").Icon("add").Text("Create"),
			),
		),
	)
}

func cuisinesDialog(id string, onSave func(ctx app.Context, e app.Event)) *cuisinesDialogCompo {
	return &cuisinesDialogCompo{ID: id, OnSave: onSave}
}

func (c *cuisinesDialogCompo) NewCuisine(ctx app.Context, e app.Event) {
	e.PreventDefault()

	c.group = CurrentGroup(ctx)
	input := app.Window().GetElementByID(c.ID + "-cuisines-dialog-name-input")
	name := input.Get("value").String()
	c.group.Cuisines = append(c.group.Cuisines, name)
	newCuisine = name
	c.Save(ctx, e)
	input.Call("blur")
	ctx.Defer(func(ctx app.Context) {
		input.Set("value", "")
	})
}

func (c *cuisinesDialogCompo) Save(ctx app.Context, e app.Event) {
	_, err := UpdateGroupCuisinesAPI.Call(c.group)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	SetCurrentGroup(c.group, ctx)
	app.Window().GetElementByID(c.ID + "-cuisines-dialog").Call("close")
	c.OnSave(ctx, e)
}

func (c *cuisinesDialogCompo) Cancel(ctx app.Context, e app.Event) {
	app.Window().GetElementByID(c.ID + "-cuisines-dialog").Call("close")
}

func (c *cuisinesDialogCompo) DeleteUnusedCuisines(ctx app.Context, e app.Event) {
	c.group = CurrentGroup(ctx)
	
	meals, err := GetMealsAPI.Call(c.group.ID)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	usedCuisines := make(map[string]bool)
	for _, meal := range meals {
		for _, cuisine := range meal.Cuisine {
			usedCuisines[cuisine] = true
		}
	}
	newCuisines := []string{}
	for cuisine, val := range usedCuisines {
		if val {
			newCuisines = append(newCuisines, cuisine)
		}
	}
	c.group.Cuisines = newCuisines
	c.Save(ctx, e)
}