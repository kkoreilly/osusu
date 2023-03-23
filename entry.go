package main

import (
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

// SetCurrentEntry sets the current entry state value to the given entry
func SetCurrentEntry(ctx app.Context, entry Entry) {
	ctx.SetState("currentEntry", entry)
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
		Title:                  "Entry",
		Description:            "Edit a meal entry",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			e.entry = GetCurrentEntry(ctx)
			e.person = GetCurrentPerson(ctx)
		},
		TitleElement: "Edit Entry (" + e.person.Name + ")",
		Elements: []app.UI{
			app.Label().ID("entry-page-date-label").Class("input-label").For("entry-page-last-done-input").Text("When did you eat this?"),
			app.Input().ID("entry-page-date-input").Class("input").Type("date").Value(e.entry.Date.Format("2006-01-02")),
			NewRadioChips("entry-page-type", "What meal did you eat this for?", "Dinner", &e.entry.Type, mealTypes...),
			NewRadioChips("entry-page-source", "How did you get this meal?", "Cooking", &e.entry.Source, mealSources...),
			NewRangeInputPersonMap("entry-page-taste", "How did this taste for you?", &e.entry.Taste, e.person),
			NewRangeInputPersonMap("entry-page-cost", "How much did this cost for you?", &e.entry.Cost, e.person),
			NewRangeInputPersonMap("entry-page-effort", "How much effort did this take for you?", &e.entry.Effort, e.person),
			NewRangeInputPersonMap("entry-page-healthiness", "How healthy was this for you?", &e.entry.Healthiness, e.person),
		},
	}
}
