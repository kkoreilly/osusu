package page

import (
	"github.com/kkoreilly/osusu/api"
	"github.com/kkoreilly/osusu/compo"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type EntryPage struct {
	app.Compo
	entry      osusu.Entry
	isEntryNew bool
	user       osusu.User
}

func (e *EntryPage) Render() app.UI {
	titleText := "Edit Entry"
	saveButtonIcon := "save"
	saveButtonText := "Save"
	if e.isEntryNew {
		titleText = "Create Entry"
		saveButtonIcon = "add"
		saveButtonText = "Create"
	}
	return &Page{
		ID:                     "entry",
		Title:                  titleText,
		Description:            "View, edit, or create a meal entry.",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			e.entry = osusu.CurrentEntry(ctx)
			e.isEntryNew = osusu.IsEntryNew(ctx)
			if e.isEntryNew {
				CurrentPage.Title = "Create Entry"
				CurrentPage.UpdatePageTitle(ctx)
			}
			e.user = osusu.CurrentUser(ctx)
			// e.entry = e.entry.FixMissingData(e.user)
		},
		TitleElement: titleText,
		Elements: []app.UI{
			app.Form().ID("entry-page-form").Class("form").OnSubmit(e.OnSubmit).Body(
				compo.DateInput().ID("entry-page-date").Label("When did you eat this?").Value(&e.entry.Date),
				compo.RadioChips().ID("entry-page-type").Label("What meal did you eat this for?").Default("Dinner").Value(&e.entry.Type).Options(osusu.MealCategories...),
				compo.RadioChips().ID("entry-page-source").Label("How did you get this meal?").Default("Cooking").Value(&e.entry.Source).Options(osusu.MealSources...),
				compo.RangeInputUserMap(&e.entry.Taste, e.user).ID("entry-page-taste").Label("How tasty think this was?"),
				compo.RangeInputUserMap(&e.entry.Cost, e.user).ID("entry-page-cost").Label("How expensive was this?"),
				compo.RangeInputUserMap(&e.entry.Effort, e.user).ID("entry-page-effort").Label("How much effort did this take?"),
				compo.RangeInputUserMap(&e.entry.Healthiness, e.user).ID("entry-page-healthiness").Label("How healthy was this?"),
				compo.ButtonRow().ID("entry-page").Buttons(
					compo.Button().ID("entry-page-delete").Class("danger").Icon("delete").Text("Delete").OnClick(e.InitialDelete).Hidden(e.isEntryNew),
					compo.Button().ID("entry-page-cancel").Class("secondary").Icon("cancel").Text("Cancel").OnClick(ReturnToReturnURL),
					compo.Button().ID("entry-page-save").Class("primary").Type("submit").Icon(saveButtonIcon).Text(saveButtonText),
				),
			),
			app.Dialog().ID("entry-page-confirm-delete").Class("modal").Body(
				app.P().ID("entry-page-confirm-delete-text").Class("confirm-delete-text").Text("Are you sure you want to delete this entry?"),
				compo.ButtonRow().ID("entry-page-confirm-delete").Buttons(
					compo.Button().ID("entry-page-confirm-delete-delete").Class("danger").Icon("delete").Text("Yes, Delete").OnClick(e.ConfirmDelete),
					compo.Button().ID("entry-page-confirm-delete-cancel").Class("secondary").Icon("cancel").Text("No, Cancel").OnClick(e.CancelDelete),
				),
			),
		},
	}
}

func (e *EntryPage) OnSubmit(ctx app.Context, event app.Event) {
	event.PreventDefault()

	if e.isEntryNew {
		entry, err := api.CreateEntryAPI.Call(e.entry)
		if err != nil {
			CurrentPage.ShowErrorStatus(err)
			return
		}
		e.entry = entry
		osusu.SetCurrentEntry(e.entry, ctx)
		ReturnToReturnURL(ctx, event)
		return
	}
	_, err := api.UpdateEntryAPI.Call(e.entry)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	osusu.SetCurrentEntry(e.entry, ctx)
	ReturnToReturnURL(ctx, event)
}

func (e *EntryPage) InitialDelete(ctx app.Context, event app.Event) {
	event.PreventDefault()
	app.Window().GetElementByID("entry-page-confirm-delete").Call("showModal")
}

func (e *EntryPage) ConfirmDelete(ctx app.Context, event app.Event) {
	event.PreventDefault()

	_, err := api.DeleteEntryAPI.Call(e.entry.ID)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	osusu.SetCurrentEntry(osusu.Entry{}, ctx)

	ReturnToReturnURL(ctx, event)
}

func (e *EntryPage) CancelDelete(ctx app.Context, event app.Event) {
	event.PreventDefault()
	app.Window().GetElementByID("entry-page-confirm-delete").Call("close")
}
