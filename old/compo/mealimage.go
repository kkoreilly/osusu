package compo

import (
	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// A MealImage is a component for a meal, recipe, or entry with an image of the meal and information about it
type MealImage struct {
	app.Compo
	ID            string
	Class         string
	Selected      bool // whether the meal image is currently selected (ie: for a context menu)
	Image         string
	MainText      string
	SecondaryText string
	Score         osusu.Score
	OnClick       app.EventHandler
	OnClickScope  []any // the on click event scope value; use this to trigger updates to the meal image on click event when certain value(s) change
}

func (m *MealImage) Render() app.UI {
	return app.Div().ID(m.ID+"-meal-image-container").Class("meal-image-container", m.Class).DataSet("no-image", m.Image == "").DataSet("selected", m.Selected).OnClick(m.OnClick, m.OnClickScope...).Body(
		app.Img().ID(m.ID+"-meal-image").Class("meal-image").Src(m.Image),
		app.Div().ID(m.ID+"-meal-image-info-container").Class("meal-image-info-container").Body(
			app.Span().ID(m.ID+"-meal-image-main-text").Class("meal-image-main-text").Text(m.MainText),
			app.Span().ID(m.ID+"-meal-image-secondary-text").Class("meal-image-secondary-text").Text(m.SecondaryText),
			app.Div().ID(m.ID+"-meal-image-score-container").Class("meal-image-score-container").Body(
				&MealScore{ID: m.ID + "-meal-image-total", Class: "meal-image-score", Score: m.Score.Total, Label: "Total"},
				&MealScore{ID: m.ID + "-meal-image-taste", Class: "meal-image-score", Score: m.Score.Taste, Label: "Taste"},
				&MealScore{ID: m.ID + "-meal-image-recency", Class: "meal-image-score", Score: m.Score.Recency, Label: "New"},
				&MealScore{ID: m.ID + "-meal-image-cost", Class: "meal-image-score", Score: m.Score.Cost, Label: "Cost"},
				&MealScore{ID: m.ID + "-meal-image-effort", Class: "meal-image-score", Score: m.Score.Effort, Label: "Effort"},
				&MealScore{ID: m.ID + "-meal-image-healthiness", Class: "meal-image-score", Score: m.Score.Healthiness, Label: "Health"},
			),
		),
	)
}
