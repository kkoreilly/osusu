package page

import (
	"sort"
	"strconv"
	"time"

	"github.com/kkoreilly/osusu/api"
	"github.com/kkoreilly/osusu/compo"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/kkoreilly/osusu/util/list"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Home struct {
	app.Compo
	group              osusu.Group
	user               osusu.User
	users              osusu.Users
	meals              osusu.Meals
	entries            osusu.Entries
	entriesForEachMeal map[int64]osusu.Entries // entries for each meal id
	recipes            osusu.Recipes
	options            osusu.Options
	usersOptions       map[string]bool
	cuisinesInUse      map[string]bool // key is cuisine name, value is if it's being used
	currentMeal        osusu.Meal      // the current selected meal for the context menu
}

func (h *Home) Render() app.UI {
	usersStrings := []string{}
	for _, u := range h.users {
		usersStrings = append(usersStrings, u.Name)
	}
	cuisines := []string{}
	for cuisine, val := range h.cuisinesInUse {
		if val {
			cuisines = append(cuisines, cuisine)
		}
	}
	// need to sort so options don't keep swapping
	sort.Strings(cuisines)
	// width, _ := app.Window().Size()
	// smallScreen := width <= 480
	subtitleText := ""
	switch h.options.Mode {
	case "Search":
		subtitleText = "Search for the best meals to eat given your current circumstances"
	case "Discover":
		subtitleText = "Discover new recipes recommended based on your previous ratings"
	case "History":
		subtitleText = "View the history of what meals you've eaten and how they were"
	}
	return &compo.Page{
		ID:                     "home",
		Title:                  "Home",
		Description:            "View, sort, and filter your meals.",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			compo.SetReturnURL("/home", ctx)
			h.group = osusu.CurrentGroup(ctx)
			if h.group.Name == "" {
				compo.Navigate("/groups", ctx)
			}
			h.user = osusu.CurrentUser(ctx)
			cuisines, err := api.GetGroupCuisines.Call(h.group.ID)
			if err != nil {
				compo.CurrentPage.ShowErrorStatus(err)
				return
			}
			h.group.Cuisines = cuisines
			osusu.SetCurrentGroup(h.group, ctx)

			users, err := api.GetUsers.Call(h.group.Members)
			if err != nil {
				compo.CurrentPage.ShowErrorStatus(err)
				return
			}
			h.users = users

			h.options = osusu.GetOptions(ctx)
			if h.options.Users == nil {
				h.options = osusu.DefaultOptions(h.group)
			}
			h.options = h.options.RemoveInvalidCuisines(h.group.Cuisines)
			if h.options.UserNames == nil {
				h.options.UserNames = make(map[int64]string)
			}
			for _, user := range h.users {
				h.options.UserNames[user.ID] = user.Name
			}
			switch ctx.Page().URL().Path {
			case "/search":
				h.options.Mode = "Search"
			case "/history":
				h.options.Mode = "History"
			case "/discover":
				h.options.Mode = "Discover"
			}
			osusu.SetOptions(h.options, ctx)
			h.Update()

			h.usersOptions = make(map[string]bool)
			for _, p := range h.users {
				if _, ok := h.options.Users[p.ID]; !ok {
					h.options.Users[p.ID] = true
				}
				h.usersOptions[p.Name] = h.options.Users[p.ID]
			}

			meals, err := api.GetMeals.Call(h.group.ID)
			if err != nil {
				compo.CurrentPage.ShowErrorStatus(err)
				return
			}
			h.meals = meals

			h.cuisinesInUse = map[string]bool{}
			for _, meal := range h.meals {
				for _, cuisine := range meal.Cuisine {
					h.cuisinesInUse[cuisine] = true
					// if the user has not yet set whether or not to allow this cuisine (if it is new), automatically set it to true
					_, ok := h.options.Cuisine[cuisine]
					if !ok {
						h.options.Cuisine[cuisine] = true
					}
				}
			}
			for cuisine := range h.options.Cuisine {
				if !h.cuisinesInUse[cuisine] {
					delete(h.options.Cuisine, cuisine)
				}
			}

			entries, err := api.GetEntries.Call(h.group.ID)
			if err != nil {
				compo.CurrentPage.ShowErrorStatus(err)
				return
			}
			h.entries = entries
			sort.Slice(h.entries, func(i, j int) bool {
				return h.entries[i].Date.After(h.entries[j].Date)
			})
			h.entriesForEachMeal = make(map[int64]osusu.Entries)
			for _, entry := range h.entries {
				entries := h.entriesForEachMeal[entry.MealID]
				if entries == nil {
					entries = osusu.Entries{}
				}
				entries = append(entries, entry)
				h.entriesForEachMeal[entry.MealID] = entries
			}
			h.SortMeals()

			h.RecommendRecipes()

			compo.CurrentPage.AddOnClick(h.PageOnClick)
		},
		TitleElement:    h.options.Mode,
		SubtitleElement: subtitleText,
		Elements: []app.UI{
			compo.ButtonRow().ID("home-page").Buttons(
				compo.Button().ID("home-page-new").Class("secondary").Icon("add").Text("New Meal").OnClick(h.NewMeal),
				compo.Button().ID("home-page-search").Class("primary").Icon("search").Text("Search").OnClick(h.ShowOptions),
			),
			compo.ButtonRow().ID("home-page-quick-options").Buttons(
				compo.CheckboxSelect().ID("home-page-options-category").Label("Categories:").Default(map[string]bool{"Dinner": true}).Value(&h.options.Category).Options(append(osusu.AllCategories, "Unset")...).OnChange(h.SaveQuickOptions),
				compo.CheckboxSelect().ID("home-page-options-users").Label("People:").Value(&h.usersOptions).Options(usersStrings...).OnChange(h.SaveQuickOptions),
				compo.CheckboxSelect().ID("home-page-options-source").Label("Sources:").Default(map[string]bool{"Cooking": true, "Dine-In": true, "Takeout": true}).Value(&h.options.Source).Options(osusu.AllSources...).OnChange(h.SaveQuickOptions).Hidden(h.options.Mode == "Discover"),
				compo.CheckboxSelect().ID("home-page-options-cuisine").Label("Cuisines:").Value(&h.options.Cuisine).Options(append(cuisines, "Unset")...).OnChange(h.SaveQuickOptions),
			),
			app.Div().ID("home-page-meals-container").Class("meal-images-container").Hidden(h.options.Mode != "Search").Body(
				app.Range(h.meals).Slice(func(i int) app.UI {
					si := strconv.Itoa(i)
					meal := h.meals[i]
					entries := h.entriesForEachMeal[meal.ID]

					// check if at least one category satisfies a category option (or there are no categories and unset is an option)
					gotCategory := len(meal.Category) == 0 && h.options.Category["Unset"]
					if !gotCategory {
						for _, mealCategory := range meal.Category {
							for optionCategory, value := range h.options.Category {
								if value && mealCategory == optionCategory {
									gotCategory = true
									break
								}
							}
							if gotCategory {
								break
							}
						}
						if !gotCategory {
							return app.Text("")
						}
					}

					// check if at least one cuisine satisfies a cuisine option (or there are no cuisines and unset is an option)
					gotCuisine := len(meal.Cuisine) == 0 && h.options.Cuisine["Unset"]
					if !gotCuisine {
						for _, mealCuisine := range meal.Cuisine {
							for optionCuisine, value := range h.options.Cuisine {
								if value && mealCuisine == optionCuisine {
									gotCuisine = true
									break
								}
							}
							if gotCuisine {
								break
							}
						}
						if !gotCuisine {
							return app.Text("")
						}
					}

					// check if at least one entry satisfies the source requirements if there is at least one entry.
					if len(entries) > 0 {
						// gotType := false
						gotSource := false
						for _, entry := range entries {
							// if h.options.Category[entry.Type] {
							// 	gotType = true
							// }
							if h.options.Source[entry.Source] {
								gotSource = true
							}
						}
						if !gotSource {
							return app.Text("")
						}
					}
					score := meal.Score(entries, h.options)
					// isCurrentMeal := meal.ID == h.currentMeal.ID

					// only put • between category and cuisine if both exist
					secondaryText := ""
					if len(meal.Category) != 0 && len(meal.Cuisine) != 0 {
						secondaryText = list.Slice(meal.Category) + " • " + list.Slice(meal.Cuisine)
					} else {
						secondaryText = list.Slice(meal.Category) + list.Slice(meal.Cuisine)
					}
					return compo.MealImage().ID("home-page-meal-" + si).Class("home-page-meal").Img(meal.Image).MainText(meal.Name).SecondaryText(secondaryText).Score(score).OnClick(func(ctx app.Context, e app.Event) { h.MealOnClick(ctx, e, meal) }).OnClickScope(meal.ID)
				}),
			),
			app.Div().ID("home-page-recipes-container").Class("meal-images-container").Hidden(h.options.Mode != "Discover").Body(
				app.Range(h.recipes).Slice(func(i int) app.UI {
					si := strconv.Itoa(i)
					recipe := h.recipes[i]
					// only put • between category and cuisine if both exist
					secondaryText := ""
					if len(recipe.Category) != 0 && len(recipe.Cuisine) != 0 {
						secondaryText = list.Slice(recipe.Category) + " • " + list.Slice(recipe.Cuisine)
					} else {
						secondaryText = list.Slice(recipe.Category) + list.Slice(recipe.Cuisine)
					}
					return compo.MealImage().ID("home-page-recipe-" + si).Class("home-page-recipe").Img(recipe.Image).MainText(recipe.Name).SecondaryText(secondaryText).Score(recipe.Score).OnClick(func(ctx app.Context, e app.Event) { h.RecipeOnClick(ctx, e, recipe) }).OnClickScope(recipe.URL)
				}),
			),
			app.Div().ID("home-page-entries-container").Class("meal-images-container").Hidden(h.options.Mode != "History").Body(
				app.Range(h.entries).Slice(func(i int) app.UI {
					si := strconv.Itoa(i)
					entry := h.entries[i]
					score := entry.Score(h.options)
					entryMeal := osusu.Meal{}
					for _, meal := range h.meals {
						if meal.ID == entry.MealID {
							entryMeal = meal
							break
						}
					}
					return compo.MealImage().ID("home-page-entry-" + si).Class("home-page-entry").Img(entryMeal.Image).MainText(entryMeal.Name).SecondaryText(entry.Date.Format("Monday, January 2, 2006")).Score(score).OnClick(func(ctx app.Context, e app.Event) { h.EntryOnClick(ctx, e, entry) }).OnClickScope(entry.ID)
				}),
			),
			// MealImage().ID("test").Img("https://static01.nyt.com/images/2021/02/17/dining/17tootired-grilled-cheese/17tootired-grilled-cheese-articleLarge.jpg?quality=75&auto=webp&disable=upscale").MainText("Grilled Cheese").Score(Score{Total: 76}),
			// app.Table().ID("home-page-meals-table").Body(
			// 	app.THead().ID("home-page-meals-table-header").Body(
			// 		app.Tr().ID("home-page-meals-table-header-row").Body(
			// 			app.If(h.options.Mode == "History",
			// 				app.Th().ID("home-page-meals-table-header-date").Text("Date"),
			// 				app.Th().ID("home-page-meals-table-header-name").Text("Meal"),
			// 			).Else(
			// 				app.Th().ID("home-page-meals-table-header-name").Text("Name"),
			// 			),
			// 			app.Th().ID("home-page-meals-table-header-total").Class("table-header-score").Text("Total"),
			// 			app.Th().ID("home-page-meals-table-header-taste").Class("table-header-score").Text("Taste"),
			// 			app.If(h.options.Mode != "History",
			// 				app.Th().ID("home-page-meals-table-header-recency").Class("table-header-score").Text("New"),
			// 			),
			// 			app.Th().ID("home-page-meals-table-header-cost").Class("table-header-score").Text("Cost"),
			// 			app.Th().ID("home-page-meals-table-header-effort").Class("table-header-score").Text("Effort"),
			// 			app.Th().ID("home-page-meals-table-header-healthiness").Class("table-header-score").Text("Health"),
			// 			app.If(!smallScreen,
			// 				app.If(h.options.Mode == "History",
			// 					app.Th().ID("home-page-meals-table-header-type").Text("Type"),
			// 					app.Th().ID("home-page-meals-table-header-source").Text("Source"),
			// 				).Else(
			// 					app.Th().ID("home-page-meals-table-header-cuisines").Text("Cuisines"),
			// 					app.Th().ID("home-page-meals-table-header-description").Text("Description"),
			// 				),
			// 			),
			// 		),
			// 	),
			// 	app.TBody().ID("home-page-meals-table-body").Body(
			// 		app.If(h.options.Mode == "Search",
			// 			app.Range(h.meals).Slice(func(i int) app.UI {
			// 				meal := h.meals[i]
			// 				si := strconv.Itoa(i)
			// 				entries := h.entriesForEachMeal[meal.ID]

			// 				// check if at least one cuisine satisfies a cuisine requirement (or there is no cuisine set)
			// 				gotCuisine := len(meal.Cuisine) == 0
			// 				for _, mealCuisine := range meal.Cuisine {
			// 					for optionCuisine, value := range h.options.Cuisine {
			// 						if value && mealCuisine == optionCuisine {
			// 							gotCuisine = true
			// 						}
			// 					}
			// 				}
			// 				if !gotCuisine {
			// 					return app.Text("")
			// 				}

			// 				// check if at least one entry satisfies the type and source requirements if there is at least one entry.
			// 				if len(entries) > 0 {
			// 					gotType := h.options.Type == "Any"
			// 					gotSource := false
			// 					for _, entry := range entries {
			// 						if entry.Type == h.options.Type {
			// 							gotType = true
			// 						}
			// 						if h.options.Source[entry.Source] {
			// 							gotSource = true
			// 						}
			// 					}
			// 					if !(gotType && gotSource) {
			// 						return app.Text("")
			// 					}
			// 				}

			// 				score := meal.Score(entries, h.options)
			// 				colorH := strconv.Itoa((score.Total * 12) / 10)
			// 				scoreText := strconv.Itoa(score.Total)
			// 				// missingData := entries.MissingData(h.user)
			// 				isCurrentMeal := meal.ID == h.currentMeal.ID
			// 				return app.Tr().ID("home-page-meal-"+si).Class("home-page-meal").DataSet("current-meal", isCurrentMeal).Style("--color-h", colorH).Style("--score", scoreText+"%").
			// 					OnClick(func(ctx app.Context, e app.Event) { h.MealOnClick(ctx, e, meal) }).Body(
			// 					app.Td().ID("home-page-meal-name-"+si).Class("home-page-meal-name").Text(meal.Name),
			// 					MealScore("home-page-meal-total-"+si, "home-page-meal-total", score.Total, "Total"),
			// 					MealScore("home-page-meal-taste-"+si, "home-page-meal-taste", score.Taste, "Taste"),
			// 					MealScore("home-page-meal-recency-"+si, "home-page-meal-recency", score.Recency, "Recency"),
			// 					MealScore("home-page-meal-cost-"+si, "home-page-meal-cost", score.Cost, "Cost"),
			// 					MealScore("home-page-meal-effort-"+si, "home-page-meal-effort", score.Effort, "Effort"),
			// 					MealScore("home-page-meal-healthiness-"+si, "home-page-meal-healthiness", score.Healthiness, "Healthiness"),
			// 					app.If(!smallScreen,
			// 						app.Td().ID("home-page-meal-cuisines-"+si).Class("home-page-meal-cuisines").Text(ListString(meal.Cuisine)),
			// 						app.Td().ID("home-page-meal-description-"+si).Class("home-page-meal-description").Text(meal.Description),
			// 					),
			// 				)
			// 			}),
			// 		).ElseIf(h.options.Mode == "History",
			// 			app.Range(h.entries).Slice(func(i int) app.UI {
			// 				si := strconv.Itoa(i)
			// 				entry := h.entries[i]
			// 				score := entry.Score(h.options)
			// 				entryMeal := Meal{}
			// 				for _, meal := range h.meals {
			// 					if meal.ID == entry.MealID {
			// 						entryMeal = meal
			// 					}
			// 				}
			// 				return app.Tr().ID("home-page-entry-"+si).Class("home-page-entry home-page-meal").OnClick(func(ctx app.Context, e app.Event) {
			// 					h.EntryOnClick(ctx, e, entry)
			// 				}).Body(
			// 					app.Td().ID("home-page-entry-date-"+si).Class("home-page-entry-date home-page-meal-name").Text(entry.Date.Format("Jan 2, 2006")),
			// 					app.Td().ID("home-page-meal-name-"+si).Class("home-page-meal-name").Text(entryMeal.Name),
			// 					MealScore("home-page-entry-total-"+si, "home-page-entry-total", score.Total, "Total"),
			// 					MealScore("home-page-entry-taste-"+si, "home-page-entry-taste", score.Taste, "Taste"),
			// 					// MealScore("home-page-meal-recency-"+si, "home-page-meal-recency", score.Recency),
			// 					MealScore("home-page-entry-cost-"+si, "home-page-entry-cost", score.Cost, "Cost"),
			// 					MealScore("home-page-entry-effort-"+si, "home-page-entry-effort", score.Effort, "Effort"),
			// 					MealScore("home-page-entry-healthiness-"+si, "home-page-entry-healthiness", score.Healthiness, "Healthiness"),
			// 					app.If(!smallScreen,
			// 						app.Td().ID("home-page-entry-type-"+si).Class("home-page-entry-type").Text(entry.Type),
			// 						app.Td().ID("home-page-entry-source-"+si).Class("home-page-entry-source").Text(entry.Source),
			// 					),
			// 				)
			// 			}),
			// 		).ElseIf(h.options.Mode == "Discover",
			// 			app.Range(h.recipes).Slice(func(i int) app.UI {
			// 				si := strconv.Itoa(i)
			// 				recipe := h.recipes[i]
			// 				// return MealImage().ID("home-page-recipe-" + si).Class("home-page-recipe-image").Img(recipe.Image).MainText(recipe.Name).Score(recipe.Score)
			// 				return app.Tr().ID("home-page-recipe-"+si).Class("home-page-recipe home-page-meal").Style("--img", "url("+recipe.Image+")").OnClick(func(ctx app.Context, e app.Event) { h.RecipeOnClick(ctx, e, recipe) }).Body(
			// 					app.Td().ID("home-page-recipe-name-"+si).Class("home-page-meal-name").Text(recipe.Name),
			// 					MealScore("home-page-recipe-total-"+si, "home-page-meal-total", recipe.Score.Total, "Total"),
			// 					MealScore("home-page-recipe-taste-"+si, "home-page-meal-taste", recipe.Score.Taste, "Taste"),
			// 					MealScore("home-page-recipe-recency-"+si, "home-page-meal-recency", recipe.Score.Recency, "Recency"),
			// 					MealScore("home-page-recipe-cost-"+si, "home-page-meal-cost", recipe.Score.Cost, "Cost"),
			// 					MealScore("home-page-recipe-effort-"+si, "home-page-meal-effort", recipe.Score.Effort, "Effort"),
			// 					MealScore("home-page-recipe-healthiness-"+si, "home-page-meal-healthiness", recipe.Score.Healthiness, "Healthiness"),
			// 					app.If(!smallScreen,
			// 						app.Td().ID("home-page-recipe-cuisines-"+si).Class("home-page-meal-cuisines").Text(ListString(recipe.Cuisine)),
			// 						app.Td().ID("home-page-recipe-description-"+si).Class("home-page-meal-description").Text(recipe.Description),
			// 					),
			// 				)
			// 			}),
			// 		),
			// 	),
			// ),
			app.Dialog().ID("home-page-meal-dialog").OnClick(h.MealDialogOnClick).Body(
				compo.Button().ID("home-page-meal-dialog-new-entry").Class("primary").Icon("add").Text("New Entry").OnClick(h.NewEntry),
				compo.Button().ID("home-page-meal-dialog-view-entries").Class("secondary").Icon("visibility").Text("View Entries").OnClick(h.ViewEntries),
				compo.Button().ID("home-page-meal-dialog-edit-meal").Class("secondary").Icon("edit").Text("Edit Meal").OnClick(h.EditMeal),
			),
			app.Dialog().ID("home-page-options").Class("modal").OnClick(h.OptionsOnClick).Body(
				app.Div().ID("home-page-options-container").OnClick(h.OptionsContainerOnClick).Body(
					app.Form().ID("home-page-options-form").Class("form").OnSubmit(h.SaveOptions).Body(
						// RadioChips().ID("home-page-options-mode").Label("What mode do you want to use?").Default("Search").Value(&h.options.Mode).Options("Search", "History", "Discover"),
						// RadioChips().ID("home-page-options-type").Label("What meal are you eating?").Default("Dinner").Value(&h.options.Type).Options(mealTypes...),
						// CheckboxChips().ID("home-page-options-users").Label("Who are you eating with?").Value(&h.usersOptions).Options(usersStrings...),
						// CheckboxChips().ID("home-page-options-source").Label("What meal sources are okay?").Default(map[string]bool{"Cooking": true, "Dine-In": true, "Takeout": true}).Value(&h.options.Source).Options(mealSources...),
						// CheckboxChips().ID("home-page-options-cuisine").Label("What cuisines are okay?").Value(&h.options.Cuisine).Options(cuisines...),
						compo.RangeInput().ID("home-page-options-taste").Label("How important is taste?").Value(&h.options.TasteWeight),
						compo.RangeInput().ID("home-page-options-recency").Label("How important is recency?").Value(&h.options.RecencyWeight),
						compo.RangeInput().ID("home-page-options-cost").Label("How important is cost?").Value(&h.options.CostWeight),
						compo.RangeInput().ID("home-page-options-effort").Label("How important is effort?").Value(&h.options.EffortWeight),
						compo.RangeInput().ID("home-page-options-healthiness").Label("How important is healthiness?").Value(&h.options.HealthinessWeight),
					),
				),
			),
		},
	}
}

