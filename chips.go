package main

import (
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// RadioChips is a component that has multiple chips, of which only one can be selected
type RadioChips struct {
	app.Compo
	ID      string
	Label   string
	Default string
	Value   *string
	Options []string
}

// Render returns the UI of the RadioChips component, which has multiple radio button chips
func (r *RadioChips) Render() app.UI {
	return app.Div().ID(r.ID+"-chips-outer-container").Class("chips-outer-container").Body(
		app.Label().ID(r.ID+"-chips-label").Class("input-label").For(r.ID+"-chips-container").Text(r.Label),
		app.Div().ID(r.ID+"-chips-container").Class("chips-container").Body(
			app.Range(r.Options).Slice(func(i int) app.UI {
				si := strconv.Itoa(i)
				val := r.Options[i]
				return app.Label().ID(r.ID+"-chip-label-"+si).Class("chip-label").For(r.ID+"-chip-input-"+si).DataSet("checked", val == *r.Value).Body(
					app.Input().ID(r.ID+"-chip-input-"+si).Class("chip-input").Type("radio").Name(r.ID).Checked(val == *r.Value).OnChange(func(ctx app.Context, e app.Event) {
						if e.Get("target").Get("checked").Bool() {
							*r.Value = val
						}
					}),
					app.Text(val),
				)
			}),
		),
	)
}

// OnInit is called when the component is loaded, and it sets the value to the default value
func (r *RadioChips) OnInit() {
	*r.Value = r.Default
}

// NewRadioChips makes a new RadioChips component with the given values
func NewRadioChips(id string, label string, defaultValue string, value *string, options ...string) *RadioChips {
	return &RadioChips{ID: id, Label: label, Default: defaultValue, Value: value, Options: options}
}

// CheckboxChips is a component that has multiple chips, of which any number can be selected
type CheckboxChips struct {
	app.Compo
	ID      string
	Label   string
	Default map[string]bool
	Value   *map[string]bool
	Options []string
}

// Render returns the UI of the CheckboxChips component, which has multiple checkbox chips
func (c *CheckboxChips) Render() app.UI {
	return app.Div().ID(c.ID+"-chips-outer-container").Class("chips-outer-container").Body(
		app.Label().ID(c.ID+"-chips-label").Class("input-label").For(c.ID+"-chips-container").Text(c.Label),
		app.Div().ID(c.ID+"-chips-container").Class("chips-container").Body(
			app.Range(c.Options).Slice(func(i int) app.UI {
				si := strconv.Itoa(i)
				val := c.Options[i]
				return app.Label().ID(c.ID+"-chip-label-"+si).Class("chip-label").For(c.ID+"-chip-input-"+si).DataSet("checked", (*c.Value)[val]).Body(
					app.Input().ID(c.ID+"-chip-input-"+si).Class("chip-input").Type("checkbox").Name(c.ID).Checked((*c.Value)[val]).OnChange(func(ctx app.Context, e app.Event) {
						(*c.Value)[val] = e.Get("target").Get("checked").Bool()
					}),
					app.Text(val),
				)
			}),
		),
	)
}

// OnInit is called when the component is loaded, and it sets the value to the default value
func (c *CheckboxChips) OnInit() {
	*c.Value = c.Default
}

// NewCheckboxChips makes a new CheckboxChips component with the given values
func NewCheckboxChips(id string, label string, defaultValue map[string]bool, value *map[string]bool, options ...string) *CheckboxChips {
	return &CheckboxChips{ID: id, Label: label, Default: defaultValue, Value: value, Options: options}
}

// CheckboxChip is a component that has one chip that can be either selected or not
type CheckboxChip struct {
	app.Compo
	ID      string
	Label   string
	Default bool
	Value   *bool
}

// Render returns the UI of the CheckboxChip component, which has one checkbox chip
func (c *CheckboxChip) Render() app.UI {
	return app.Label().ID(c.ID+"-chip-label").Class("chip-label").For(c.ID+"-chip-input").DataSet("checked", *c.Value).Body(
		app.Input().ID(c.ID+"-chip-input").Class("chip-input").Type("checkbox").Name(c.ID).Checked(*c.Value).OnChange(func(ctx app.Context, e app.Event) {
			*c.Value = e.Get("target").Get("checked").Bool()
		}),
		app.Text(c.Label),
	)
}

// OnInit is called when the component is loaded, and it sets the value to the default value
func (c *CheckboxChip) OnInit() {
	*c.Value = c.Default
}

// NewCheckboxChip makes a new CheckboxChip component with the given values
func NewCheckboxChip(id string, label string, defaultValue bool, value *bool) *CheckboxChip {
	return &CheckboxChip{ID: id, Label: label, Default: defaultValue, Value: value}
}
