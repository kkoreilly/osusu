package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Meal is a struct that represents the data of a meal
type Meal struct {
	ID          int64
	GroupID     int64
	Name        string
	Description string
	Cuisine     []string
}

// Meals is a slice that represents multiple meals
type Meals []Meal

// Score produces a score object for the meal based on the given entries and options
func (m Meal) Score(entries Entries, options Options) Score {
	scores := []Score{}
	for _, entry := range entries {
		scores = append(scores, entry.Score(options))
	}
	return AverageScore(scores...)
	// entriesSum := 0
	// var latestDate time.Time
	// for _, entry := range entries {
	// 	entriesSum += entry.Score(options)
	// 	if entry.Date.After(latestDate) {
	// 		latestDate = entry.Date
	// 	}
	// }
	// recencyScore := int(2 * time.Now().Truncate(time.Hour*24).UTC().Sub(latestDate) / (time.Hour * 24))
	// if recencyScore > 100 {
	// 	recencyScore = 100
	// }
	// // add up all of the weights except recency and multiply all of the scores except for recency by them to make the other weights affect how much recency matters
	// weightsToal := options.CostWeight + options.EffortWeight + options.HealthinessWeight + options.TasteWeight
	// sum := weightsToal*entriesSum + options.RecencyWeight*recencyScore
	// den := weightsToal*len(entries) + options.RecencyWeight
	// if den == 0 {
	// 	return 0
	// }
	// return sum / den
}

// RemoveInvalidCuisines returns the the meal with all invalid cuisines removed, using the given cuisine options
func (m Meal) RemoveInvalidCuisines(cuisines []string) Meal {
	res := []string{}
	for _, mealCuisine := range m.Cuisine {
		for _, cuisineOption := range cuisines {
			if mealCuisine == cuisineOption {
				res = append(res, mealCuisine)
			}
		}
	}
	m.Cuisine = res
	return m
}

// CurrentMeal gets and returns the current meal state value, using the given context
func CurrentMeal(ctx app.Context) Meal {
	var meal Meal
	ctx.GetState("currentMeal", &meal)
	return meal
}

// SetCurrentMeal sets the current meal state value to the given meal, using the given context
func SetCurrentMeal(meal Meal, ctx app.Context) {
	ctx.SetState("currentMeal", meal, app.Persist)
}

// IsMealNew gets the state value specifying whether the current meal is new
func IsMealNew(ctx app.Context) bool {
	var isMealNew bool
	ctx.GetState("isMealNew", &isMealNew)
	return isMealNew
}

// SetIsMealNew sets the state value specifying whether the current meal is new
func SetIsMealNew(isMealNew bool, ctx app.Context) {
	ctx.SetState("isMealNew", isMealNew, app.Persist)
}

var (
	mealTypes   = []string{"Breakfast", "Lunch", "Dinner"}
	mealSources = []string{"Cooking", "Dine-In", "Takeout"}
)

type meal struct {
	app.Compo
	group     Group
	user      User
	meal      Meal
	isMealNew bool
	cuisine   map[string]bool
}

func (m *meal) Render() app.UI {
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
	return &Page{
		ID:                     "meal",
		Title:                  titleText,
		Description:            "Edit, view, or create a meal.",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			m.group = CurrentGroup(ctx)
			m.user = CurrentUser(ctx)
			m.meal = CurrentMeal(ctx)
			m.isMealNew = IsMealNew(ctx)

			if m.isMealNew {
				CurrentPage.Title = "Create Meal"
				CurrentPage.UpdatePageTitle(ctx)
			}

			cuisines, err := GetGroupCuisinesAPI.Call(m.group.ID)
			if err != nil {
				CurrentPage.ShowErrorStatus(err)
				return
			}
			m.group.Cuisines = cuisines
			SetCurrentGroup(m.group, ctx)

			m.meal = m.meal.RemoveInvalidCuisines(m.group.Cuisines)

			if m.meal.Cuisine == nil {
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
				TextInput().ID("meal-page-name").Label("Name:").Value(&m.meal.Name).AutoFocus(true),
				Textarea().ID("meal-page-description").Label("Description:").Value(&m.meal.Description),
				CheckboxChips().ID("meal-page-cuisine").Label("Cuisines:").Value(&m.cuisine).Options(append(cuisines, "+")...).OnChange(m.CuisinesOnChange),
				newCuisinesDialog("meal-page", m.CuisinesDialogOnSave),
				ButtonRow().ID("meal-page").Buttons(
					Button().ID("meal-page-cancel").Class("secondary").Icon("cancel").Text("Cancel").OnClick(NavigateEvent("/home")),
					Button().ID("meal-page-save").Class("primary").Type("submit").Icon(saveButtonIcon).Text(saveButtonText),
				),
			),
		},
	}
}

func (m *meal) CuisinesOnChange(ctx app.Context, event app.Event, val string) {
	if val == "+" {
		m.cuisine[val] = false
		event.Get("target").Set("checked", false)
		app.Window().GetElementByID("meal-page-cuisines-dialog").Call("showModal")
	}
}

func (m *meal) CuisinesDialogOnSave(ctx app.Context, event app.Event) {
	m.user = CurrentUser(ctx)
	m.meal = m.meal.RemoveInvalidCuisines(m.group.Cuisines)
}

func (m *meal) OnSubmit(ctx app.Context, event app.Event) {
	event.PreventDefault()

	m.meal.Cuisine = []string{}
	for cuisine, value := range m.cuisine {
		if value {
			m.meal.Cuisine = append(m.meal.Cuisine, cuisine)
		}
	}

	if m.isMealNew {
		m.meal.GroupID = m.group.ID
		meal, err := CreateMealAPI.Call(m.meal)
		if err != nil {
			CurrentPage.ShowErrorStatus(err)
			return
		}
		m.meal = meal
		SetCurrentMeal(m.meal, ctx)
		entries, err := GetEntriesForMealAPI.Call(m.meal.ID)
		if err != nil {
			CurrentPage.ShowErrorStatus(err)
			return
		}
		entry := NewEntry(m.group, m.user, m.meal, entries)
		SetIsEntryNew(true, ctx)
		SetCurrentEntry(entry, ctx)
		Navigate("/entry", ctx)
		return
	}
	_, err := UpdateMealAPI.Call(m.meal)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}

	SetCurrentMeal(m.meal, ctx)

	Navigate("/home", ctx)
}

func (m *meal) ViewEntries(ctx app.Context, event app.Event) {
	event.PreventDefault()

	m.meal.Cuisine = []string{}
	for cuisine, value := range m.cuisine {
		if value {
			m.meal.Cuisine = append(m.meal.Cuisine, cuisine)
		}
	}

	_, err := UpdateMealAPI.Call(m.meal)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	SetCurrentMeal(m.meal, ctx)

	Navigate("/entries", ctx)
}
