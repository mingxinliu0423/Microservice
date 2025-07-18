package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:		[]string{"https://*", "http://*"},
		AllowedMethods: 	[]string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:		[]string{"Accpet", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: 	[]string("Link"),
		AllowCredentials:	true,
		MaxAge:				300,
	}))

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to theAuth service!"))
	})

	mux.Post("/authenticate", app.Authenticate)

	return mux
}