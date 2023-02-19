package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Meal is a struct that represents the data of a meal
type Meal struct {
	ID          int
	Name        string
	Cost        int
	Effort      int
	Healthiness int
	Taste       map[int]int // key is person id, value is taste rating
	LastDone    time.Time
	UserID      int
}

// Meals is a slice that represents multiple meals
type Meals []Meal

// Score produces a score from 0 to 100 for the meal based on its attributes and the given options
func (m Meal) Score(options Options) int {
	// average of all attributes
	var tasteSum int
	for i, v := range m.Taste {
		use := options.People[i]
		// invert the person's rating if they are not participating
		if use {
			tasteSum += v
		} else {
			tasteSum += 100 - v
		}
	}
	recencyScore := int(2 * time.Now().Truncate(time.Hour*24).UTC().Sub(m.LastDone) / (time.Hour * 24))
	if recencyScore > 100 {
		recencyScore = 100
	}
	sum := options.CostWeight*(100-m.Cost) + options.EffortWeight*(100-m.Effort) + options.HealthinessWeight*m.Healthiness + options.TasteWeight*tasteSum + options.RecencyWeight*recencyScore
	den := options.CostWeight + options.EffortWeight + options.HealthinessWeight + len(m.Taste)*options.TasteWeight + options.RecencyWeight
	if den == 0 {
		return 0
	}
	return sum / den
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
	req, err := NewRequest(http.MethodGet, "/api/getMeals?u="+strconv.Itoa(user.ID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
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

// CreateMealRequest sends an HTTP request to the server to create a meal associated with the given user and returns the created meal if successful and an error if not
func CreateMealRequest(user User) (Meal, error) {
	userID := user.ID
	req, err := NewRequest(http.MethodPost, "/api/createMeal", bytes.NewBufferString(strconv.Itoa(userID)))
	if err != nil {
		return Meal{}, err
	}
	req.Header.Set("Content-Type", "text/plain")
	resp, err := http.DefaultClient.Do(req)
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
	req, err := NewRequest(http.MethodPost, "/api/updateMeal", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
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

// DeleteMealRequest sends an HTTP request to the server to delete the given meal
func DeleteMealRequest(meal Meal) error {
	id := meal.ID
	req, err := NewRequest(http.MethodDelete, "/api/deleteMeal", bytes.NewBufferString(strconv.Itoa(id)))
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