func (h *Home) OptionsOnClick(ctx app.Context, e app.Event) {
	// if the options dialog on click event is triggered, close the options because the dialog includes the whole page and a separate event will cancel this if they actually clicked on the dialog
	h.SaveOptions(ctx, e)
}

func (h *Home) OptionsContainerOnClick(ctx app.Context, e app.Event) {
	// cancel the closing of the dialog if they actually click on the dialog
	e.Call("stopPropagation")
}

func (h *Home) NewMeal(ctx app.Context, e app.Event) {
	osusu.SetIsMealNew(true, ctx)
	osusu.SetCurrentMeal(osusu.Meal{}, ctx)
	compo.Navigate("/meal", ctx)
}

func (h *Home) PageOnClick(ctx app.Context, e app.Event) {
	// close meal dialog on page click (this will be stopped by another event if someone clicks on the meal dialog itself)
	mealDialog := app.Window().GetElementByID("home-page-meal-dialog")
	// need to check for null because people can click on page before dialog exists
	if !mealDialog.IsNull() {
		mealDialog.Call("close")
		h.currentMeal = osusu.Meal{}
	}
}

func (h *Home) MealOnClick(ctx app.Context, e app.Event, meal osusu.Meal) {
	e.Call("stopPropagation")
	h.currentMeal = meal
	osusu.SetCurrentMeal(meal, ctx)
	dialog := app.Window().GetElementByID("home-page-meal-dialog")
	if dialog.Get("open").Bool() {
		ctx.Dispatch(func(ctx app.Context) {
			dialog.Call("close")
		})
		ctx.Defer(func(ctx app.Context) {
			time.Sleep(250 * time.Millisecond)
			h.UpdateMealDialogPosition(ctx, e, dialog)
			dialog.Call("show")
		})
		return
	}
	h.UpdateMealDialogPosition(ctx, e, dialog)
	dialog.Call("show")
}

