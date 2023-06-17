package compo

import (
	"github.com/kkoreilly/osusu/api"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// OptionsCompo is a pop-up dialog component that displays sorting and filtering options for meals, recipes, and entries
type OptionsCompo struct {
	app.Compo
	id         string
	class      string
	options    *osusu.Options
	exclude    []string
	excludeMap map[string]bool // a map version of exclude with every value equal to true
	onSave     app.EventHandler
}

// Options returns a new options pop-dialog component that displays sorting and filtering options for meals, recipes, and entries
func Options() *OptionsCompo {
	return &OptionsCompo{}
}

// Render returns the UI of the options component
func (o *OptionsCompo) Render() app.UI {
	return app.Dialog().ID(o.id+"-options").Class("modal options", o.class).OnClick(o.onClick).Body(
		app.Div().ID(o.id + "-options-container").Class("options-container").OnClick(o.containerOnClick).Body(
			app.Form().ID(o.id+"-options-form").Class("form options-form").OnSubmit(o.SaveOptions).Body(
				app.If(!o.excludeMap["taste"], RangeInput().ID(o.id+"-options-taste").Class("options-taste-input").Label("How important is taste?").Value(&o.options.TasteWeight)),
				app.If(!o.excludeMap["recency"], RangeInput().ID(o.id+"-options-recency").Class("options-recency-input").Label("How important is recency?").Value(&o.options.RecencyWeight)),
				app.If(!o.excludeMap["cost"], RangeInput().ID(o.id+"-options-cost").Class("options-cost-input").Label("How important is cost?").Value(&o.options.CostWeight)),
				app.If(!o.excludeMap["effort"], RangeInput().ID(o.id+"-options-effort").Class("options-effort-input").Label("How important is effort?").Value(&o.options.EffortWeight)),
				app.If(!o.excludeMap["healthiness"], RangeInput().ID(o.id+"-options-healthiness").Class("options-healthiness-input").Label("How important is healthiness?").Value(&o.options.HealthinessWeight)),
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

// SaveOptions saves the options to local storage and then closes the dialog
func (o *OptionsCompo) SaveOptions(ctx app.Context, e app.Event) {
	osusu.SetOptions(*o.options, ctx)

	app.Window().GetElementByID(o.id + "-options").Call("close")

	if o.onSave != nil {
		o.onSave(ctx, e)
	}
}

// OnInit is called when the options component is loaded
func (o *OptionsCompo) OnInit() {
	o.excludeMap = map[string]bool{}
	if o.exclude != nil {
		for _, option := range o.exclude {
			o.excludeMap[option] = true
		}
	}
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

// Exclude excludes the given options from the options component such that they are not displayed
func (o *OptionsCompo) Exclude(exclude ...string) *OptionsCompo {
	o.exclude = exclude
	return o
}

// OnSave sets the function to be called when the options are saved
func (o *OptionsCompo) OnSave(onSave app.EventHandler) *OptionsCompo {
	o.onSave = onSave
	return o
}

// QuickOptionsCompo is a component that displays quick dropdown sorting and filtering options for meals, recipes, and entries
type QuickOptionsCompo struct {
	app.Compo
	id           string
	options      *osusu.Options
	exclude      []string
	excludeMap   map[string]bool // a map version of exclude with every value equal to true
	group        osusu.Group
	meals        osusu.Meals
	onSave       app.EventHandler
	users        osusu.Users
	usersOptions map[string]bool
	usersStrings []string
	cuisines     []string
}

// QuickOptions returns a new component that displays quick dropdown sorting and filtering options for meals, recipes, and entries
func QuickOptions() *QuickOptionsCompo {
	return &QuickOptionsCompo{}
}

// Render returns the UI of the quick options component
func (q *QuickOptionsCompo) Render() app.UI {
	return &ButtonRow{ID: q.id + "-quick-options", Buttons: []app.UI{
		CheckboxSelect().ID(q.id + "-quick-options-category").Label("Categories:").Default(map[string]bool{"Dinner": true}).Value(&q.options.Category).Options(append(osusu.AllCategories, "Unset")...).OnChange(q.SaveOptions).Hidden(q.excludeMap["category"]),
		CheckboxSelect().ID(q.id + "-quick-options-users").Label("People:").Value(&q.usersOptions).Options(q.usersStrings...).OnChange(q.SaveOptions).Hidden(q.excludeMap["users"]),
		CheckboxSelect().ID(q.id + "-quick-options-source").Label("Sources:").Default(map[string]bool{"Cooking": true, "Dine-In": true, "Takeout": true}).Value(&q.options.Source).Options(osusu.AllSources...).OnChange(q.SaveOptions).Hidden(q.excludeMap["source"]),
		CheckboxSelect().ID(q.id + "-quick-options-cuisine").Label("Cuisines:").Value(&q.options.Cuisine).Options(q.cuisines...).OnChange(q.SaveOptions).Hidden(q.excludeMap["cuisine"]),
	}}
}

// SaveOptions saves the quick options to local storage
func (q *QuickOptionsCompo) SaveOptions(ctx app.Context, e app.Event, val string) {
	for _, u := range q.users {
		q.options.Users[u.ID] = q.usersOptions[u.Name]
	}

	osusu.SetOptions(*q.options, ctx)

	if q.onSave != nil {
		q.onSave(ctx, e)
	}
}

// OnInit is called when the quick options component is loaded
func (q *QuickOptionsCompo) OnInit() {
	q.excludeMap = map[string]bool{}
	if q.exclude != nil {
		for _, option := range q.exclude {
			q.excludeMap[option] = true
		}
	}

	users, err := api.GetUsers.Call(q.group.Members)
	if err != nil {
		CurrentPage.ShowErrorStatus(err)
		return
	}
	q.users = users

	q.usersOptions = make(map[string]bool)
	for _, p := range q.users {
		if _, ok := q.options.Users[p.ID]; !ok {
			q.options.Users[p.ID] = true
		}
		q.usersOptions[p.Name] = q.options.Users[p.ID]
	}

	q.usersStrings = []string{}
	for _, u := range q.users {
		q.usersStrings = append(q.usersStrings, u.Name)
	}

	// if we have no meals to draw cuisines from (ie: we are probably in discover mode), then just use base cuisines
	if q.meals == nil {
		q.cuisines = osusu.BaseCuisines
	} else {
		cuisinesInUse := map[string]bool{}
		for _, meal := range q.meals {
			for _, cuisine := range meal.Cuisine {
				cuisinesInUse[cuisine] = true
				// if the user has not yet set whether or not to allow this cuisine (if it is new), automatically set it to true
				_, ok := q.options.Cuisine[cuisine]
				if !ok {
					q.options.Cuisine[cuisine] = true
				}
			}
		}
		for cuisine := range q.options.Cuisine {
			if !cuisinesInUse[cuisine] {
				delete(q.options.Cuisine, cuisine)
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

// ID sets the id of the quick options component
func (q *QuickOptionsCompo) ID(id string) *QuickOptionsCompo {
	q.id = id
	return q
}

// Options sets the actual options value of the quick options component
func (q *QuickOptionsCompo) Options(options *osusu.Options) *QuickOptionsCompo {
	q.options = options
	return q
}

// Exclude excludes the given options from the quick options component such that they are not displayed
func (q *QuickOptionsCompo) Exclude(exclude ...string) *QuickOptionsCompo {
	q.exclude = exclude
	return q
}

// Group sets the group that the person viewing the quick options component is currently in
func (q *QuickOptionsCompo) Group(group osusu.Group) *QuickOptionsCompo {
	q.group = group
	return q
}

// Meals sets the meals that the quick options component controls. This affects the cuisine options that are displayed.
// To have all base cuisine options displayed, do not call this function, or call it with nil.
func (q *QuickOptionsCompo) Meals(meals osusu.Meals) *QuickOptionsCompo {
	q.meals = meals
	return q
}

// OnSave sets the function be called when the quick options are saved
func (q *QuickOptionsCompo) OnSave(onSave app.EventHandler) *QuickOptionsCompo {
	q.onSave = onSave
	return q
}
