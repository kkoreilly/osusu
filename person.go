package main

// // Person represents the data of a person
// type Person struct {
// 	ID     int64
// 	UserID int64
// 	Name   string
// }
//
// type person struct {
// 	app.Compo
// 	person Person
// }

// func (p *person) Render() app.UI {
// 	return &Page{
// 		ID:                     "person",
// 		Title:                  "Person",
// 		Description:            "Confirm person",
// 		AuthenticationRequired: true,
// 		OnNavFunc: func(ctx app.Context) {
// 			p.person = GetCurrentPerson(ctx)
// 			app.Window().GetElementByID("person-page-name-input").Set("value", p.person.Name)
// 		},
// 		TitleElement: "Confirm Person",
// 		Elements: []app.UI{
// 			app.Form().ID("person-page-form").Class("form").OnSubmit(p.OnSubmit).Body(
// 				app.Label().ID("person-page-name-label").Class("input-label").For("person-page-name-input").Text("Name:"),
// 				app.Input().ID("person-page-name-input").Class("input").Type("text").Placeholder("Person Name").AutoFocus(true),
// 				app.Div().ID("person-page-action-button-row").Class("action-button-row").Body(
// 					app.Input().ID("person-page-delete-button").Class("action-button", "danger-action-button").Type("button").Value("Delete").OnClick(p.InitialDelete),
// 					app.Button().ID("person-page-back-button").Class("action-button", "secondary-action-button").Type("button").OnClick(NavigateEvent("/people")).Text("Back"),
// 					app.Input().ID("person-page-continue-button").Class("action-button", "primary-action-button").Type("submit").Value("Continue"),
// 				),
// 			),
// 			app.Dialog().ID("person-page-confirm-delete").Body(
// 				app.P().ID("person-page-confirm-delete-text").Class("confirm-delete-text").Text("Are you sure you want to delete this person?"),
// 				app.Div().ID("person-page-confirm-delete-action-button-row").Class("action-button-row").Body(
// 					app.Button().ID("person-page-confirm-delete-delete").Class("action-button", "danger-action-button").Text("Yes, Delete").OnClick(p.ConfirmDelete),
// 					app.Button().ID("person-page-confirm-delete-cancel").Class("action-button", "secondary-action-button").Text("No, Cancel").OnClick(p.CancelDelete),
// 				),
// 			),
// 		},
// 	}
// }

// func (p *person) OnSubmit(ctx app.Context, e app.Event) {
// 	e.PreventDefault()

// 	originalName := p.person.Name
// 	p.person.Name = app.Window().GetElementByID("person-page-name-input").Get("value").String()

// 	if p.person.Name != originalName {
// 		_, err := UpdatePersonAPI.Call(p.person)
// 		if err != nil {
// 			CurrentPage.ShowErrorStatus(err)
// 			return
// 		}
// 		SetCurrentPerson(p.person, ctx)
// 	}

// 	Navigate("/home", ctx)
// }

// func (p *person) InitialDelete(ctx app.Context, event app.Event) {
// 	event.PreventDefault()
// 	app.Window().GetElementByID("person-page-confirm-delete").Call("showModal")
// }

// func (p *person) ConfirmDelete(ctx app.Context, event app.Event) {
// 	event.PreventDefault()

// 	_, err := DeletePersonAPI.Call(p.person.ID)
// 	if err != nil {
// 		CurrentPage.ShowErrorStatus(err)
// 		return
// 	}
// 	SetCurrentPerson(Person{}, ctx)
// 	Navigate("/people", ctx)
// }

// func (p *person) CancelDelete(ctx app.Context, event app.Event) {
// 	event.PreventDefault()
// 	app.Window().GetElementByID("person-page-confirm-delete").Call("close")
// }
