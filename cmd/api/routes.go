package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(app.rateLimit)

	router.NotFound(app.notFoundResponse)
	router.MethodNotAllowed(app.methodNotAllowedResponse)

	router.Get("/v1/healthcheck", app.checkHealth)
	router.Post("/v1/movies", app.createMovie)
	router.Get("/v1/movies", app.listAllMovies)
	router.Get("/v1/movies/{id}", app.showMovie)
	router.Patch("/v1/movies/{id}", app.updateMovie)
	router.Delete("/v1/movies/{id}", app.deleteMovie)

	router.Post("/v1/users", app.registerUser)

	return router
}
