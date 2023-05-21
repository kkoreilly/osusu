package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

// ButtonRowCompo is a component that contains a row of multiple buttons
type ButtonRowCompo struct {
	app.Compo
	id      string
	buttons []app.UI
}

// Render returns the UI of the button row component
func (b *ButtonRowCompo) Render() app.UI {
	return app.Div().ID(b.id + "-button-row").Class("button-row").Body(
		b.buttons...,
	)
}

// ButtonRow returns a new button row component
func ButtonRow() *ButtonRowCompo {
	return &ButtonRowCompo{}
}

// ID sets the ID of the button row component
func (b *ButtonRowCompo) ID(id string) *ButtonRowCompo {
	b.id = id
	return b
}

// Buttons sets the buttons in the button row component
func (b *ButtonRowCompo) Buttons(buttons ...app.UI) *ButtonRowCompo {
	b.buttons = buttons
	return b
}

// ButtonCompo is a button component with text and an optional icon
type ButtonCompo struct {
	app.Compo
	id           string
	ButtonClass  string // can sometimes change so needs to be exported
	typ          string
	DataSet      string
	ButtonIcon   string // can sometimes change so needs to be exported
	ButtonText   string // can sometimes change so needs to be exported
	onClick      app.EventHandler
	ButtonHidden bool // can sometimes change so needs to be exported
}

// Render returns the UI of the button component
func (b *ButtonCompo) Render() app.UI {
	if b.typ == "" {
		b.typ = "button"
	}
	return app.Button().ID(b.id+"-button").Class(b.ButtonClass+"-button", "button").Type(b.typ).OnClick(b.onClick).Hidden(b.ButtonHidden).Body(
		app.Span().ID(b.id+"-button-icon").Class(b.ButtonClass+"-button-icon", "button-icon", "material-symbols-outlined").Text(b.ButtonIcon).Hidden(b.ButtonIcon == ""),
		app.Span().ID(b.id+"-button-text").Class(b.ButtonClass+"-button-text", "button-text").Text(b.ButtonText).Hidden(b.ButtonText == ""),
	)
}

// Button returns a new button component
func Button() *ButtonCompo {
	return &ButtonCompo{}
}

// ID sets the ID of the button component to the given value
func (b *ButtonCompo) ID(id string) *ButtonCompo {
	b.id = id
	return b
}

// Class sets the class of the button component to the given value (ex: primary, secondary, tertiary)
func (b *ButtonCompo) Class(class string) *ButtonCompo {
	b.ButtonClass = class
	return b
}

// Type sets the button type of the button component to the given value (ex: button, submit).
// If this is not called, the default value is button.
func (b *ButtonCompo) Type(typ string) *ButtonCompo {
	b.typ = typ
	return b
}

// Icon sets the icon of the button component to the given value
func (b *ButtonCompo) Icon(icon string) *ButtonCompo {
	b.ButtonIcon = icon
	return b
}

// Text sets the text of the button component to the given value
func (b *ButtonCompo) Text(text string) *ButtonCompo {
	b.ButtonText = text
	return b
}

// OnClick sets the on click function of the button component
func (b *ButtonCompo) OnClick(h app.EventHandler) *ButtonCompo {
	b.onClick = h
	return b
}

// Hidden sets whether the button is hidden
func (b *ButtonCompo) Hidden(hidden bool) *ButtonCompo {
	b.ButtonHidden = hidden
	return b
}
