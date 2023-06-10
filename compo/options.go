package compo

import (
	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// OptionsCompo is a component that displays sorting and filtering options for meals, recipes, and entries
type OptionsCompo struct {
	app.Compo
	id      string
	class   string
	options *osusu.Options
	onSave  app.EventHandler
}

// Options returns a new options component that displays sorting and filtering options for meals, recipes, and entries
func Options() *OptionsCompo {
	return &OptionsCompo{}
}

// Render returns the UI of the options component
func (o *OptionsCompo) Render() app.UI {
	return app.Dialog().ID(o.id+"-options").Class("modal options", o.class).OnClick(o.onClick).Body(
		app.Div().ID(o.id + "-options-container").Class("options-container").OnClick(o.containerOnClick).Body(
			app.Form().ID(o.id+"-options-form").Class("form options-form").OnSubmit(o.SaveOptions).Body(
				RangeInput().ID(o.id+"-options-taste").Class("options-taste-input").Label("How important is taste?").Value(&o.options.TasteWeight),
				RangeInput().ID(o.id+"-options-recency").Class("options-recency-input").Label("How important is recency?").Value(&o.options.RecencyWeight),
				RangeInput().ID(o.id+"-options-cost").Class("options-cost-input").Label("How important is cost?").Value(&o.options.CostWeight),
				RangeInput().ID(o.id+"-options-effort").Class("options-effort-input").Label("How important is effort?").Value(&o.options.EffortWeight),
				RangeInput().ID(o.id+"-options-healthiness").Class("options-healthiness-input").Label("How important is healthiness?").Value(&o.options.HealthinessWeight),
			),
		),
	)
}

func (o *OptionsCompo) onClick(ctx app.Context, e app.Event) {
	// if the options dialog on click event is triggered, save and close the options because the dialog includes the whole page and a separate event will cancel this if they actually clicked on the dialog
	o.SaveOptions(ctx, e)
}

func (o *OptionsCompo) containerOnClick(ctx app.Context, e app.Event) {
	// cancel the closing of the dialog if they actually click on the dialog
	e.Call("stopPropagation")
}

// SaveOptions saves and closes the options
func (o *OptionsCompo) SaveOptions(ctx app.Context, e app.Event) {
	e.PreventDefault()

	// for _, u := range o.users {
	// 	o.options.Users[u.ID] = o.usersOptions[u.Name]
	// }

	osusu.SetOptions(*o.options, ctx)

	app.Window().GetElementByID(o.id + "-options").Call("close")

	o.onSave(ctx, e)
}

// ID sets the ID of the options component
func (o *OptionsCompo) ID(id string) *OptionsCompo {
	o.id = id
	return o
}

// Class sets the class of the options component
func (o *OptionsCompo) Class(class string) *OptionsCompo {
	o.class = class
	return o
}

// Options sets the actual options value of the options component that will be displayed and updated
func (o *OptionsCompo) Options(options *osusu.Options) *OptionsCompo {
	o.options = options
	return o
}

// OnSave sets the function to be called when the options are saved
func (o *OptionsCompo) OnSave(onSave app.EventHandler) *OptionsCompo {
	o.onSave = onSave
	return o
}