func (h *Home) EntryOnClick(ctx app.Context, e app.Event, entry osusu.Entry) {
	osusu.SetIsEntryNew(false, ctx)
	osusu.SetCurrentEntry(entry, ctx)
	compo.Navigate("/entry", ctx)
}

func (h *Home) RecipeOnClick(ctx app.Context, e app.Event, recipe osusu.Recipe) {
	osusu.SetCurrentRecipe(recipe, ctx)
	compo.Navigate("/recipe", ctx)
}

func (h *Home) UpdateMealDialogPosition(ctx app.Context, e app.Event, dialog app.Value) {
	pageX, pageY := e.Get("pageX").Int(), e.Get("pageY").Int()
	clientX, clientY := e.Get("clientX").Int(), e.Get("clientY").Int()
	clientWidth, clientHeight := dialog.Get("clientWidth").Int(), dialog.Get("clientHeight").Int()
	pageWidth, pageHeight := app.Window().Size()
	translateX, translateY := "0%", "0%"
	if clientX+clientWidth >= pageWidth {
		translateX = "-100%"
	}
	if clientY+clientHeight >= pageHeight {
		translateY = "-100%"
	}
	dialog.Get("style").Set("top", strconv.Itoa(pageY)+"px")
	dialog.Get("style").Set("left", strconv.Itoa(pageX)+"px")
	dialog.Get("style").Set("transform", "translate("+translateX+", "+translateY)
}

