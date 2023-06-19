package compo

import (
	"time"

	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Input is a component that includes an input field and an associated label
type Input[T any] struct {
	app.Compo
	ID                 string
	Class              string
	IsTextarea         bool   // whether the input is a text area instead of an input
	Type               string // the html type of the input (ex: text, password, range)
	Label              string
	Placeholder        string
	Value              *T
	ValueFunc          func(v app.Value) T // the (required) function used to convert the html value to the actual value
	DisplayFunc        func(T) any         // the (optional) function used to convert the actual value to the html value displayed
	AutoFocus          bool
	ButtonIcon         string           // if set to something other than "", an icon button will be rendered in the input with the given icon
	ButtonOnClick      app.EventHandler // the function called when the icon button created by ButtonIcon is clicked
	ButtonOnClickScope []any            // this should be set to cause the ButtonOnClick function to update when the given value(s) change
	showPassword       bool             // whether to show the password (only relevant for password inputs); used internally only -- do not set
}

// Render returns the UI of the input component, which includes a label and an input associated with it
func (i *Input[T]) Render() app.UI {
	var value any = *i.Value
	if i.DisplayFunc != nil {
		value = i.DisplayFunc(*i.Value)
	}
	inputType := i.Type
	buttonIcon := i.ButtonIcon
	// special case: if we are currently showing the password, we have to override the input type and the button icon to match
	if i.showPassword {
		inputType = "text"
		buttonIcon = "visibility_off"
	}
	var input app.UI = app.Input().ID(i.ID+"-input").Class("input", i.Class).Type(inputType).Placeholder(i.Placeholder).AutoFocus(i.AutoFocus).Value(value).OnChange(func(ctx app.Context, e app.Event) {
		*i.Value = i.ValueFunc(e.Get("target"))
	})
	if i.IsTextarea {
		input = app.Textarea().ID(i.ID+"-input").Class("input", i.Class).Placeholder(i.Placeholder).AutoFocus(i.AutoFocus).Text(value).OnChange(func(ctx app.Context, e app.Event) {
			*i.Value = i.ValueFunc(e.Get("target"))
		})
	}
	return app.Div().ID(i.ID+"-input-container").Class("input-container").DataSet("has-button", i.ButtonIcon != "").Body(
		app.Label().ID(i.ID+"-input-label").Class("input-label").For(i.ID+"-input").Text(i.Label),
		input,
		&Button{ID: i.ID + "-input-button", Class: "input", Icon: buttonIcon, OnClick: i.ButtonOnClick, OnClickScope: i.ButtonOnClickScope, Hidden: i.ButtonIcon == ""},
	)
}

// TextInput converts the given input component into a text input component
func TextInput(input *Input[string]) *Input[string] {
	input.Type = "text"
	input.ValueFunc = ValueFuncString
	return input
}

// PasswordInput converts the given input component into a password input component
func PasswordInput(input *Input[string]) *Input[string] {
	input.Type = "password"
	input.ValueFunc = ValueFuncString
	input.ButtonIcon = "visibility"
	input.ButtonOnClick = func(ctx app.Context, e app.Event) {
		input.showPassword = !input.showPassword
	}
	return input
}

// RangeInput converts the given input component into a range input component
func RangeInput(input *Input[int]) *Input[int] {
	input.Class = "input-range"
	input.Type = "range"
	input.ValueFunc = ValueFuncInt
	return input
}

// DateInput converts the given input component into a date input component
func DateInput(input *Input[time.Time]) *Input[time.Time] {
	input.Type = "date"
	input.ValueFunc = ValueFuncDate
	input.DisplayFunc = DisplayFuncDate
	return input
}

// TextareaInput converts the given input component into a textarea input component
func TextareaInput(input *Input[string]) *Input[string] {
	input.Class = "input-textarea"
	input.IsTextarea = true
	input.ValueFunc = ValueFuncString
	return input
}

// RangeInputUserMap converts the input component into a range input component that has its values associated with the entry in the user map corresponding to the given user
func RangeInputUserMap(input *Input[int], value *osusu.UserMap, user osusu.User) *Input[int] {
	val := (*value)[user.ID]
	input.Class = "input-range"
	input.Type = "range"
	input.Value = &val
	input.ValueFunc = func(v app.Value) int {
		res := v.Get("valueAsNumber").Int()
		(*value)[user.ID] = res
		return res
	}
	return input
}

// ValueFuncString is a basic value function for a string input
func ValueFuncString(v app.Value) string {
	return v.Get("value").String()
}

// ValueFuncInt is a basic value function for an int input
func ValueFuncInt(v app.Value) int {
	return v.Get("valueAsNumber").Int()
}

// ValueFuncDate is a basic value function for a date input
func ValueFuncDate(v app.Value) time.Time {
	return time.UnixMilli(int64(v.Get("valueAsNumber").Int())).UTC()
}

// DisplayFuncDate is a basic display function for a date input
func DisplayFuncDate(v time.Time) any {
	return v.Format("2006-01-02")
}
