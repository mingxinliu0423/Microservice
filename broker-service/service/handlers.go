package main

import (
	"net/http"
)

// Broker is the HTTP handler for the main endpoint.
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	// Create a payload for the response
	payload := jsonResponse{
		Success: true,
		Message: "Request successfully routed to broker service.",
	}

	// Write the JSON response using the writeJSON method
	_ = app.writeJSON(w, http.StatusOK, payload)
}
