package compo

import (
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// MealScore returns a table cell with a score pie circle containing score information for a meal or entry
func MealScore(id string, class string, score int, label string) app.UI {
	return app.Div().ID(id).Class("meal-score", class).Style("--score", strconv.Itoa(score)).Style("--color-l", strconv.Itoa(score/4+45)+"%").Body(
		app.Span().ID(id+"-label").Class("meal-score-label", class+"-label").Text(label),
		app.Div().ID(id+"-circle").Class("meal-score-circle", "pie", class+"-circle").Text(score),
	)
}
