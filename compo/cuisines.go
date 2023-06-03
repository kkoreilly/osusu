package compo

import (
	"github.com/kkoreilly/osusu/api"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type CuisinesDialogCompo struct {
	app.Compo
	ID             string
	OnSave         func(ctx app.Context, e app.Event)
	group          osusu.Group
	newCuisineName string
}

// NewCuisinecCreated is whether a new cuisine was created by the cuisines dialog
var NewCuisineCreated bool

// NewCuisine is the newly created cuisine (only applicable if NewCuisineCreated is true)
var NewCuisine string

func (c *CuisinesDialogCompo) Render() app.UI {
	return app.Div().ID(c.ID+"-cuisines-dialog-container").Class("cuisines-dialog-container").Body(
		app.Dialog().ID(c.ID+"-cuisines-dialog").Class("cuisines-dialog", "modal").Body(
			app.Form().ID(c.ID+"-cuisines-dialog-form").Class("form").OnSubmit(c.NewCuisine).Body(
				TextInput().ID(c.ID+"-cuisines-dialog-name").Label("Create New Cuisine:").Value(&c.newCuisineName),
				ButtonRow().ID(c.ID+"-cuisines-dialog-button-row").Buttons(
					Button().ID(c.ID+"-cuisines-dialog-delete").Class("danger").Icon("delete").Text("Delete Unused Cuisines").OnClick(c.InitialDelete),
					Button().ID(c.ID+"-cuisines-dialog-cancel").Class("secondary").Icon("cancel").Text("Cancel").OnClick(c.Cancel),
					Button().ID(c.ID+"-cuisines-dialog-new").Class("primary").Type("submit").Icon("add").Text("Create"),
				),
			),
		),
		app.Dialog().ID(c.ID+"-cuisines-dialog-confirm-delete").Class("cuisines-dialog-confirm-delete", "modal").Body(
			app.P().ID(c.ID+"-cuisines-dialog-confirm-delete-text").Class("confirm-delete-text").Text("Are you sure you want to delete all unused cuisines?"),
			ButtonRow().ID(c.ID+"-cuisines-dialog-confirm-delete").Buttons(
				Button().ID(c.ID+"-cuisines-dialog-confirm-delete-delete").Class("danger").Icon("delete").Text("Yes, Delete").OnClick(c.DeleteUnusedCuisines),
				Button().ID(c.ID+"-cuisines-dialog-confirm-delete-cancel").Class("secondary").Icon("cancel").Text("No, Cancel").OnClick(c.CancelDelete),
			),
		),
	)
}
func CuisinesDialog(id string, onSave func(ctx app.Context, e app.Event)) *CuisinesDialogCompo {
	return &CuisinesDialogCompo{ID: id, OnSave: onSave}
}

func (c *CuisinesDialogCompo) NewCuisine(ctx app.Context, e app.Event) {
	e.PreventDefault()

	NewCuisineCreated = true
	c.group = osusu.CurrentGroup(ctx)
	input := app.Window().GetElementByID(c.ID + "-cuisines-dialog-name-input")
	name := input.Get("value").String()
	c.group.Cuisines = append(c.group.Cuisines, name)
	NewCuisine = name
	c.Save(ctx, e)
	input.Call("blur")
	ctx.Defer(func(ctx app.Context) {
		input.Set("value", "")
	})
}

func (c *CuisinesDialogCompo) Save(ctx app.Context, e app.Event) {
	_, err := api.UpdateGroupCuisines.Call(c.group)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	osusu.SetCurrentGroup(c.group, ctx)
	app.Window().GetElementByID(c.ID + "-cuisines-dialog").Call("close")
	c.OnSave(ctx, e)
}

func (c *CuisinesDialogCompo) Cancel(ctx app.Context, e app.Event) {
	app.Window().GetElementByID(c.ID + "-cuisines-dialog").Call("close")
}

func (c *CuisinesDialogCompo) InitialDelete(ctx app.Context, e app.Event) {
	app.Window().GetElementByID(c.ID + "-cuisines-dialog-confirm-delete").Call("showModal")
}

func (c *CuisinesDialogCompo) DeleteUnusedCuisines(ctx app.Context, e app.Event) {
	NewCuisineCreated = false

	c.group = osusu.CurrentGroup(ctx)

	meals, err := api.GetMeals.Call(c.group.ID)
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
	app.Window().GetElementByID(c.ID + "-cuisines-dialog-confirm-delete").Call("close")
}

func (c *CuisinesDialogCompo) CancelDelete(ctx app.Context, e app.Event) {
	app.Window().GetElementByID(c.ID + "-cuisines-dialog-confirm-delete").Call("close")
}
