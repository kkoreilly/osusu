// Package api contains the APIs used to communicate between the client and the server
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/kkoreilly/osusu/osusu"
	"github.com/kkoreilly/osusu/server"
)

// An API is a structural representation of a requestable API that accepts the given method on the given path and executes the given function.
// APIs are requested with the Call function.
type API[I, O any] struct {
	Method string
	Path   string
	Func   func(data I) (O, error)
}

// NewAPI creates and returns a new API with the given values
func NewAPI[I, O any](method string, path string, serverFunc func(data I) (O, error)) *API[I, O] {
	server.HandleFunc(method, path, func(w http.ResponseWriter, r *http.Request) {
		var data I
		if method == http.MethodGet {
			urlData := r.URL.Query().Get("data")
			err := json.Unmarshal([]byte(urlData), &data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		res, err := serverFunc(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(json)
	})
	return &API[I, O]{
		Method: method,
		Path:   path,
		Func:   serverFunc,
	}
}

// Call requests the API with the given input data and returns the result
func (a *API[I, O]) Call(data I) (O, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return *new(O), err
	}
	var req *http.Request
	if a.Method == http.MethodGet {
		req, err = osusu.NewRequest(a.Method, a.Path+"?data="+string(jsonData), nil)
		if err != nil {
			return *new(O), err
		}
	} else {
		req, err = osusu.NewRequest(a.Method, a.Path, bytes.NewBuffer(jsonData))
		if err != nil {
			return *new(O), err
		}
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return *new(O), err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return *new(O), err
		}
		return *new(O), errors.New(string(body))
	}
	var res O
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return *new(O), err
	}
	return res, nil
}
