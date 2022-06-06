package main

import (
	"net/http"

	"github.com/stephenjlovell/go-micro/logger/data"

	helpers "github.com/stephenjlovell/json-helpers"
)

type LoggerPayload struct {
	Name string `json:name`
	Data string `json:data`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	requestPayload := &LoggerPayload{}
	_ = helpers.ReadJSON(w, r, requestPayload)
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}
	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	helpers.WriteJSON(w, http.StatusAccepted, helpers.JsonResponse{
		Error:   false,
		Message: "logged",
	})
}
