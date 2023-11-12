package compo

import (
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// MealScore is a component with a score pie circle containing score information for a meal, recipe, or entry
type MealScore struct {
	app.Compo
	ID    string
	Class string
	Score int
	Label string // the label for the score value (ex: cost, health)
}

func (m *MealScore) Render() app.UI {
	return app.Div().ID(m.ID).Class("meal-score", m.Class).Style("--score", strconv.Itoa(m.Score)).Style("--color-l", strconv.Itoa(m.Score/4+45)+"%").Body(
		app.Span().ID(m.ID+"-label").Class("meal-score-label", m.Class+"-label").Text(m.Label),
		app.Div().ID(m.ID+"-circle").Class("meal-score-circle", "pie", m.Class+"-circle").Text(m.Score),
	)
}
