package page

import (
	"sort"
	"strconv"

	"github.com/kkoreilly/osusu/api"
	"github.com/kkoreilly/osusu/compo"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/kkoreilly/osusu/util/list"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Discover struct {
	app.Compo
	group       osusu.Group
	user        osusu.User
	meals       osusu.Meals
	entries     osusu.Entries
	mealEntries map[int64]osusu.Entries
	recipes     osusu.Recipes
	options     osusu.Options
}

func (d *Discover) Render() app.UI {
	return &compo.Page{
		ID:                     "discover",
		Title:                  "Discover",
		Description:            "Discover new recipes recommended based on your previous ratings",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			compo.SetReturnURL("/discover", ctx)
			d.group = osusu.CurrentGroup(ctx)
			if d.group.Name == "" {
				compo.Navigate("/groups", ctx)
			}
			d.user = osusu.CurrentUser(ctx)
			cuisines, err := api.GetGroupCuisines.Call(d.group.ID)
			if err != nil {
				compo.CurrentPage.ShowErrorStatus(err)
				return
			}
			d.group.Cuisines = cuisines
			osusu.SetCurrentGroup(d.group, ctx)

			d.options = osusu.GetOptions(ctx)
			if d.options.Users == nil {
				d.options = osusu.DefaultOptions(d.group)
			}

			meals, err := api.GetMeals.Call(d.group.ID)
			if err != nil {
				compo.CurrentPage.ShowErrorStatus(err)
				return
			}
			d.meals = meals

			entries, err := api.GetEntries.Call(d.group.ID)
			if err != nil {
				compo.CurrentPage.ShowErrorStatus(err)
				return
			}
			d.entries = entries
			sort.Slice(d.entries, func(i, j int) bool {
				return d.entries[i].Date.After(d.entries[j].Date)
			})
			d.mealEntries = make(map[int64]osusu.Entries)
			for _, entry := range d.entries {
				entries := d.mealEntries[entry.MealID]
				if entries == nil {
					entries = osusu.Entries{}
				}
				entries = append(entries, entry)
				d.mealEntries[entry.MealID] = entries
			}

			d.RecommendRecipes()
		},
		TitleElement:    "Discover",
		SubtitleElement: "Discover new recipes recommended based on your previous ratings",
		Elements: []app.UI{
			&compo.ButtonRow{ID: "discover-page", Buttons: []app.UI{
				&compo.Button{ID: "discover-page-sort", Class: "primary", Icon: "sort", Text: "Sort", OnClick: d.ShowOptions},
			}},
			&compo.QuickOptions{ID: "discover-page", Options: &d.options, Exclude: []string{"source"}, Group: d.group, Meals: nil, OnSave: func(ctx app.Context, e app.Event) { d.RecommendRecipes() }},
			app.P().ID("discover-page-no-recipes-shown").Class("centered-text").Text("No recipes satisfy your filters. Please try changing them.").Hidden(len(d.recipes) != 0),
			app.Div().ID("discover-page-recipes-container").Class("meal-images-container").Body(
				app.Range(d.recipes).Slice(func(i int) app.UI {
					si := strconv.Itoa(i)
					recipe := d.recipes[i]
					// only put • between category and cuisine if both exist
					secondaryText := ""
					if len(recipe.Category) != 0 && len(recipe.Cuisine) != 0 {
						secondaryText = list.Slice(recipe.Category) + " • " + list.Slice(recipe.Cuisine)
					} else {
						secondaryText = list.Slice(recipe.Category) + list.Slice(recipe.Cuisine)
					}
					return &compo.MealImage{ID: "discover-page-recipe-" + si, Class: "discover-page-recipe", Image: recipe.Image, MainText: recipe.Name, SecondaryText: secondaryText, Score: recipe.Score, OnClick: func(ctx app.Context, e app.Event) { d.RecipeOnClick(ctx, e, recipe) }, OnClickScope: []any{recipe.URL}}
				}),
			),
			&compo.Options{ID: "discover-page", Options: &d.options, OnSave: func(ctx app.Context, e app.Event) { d.RecommendRecipes() }},
		},
	}
}

func (d *Discover) RecipeOnClick(ctx app.Context, e app.Event, recipe osusu.Recipe) {
	osusu.SetCurrentRecipe(recipe, ctx)
	compo.Navigate("/recipe", ctx)
}

func (d *Discover) ShowOptions(ctx app.Context, e app.Event) {
	app.Window().GetElementByID("discover-page-options").Call("showModal")
}

// RecommendRecipes loads recipe recommendations for discover mode
func (d *Discover) RecommendRecipes() {
	wordScoreMap := osusu.WordScoreMap(d.meals, d.mealEntries, d.options)
	usedSources := map[string]bool{}
	for _, meal := range d.meals {
		if meal.Source != "" {
			usedSources[meal.Source] = true
		}
	}
	recipes, err := api.RecommendRecipes.Call(osusu.RecommendRecipesData{WordScoreMap: wordScoreMap, Options: d.options, UsedSources: usedSources, N: 0})
	if err != nil {
		compo.CurrentPage.ShowErrorStatus(err)
		return
	}
	d.recipes = recipes
}
