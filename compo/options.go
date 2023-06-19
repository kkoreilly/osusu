package compo

import (
	"github.com/kkoreilly/osusu/api"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Options is a pop-up dialog component that displays sorting and filtering options for meals, recipes, and entries
type Options struct {
	app.Compo
	ID         string
	Class      string
	Options    *osusu.Options
	Exclude    []string        // options that will not be displayed
	excludeMap map[string]bool // a map version of exclude with every value equal to true
	OnSave     app.EventHandler
}

// Render returns the UI of the options component
func (o *Options) Render() app.UI {
	return app.Dialog().ID(o.ID+"-options").Class("modal options", o.Class).OnClick(o.onClick).Body(
		app.Div().ID(o.ID + "-options-container").Class("options-container").OnClick(o.containerOnClick).Body(
			app.Form().ID(o.ID+"-options-form").Class("form options-form").OnSubmit(o.SaveOptions).Body(
				app.If(!o.excludeMap["taste"], RangeInput(&Input[int]{ID: o.ID + "-options-taste", Class: "options-taste-input", Label: "How important is taste?", Value: &o.Options.TasteWeight})),
				app.If(!o.excludeMap["recency"], RangeInput(&Input[int]{ID: o.ID + "-options-recency", Class: "options-recency-input", Label: "How important is recency?", Value: &o.Options.RecencyWeight})),
				app.If(!o.excludeMap["cost"], RangeInput(&Input[int]{ID: o.ID + "-options-cost", Class: "options-cost-input", Label: "How important is cost?", Value: &o.Options.CostWeight})),
				app.If(!o.excludeMap["effort"], RangeInput(&Input[int]{ID: o.ID + "-options-effort", Class: "options-effort-input", Label: "How important is effort?", Value: &o.Options.EffortWeight})),
				app.If(!o.excludeMap["healthiness"], RangeInput(&Input[int]{ID: o.ID + "-options-healthiness", Class: "options-healthiness-input", Label: "How important is healthiness?", Value: &o.Options.HealthinessWeight})),
			),
		),
	)
}

func (o *Options) onClick(ctx app.Context, e app.Event) {
	// if the options dialog on click event is triggered, save and close the options because the dialog includes the whole page and a separate event will cancel this if they actually clicked on the dialog
	o.SaveOptions(ctx, e)
}

func (o *Options) containerOnClick(ctx app.Context, e app.Event) {
	// cancel the closing of the dialog if they actually click on the dialog
	e.Call("stopPropagation")
}

// SaveOptions saves the options to local storage and then closes the dialog
func (o *Options) SaveOptions(ctx app.Context, e app.Event) {
	osusu.SetOptions(*o.Options, ctx)

	app.Window().GetElementByID(o.ID + "-options").Call("close")

	if o.OnSave != nil {
		o.OnSave(ctx, e)
	}
}

// OnInit is called when the options component is loaded
func (o *Options) OnInit() {
	o.excludeMap = map[string]bool{}
	if o.Exclude != nil {
		for _, option := range o.Exclude {
			o.excludeMap[option] = true
		}
	}
}

// QuickOptions is a component that displays quick dropdown sorting and filtering options for meals, recipes, and entries
type QuickOptions struct {
	app.Compo
	ID           string
	Options      *osusu.Options
	Exclude      []string        // options that will not be displayed
	excludeMap   map[string]bool // a map version of exclude with every value equal to true
	Group        osusu.Group
	Meals        osusu.Meals // these are used to determine the cuisine options that are displayed; to have all base cuisine options displayed, do not set this, or set this to nil
	OnSave       app.EventHandler
	users        osusu.Users
	usersOptions map[string]bool
	usersStrings []string
	cuisines     []string
}

// Render returns the UI of the quick options component
func (q *QuickOptions) Render() app.UI {
	return &ButtonRow{ID: q.ID + "-quick-options", Buttons: []app.UI{
		&Chips[map[string]bool]{ID: q.ID + "-quick-options-category", IsSelect: true, Type: "checkbox", Label: "Categories:", Default: map[string]bool{"Dinner": true}, Value: &q.Options.Category, Options: append(osusu.AllCategories, "Unset"), OnChange: q.SaveOptions, Hidden: q.excludeMap["category"]},
		&Chips[map[string]bool]{ID: q.ID + "-quick-options-users", IsSelect: true, Type: "checkbox", Label: "People:", Value: &q.usersOptions, Options: q.usersStrings, OnChange: q.SaveOptions, Hidden: q.excludeMap["users"]},
		&Chips[map[string]bool]{ID: q.ID + "-quick-options-source", IsSelect: true, Type: "checkbox", Label: "Sources:", Default: map[string]bool{"Cooking": true, "Dine-In": true, "Takeout": true}, Value: &q.Options.Source, Options: osusu.AllSources, OnChange: q.SaveOptions, Hidden: q.excludeMap["source"]},
		&Chips[map[string]bool]{ID: q.ID + "-quick-options-cuisine", IsSelect: true, Type: "checkbox", Label: "Cuisines:", Value: &q.Options.Cuisine, Options: q.cuisines, OnChange: q.SaveOptions, Hidden: q.excludeMap["cuisine"]},
	}}
}

// SaveOptions saves the quick options to local storage
func (q *QuickOptions) SaveOptions(ctx app.Context, e app.Event, val string) {
	for _, u := range q.users {
		q.Options.Users[u.ID] = q.usersOptions[u.Name]
	}

	osusu.SetOptions(*q.Options, ctx)

	if q.OnSave != nil {
		q.OnSave(ctx, e)
	}
}

// OnInit is called when the quick options component is loaded
func (q *QuickOptions) OnInit() {
	q.excludeMap = map[string]bool{}
	if q.Exclude != nil {
		for _, option := range q.Exclude {
			q.excludeMap[option] = true
		}
	}

	users, err := api.GetUsers.Call(q.Group.Members)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	q.users = users

	q.usersOptions = make(map[string]bool)
	for _, p := range q.users {
		if _, ok := q.Options.Users[p.ID]; !ok {
			q.Options.Users[p.ID] = true
		}
		q.usersOptions[p.Name] = q.Options.Users[p.ID]
	}

	q.usersStrings = []string{}
	for _, u := range q.users {
		q.usersStrings = append(q.usersStrings, u.Name)
	}

	// if we have no meals to draw cuisines from (ie: we are probably in discover mode), then just use base cuisines
	if q.Meals == nil {
		q.cuisines = osusu.BaseCuisines
	} else {
		cuisinesInUse := map[string]bool{}
		for _, meal := range q.Meals {
			for _, cuisine := range meal.Cuisine {
				cuisinesInUse[cuisine] = true
				// if the user has not yet set whether or not to allow this cuisine (if it is new), automatically set it to true
				_, ok := q.Options.Cuisine[cuisine]
				if !ok {
					q.Options.Cuisine[cuisine] = true
				}
			}
		}
		for cuisine := range q.Options.Cuisine {
			if !cuisinesInUse[cuisine] {
				delete(q.Options.Cuisine, cuisine)
			}
		}

		q.cuisines = []string{}
		for cuisine, val := range cuisinesInUse {
			if val {
				q.cuisines = append(q.cuisines, cuisine)
			}
		}
	}
}
