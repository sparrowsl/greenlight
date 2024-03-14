package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)

	router.NotFound(app.notFoundResponse)
	router.MethodNotAllowed(app.methodNotAllowedResponse)

	router.Get("/v1/healthcheck", app.checkHealth)
	router.Post("/v1/movies", app.createMovie)
	router.Get("/v1/movies/{id}", app.showMovie)
	router.Put("/v1/movies/{id}", app.updateMovie)
	router.Delete("/v1/movies/{id}", app.deleteMovie)

	return router
}
