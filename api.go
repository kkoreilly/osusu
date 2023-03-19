package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type API[I, O any] struct {
	Method string
	Path   string
	Func   func(data I) (O, error)
}

func NewAPI[I, O any](method string, path string, serverFunc func(data I) (O, error)) *API[I, O] {
	HandleFunc(method, path, func(w http.ResponseWriter, r *http.Request) {
		var data I
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
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

func (a *API[I, O]) Call(data I) (O, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return *new(O), err
	}
	req, err := NewRequest(a.Method, a.Path, bytes.NewBuffer(jsonData))
	if err != nil {
		return *new(O), err
	}
	req.Header.Set("Content-Type", "application/json")
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
