package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// People represents the data of multiple people
type People []Person

type people struct {
	app.Compo
	people People
}

func (p *people) Render() app.UI {
	return app.Div().Body(
		app.H1().ID("people-page-title").Class("page-title").Text("Who Are You?"),
		app.Div().ID("people-page-action-button-row").Class("action-button-row").Body(
			app.Button().ID("people-page-new-person").Class("action-button", "blue-action-button").Text("New Person").OnClick(p.New),
		),
		app.Div().ID("people-page-people-container").Body(
			app.Range(p.people).Slice(func(i int) app.UI {
				return app.Button().ID("people-page-person-" + strconv.Itoa(i)).Class("people-page-person").Text(p.people[i].Name).
					OnClick(func(ctx app.Context, e app.Event) { p.PersonOnClick(ctx, e, p.people[i]) })
			}),
		),
	)
}

func (p *people) OnNav(ctx app.Context) {
	people, err := GetPeopleRequest(GetCurrentUser(ctx))
	if err != nil {
		log.Println(err)
		return
	}
	p.people = people
}

func (p *people) New(ctx app.Context, e app.Event) {
	person, err := CreatePersonRequest(GetCurrentUser(ctx))
	if err != nil {
		log.Println(err)
		return
	}
	SetCurrentPerson(person, ctx)
	ctx.Navigate("/person")
}

func (p *people) PersonOnClick(ctx app.Context, e app.Event, person Person) {
	SetCurrentPerson(person, ctx)
	ctx.Navigate("/person")
}

// SetCurrentPerson sets the current person state value
func SetCurrentPerson(person Person, ctx app.Context) {
	ctx.SetState("currentPerson", person, app.Persist)
}

// GetCurrentPerson gets the current person state value
func GetCurrentPerson(ctx app.Context) Person {
	var person Person
	ctx.GetState("currentPerson", &person)
	return person
}

// // UpdatePeopleRequest sends a request to the server to update the people for the user
// func UpdatePeopleRequest(user User) error {
// 	jsonData, err := json.Marshal(user)
// 	if err != nil {
// 		return err
// 	}
// 	resp, err := http.Post("/api/updatePeople", "application/json", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != http.StatusOK {
// 		body, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			return err
// 		}
// 		return fmt.Errorf("Error %s: %v", resp.Status, string(body))
// 	}
// 	return nil
// }

// GetPeopleRequest sends an HTTP request to the server to get the people for the given user
func GetPeopleRequest(user User) (People, error) {
	resp, err := http.Get("/api/getPeople?u=" + strconv.Itoa(user.ID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Error %s: %v", resp.Status, string(body))
	}
	var people People
	err = json.NewDecoder(resp.Body).Decode(&people)
	if err != nil {
		return nil, err
	}
	return people, nil
}

// CreatePersonRequest sends an HTTP request to the server to create a person associated with the given user and returns the created person if successful and an error if not
func CreatePersonRequest(user User) (Person, error) {
	userID := user.ID
	resp, err := http.Post("/api/createPerson", "text/plain", bytes.NewBufferString(strconv.Itoa(userID)))
	if err != nil {
		return Person{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return Person{}, err
		}
		return Person{}, fmt.Errorf("Error %s: %v", resp.Status, string(body))
	}
	// return gotten person
	var res Person
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return Person{}, err
	}
	return res, nil
}

// UpdatePersonRequest sends an HTTP request to the server to update the given person
func UpdatePersonRequest(person Person) error {
	jsonData, err := json.Marshal(person)
	if err != nil {
		return err
	}
	resp, err := http.Post("/api/updatePerson", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Error %s: %v", resp.Status, string(body))
	}
	return nil
}

// DeletePersonRequest sends an HTTP request to the server to delete the given person
func DeletePersonRequest(person Person) error {
	id := person.ID
	req, err := http.NewRequest(http.MethodDelete, "/api/deletePerson", bytes.NewBufferString(strconv.Itoa(id)))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Error %s: %v", resp.Status, string(body))
	}
	return nil
}
