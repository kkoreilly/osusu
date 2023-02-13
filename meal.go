package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Meal is a struct that represents the data of a meal
type Meal struct {
	ID          int
	Name        string
	Cost        int
	Effort      int
	Healthiness int
	Owner       int
}

// Meals is a slice that represents multiple meals
type Meals []Meal

// DefaultMeal makes and returns a meal with default values
func DefaultMeal() Meal {
	return Meal{
		Name:        "",
		Cost:        50,
		Effort:      50,
		Healthiness: 50,
	}
}

// SetCurrentMeal sets the current meal state value to the given meal, using the given context
func SetCurrentMeal(meal Meal, ctx app.Context) {
	ctx.SetState("currentMeal", meal, app.Persist)
}

// GetCurrentMeal gets and returns the current meal state value, using the given context
func GetCurrentMeal(ctx app.Context) Meal {
	var meal Meal
	ctx.GetState("currentMeal", &meal)
	return meal
}

// GetMealsRequest sends an HTTP request to the server to get the meals for the given user
func GetMealsRequest(user User) (Meals, error) {
	resp, err := http.Get("/api/getMeals?u=" + strconv.Itoa(user.ID))
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
	var meals Meals
	err = json.NewDecoder(resp.Body).Decode(&meals)
	if err != nil {
		return nil, err
	}
	return meals, nil
}

// CreateMealRequest sends an HTTP request to the server to create a meal with the given user as the owner and returns the created meal if successful and an error if not
func CreateMealRequest(user User) (Meal, error) {
	owner := user.ID
	resp, err := http.Post("/api/createMeal", "text/plain", bytes.NewBufferString(strconv.Itoa(owner)))
	if err != nil {
		return Meal{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return Meal{}, err
		}
		return Meal{}, fmt.Errorf("Error %s: %v", resp.Status, string(body))
	}
	// return gotten meal
	var res Meal
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return Meal{}, err
	}
	return res, nil
}

// UpdateMealRequest sends an HTTP request to the server to update the given meal
func UpdateMealRequest(meal Meal) error {
	jsonData, err := json.Marshal(meal)
	if err != nil {
		return err
	}
	resp, err := http.Post("/api/updateMeal", "application/json", bytes.NewBuffer(jsonData))
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

// DeleteMealRequest sends an HTTP request to the server to delete the given meal
func DeleteMealRequest(meal Meal) error {
	id := meal.ID
	req, err := http.NewRequest(http.MethodDelete, "/api/deleteMeal", bytes.NewBufferString(strconv.Itoa(id)))
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
