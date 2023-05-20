package main

import (
	"sort"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Entry is an entry with information about how a meal was at a certain point in time
type Entry struct {
	ID          int64
	GroupID     int64
	MealID      int64
	Date        time.Time
	Type        string
	Source      string
	Cost        UserMap
	Effort      UserMap
	Healthiness UserMap
	Taste       UserMap
}

// NewEntry returns a new entry for the given group, user, and meal with the given existing entries for the meal
func NewEntry(group Group, user User, meal Meal, entries Entries) Entry {
	newEntry := Entry{
		GroupID:     group.ID,
		MealID:      meal.ID,
		Date:        time.Now(),
		Type:        "Dinner",
		Source:      "Cooking",
		Cost:        UserMap{user.ID: 50},
		Effort:      UserMap{user.ID: 50},
		Healthiness: UserMap{user.ID: 50},
		Taste:       UserMap{user.ID: 50},
	}
	// if there are previous entries, copy the values from the latest (the entries are already sorted with the latest first by OnNav)
	// we only copy the person map values for the person creating the new entry.
	if len(entries) > 0 {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Date.After(entries[j].Date)
		})
		previousEntry := entries[0]
		previousEntry = previousEntry.FixMissingData(user)
		newEntry = Entry{
			GroupID:     group.ID,
			MealID:      meal.ID,
			Date:        time.Now(),
			Type:        previousEntry.Type,
			Source:      previousEntry.Source,
			Cost:        UserMap{user.ID: previousEntry.Cost[user.ID]},
			Effort:      UserMap{user.ID: previousEntry.Effort[user.ID]},
			Healthiness: UserMap{user.ID: previousEntry.Healthiness[user.ID]},
			Taste:       UserMap{user.ID: previousEntry.Taste[user.ID]},
		}
	}
	return newEntry
}

// Score produces a score object for the entry based on its attributes and the given options
func (e Entry) Score(options Options) Score {
	score := Score{
		Cost:        e.Cost.Sum(options.Users, true) / len(e.Cost),
		Effort:      e.Effort.Sum(options.Users, true) / len(e.Effort),
		Healthiness: e.Healthiness.Sum(options.Users, false) / len(e.Healthiness),
		Taste:       e.Taste.Sum(options.Users, false) / len(e.Taste),
	}
	// recency is irrelevant for entries, so ignore it for this calculation
	options.RecencyWeight = 0
	score.Total = score.ComputeTotal(options)
	return score
}

// // MissingData returns whether the given user is missing data in the given entry
// func (e Entry) MissingData(user User) bool {
// 	return !(e.Cost.HasValueSet(user) && e.Effort.HasValueSet(user) && e.Healthiness.HasValueSet(user) && e.Taste.HasValueSet(user))
// }

// FixMissingData fixes any missing data for the given user for the given entry by setting their values to the average of everyone else's ratings and returning the updated entry
func (e Entry) FixMissingData(user User) Entry {
	if !e.Cost.HasValueSet(user) {
		e.Cost[user.ID] = e.Cost.Average()
	}
	if !e.Effort.HasValueSet(user) {
		e.Effort[user.ID] = e.Effort.Average()
	}
	if !e.Healthiness.HasValueSet(user) {
		e.Healthiness[user.ID] = e.Healthiness.Average()
	}
	if !e.Taste.HasValueSet(user) {
		e.Taste[user.ID] = e.Taste.Average()
	}
	return e
}

// RemoveInvalid returns the entry with all invalid entries associated with nonexistent users removed
func (e Entry) RemoveInvalid(users []User) Entry {
	e.Cost.RemoveInvalid(users)
	e.Effort.RemoveInvalid(users)
	e.Healthiness.RemoveInvalid(users)
	e.Taste.RemoveInvalid(users)
	return e
}

// CurrentEntry gets and returns the current entry state value
func CurrentEntry(ctx app.Context) Entry {
	var entry Entry
	ctx.GetState("currentEntry", &entry)
	return entry
}

// SetCurrentEntry sets the current entry state value to the given entry
func SetCurrentEntry(entry Entry, ctx app.Context) {
	ctx.SetState("currentEntry", entry, app.Persist)
}

// IsEntryNew returns the state value specifying whether the current entry is new
func IsEntryNew(ctx app.Context) bool {
	var isEntryNew bool
	ctx.GetState("isEntryNew", &isEntryNew)
	return isEntryNew
}

