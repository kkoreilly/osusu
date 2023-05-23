package main

import (
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// A MealImageCompo is a component for a meal, recipe, or entry with an image of the meal and information about it
type MealImageCompo struct {
	app.Compo
	id       string
	class    string
	img      string
	mainText string
	score    Score
}

func (m *MealImageCompo) Render() app.UI {
	return app.Div().ID(m.id+"-meal-image-container").Class("meal-image-container", m.class).Body(
		app.Img().ID(m.id+"-meal-image").Class("meal-image").Src(m.img),
		app.Span().ID(m.id+"-meal-image-main-text").Class("meal-image-main-text").Text(m.mainText+" - "+strconv.Itoa(m.score.Total)),
		// app.Span().ID(m.id+"-meal-image-score-text").Class("meal-image-score-text").Text("Score: "+strconv.Itoa(s))
	)
}

// MealImage returns a new meal image component
func MealImage() *MealImageCompo {
	return &MealImageCompo{}
}

// ID sets the id of the meal image component
func (m *MealImageCompo) ID(id string) *MealImageCompo {
	m.id = id
	return m
}

// Class sets the class of the meal image component
func (m *MealImageCompo) Class(class string) *MealImageCompo {
	m.class = class
	return m
}

// Img sets the image of the meal image component
func (m *MealImageCompo) Img(img string) *MealImageCompo {
	m.img = img
	return m
}

// MainText sets the main text of the meal image component
func (m *MealImageCompo) MainText(mainText string) *MealImageCompo {
	m.mainText = mainText
	return m
}

// Score sets the score of the meal image component
func (m *MealImageCompo) Score(score Score) *MealImageCompo {
	m.score = score
	return m
}
