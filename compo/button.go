package compo

import "github.com/maxence-charriere/go-app/v9/pkg/app"

// ButtonRow is a component that contains a row of multiple buttons
type ButtonRow struct {
	app.Compo
	ID      string
	Buttons []app.UI // Buttons are the buttons contained within the button row component
}

// Render returns the UI of the button row component
func (b *ButtonRow) Render() app.UI {
	return app.Div().ID(b.ID + "-button-row").Class("button-row").Body(
		b.Buttons...,
	)
}

// Button is a button component with text and an optional icon
type Button struct {
	app.Compo
	ID           string
	Class        string // Class is the CSS class of the button (ex: primary, secondary, tertiary)
	Type         string // Type is the HTML type of the button (ex: button, submit)
	Icon         string
	Text         string
	OnClick      app.EventHandler
	OnClickScope []any // OnClickScope is the on click event scope value that can be set to trigger updates to the click event when certain value(s) change
	Hidden       bool
}

// Render returns the UI of the button component
func (b *Button) Render() app.UI {
	if b.Type == "" {
		b.Type = "button"
	}
	return app.Button().ID(b.ID+"-button").Class(b.Class+"-button", "button").Type(b.Type).OnClick(b.OnClick, b.OnClickScope...).Hidden(b.Hidden).Body(
		app.Span().ID(b.ID+"-button-icon").Class(b.Class+"-button-icon", "button-icon", "material-symbols-outlined").Text(b.Icon).Hidden(b.Icon == ""),
		app.Span().ID(b.ID+"-button-text").Class(b.Class+"-button-text", "button-text").Text(b.Text).Hidden(b.Text == ""),
	)
}
