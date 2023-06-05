package page

import (
	"github.com/kkoreilly/osusu/api"
	"github.com/kkoreilly/osusu/compo"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type MealPage struct {
	app.Compo
	group     osusu.Group
	user      osusu.User
	meal      osusu.Meal
	isMealNew bool
	category  map[string]bool
	cuisine   map[string]bool
}

func (m *MealPage) Render() app.UI {
	// need to copy to separate array from because append modifies the underlying array
	var cuisines = make([]string, len(m.group.Cuisines))
	copy(cuisines, m.group.Cuisines)
	titleText := "Edit Meal"
	saveButtonIcon := "save"
	saveButtonText := "Save"
	if m.isMealNew {
		titleText = "Create Meal"
		saveButtonIcon = "add"
		saveButtonText = "Create"
	}
	return &compo.Page{
		ID:                     "meal",
		Title:                  titleText,
		Description:            "Edit, view, or create a meal.",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			m.group = osusu.CurrentGroup(ctx)
			m.user = osusu.CurrentUser(ctx)
			m.meal = osusu.CurrentMeal(ctx)
			m.isMealNew = osusu.IsMealNew(ctx)

			if m.isMealNew {
				compo.CurrentPage.Title = "Create Meal"
				compo.CurrentPage.UpdatePageTitle(ctx)
			}

			cuisines, err := api.GetGroupCuisines.Call(m.group.ID)
			if err != nil {
				compo.CurrentPage.ShowErrorStatus(err)
				return
			}
			m.group.Cuisines = cuisines
			osusu.SetCurrentGroup(m.group, ctx)

			m.meal = m.meal.RemoveInvalidCuisines(m.group.Cuisines)

			// need to check that length is 0 as well because we could have data from recipe import
			if m.isMealNew && len(m.meal.Category) == 0 {
				m.meal.Category = []string{"Dinner"}
			}
			m.category = make(map[string]bool)
			for _, category := range m.meal.Category {
				m.category[category] = true
			}

			// need to check that length is 0 as well because we could have data from recipe import
			if m.isMealNew && len(m.meal.Cuisine) == 0 {
				m.meal.Cuisine = []string{"American"}
			}
			m.cuisine = make(map[string]bool)
			for _, cuisine := range m.meal.Cuisine {
				m.cuisine[cuisine] = true
			}
		},
		TitleElement: titleText,
		Elements: []app.UI{
			app.Form().ID("meal-page-form").Class("form").OnSubmit(m.OnSubmit).Body(
				compo.TextInput().ID("meal-page-name").Label("Name:").Value(&m.meal.Name).AutoFocus(true),
				compo.Textarea().ID("meal-page-description").Label("Description:").Value(&m.meal.Description),
				compo.TextInput().ID("meal-page-source").Label("Source:").Value(&m.meal.Source),
				// Button().ID("meal-page-view-source").Class("secondary").Value()
				compo.TextInput().ID("meal-page-image").Label("Image:").Value(&m.meal.Image),
				compo.CheckboxChips().ID("meal-page-category").Label("Categories:").Value(&m.category).Options(osusu.AllCategories...),
				compo.CheckboxChips().ID("meal-page-cuisine").Label("Cuisines:").Value(&m.cuisine).Options(append(cuisines, "+")...).OnChange(m.CuisinesOnChange),
				compo.CuisinesDialog("meal-page", m.CuisinesDialogOnSave),
				compo.ButtonRow().ID("meal-page").Buttons(
					compo.Button().ID("meal-page-delete").Class("danger").Icon("delete").Text("Delete").OnClick(m.DeleteMeal).Hidden(m.isMealNew),
					compo.Button().ID("meal-page-cancel").Class("secondary").Icon("cancel").Text("Cancel").OnClick(compo.NavigateEvent("/home")),
					compo.Button().ID("meal-page-save").Class("primary").Type("submit").Icon(saveButtonIcon).Text(saveButtonText),
				),
				app.Dialog().ID("meal-page-confirm-delete-meal").Class("modal").Body(
					app.P().ID("meal-page-confirm-delete-meal-text").Class("confirm-delete-text").Text("Are you sure you want to delete this meal?"),
					compo.ButtonRow().ID("meal-page-confirm-delete-meal").Buttons(
						compo.Button().ID("meal-page-confirm-delete-meal-delete").Class("danger").Icon("delete").Text("Yes, Delete").OnClick(m.ConfirmDeleteMeal),
						compo.Button().ID("meal-page-confirm-delete-meal-cancel").Class("secondary").Icon("cancel").Text("No, Cancel").OnClick(m.CancelDeleteMeal),
					),
				),
			),
		},
	}
}