// SetIsEntryNew sets the state value specifying whether the current entry is new
func SetIsEntryNew(isEntryNew bool, ctx app.Context) {
	ctx.SetState("isEntryNew", isEntryNew, app.Persist)
}

type entry struct {
	app.Compo
	entry      Entry
	isEntryNew bool
	user       User
}

func (e *entry) Render() app.UI {
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
			e.entry = CurrentEntry(ctx)
			e.isEntryNew = IsEntryNew(ctx)
			if e.isEntryNew {
				CurrentPage.Title = "Create Entry"
				CurrentPage.UpdatePageTitle(ctx)
			}
			e.user = CurrentUser(ctx)
			e.entry = e.entry.FixMissingData(e.user)
		},
		TitleElement: titleText,
		Elements: []app.UI{
			app.Form().ID("entry-page-form").Class("form").OnSubmit(e.OnSubmit).Body(
				DateInput().ID("entry-page-date").Label("When did you eat this?").Value(&e.entry.Date),
				RadioChips().ID("entry-page-type").Label("What meal did you eat this for?").Default("Dinner").Value(&e.entry.Type).Options(mealTypes...),
				RadioChips().ID("entry-page-source").Label("How did you get this meal?").Default("Cooking").Value(&e.entry.Source).Options(mealSources...),
				RangeInputUserMap(&e.entry.Taste, e.user).ID("entry-page-taste").Label("How tasty think this was?"),
				RangeInputUserMap(&e.entry.Cost, e.user).ID("entry-page-cost").Label("How expensive was this?"),
				RangeInputUserMap(&e.entry.Effort, e.user).ID("entry-page-effort").Label("How much effort did this take?"),
				RangeInputUserMap(&e.entry.Healthiness, e.user).ID("entry-page-healthiness").Label("How healthy was this?"),
				ButtonRow().ID("entry-page").Buttons(
					Button().ID("entry-page-delete").Class("danger").Icon("delete").Text("Delete").OnClick(e.InitialDelete).Hidden(e.isEntryNew),
					Button().ID("entry-page-cancel").Class("secondary").Icon("cancel").Text("Cancel").OnClick(ReturnToReturnURL),
					Button().ID("entry-page-save").Class("primary").Type("submit").Icon(saveButtonIcon).Text(saveButtonText),
				),
			),
			app.Dialog().ID("entry-page-confirm-delete").Class("modal").Body(
				app.P().ID("entry-page-confirm-delete-text").Class("confirm-delete-text").Text("Are you sure you want to delete this entry?"),
				ButtonRow().ID("entry-page-confirm-delete").Buttons(
					Button().ID("entry-page-confirm-delete-delete").Class("danger").Icon("delete").Text("Yes, Delete").OnClick(e.ConfirmDelete),
					Button().ID("entry-page-confirm-delete-cancel").Class("secondary").Icon("cancel").Text("No, Cancel").OnClick(e.CancelDelete),
				),
			),
		},
	}
}

func (e *entry) OnSubmit(ctx app.Context, event app.Event) {
	event.PreventDefault()

	if e.isEntryNew {
		entry, err := CreateEntryAPI.Call(e.entry)
		if err != nil {
			CurrentPage.ShowErrorStatus(err)
			return
		}
		e.entry = entry
		SetCurrentEntry(e.entry, ctx)
		ReturnToReturnURL(ctx, event)
		return
	}
	_, err := UpdateEntryAPI.Call(e.entry)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	SetCurrentEntry(e.entry, ctx)
	ReturnToReturnURL(ctx, event)
}

func (e *entry) InitialDelete(ctx app.Context, event app.Event) {
	event.PreventDefault()
	app.Window().GetElementByID("entry-page-confirm-delete").Call("showModal")
}

func (e *entry) ConfirmDelete(ctx app.Context, event app.Event) {
	event.PreventDefault()

	_, err := DeleteEntryAPI.Call(e.entry.ID)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	SetCurrentEntry(Entry{}, ctx)

	Navigate("/entries", ctx)
}

func (e *entry) CancelDelete(ctx app.Context, event app.Event) {
	event.PreventDefault()
	app.Window().GetElementByID("entry-page-confirm-delete").Call("close")
}
