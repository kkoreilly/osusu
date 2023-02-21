package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

// Input is a component that includes an input field and an associated label
type Input[T any] struct {
	app.Compo
	ID          string
	Label       string
	InputClass  string
	Type        string
	Placeholder string
	AutoFocus   bool
	Value       *T
	ValueFunc   func(app.Value) T
}

// Render returns the UI of the input component, which includes a label and an input associated with it
func (i *Input[T]) Render() app.UI {
	return app.Div().ID(i.ID+"-input-container").Class("input-container").Body(
		app.Label().ID(i.ID+"-input-label").Class("input-label").For(i.ID+"-input").Text(i.Label),
		app.Input().ID(i.ID+"-input").Class("input", i.InputClass).Type(i.Type).Placeholder(i.Placeholder).AutoFocus(i.AutoFocus).Value(*i.Value).OnChange(func(ctx app.Context, e app.Event) {
			*i.Value = i.ValueFunc(e.Get("target"))
		}),
	)
}

// ValueFuncString is a basic value function for a string input
func ValueFuncString(v app.Value) string {
	return v.Get("value").String()
}

// ValueFuncInt is a basic value function for an int input
func ValueFuncInt(v app.Value) int {
	return v.Get("valueAsNumber").Int()
}

// NewTextInput makes a new text input component from the given values
func NewTextInput(id string, label string, placeholder string, autoFocus bool, value *string) *Input[string] {
	return &Input[string]{ID: id, Label: label, Type: "text", Placeholder: placeholder, AutoFocus: autoFocus, Value: value, ValueFunc: ValueFuncString}
}

// NewRangeInput makes a new range input component from the given values
func NewRangeInput(id string, label string, value *int) *Input[int] {
	return &Input[int]{ID: id, Label: label, InputClass: "input-range", Type: "range", Value: value, ValueFunc: ValueFuncInt}
}

// SetType sets the input type of the input to the given value
func (i *Input[T]) SetType(inputType string) *Input[T] {
	i.Type = inputType
	return i
}
