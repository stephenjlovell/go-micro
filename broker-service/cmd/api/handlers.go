package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type SubmissionPayload struct {
	Error  bool        `json:"error"`
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	payload := &SubmissionPayload{}

	err := app.readJSON(w, r, payload)
	app.handleError(err, w)

	switch payload.Action {
	case "auth":
		app.Authenticate(w, payload.Auth)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) Authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json to send to the authentication service
	jsonData, _ := json.MarshalIndent(a, "", "\t")
	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	app.handleError(err, w)
	client := &http.Client{}
	response, err := client.Do(request)
	app.handleError(err, w)
	defer response.Body.Close()
	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("unable to authenticate"))
		return
	}

	serviceResponse := &jsonResponse{}
	err = json.NewDecoder(response.Body).Decode(serviceResponse)
	app.handleError(err, w)
	if serviceResponse.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	app.writeJSON(w, http.StatusAccepted, &jsonResponse{
		Error:   false,
		Message: "authenticated",
		Data:    serviceResponse.Data,
	})
}

func (app *Config) handleError(err error, w http.ResponseWriter) {
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}
