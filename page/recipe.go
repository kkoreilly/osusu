package page

import (
	"github.com/kkoreilly/osusu/compo"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Recipe struct {
	app.Compo
	recipe osusu.Recipe
}

func (r *Recipe) Render() app.UI {
	return &compo.Page{
		ID:                     "recipe",
		Title:                  "Add Recipe",
		Description:            "View and add a new recipe",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			r.recipe = osusu.CurrentRecipe(ctx)
		},
		TitleElement: "Add Recipe",
		Elements: []app.UI{
			app.Div().ID("recipe-page-info-container").Class("form").Body(
				app.Span().ID("recipe-page-name").Text("Name: "+r.recipe.Name),
				app.If(r.recipe.Description != "",
					app.Span().ID("recipe-page-description").Text("Description: "+r.recipe.Description),
				),
				app.If(r.recipe.TotalTime != "" && r.recipe.TotalTime != "0s",
					app.Span().ID("recipe-page-time").Text("Total Time: "+r.recipe.TotalTime),
				),
				app.If(r.recipe.Image != "",
					app.Span().ID("recipe-page-image-label").Text("Image:"),
					app.Img().ID("recipe-page-image").Src(r.recipe.Image).Alt(r.recipe.Name+" image"),
				),
				&compo.ButtonRow{ID: "recipe-page", Buttons: []app.UI{
					&compo.Button{ID: "recipe-page-back", Class: "secondary", Icon: "arrow_back", Text: "Back", OnClick: compo.ReturnToReturnURL},
					&compo.Button{ID: "recipe-page-view-recipe", Class: "tertiary", Icon: "visibility", Text: "View", OnClick: compo.NavigateEvent(r.recipe.URL)},
					&compo.Button{ID: "recipe-page-add-recipe", Class: "primary", Icon: "add", Text: "Add", OnClick: r.Add},
				}},
			),
		},
	}
}

func (r *Recipe) Add(ctx app.Context, e app.Event) {
	osusu.SetIsMealNew(true, ctx)
	meal := osusu.Meal{
		Name:        r.recipe.Name,
		Description: r.recipe.Description,
		Source:      r.recipe.URL,
		Image:       r.recipe.Image,
		Category:    r.recipe.Category,
		Cuisine:     r.recipe.Cuisine,
	}
	osusu.SetCurrentMeal(meal, ctx)
	compo.Navigate("/meal", ctx)
}