func (m *MealPage) CuisinesOnChange(ctx app.Context, event app.Event, val string) {
	if val == "+" {
		m.cuisine[val] = false
		event.Get("target").Set("checked", false)
		app.Window().GetElementByID("meal-page-cuisines-dialog").Call("showModal")
	}
}

func (m *MealPage) CuisinesDialogOnSave(ctx app.Context, event app.Event) {
	m.group = osusu.CurrentGroup(ctx)
	if compo.NewCuisineCreated {
		m.cuisine[compo.NewCuisine] = true
	}
	m.meal.RemoveInvalidCuisines(m.group.Cuisines)
}

func (m *MealPage) OnSubmit(ctx app.Context, event app.Event) {
	event.PreventDefault()

	m.meal.Category = []string{}
	for category, value := range m.category {
		if value {
			m.meal.Category = append(m.meal.Category, category)
		}
	}

	m.meal.Cuisine = []string{}
	for cuisine, value := range m.cuisine {
		if value {
			m.meal.Cuisine = append(m.meal.Cuisine, cuisine)
		}
	}

	if m.isMealNew {
		m.meal.GroupID = m.group.ID
		meal, err := api.CreateMeal.Call(m.meal)
		if err != nil {
			compo.CurrentPage.ShowErrorStatus(err)
			return
		}
		m.meal = meal
		osusu.SetCurrentMeal(m.meal, ctx)
		entries, err := api.GetEntriesForMeal.Call(m.meal.ID)
		if err != nil {
			compo.CurrentPage.ShowErrorStatus(err)
			return
		}
		entry := osusu.NewEntry(m.group, m.user, m.meal, entries)
		osusu.SetIsEntryNew(true, ctx)
		osusu.SetCurrentEntry(entry, ctx)
		compo.Navigate("/entry", ctx)
		return
	}
	_, err := api.UpdateMeal.Call(m.meal)
	if err != nil {
		compo.CurrentPage.ShowErrorStatus(err)
		return
	}

	osusu.SetCurrentMeal(m.meal, ctx)

	compo.Navigate("/home", ctx)
}

func (m *MealPage) ViewEntries(ctx app.Context, event app.Event) {
	event.PreventDefault()

	m.meal.Cuisine = []string{}
	for cuisine, value := range m.cuisine {
		if value {
			m.meal.Cuisine = append(m.meal.Cuisine, cuisine)
		}
	}

	_, err := api.UpdateMeal.Call(m.meal)
	if err != nil {
		compo.CurrentPage.ShowErrorStatus(err)
		return
	}
	osusu.SetCurrentMeal(m.meal, ctx)

	compo.Navigate("/entries", ctx)
}

func (m *MealPage) DeleteMeal(ctx app.Context, e app.Event) {
	e.PreventDefault()
	app.Window().GetElementByID("meal-page-confirm-delete-meal").Call("showModal")
}

func (m *MealPage) ConfirmDeleteMeal(ctx app.Context, e app.Event) {
	e.PreventDefault()

	_, err := api.DeleteMeal.Call(m.meal.ID)
	if err != nil {
		compo.CurrentPage.ShowErrorStatus(err)
		return
	}
	osusu.SetCurrentMeal(osusu.Meal{}, ctx)
	app.Window().GetElementByID("meal-page-confirm-delete-meal").Call("close")
	compo.ReturnToReturnURL(ctx, e)
}

func (m *MealPage) CancelDeleteMeal(ctx app.Context, e app.Event) {
	e.PreventDefault()
	app.Window().GetElementByID("meal-page-confirm-delete-meal").Call("close")
}