func (h *Home) MealDialogOnClick(ctx app.Context, e app.Event) {
	// stop the meal dialog from being closed by the page on click event
	e.Call("stopPropagation")
}

func (h *Home) NewEntry(ctx app.Context, e app.Event) {
	entry := osusu.NewEntry(h.group, h.user, h.currentMeal, h.entriesForEachMeal[h.currentMeal.ID])
	osusu.SetIsEntryNew(true, ctx)
	osusu.SetCurrentEntry(entry, ctx)
	compo.Navigate("/entry", ctx)
}

func (h *Home) ViewEntries(ctx app.Context, e app.Event) {
	compo.Navigate("/entries", ctx)
}

func (h *Home) EditMeal(ctx app.Context, e app.Event) {
	osusu.SetIsMealNew(false, ctx)
	compo.Navigate("/meal", ctx)
}

func (h *Home) ShowOptions(ctx app.Context, e app.Event) {
	app.Window().GetElementByID("home-page-options").Call("showModal")
}

func (h *Home) SaveQuickOptions(ctx app.Context, e app.Event, val string) {
	h.SaveOptions(ctx, e)
}

func (h *Home) SaveOptions(ctx app.Context, e app.Event) {
	e.PreventDefault()

	for _, u := range h.users {
		h.options.Users[u.ID] = h.usersOptions[u.Name]
	}

	osusu.SetOptions(h.options, ctx)

	app.Window().GetElementByID("home-page-options").Call("close")

	if h.options.Mode == "Search" {
		h.SortMeals()
	} else if h.options.Mode == "Discover" {
		h.RecommendRecipes()
	}
}

