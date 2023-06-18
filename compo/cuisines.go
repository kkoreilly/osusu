package compo

import (
	"github.com/kkoreilly/osusu/api"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// CuisinesDialog is a component that allows the creation of new cuisines
type CuisinesDialog struct {
	app.Compo
	ID             string
	Group          *osusu.Group
	Cuisine        map[string]bool // the cuisine value that is being controlled on the page with the cuisines dialog; this will be used to set the newly created cuisine to true
	OnSave         func(ctx app.Context, e app.Event)
	newCuisineName string
}

func (c *CuisinesDialog) Render() app.UI {
	return app.Div().ID(c.ID+"-cuisines-dialog-container").Class("cuisines-dialog-container").Body(
		app.Dialog().ID(c.ID+"-cuisines-dialog").Class("cuisines-dialog", "modal").Body(
			app.Form().ID(c.ID+"-cuisines-dialog-form").Class("form").OnSubmit(c.NewCuisine).Body(
				TextInput().ID(c.ID+"-cuisines-dialog-name").Label("Create New Cuisine:").Value(&c.newCuisineName),
				&ButtonRow{ID: c.ID + "-cuisines-dialog-button-row", Buttons: []app.UI{
					&Button{ID: c.ID + "-cuisines-dialog-delete", Class: "danger", Icon: "delete", Text: "Delete Unused Cuisines", OnClick: c.InitialDelete},
					&Button{ID: c.ID + "-cuisines-dialog-cancel", Class: "secondary", Icon: "cancel", Text: "Cancel", OnClick: c.Cancel},
					&Button{ID: c.ID + "-cuisines-dialog-new", Class: "primary", Type: "submit", Icon: "add", Text: "Create"},
				}},
			),
		),
		app.Dialog().ID(c.ID+"-cuisines-dialog-confirm-delete").Class("cuisines-dialog-confirm-delete", "modal").Body(
			app.P().ID(c.ID+"-cuisines-dialog-confirm-delete-text").Class("confirm-delete-text").Text("Are you sure you want to delete all unused cuisines?"),
			&ButtonRow{ID: c.ID + "-cuisines-dialog-confirm-delete", Buttons: []app.UI{
				&Button{ID: c.ID + "-cuisines-dialog-confirm-delete-delete", Class: "danger", Icon: "delete", Text: "Yes, Delete", OnClick: c.DeleteUnusedCuisines},
				&Button{ID: c.ID + "-cuisines-dialog-confirm-delete-cancel", Class: "secondary", Icon: "cancel", Text: "No, Cancel", OnClick: c.CancelDelete},
			}},
		),
	)
}

func (c *CuisinesDialog) NewCuisine(ctx app.Context, e app.Event) {
	e.PreventDefault()

	input := app.Window().GetElementByID(c.ID + "-cuisines-dialog-name-input")
	name := input.Get("value").String()
	c.Group.Cuisines = append(c.Group.Cuisines, name)
	c.Cuisine[name] = true
	c.Save(ctx, e)
	input.Call("blur")
	ctx.Defer(func(ctx app.Context) {
		input.Set("value", "")
	})
}

func (c *CuisinesDialog) Save(ctx app.Context, e app.Event) {
	_, err := api.UpdateGroupCuisines.Call(*c.Group)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	osusu.SetCurrentGroup(*c.Group, ctx)
	app.Window().GetElementByID(c.ID + "-cuisines-dialog").Call("close")
	if c.OnSave != nil {
		c.OnSave(ctx, e)
	}
}

func (c *CuisinesDialog) Cancel(ctx app.Context, e app.Event) {
	app.Window().GetElementByID(c.ID + "-cuisines-dialog").Call("close")
}

func (c *CuisinesDialog) InitialDelete(ctx app.Context, e app.Event) {
	app.Window().GetElementByID(c.ID + "-cuisines-dialog-confirm-delete").Call("showModal")
}

func (c *CuisinesDialog) DeleteUnusedCuisines(ctx app.Context, e app.Event) {
	meals, err := api.GetMeals.Call(c.Group.ID)
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
	c.Group.Cuisines = newCuisines
	c.Save(ctx, e)
	app.Window().GetElementByID(c.ID + "-cuisines-dialog-confirm-delete").Call("close")
}

func (c *CuisinesDialog) CancelDelete(ctx app.Context, e app.Event) {
	app.Window().GetElementByID(c.ID + "-cuisines-dialog-confirm-delete").Call("close")
}
