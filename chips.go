package main

import (
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// ChipsCompo is a component that contains multiple selectable chips. It can either contain checkbox or radio chips.
// If it contains checkbox chips, it will be of type map[string]bool, and if it has radio chips, it will be of type string.
type ChipsCompo[T string | map[string]bool] struct {
	app.Compo
	id         string
	class      string
	typ        string
	label      string
	defaultVal T
	value      *T
	OptionsVal []string // can change so needs to be exported
	onChange   func(ctx app.Context, e app.Event, val string)
}

// Render returns the UI of the chips component
func (c *ChipsCompo[T]) Render() app.UI {
	return app.Div().ID(c.id+"-chips-outer-container").Class("chips-outer-container", c.class).Body(
		app.Label().ID(c.id+"-chips-label").Class("input-label").Text(c.label),
		app.Div().ID(c.id+"-chips-container").Class("chips-container").Body(
			app.Range(c.OptionsVal).Slice(func(i int) app.UI {
				si := strconv.Itoa(i)
				optionVal := c.OptionsVal[i]
				var checked bool
				if actualVal, ok := any(*c.value).(string); ok {
					checked = optionVal == actualVal
				} else if actualVal, ok := any(*c.value).(map[string]bool); ok {
					checked = actualVal[optionVal]
				}
				return app.Label().ID(c.id+"-chip-label-"+si).Class("chip-label").For(c.id+"-chip-input-"+si).DataSet("checked", checked).Body(
					app.Input().ID(c.id+"-chip-input-"+si).Class("chip-input").Type(c.typ).Name(c.id).Checked(checked).OnChange(func(ctx app.Context, e app.Event) {
						// need to get val again to get updated value
						optionVal := c.OptionsVal[i]
						if val, ok := any(optionVal).(T); ok {
							if e.Get("target").Get("checked").Bool() {
								*c.value = val
							}
						} else if actualVal, ok := any(*c.value).(map[string]bool); ok {
							actualVal[optionVal] = e.Get("target").Get("checked").Bool()
						}
						if c.onChange != nil {
							c.onChange(ctx, e, optionVal)
						}
					}),
					app.Text(optionVal),
				)
			}),
		),
	)
}

// OnNav is called when the chips component is loaded. It loads the default value, if it set.
func (c *ChipsCompo[T]) OnNav(ctx app.Context) {
	if val, ok := any(c.defaultVal).(string); ok {
		if val != "" {
			c.value = &c.defaultVal
		}
	} else if val, ok := any(c.defaultVal).(map[string]bool); ok {
		if val != nil {
			c.value = &c.defaultVal
		}
	}
}

// Chips returns a new chips component
func Chips[T string | map[string]bool]() *ChipsCompo[T] {
	return &ChipsCompo[T]{}

}

// RadioChips returns a new radio chips component
func RadioChips() *ChipsCompo[string] {
	return Chips[string]().Type("radio")
}

// CheckboxChips returns a new checkbox chips component
func CheckboxChips() *ChipsCompo[map[string]bool] {
	return Chips[map[string]bool]().Type("checkbox")
}

// ID sets the id of the chips component to the given value
func (c *ChipsCompo[T]) ID(id string) *ChipsCompo[T] {
	c.id = id
	return c
}

// Class sets the class of the chips component to the given value
func (c *ChipsCompo[T]) Class(class string) *ChipsCompo[T] {
	c.class = class
	return c
}

// Type sets the type of the chips component to the given value
func (c *ChipsCompo[T]) Type(typ string) *ChipsCompo[T] {
	c.typ = typ
	return c
}

// Label sets the label of the chips component to the given value
func (c *ChipsCompo[T]) Label(label string) *ChipsCompo[T] {
	c.label = label
	return c
}

// Default sets the default value of the chips component to the given value
func (c *ChipsCompo[T]) Default(defaultVal T) *ChipsCompo[T] {
	c.defaultVal = defaultVal
	return c
}

// Value sets the value of the chips component to the given value
func (c *ChipsCompo[T]) Value(value *T) *ChipsCompo[T] {
	c.value = value
	return c
}

// Options sets the options for the chips component to the given value
func (c *ChipsCompo[T]) Options(options ...string) *ChipsCompo[T] {
	c.OptionsVal = options
	return c
}

func (c *ChipsCompo[T]) OnChange(onChange func(ctx app.Context, e app.Event, val string)) *ChipsCompo[T] {
	c.onChange = onChange
	return c
}

// CheckboxChipCompo is a component that has one chip that can be either selected or not
type CheckboxChipCompo struct {
	app.Compo
	id         string
	label      string
	defaultVal bool
	value      *bool
}

// Render returns the UI of the CheckboxChip component
func (c *CheckboxChipCompo) Render() app.UI {
	return app.Label().ID(c.id+"-chip-label").Class("chip-label").For(c.id+"-chip-input").DataSet("checked", *c.value).Body(
		app.Input().ID(c.id+"-chip-input").Class("chip-input").Type("checkbox").Name(c.id).Checked(*c.value).OnChange(func(ctx app.Context, e app.Event) {
			*c.value = e.Get("target").Get("checked").Bool()
		}),
		app.Text(c.label),
	)
}

// OnInit is called when the component is loaded, and it sets the value to the default value
func (c *CheckboxChipCompo) OnInit() {
	*c.value = c.defaultVal
}

// CheckboxChip returns a new checkbox chip component
func CheckboxChip() *CheckboxChipCompo {
	return &CheckboxChipCompo{}
}

// ID sets the ID of the checkbox chip component
func (c *CheckboxChipCompo) ID(id string) *CheckboxChipCompo {
	c.id = id
	return c
}

// Label sets the label of the checkbox chip component
func (c *CheckboxChipCompo) Label(label string) *CheckboxChipCompo {
	c.label = label
	return c
}

// Default sets the default value of the checkbox chip component
func (c *CheckboxChipCompo) Default(defaultVal bool) *CheckboxChipCompo {
	c.defaultVal = defaultVal
	return c
}

// Value sets the value of the checkbox chip component
func (c *CheckboxChipCompo) Value(value *bool) *CheckboxChipCompo {
	c.value = value
	return c
}
