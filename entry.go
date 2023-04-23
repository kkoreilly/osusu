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

// Score produces a score from 0 to 100 for the entry based on its attributes and the given options
func (e Entry) Score(options Options) int {
	// average of all attributes
	sum := options.CostWeight*e.Cost.Sum(options.Users, true) + options.EffortWeight*e.Effort.Sum(options.Users, true) + options.HealthinessWeight*e.Healthiness.Sum(options.Users, false) + options.TasteWeight*e.Taste.Sum(options.Users, false)
	den := len(e.Cost)*options.CostWeight + len(e.Effort)*options.EffortWeight + len(e.Healthiness)*options.HealthinessWeight + len(e.Taste)*options.TasteWeight
	if den == 0 {
		return 0
	}
	return sum / den
}

// MissingData returns whether the given user is missing data in the given entry
func (e Entry) MissingData(user User) bool {
	return !(e.Cost.HasValueSet(user) && e.Effort.HasValueSet(user) && e.Healthiness.HasValueSet(user) && e.Taste.HasValueSet(user))
}

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

// SetIsEntryNew sets the state value specifying whether the current entry is new
func SetIsEntryNew(isEntryNew bool, ctx app.Context) {
	ctx.SetState("isEntryNew", isEntryNew, app.Persist)
}

// GetIsEntryNew returns the state value specifying whether the current entry is new
func GetIsEntryNew(ctx app.Context) bool {
	var isEntryNew bool
	ctx.GetState("isEntryNew", &isEntryNew)
	return isEntryNew
}

type entry struct {
	app.Compo
	entry      Entry
	isEntryNew bool
	user       User
}

func (e *entry) Render() app.UI {
	titleText := "Edit Entry"
	saveButtonText := "Save"
	if e.isEntryNew {
		titleText = "Create Entry"
		saveButtonText = "Create"
	}
	return &Page{
		ID:                     "entry",
		Title:                  titleText,
		Description:            "View, edit, or create a meal entry.",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			e.entry = GetCurrentEntry(ctx)
			e.isEntryNew = GetIsEntryNew(ctx)
			if e.isEntryNew {
				CurrentPage.Title = "Create Entry"
				CurrentPage.UpdatePageTitle(ctx)
			}
			// people, err := GetPeopleAPI.Call(GetCurrentUser(ctx).ID)
			// if err != nil {
			// 	CurrentPage.ShowErrorStatus(err)
			// 	return
			// }
			// e.entry = e.entry.RemoveInvalid(people)
			e.user = GetCurrentUser(ctx)
			e.entry = e.entry.FixMissingData(e.user)
		},
		TitleElement: titleText,
		Elements: []app.UI{
			app.Form().ID("entry-page-form").Class("form").OnSubmit(e.OnSubmit).Body(
				app.Label().ID("entry-page-date-label").Class("input-label").For("entry-page-date-input").Text("When did you eat this?"),
				app.Input().ID("entry-page-date-input").Class("input").Type("date").Value(e.entry.Date.Format("2006-01-02")),
				NewRadioChips("entry-page-type", "What meal did you eat this for?", "Dinner", &e.entry.Type, mealTypes...),
				NewRadioChips("entry-page-source", "How did you get this meal?", "Cooking", &e.entry.Source, mealSources...),
				NewRangeInputUserMap("entry-page-taste", "How tasty do you think this was?", &e.entry.Taste, e.user),
				NewRangeInputUserMap("entry-page-cost", "How expensive do you think this was?", &e.entry.Cost, e.user),
				NewRangeInputUserMap("entry-page-effort", "How much effort do you think this took?", &e.entry.Effort, e.user),
				NewRangeInputUserMap("entry-page-healthiness", "How healthy do you think this was?", &e.entry.Healthiness, e.user),
				app.Div().ID("entry-page-action-button-row").Class("action-button-row").Body(
					app.If(!e.isEntryNew,
						app.Button().ID("entry-page-delete-button").Class("action-button", "danger-action-button").Type("button").Text("Delete").OnClick(e.InitialDelete),
					),
					app.Button().ID("entry-page-cancel-button").Class("action-button", "secondary-action-button").Type("button").Text("Cancel").OnClick(ReturnToReturnURL),
					app.Button().ID("entry-page-save-button").Class("action-button", "primary-action-button").Type("submit").Text(saveButtonText),
				),
			),
			app.Dialog().ID("entry-page-confirm-delete").Class("modal").Body(
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

	Navigate("/entries", ctx)
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
