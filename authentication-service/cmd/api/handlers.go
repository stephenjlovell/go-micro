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

// TODO: duplicated in logger-service, broker-service
type LoggerPayload struct {
	Name string `json:name`
	Data string `json:data`
}

const (
	loggerServiceUrl = "http://logger-service/log"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	authPayload, err := app.readAuthPayload(w, r)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	user, err := app.doAuthenticate(w, r, authPayload)
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
	_ = app.doLogging(outcome, authPayload.Email)
}

// log result of authentication attempt
func (app *Config) doLogging(outcome, data string) (err error) {
	response, err := app.doServiceRequest("POST", loggerServiceUrl, LoggerPayload{
		Name: fmt.Sprintf("login attempt %s", outcome),
		Data: data,
	})
	defer func() {
		closeErr := response.Body.Close()
		if err == nil {
			err = closeErr
		}
	}()
	return
}

func (app *Config) readAuthPayload(w http.ResponseWriter, r *http.Request) (*AuthPayload, error) {
	payload := &AuthPayload{}
	err := app.readJSON(w, r, payload)
	return payload, err
}

func (app *Config) doAuthenticate(w http.ResponseWriter, r *http.Request, payload *AuthPayload) (*data.User, error) {
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

// TODO: duplicated in broker-service
func (app *Config) doServiceRequest(method, url string, data any) (*http.Response, error) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")
	// call the service
	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	return client.Do(request)
}
