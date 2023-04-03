package main

import (
	"log"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Entry is an entry with information about how a meal was at a certain point in time
type Entry struct {
	ID          int
	UserID      int
	MealID      int
	Date        time.Time
	Type        string
	Source      string
	Cost        PersonMap
	Effort      PersonMap
	Healthiness PersonMap
	Taste       PersonMap
}

// Score produces a score from 0 to 100 for the entry based on its attributes and the given options
func (e Entry) Score(options Options) int {
	// average of all attributes
	sum := options.CostWeight*e.Cost.Sum(options.People, true) + options.EffortWeight*e.Effort.Sum(options.People, true) + options.HealthinessWeight*e.Healthiness.Sum(options.People, false) + options.TasteWeight*e.Taste.Sum(options.People, false)
	den := len(e.Cost)*options.CostWeight + len(e.Effort)*options.EffortWeight + len(e.Healthiness)*options.HealthinessWeight + len(e.Taste)*options.TasteWeight
	if den == 0 {
		return 0
	}
	return sum / den
}

// MissingData returns whether the given person is missing data in the given entry
func (e Entry) MissingData(person Person) bool {
	return !(e.Cost.HasValueSet(person) && e.Effort.HasValueSet(person) && e.Healthiness.HasValueSet(person) && e.Taste.HasValueSet(person))
}

// SetCurrentEntry sets the current entry state value to the given entry
func SetCurrentEntry(entry Entry, ctx app.Context) {
	ctx.SetState("currentEntry", entry, app.Persist)
}

// GetCurrentEntry gets and returns the current entry state value
func GetCurrentEntry(ctx app.Context) Entry {
	var entry Entry
	ctx.GetState("currentEntry", &entry)
	return entry
}

type entry struct {
	app.Compo
	entry  Entry
	person Person
}

func (e *entry) Render() app.UI {
	return &Page{
		ID:                     "entry",
		Title:                  "Edit Entry",
		Description:            "Edit a meal entry",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			e.entry = GetCurrentEntry(ctx)
			e.person = GetCurrentPerson(ctx)
			// if the person is missing data for this entry, set their value to the average of all other people's ratings
			if !e.entry.Cost.HasValueSet(e.person) {
				e.entry.Cost[e.person.ID] = e.entry.Cost.Average()
			}
			if !e.entry.Effort.HasValueSet(e.person) {
				e.entry.Effort[e.person.ID] = e.entry.Effort.Average()
			}
			if !e.entry.Healthiness.HasValueSet(e.person) {
				e.entry.Healthiness[e.person.ID] = e.entry.Healthiness.Average()
			}
			if !e.entry.Taste.HasValueSet(e.person) {
				e.entry.Taste[e.person.ID] = e.entry.Taste.Average()
			}
		},
		TitleElement: "Edit Entry (" + e.person.Name + ")",
		Elements: []app.UI{
			app.Form().ID("entry-page-form").Class("form").OnSubmit(e.OnSubmit).Body(
				app.Label().ID("entry-page-date-label").Class("input-label").For("entry-page-date-input").Text("When did you eat this?"),
				app.Input().ID("entry-page-date-input").Class("input").Type("date").Value(e.entry.Date.Format("2006-01-02")),
				NewRadioChips("entry-page-type", "What meal did you eat this for?", "Dinner", &e.entry.Type, mealTypes...),
				NewRadioChips("entry-page-source", "How did you get this meal?", "Cooking", &e.entry.Source, mealSources...),
				NewRangeInputPersonMap("entry-page-taste", "How tasty do you think this was?", &e.entry.Taste, e.person),
				NewRangeInputPersonMap("entry-page-cost", "How expensive do you think this was?", &e.entry.Cost, e.person),
				NewRangeInputPersonMap("entry-page-effort", "How much effort do you think this took?", &e.entry.Effort, e.person),
				NewRangeInputPersonMap("entry-page-healthiness", "How healthy do you think this was?", &e.entry.Healthiness, e.person),
				app.Div().ID("entry-page-action-button-row").Class("action-button-row").Body(
					app.Input().ID("entry-page-delete-button").Class("action-button", "danger-action-button").Type("button").Value("Delete").OnClick(e.InitialDelete),
					app.A().ID("entry-page-cancel-button").Class("action-button", "secondary-action-button").Href("/entries").Text("Cancel"),
					app.Input().ID("entry-page-save-button").Class("action-button", "primary-action-button").Type("submit").Value("Save"),
				),
			),
			app.Dialog().ID("entry-page-confirm-delete").Body(
				app.P().ID("entry-page-confirm-delete-text").Class("confirm-delete-text").Text("Are you sure you want to delete this entry?"),
				app.Div().ID("entry-page-confirm-delete-action-button-row").Class("action-button-row").Body(
					app.Button().ID("entry-page-confirm-delete-delete").Class("action-button", "danger-action-button").Text("Yes, Delete").OnClick(e.ConfirmDelete),
					app.Button().ID("entry-page-confirm-delete-cancel").Class("action-button", "secondary-action-button").Text("No, Cancel").OnClick(e.CancelDelete),
				),
			),
		},
	}
}

func (e *entry) OnSubmit(ctx app.Context, event app.Event) {
	event.PreventDefault()

	e.entry.Date = time.UnixMilli(int64((app.Window().GetElementByID("entry-page-date-input").Get("valueAsNumber").Int()))).UTC()

	_, err := UpdateEntryAPI.Call(e.entry)
	if err != nil {
		log.Println(err)
		return
	}
	SetCurrentEntry(e.entry, ctx)

	ctx.Navigate("/entries")
}

func (e *entry) InitialDelete(ctx app.Context, event app.Event) {
	event.PreventDefault()
	app.Window().GetElementByID("entry-page-confirm-delete").Call("showModal")
}

func (e *entry) ConfirmDelete(ctx app.Context, event app.Event) {
	event.PreventDefault()

	_, err := DeleteEntryAPI.Call(e.entry.ID)
	if err != nil {
		log.Println(err)
		return
	}
	SetCurrentEntry(Entry{}, ctx)

	ctx.Navigate("/entries")
}

func (e *entry) CancelDelete(ctx app.Context, event app.Event) {
	event.PreventDefault()
	app.Window().GetElementByID("entry-page-confirm-delete").Call("close")
}
