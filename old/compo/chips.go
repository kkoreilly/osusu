package compo

import (
	"strconv"

	"github.com/kkoreilly/osusu/util/cond"
	"github.com/kkoreilly/osusu/util/list"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Chips is a component that contains multiple selectable chips. It can either contain checkbox or radio chips.
// If it contains checkbox chips, it must be of type map[string]bool, and if it has radio chips, it must be of type string.
type Chips[T string | map[string]bool] struct {
	app.Compo
	ID         string
	Class      string
	IsSelect   bool   // whether the chips are displayed as a select (dropdown) element
	Type       string // the type of the chips (radio or checkbox)
	Label      string
	Default    T
	Value      *T
	Options    []string
	OnChange   func(ctx app.Context, e app.Event, val string)
	Hidden     bool
	selectOpen bool
}

// Render returns the UI of the chips component
func (c *Chips[T]) Render() app.UI {
	// get formatted value text and value icon for value element only if we are going to display it
	var valueText, valueIcon string
	if c.IsSelect {
		if actualVal, ok := any(*c.Value).(string); ok {
			valueText = actualVal
		} else if actualVal, ok := any(*c.Value).(map[string]bool); ok {
			valueText = list.MapNum(actualVal, 3)
		}
		if c.selectOpen {
			valueIcon = "expand_less"
		} else {
			valueIcon = "expand_more"
		}
	}
	return app.Div().ID(c.ID+"-chips-outer-container").Class("chips-outer-container", c.Class, cond.If(c.IsSelect, "select")).DataSet("open", c.selectOpen).Hidden(c.Hidden).Body(
		app.Span().ID(c.ID+"-chips-label").Class("input-label").Text(c.Label),
		// current value of the chips, used in select
		&Button{ID: c.ID + "-chips-value", Class: "tertiary", Icon: valueIcon, Text: valueText, Hidden: !c.IsSelect, OnClick: c.ToggleSelect},
		app.Div().ID(c.ID+"-chips-container").Class("chips-container").Hidden(c.IsSelect && !c.selectOpen).Body(
			app.Range(c.Options).Slice(func(i int) app.UI {
				si := strconv.Itoa(i)
				optionVal := c.Options[i]
				var checked bool
				if actualVal, ok := any(*c.Value).(string); ok {
					checked = optionVal == actualVal
				} else if actualVal, ok := any(*c.Value).(map[string]bool); ok {
					checked = actualVal[optionVal]
				}
				return app.Label().ID(c.ID+"-chip-label-"+si).Class("chip-label").For(c.ID+"-chip-input-"+si).DataSet("checked", checked).Body(
					app.Input().ID(c.ID+"-chip-input-"+si).Class("chip-input").Type(c.Type).Name(c.ID).Checked(checked).OnChange(func(ctx app.Context, e app.Event) {
						// need to get val again to get updated value
						optionVal := c.Options[i]
						if val, ok := any(optionVal).(T); ok {
							if e.Get("target").Get("checked").Bool() {
								*c.Value = val
							}
						} else if actualVal, ok := any(*c.Value).(map[string]bool); ok {
							actualVal[optionVal] = e.Get("target").Get("checked").Bool()
						}
						if c.OnChange != nil {
							c.OnChange(ctx, e, optionVal)
						}
					}),
					app.Text(optionVal),
				)
			}),
		),
	)
}

// OnNav is called when the chips component is loaded. It loads the default value, if it set.
func (c *Chips[T]) OnInit() {
	if actualDefaultVal, ok := any(c.Default).(string); ok {
		// only use default if it is set and actual value is unset to prevent overriding existing info
		actualVal := any(*c.Value).(string)
		if actualDefaultVal != "" && actualVal == "" {
			*c.Value = c.Default
		}
	} else if actualDefaultVal, ok := any(c.Default).(map[string]bool); ok {
		// only use default if it is set and actual value is unset to prevent overriding existing info
		actualVal := any(*c.Value).(map[string]bool)
		if actualDefaultVal != nil && actualVal == nil {
			*c.Value = c.Default
		}
	}
	if c.IsSelect {
		CurrentPage.AddOnClick(func(ctx app.Context, e app.Event) {
			id := e.Get("target").Get("id").String()
			class := e.Get("target").Get("className").String()
			// never close if we just opened the select (clicked on the value)
			// no point in closing if select isn't even open
			// if type is radio, then always close because they are just selecting one option and that is standard behavior for dropdown
			// otherwise (if type is checkbox), keep open (return false) if they click anywhere inside the select, so they can change multiple options.
			if id != c.ID+"-chips-value-button" && id != c.ID+"-chips-value-button-icon" && id != c.ID+"-chips-value-button-text" && c.selectOpen && (c.Type == "radio" || (class != "chips-container" && class != "chip-label" && class != "chip-input")) {
				c.selectOpen = false
				c.Update()
			}
		})
	}
}

// ToggleSelect toggles whether the select is open
func (c *Chips[T]) ToggleSelect(ctx app.Context, e app.Event) {
	c.selectOpen = !c.selectOpen
}

// CheckboxChip is a component that has one chip that can be either selected or not
type CheckboxChip struct {
	app.Compo
	ID      string
	Label   string
	Default bool
	Value   *bool
}

// Render returns the UI of the CheckboxChip component
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
