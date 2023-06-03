package compo

import (
	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// A MealImageCompo is a component for a meal, recipe, or entry with an image of the meal and information about it
type MealImageCompo struct {
	app.Compo
	// everything can change so they need to be exported
	IDValue            string
	ClassValue         string
	ImgValue           string
	MainTextValue      string
	SecondaryTextValue string
	ScoreValue         osusu.Score
	OnClickValue       app.EventHandler
	OnClickScopeValue  []any
}

func (m *MealImageCompo) Render() app.UI {
	return app.Div().ID(m.IDValue+"-meal-image-container").Class("meal-image-container", m.ClassValue).DataSet("no-image", m.ImgValue == "").OnClick(m.OnClickValue, m.OnClickScopeValue...).Body(
		app.Img().ID(m.IDValue+"-meal-image").Class("meal-image").Src(m.ImgValue),
		app.Span().ID(m.IDValue+"-meal-image-main-text").Class("meal-image-main-text").Text(m.MainTextValue),
		app.Span().ID(m.IDValue+"-meal-image-secondary-text").Class("meal-image-secondary-text").Text(m.SecondaryTextValue),
		app.Div().ID(m.IDValue+"-meal-image-info-container").Class("meal-image-info-container").Body(
			MealScore(m.IDValue+"-meal-image-total", "meal-image-score", m.ScoreValue.Total, "Total"),
			MealScore(m.IDValue+"-meal-image-taste", "meal-image-score", m.ScoreValue.Taste, "Taste"),
			MealScore(m.IDValue+"-meal-image-recency", "meal-image-score", m.ScoreValue.Recency, "New"),
			MealScore(m.IDValue+"-meal-image-cost", "meal-image-score", m.ScoreValue.Cost, "Cost"),
			MealScore(m.IDValue+"-meal-image-effort", "meal-image-score", m.ScoreValue.Effort, "Effort"),
			MealScore(m.IDValue+"-meal-image-healthiness", "meal-image-score", m.ScoreValue.Healthiness, "Health"),
		),

		// app.Span().ID(m.id+"-meal-image-score-text").Class("meal-image-score-text").Text("Score: "+strconv.Itoa(m.score.Total)),
	)
}

// MealImage returns a new meal image component
func MealImage() *MealImageCompo {
	return &MealImageCompo{}
}

// ID sets the id of the meal image component
func (m *MealImageCompo) ID(id string) *MealImageCompo {
	m.IDValue = id
	return m
}

// Class sets the class of the meal image component
func (m *MealImageCompo) Class(class string) *MealImageCompo {
	m.ClassValue = class
	return m
}

// Img sets the image of the meal image component
func (m *MealImageCompo) Img(img string) *MealImageCompo {
	m.ImgValue = img
	return m
}

// MainText sets the main text of the meal image component
func (m *MealImageCompo) MainText(mainText string) *MealImageCompo {
	m.MainTextValue = mainText
	return m
}

// SecondaryText sets the secondary text of the meal image component
func (m *MealImageCompo) SecondaryText(secondaryText string) *MealImageCompo {
	m.SecondaryTextValue = secondaryText
	return m
}

// Score sets the score of the meal image component
func (m *MealImageCompo) Score(score osusu.Score) *MealImageCompo {
	m.ScoreValue = score
	return m
}

// OnClick sets the on click event for the meal image component
func (m *MealImageCompo) OnClick(onClick app.EventHandler) *MealImageCompo {
	m.OnClickValue = onClick
	return m
}

// OnClickScopeValue sets the on click event scope value for the meal image component.
// Use this to trigger updates to the meal image on click event when a given value changes (ex: a meal id)
func (m *MealImageCompo) OnClickScope(scope ...any) *MealImageCompo {
	m.OnClickScopeValue = scope
	return m
}
