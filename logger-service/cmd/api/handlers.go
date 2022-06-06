package main

import (
	"net/http"

	"github.com/stephenjlovell/go-micro/logger/data"
)

type LoggerPayload struct {
	Name string `json:name`
	Data string `json:data`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	requestPayload := &LoggerPayload{}
	_ = app.readJSON(w, r, requestPayload)
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}
	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted, jsonResponse{
		Error:   false,
		Message: "logged",
	})
}
