package main

import (
	"log-service/data"
	"net/http"
)

// JSONPayload represents the expected JSON payload for writing a log entry
type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// WriteLog handles the HTTP request to write a log entry
func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// Read JSON into var
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	// Insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Success: true,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}
