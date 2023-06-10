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
	group        osusu.Group
	user         osusu.User
	users        osusu.Users
	meals        osusu.Meals
	entries      osusu.Entries
	mealEntries  map[int64]osusu.Entries
	recipes      osusu.Recipes
	options      osusu.Options
	usersOptions map[string]bool
}

func (d *Discover) Render() app.UI {
	usersStrings := []string{}
	for _, u := range d.users {
		usersStrings = append(usersStrings, u.Name)
	}
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

			users, err := api.GetUsers.Call(d.group.Members)
			if err != nil {
				compo.CurrentPage.ShowErrorStatus(err)
				return
			}
			d.users = users

			d.options = osusu.GetOptions(ctx)
			if d.options.Users == nil {
				d.options = osusu.DefaultOptions(d.group)
			}
			d.options.Mode = "Discover"
			osusu.SetOptions(d.options, ctx)
			if d.options.UserNames == nil {
				d.options.UserNames = make(map[int64]string)
			}
			for _, user := range d.users {
				d.options.UserNames[user.ID] = user.Name
			}

			d.usersOptions = make(map[string]bool)
			for _, p := range d.users {
				if _, ok := d.options.Users[p.ID]; !ok {
					d.options.Users[p.ID] = true
				}
				d.usersOptions[p.Name] = d.options.Users[p.ID]
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
			compo.ButtonRow().ID("discover-page").Buttons(
				compo.Button().ID("discover-page-search").Class("primary").Icon("search").Text("Search").OnClick(d.ShowOptions),
			),
			compo.ButtonRow().ID("discover-page-quick-options").Buttons(
				compo.CheckboxSelect().ID("discover-page-options-category").Label("Categories:").Default(map[string]bool{"Dinner": true}).Value(&d.options.Category).Options(append(osusu.AllCategories, "Unset")...).OnChange(d.SaveQuickOptions),
				compo.CheckboxSelect().ID("discover-page-options-users").Label("People:").Value(&d.usersOptions).Options(usersStrings...).OnChange(d.SaveQuickOptions),
				compo.CheckboxSelect().ID("discover-page-options-cuisine").Label("Cuisines:").Value(&d.options.Cuisine).Options(osusu.BaseCuisines...).OnChange(d.SaveQuickOptions),
			),
			app.Div().ID("discover-page-recipes-container").Class("meal-images-container").Hidden(d.options.Mode != "Discover").Body(
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
					return compo.MealImage().ID("discover-page-recipe-" + si).Class("discover-page-recipe").Img(recipe.Image).MainText(recipe.Name).SecondaryText(secondaryText).Score(recipe.Score).OnClick(func(ctx app.Context, e app.Event) { d.RecipeOnClick(ctx, e, recipe) }).OnClickScope(recipe.URL)
				}),
			),
			compo.Options().ID("discover-page").Options(&d.options).OnSave(func(ctx app.Context, e app.Event) { d.RecommendRecipes() }),
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

func (d *Discover) SaveQuickOptions(ctx app.Context, e app.Event, val string) {
	d.SaveOptions(ctx, e)
}

func (d *Discover) SaveOptions(ctx app.Context, e app.Event) {
	e.PreventDefault()

	for _, u := range d.users {
		d.options.Users[u.ID] = d.usersOptions[u.Name]
	}

	osusu.SetOptions(d.options, ctx)

	app.Window().GetElementByID("discover-page-options").Call("close")

	d.RecommendRecipes()
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
