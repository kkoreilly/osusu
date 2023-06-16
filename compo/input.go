package compo

import (
	"log"
	"time"

	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// InputCompo is a component that includes an input field and an associated label
type InputCompo[T any] struct {
	app.Compo
	id            string
	class         string
	isTextarea    bool   // whether the input is a text area instead of an input
	InputType     string // can sometimes change (with password field and show password) so needs to be exported
	label         string
	placeholder   string
	value         *T
	valueFunc     func(app.Value) T
	displayFunc   func(T) any
	autoFocus     bool
	ButtonIconVal string // can sometimes change so needs to be exported
	buttonOnClick app.EventHandler
}

// Render returns the UI of the input component, which includes a label and an input associated with it
func (i *InputCompo[T]) Render() app.UI {
	var value any = *i.value
	if i.displayFunc != nil {
		value = i.displayFunc(*i.value)
	}
	var input app.UI = app.Input().ID(i.id+"-input").Class("input", i.class).Type(i.InputType).Placeholder(i.placeholder).AutoFocus(i.autoFocus).Value(value).OnChange(func(ctx app.Context, e app.Event) {
		*i.value = i.valueFunc(e.Get("target"))
		log.Println(*i.value)
	})
	if i.isTextarea {
		input = app.Textarea().ID(i.id+"-input").Class("input", i.class).Placeholder(i.placeholder).AutoFocus(i.autoFocus).Text(value).OnChange(func(ctx app.Context, e app.Event) {
			*i.value = i.valueFunc(e.Get("target"))
		})
	}

	return app.Div().ID(i.id+"-input-container").Class("input-container").DataSet("has-button", i.ButtonIconVal != "").Body(
		app.Label().ID(i.id+"-input-label").Class("input-label").For(i.id+"-input").Text(i.label),
		input,
		Button().ID(i.id+"-input-button").Class("input").Icon(i.ButtonIconVal).OnClick(i.buttonOnClick).Hidden(i.ButtonIconVal == ""),
	)
}

// Input returns a new input component
func Input[T any]() *InputCompo[T] {
	return &InputCompo[T]{}
}

// TextInput returns a new string input component with type text and value func ValueFuncString
func TextInput() *InputCompo[string] {
	return Input[string]().Type("text").ValueFunc(ValueFuncString)
}

// RangeInput returns a new range input component with type range and value func ValueFuncInt
func RangeInput() *InputCompo[int] {
	return Input[int]().Class("input-range").Type("range").ValueFunc(ValueFuncInt)
}

// DateInput returns a new date input component
func DateInput() *InputCompo[time.Time] {
	return Input[time.Time]().Type("date").ValueFunc(ValueFuncDate).DisplayFunc(DisplayFuncDate)
}

// Textarea returns a new textarea input component
func Textarea() *InputCompo[string] {
	return Input[string]().Class("input-textarea").IsTextarea(true).ValueFunc(ValueFuncString)
}

// RangeInputUserMap returns a new range input component that has its values associated with the entry in the user map corresponding to the given user
func RangeInputUserMap(value *osusu.UserMap, user osusu.User) *InputCompo[int] {
	val := (*value)[user.ID]
	return Input[int]().Class("input-range").Type("range").Value(&val).ValueFunc(func(v app.Value) int {
		res := v.Get("valueAsNumber").Int()
		(*value)[user.ID] = res
		return res
	})
}

// ID sets the ID of the input component to the given value
func (i *InputCompo[T]) ID(id string) *InputCompo[T] {
	i.id = id
	return i
}

// Class adds the given value to the class of the input component
func (i *InputCompo[T]) Class(class string) *InputCompo[T] {
	i.class += class + " "
	return i
}

// IsTextarea sets whether the input component is a textarea. The default value if this function is not called is false.
func (i *InputCompo[T]) IsTextarea(isTextarea bool) *InputCompo[T] {
	i.isTextarea = isTextarea
	return i
}

// Type sets the input type of the input (ex: text, password, range)
func (i *InputCompo[T]) Type(typ string) *InputCompo[T] {
	i.InputType = typ
	return i
}

// Label sets the label of the input
func (i *InputCompo[T]) Label(label string) *InputCompo[T] {
	i.label = label
	return i
}

// Placeholder sets the placeholder of the input
func (i *InputCompo[T]) Placeholder(placeholder string) *InputCompo[T] {
	i.placeholder = placeholder
	return i
}

// Value sets the value of the input component to stay equal with the given pointer
func (i *InputCompo[T]) Value(value *T) *InputCompo[T] {
	i.value = value
	return i
}

// ValueFunc sets the function used to convert the value of the input to a usable value
func (i *InputCompo[T]) ValueFunc(valueFunc func(app.Value) T) *InputCompo[T] {
	i.valueFunc = valueFunc
	return i
}

// DisplayFunc sets the function used to convert the set value to the value actually displayed in the input
func (i *InputCompo[T]) DisplayFunc(displayFunc func(T) any) *InputCompo[T] {
	i.displayFunc = displayFunc
	return i
}

// AutoFocus sets whether to automatically focus the input on page load
func (i *InputCompo[T]) AutoFocus(autoFocus bool) *InputCompo[T] {
	i.autoFocus = autoFocus
	return i
}

// ButtonIcon, if called with a value other than "", causes an icon button to be rendered inside of the input with the given button icon.
// ButtonOnClick should be called to set the function called when the icon button is clicked.
func (i *InputCompo[T]) ButtonIcon(buttonIcon string) *InputCompo[T] {
	i.ButtonIconVal = buttonIcon
	return i
}

// ButtonOnClick sets the function that is called when the icon button specified with ButtonIcon is clicked on.
func (i *InputCompo[T]) ButtonOnClick(buttonOnClick app.EventHandler) *InputCompo[T] {
	i.buttonOnClick = buttonOnClick
	return i
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
