package osusu

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
	Category    string
	Source      string
	Cost        UserMap
	Effort      UserMap
	Healthiness UserMap
	Taste       UserMap
}

// NewEntry returns a new entry for the given group, user, and meal with the given existing entries for the meal
func NewEntry(group Group, user User, meal Meal, entries Entries) Entry {
	category := ""
	if len(meal.Category) != 0 {
		category = meal.Category[0]
	}
	newEntry := Entry{
		GroupID:     group.ID,
		MealID:      meal.ID,
		Date:        time.Now(),
		Category:    category,
		Source:      "Cooking",
		Cost:        UserMap{user.ID: 50},
		Effort:      UserMap{user.ID: 50},
		Healthiness: UserMap{user.ID: 50},
		Taste:       UserMap{user.ID: 50},
	}
	// if there are previous entries, copy the values from the latest entry where the person has rated things, if there is one.
	// we only copy the person map values for the person creating the new entry.
	if len(entries) > 0 {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Date.After(entries[j].Date)
		})
		for _, previousEntry := range entries {
			if previousEntry.Cost[user.ID] != 0 {
				newEntry = Entry{
					GroupID:     group.ID,
					MealID:      meal.ID,
					Date:        time.Now(),
					Category:    previousEntry.Category,
					Source:      previousEntry.Source,
					Cost:        UserMap{user.ID: previousEntry.Cost[user.ID]},
					Effort:      UserMap{user.ID: previousEntry.Effort[user.ID]},
					Healthiness: UserMap{user.ID: previousEntry.Healthiness[user.ID]},
					Taste:       UserMap{user.ID: previousEntry.Taste[user.ID]},
				}
			}
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