// RecommendRecipes loads recipe recommendations for discover mode
func (h *Home) RecommendRecipes() {
	wordScoreMap := osusu.WordScoreMap(h.meals, h.entriesForEachMeal, h.options)
	usedSources := map[string]bool{}
	for _, meal := range h.meals {
		if meal.Source != "" {
			usedSources[meal.Source] = true
		}
	}
	recipes, err := api.RecommendRecipes.Call(osusu.RecommendRecipesData{WordScoreMap: wordScoreMap, Options: h.options, UsedSources: usedSources, N: 0})
	if err != nil {
		compo.CurrentPage.ShowErrorStatus(err)
		return
	}
	h.recipes = recipes
}

func (h *Home) SortMeals() {
	sort.Slice(h.meals, func(i, j int) bool {
		// // prioritize meals with missing data, then score
		// mealI := h.meals[i]
		// entriesI := h.entriesForEachMeal[mealI.ID]
		// iMissingData := entriesI.MissingData(h.user)
		// mealJ := h.meals[j]
		// entriesJ := h.entriesForEachMeal[mealJ.ID]
		// jMissingData := entriesJ.MissingData(h.user)
		// if iMissingData && !jMissingData {
		// 	return true
		// }
		// if !iMissingData && jMissingData {
		// 	return false
		// }
		// sort by recency in history mode, score otherwise
		if h.options.Mode == "History" {
			return h.meals[i].LatestDate(h.entriesForEachMeal[h.meals[i].ID]).After(h.meals[j].LatestDate(h.entriesForEachMeal[h.meals[j].ID]))
		}
		return h.meals[i].Score(h.entriesForEachMeal[h.meals[i].ID], h.options).Total > h.meals[j].Score(h.entriesForEachMeal[h.meals[j].ID], h.options).Total

	})
}
