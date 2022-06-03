package main

import (
	"authentication/data"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

const (
	loggerServiceUrl = "http://logger-service/log"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	user, err := app.doAuthenticate(w, r)
	// send response
	var outcome string
	if err != nil {
		outcome = "failed"
		app.errorJSON(w, err, http.StatusBadRequest)
	} else {
		outcome = "succeded"
		app.writeJSON(w, http.StatusAccepted, jsonResponse{
			Error:   false,
			Message: "Logged in",
			Data:    user,
		})
	}
	// log result of authentication attempt
	response, _ := app.doServiceRequest("POST", loggerServiceUrl, LoggerPayload{
		Name: fmt.Sprintf("login attempt %s", outcome),
		Data: "",
	})
	defer response.Body.Close()
}

func (app *Config) doAuthenticate(w http.ResponseWriter, r *http.Request) (*data.User, error) {
	payload := &AuthPayload{}
	err := app.readJSON(w, r, payload)
	if err != nil {
		return nil, err
	}
	user, err := app.Models.User.GetByEmail(payload.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	valid, err := user.PasswordMatches(payload.Password)
	if err != nil || !valid {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

// TODO: duplicated in logger-service
type LoggerPayload struct {
	Name string `json:name`
	Data string `json:data`
}

// TODO: duplicated in broker-service
func (app *Config) doServiceRequest(method, url string, data any) (*http.Response, error) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")
	// call the service
	request, err := http.NewRequest("POST", loggerServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	return client.Do(request)
}
