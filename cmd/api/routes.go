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
	router.Use(app.authenticate)

	router.NotFound(app.notFoundResponse)
	router.MethodNotAllowed(app.methodNotAllowedResponse)

	router.Get("/v1/healthcheck", app.checkHealth)

	router.Group(func(r chi.Router) {
		r.Use(app.requireActivatedUser)

		r.Post("/v1/movies", app.createMovie)
		r.Get("/v1/movies", app.listAllMovies)
		r.Get("/v1/movies/{id}", app.showMovie)
		r.Patch("/v1/movies/{id}", app.updateMovie)
		r.Delete("/v1/movies/{id}", app.deleteMovie)
	})

	router.Put("/v1/users/activated", app.activateUser)
	router.Post("/v1/users", app.registerUser)

	router.Post("/v1/tokens/authentication", app.createAuthenticationToken)

	return router
}
