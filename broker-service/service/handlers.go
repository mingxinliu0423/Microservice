package main

import (
	"broker/logs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RequestPayload represents the structure of the JSON payload
// It contains an "action" field to determine the action to perform
// and an optional "auth" field for authentication data
type RequestPayload struct {
	Action string      `json:"action"`         // Action to perform
	Auth   AuthPayload `json:"auth,omitempty"` // Optional authentication data
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

// AuthPayload represents the structure of the authentication data
// It contains fields for email and password
type AuthPayload struct {
	Email    string `json:"email"`    // User's email
	Password string `json:"password"` // User's password
}

// LogPayload represents the structure of the payload for logging
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

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

// HandleSubmission is the main point of entry into the broker. It accepts a JSON
// payload and performs an action based on the value of "action" in that JSON.
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItem(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

// authenticate calls the authentication microservice and sends back the appropriate response
func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if !jsonFromService.Success {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Success = true
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

// logItem sends a log entry to the logger service
func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Success = false
	payload.Message = "logged"

	app.writeJSON(w, http.StatusAccepted, payload)

}

// sendMail sends an email using the mail service.
func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	// Marshal the message payload into JSON format
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// Define the URL of the mail service
	mailServiceURL := "http://mail-service/send"

	// Create a POST request to the mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Set the request content type to JSON
	request.Header.Set("Content-Type", "application/json")

	// Create a new HTTP client and send the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// Check if the response status code is not the expected status
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mail service"))
		return
	}

	// Send a success response back to the client
	var payload jsonResponse
	payload.Success = true
	payload.Message = "Message sent to " + msg.To

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	// Read and parse JSON request payload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		// Return an error response if JSON parsing fails
		app.errorJSON(w, err)
		return
	}

	// Establish a connection to the gRPC server
	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		// Return an error response if the connection fails
		app.errorJSON(w, err)
		return
	}
	defer conn.Close() // Ensure the connection is closed when done

	// Create a new gRPC client for the log service
	c := logs.NewLogServiceClient(conn)
	// Set a timeout context for the gRPC request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // Ensure the context is canceled when done

	// Send the log request to the gRPC server
	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		// Return an error response if the gRPC request fails
		app.errorJSON(w, err)
		return
	}

	// Create and send a success response
	var payload jsonResponse
	payload.Success = true
	payload.Message = "logged"

	app.writeJSON(w, http.StatusAccepted, payload)
}
