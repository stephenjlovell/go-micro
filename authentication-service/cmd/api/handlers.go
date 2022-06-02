package main

import (
	"errors"
	"net/http"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	payload := &AuthPayload{}

	err := app.readJSON(w, r, payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(payload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
	}

	valid, err := user.PasswordMatches(payload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
	}

	app.writeJSON(w, http.StatusOK, jsonResponse{
		Error:   false,
		Message: "Logged in",
		Data:    user,
	})
}
