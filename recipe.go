package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type recipe struct {
	app.Compo
	recipe Recipe
}

func (r *recipe) Render() app.UI {
	return &Page{
		ID:                     "recipe",
		Title:                  "Add Recipe",
		Description:            "View and add a new recipe",
		AuthenticationRequired: true,
		OnNavFunc: func(ctx app.Context) {
			r.recipe = CurrentRecipe(ctx)
		},
		TitleElement: "Add Recipe",
		Elements: []app.UI{
			app.Div().ID("recipe-page-info-container").Class("form").Body(
				app.Span().ID("recipe-page-name").Text("Name: "+r.recipe.Name),
				app.Span().ID("recipe-page-source").Text("Source: "+r.recipe.Source),
				app.If(r.recipe.TotalTime != "" && r.recipe.TotalTime != "0s",
					app.Span().ID("recipe-page-time").Text("Total Time: "+r.recipe.TotalTime),
				),
				app.If(r.recipe.Description != "",
					app.Span().ID("recipe-page-description").Text("Description: "+r.recipe.Description),
				),
				app.If(r.recipe.Image != "",
					app.Span().ID("recipe-page-image-label").Text("Image:"),
					app.Img().ID("recipe-page-image").Src(r.recipe.Image).Alt(r.recipe.Name+" image"),
				),
				ButtonRow().ID("recipe-page").Buttons(
					Button().ID("recipe-page-back").Class("secondary").Icon("arrow_back").Text("Back").OnClick(ReturnToReturnURL),
					Button().ID("recipe-page-view-recipe").Class("tertiary").Icon("visibility").Text("View").OnClick(NavigateEvent(r.recipe.URL)),
					Button().ID("recipe-page-add-recipe").Class("primary").Icon("add").Text("Add"),
				),
			),
		},
	}
}
