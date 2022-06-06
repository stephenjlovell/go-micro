package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	helpers "github.com/stephenjlovell/json-helpers"
)

const (
	loggerServiceUrl         = "http://logger-service/log"
	authenticationServiceUrl = "http://authentication-service/authenticate"
)

type LoggerPayload struct {
	Name string `json:name`
	Data string `json:data`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type SubmissionPayload struct {
	Error  bool          `json:"error"`
	Action string        `json:"action"`
	Auth   AuthPayload   `json:"auth,omitempty"`
	Log    LoggerPayload `json:"log,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := helpers.JsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	_ = helpers.WriteJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	payload := &SubmissionPayload{}
	err := helpers.ReadJSON(w, r, payload)
	if err != nil {
		_ = helpers.ErrorJSON(w, err)
		return
	}
	errCode := http.StatusUnprocessableEntity
	switch payload.Action {
	case "auth":
		err = app.authenticate(w, payload.Auth)
		errCode = http.StatusUnauthorized
	case "log":
		err = app.logItem(w, payload.Log)
	default:
		err = errors.New("unknown action")
	}
	if err != nil {
		_ = helpers.ErrorJSON(w, err, errCode)
		return
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) error {
	response, err := helpers.DoRequest("POST", authenticationServiceUrl, a)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		return errors.New("invalid credentials")
	} else if response.StatusCode != http.StatusAccepted {
		return errors.New("unable to authenticate")
	}
	// send response
	serviceResponse := &helpers.JsonResponse{}
	err = json.NewDecoder(response.Body).Decode(serviceResponse)
	if err != nil || serviceResponse.Error {
		return errors.New("unable to authenticate")
	}
	return helpers.WriteJSON(w, http.StatusAccepted, &helpers.JsonResponse{
		Error:   false,
		Message: "authenticated",
		Data:    serviceResponse.Data,
	})
}

func (app *Config) logItem(w http.ResponseWriter, l LoggerPayload) error {
	response, err := helpers.DoRequest("POST", loggerServiceUrl, l)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	// make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		return errors.New(fmt.Sprintf("request failed: %s", err))
	}
	// send response
	return helpers.WriteJSON(w, http.StatusAccepted, &helpers.JsonResponse{
		Error:   false,
		Message: "logged",
	})
}
